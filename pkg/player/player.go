package player

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/watch"
)

const (
	defaultChips         = 100
	defaultActionTimeout = 5 * time.Second
)

type ErrNotEnoughChips struct {
	Have int
	Want int
}

func (e ErrNotEnoughChips) Error() string {
	return fmt.Sprintf("not enough chips, have: %d, want: %d", e.Have, e.Want)
}

type Player struct {
	name          string
	id            uuid.UUID
	actionTimeout time.Duration
	holeCards     [2]*card.Card
	chips         int
	watcher       watch.Interface
	status        Status

	// for action handling
	done             sync.Once
	active           chan bool
	actionChan       chan Action
	availableActions map[ActionType]Action
}

func New(opts ...Option) *Player {
	id := uuid.New()
	p := &Player{
		id:            id,
		actionTimeout: defaultActionTimeout,
		status:        StatusUnknown,
		active:        make(chan bool, 1),
		actionChan:    make(chan Action, 1),
	}
	p.Apply(opts...)
	if len(p.name) == 0 {
		p.name = base64.StdEncoding.EncodeToString([]byte(id.String()))[:7]
	}
	if p.chips == 0 {
		p.chips = defaultChips
	}
	return p
}

type Option func(*Player)

func WithName(name string) Option {
	return func(p *Player) {
		p.name = name
	}
}

func WithChips(chips int) Option {
	return func(p *Player) {
		p.chips = chips
	}
}

func WithActionTimeout(timeout time.Duration) Option {
	return func(p *Player) {
		p.actionTimeout = timeout
	}
}

func WithWatcher(watcher watch.Interface) Option {
	return func(p *Player) {
		p.watcher = watcher
	}
}

func (p *Player) Watch() <-chan watch.Event {
	return p.watcher.ResultChan()
}

func (p *Player) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(p)
	}
}

// TODO(@yshngg): Implement BestFivePockerHand method
func (p *Player) BestFivePockerHand(communityCards [5]*card.Card) [5]*card.Card {
	var bestFive [5]*card.Card

	return bestFive
}

func (p *Player) Bet(chips int) error {
	if p.chips < chips {
		return ErrNotEnoughChips{p.chips, chips}
	}
	p.chips -= chips
	return nil
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) ID() uuid.UUID {
	return p.id
}

func (p *Player) HoleCards() [2]*card.Card {
	return p.holeCards
}

func (p *Player) SetHoleCards(cards [2]*card.Card) {
	p.holeCards = cards
}

func (p *Player) Chips() int {
	return p.chips
}

func (p *Player) Status() Status {
	return p.status
}

func (p *Player) Active() <-chan bool {
	return p.active
}

func (p *Player) TakeAction(ctx context.Context, action Action) error {
	if p.availableActions == nil {
		return fmt.Errorf("not have available actions")
	}
	require, ok := p.availableActions[action.Type]
	if !ok {
		return fmt.Errorf("not available action: %v, available action: %v", action, p.availableActions)
	}

	switch action.Type {
	case ActionCheck, ActionFold:
		p.actionChan <- action
	case ActionBet, ActionCall, ActionRaise:
		// Equivalent to: !(require.Chips <= action.Chips <= p.chips)
		if require.Chips > action.Chips || action.Chips > p.chips {
			return fmt.Errorf("not enough chips: %d, can not take the action: %v", p.chips, action)
		}
		p.chips -= action.Chips
		p.actionChan <- action
		return nil
	case ActionAllIn:
		if p.chips <= 0 {
			return fmt.Errorf("not enough chips: %d", p.chips)
		}
		action.Chips = p.chips
		p.chips = 0
		p.actionChan <- action
		return nil
	}
	return nil
}

func (p *Player) WaitForAction(ctx context.Context, availableActions map[ActionType]Action) Action {
	ctx, cancel := context.WithTimeoutCause(ctx, p.actionTimeout, fmt.Errorf("action timeout"))
	defer cancel()

	// drain active and action channels
Drain:
	for {
		select {
		case <-p.active:
		case <-p.actionChan:
		default:
			if len(p.active) == 0 && len(p.actionChan) == 0 {
				break Drain
			}
		}
	}

	p.active <- true
	p.availableActions = availableActions
	defer func() {
		p.availableActions = nil
	}()

	action := Action{Type: ActionFold}
	select {
	case <-ctx.Done():
		<-p.active
		if _, ok := p.availableActions[ActionCheck]; ok {
			action = Action{Type: ActionCheck}
		}

	case action = <-p.actionChan:
	}

	p.status = action.Type.IntoStatus()
	return action
}

func (p *Player) Done() {
	if p.watcher != nil {
		p.watcher.Stop()
	}
	p.done.Do(func() {
		close(p.active)
		close(p.actionChan)
	})
}
