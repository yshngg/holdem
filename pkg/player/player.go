package player

import (
	"context"
	"crypto/md5"
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
	// human readable identity
	name string
	// machine readable identity
	id            string
	actionTimeout time.Duration
	holeCards     [2]*card.Card
	chips         int
	watcher       watch.Interface
	status        StatusType

	// for action handling
	once sync.Once

	// activeChan indicates
	activeChan chan []Action
	actionChan chan Action
	// availableActions map[ActionType]Action
}

func New(opts ...Option) *Player {
	p := &Player{
		status: StatusIdle,
	}
	for _, opt := range opts {
		opt(p)
	}
	if p.actionTimeout <= 0 {
		p.actionTimeout = defaultActionTimeout
	}
	if len(p.id) == 0 {
		p.id = uuid.New().String()
	}
	if len(p.name) == 0 {
		h := md5.New()
		sum := h.Sum([]byte(p.id))
		p.name = base64.StdEncoding.EncodeToString(sum)[:7]
	}
	if p.chips == 0 {
		p.chips = defaultChips
	}
	return p
}

type Option func(*Player)

func WithID(id string) Option {
	return func(p *Player) {
		p.id = id
	}
}

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

// WithStatus only used for testing
func WithStatus(status StatusType) Option {
	return func(p *Player) {
		p.status = status
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

func (p *Player) StopWatch() {
	if p.watcher != nil {
		p.watcher.Stop()
	}
}

func (p *Player) Watch() <-chan watch.Event {
	out := watch.Filter(p.watcher, func(in watch.Event) (watch.Event, bool) {
		// TODO(@yshngg): Implement filtering logic
		return in, true
	})
	return out.Watch()
}

// func (p *Player) Apply(opts ...Option) {
// 	for _, opt := range opts {
// 		opt(p)
// 	}
// }

// TODO(@yshngg): Implement BestFiveCard method
func (p *Player) BestFiveCard(communityCards ...*card.Card) [5]*card.Card {
	var bestFive [5]*card.Card

	return bestFive
}

func (p *Player) Check(ctx context.Context) error {
	return p.takeAction(ctx, Action{Type: ActionCheck})
}

func (p *Player) Fold(ctx context.Context) error {
	return p.takeAction(ctx, Action{Type: ActionFold})
}

func (p *Player) Bet(ctx context.Context, chips int) error {
	return p.takeAction(ctx, Action{Type: ActionBet, Chips: chips})
}

func (p *Player) Raise(ctx context.Context, chips int) error {
	return p.takeAction(ctx, Action{Type: ActionRaise, Chips: chips})
}

func (p *Player) Call(ctx context.Context) error {
	return p.takeAction(ctx, Action{Type: ActionCall})
}

func (p *Player) AllIn(ctx context.Context) error {
	return p.takeAction(ctx, Action{Type: ActionAllIn})
}

func (p *Player) ShowHoleCards(ctx context.Context) error {
	return p.takeAction(ctx, Action{Type: ActionShowHoleCards})
}

func (p *Player) HideHoleCards(ctx context.Context) error {
	return p.takeAction(ctx, Action{Type: ActionHideHoleCards})
}

func (p *Player) Ready() error {
	if p.status != StatusIdle {
		return fmt.Errorf("player is not idle, cannot ready")
	}
	p.activeChan = make(chan []Action, 1)
	p.actionChan = make(chan Action)
	p.status = StatusReady
	return nil
}

// Reset reset the player's status.
func (p *Player) Reset() error {
	p.activeChan = make(chan []Action, 1)
	p.actionChan = make(chan Action)
	p.status = StatusReady
	p.holeCards = [2]*card.Card{}
	return nil
}

func (p *Player) CancelReady() error {
	if p.status != StatusReady {
		return fmt.Errorf("player is not ready, cannot cancel")
	}
	p.status = StatusIdle
	p.activeChan = nil
	p.actionChan = nil
	return nil
}

// Gone is used to release resources.
func (p *Player) Gone() error {
	p.once.Do(func() {
		if p.activeChan != nil {
			close(p.activeChan)
		}
		if p.actionChan != nil {
			close(p.actionChan)
		}
	})
	return nil
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) ID() string {
	return p.id
}

func (p *Player) HoleCards() [2]*card.Card {
	return p.holeCards
}

func (p *Player) SetHoleCards(cards [2]*card.Card) error {
	p.holeCards = cards
	if p.status != StatusReady {
		return fmt.Errorf("player is not ready, cannot wait to act")
	}
	return nil
}

func (p *Player) Chips() int {
	return p.chips
}

func (p *Player) Status() StatusType {
	return p.status
}

func (p *Player) Active() <-chan []Action {
	return p.activeChan
}

func (p *Player) takeAction(ctx context.Context, action Action) error {
	ctx, cancel := context.WithTimeoutCause(ctx, p.actionTimeout, fmt.Errorf("action timeout"))
	defer cancel()

	// if p.availableActions == nil {
	// 	return fmt.Errorf("not have available actions")
	// }
	// require, ok := p.availableActions[action.Type]
	// if !ok {
	// 	return fmt.Errorf("not available action: %v, available action: %v", action, p.availableActions)
	// }

	switch action.Type {
	case ActionCheck, ActionFold:
	case ActionBet, ActionRaise:
		// Equivalent to: !(require.Chips <= action.Chips <= p.chips)
		if require.Chips > action.Chips || action.Chips > p.chips {
			return fmt.Errorf("not enough chips: %d, can not take the action: %v", p.chips, action)
		}
		p.chips -= action.Chips
	case ActionCall:
		// Equivalent to: !(require.Chips <= p.chips)
		if require.Chips > p.chips {
			return fmt.Errorf("not enough chips: %d, can not take the action: %v", p.chips, action)
		}
		action.Chips = require.Chips
		p.chips -= require.Chips
	case ActionAllIn:
		if p.chips <= 0 {
			return fmt.Errorf("not enough chips: %d", p.chips)
		}
		action.Chips = p.chips
		p.chips = 0
	default:
		return fmt.Errorf("invalid action type: %v", action.Type)
	}

	select {
	case p.actionChan <- action:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("failed to take action: %v, reason: %w", action, ctx.Err())
	}
}

// WaitForAction wait for the player to take action
// action[0] is the default action.
// TODO(@yshngg): check correct of function call params
func (p *Player) WaitForAction(ctx context.Context, available []Action) (*Action, error) {
	ctx, cancel := context.WithTimeoutCause(ctx, p.actionTimeout, fmt.Errorf("action timeout"))
	defer cancel()

	if p.status != StatusWaiting {
		return nil, fmt.Errorf("player %s [id: %s] does not wait to act, status: %s", p.name, p.id, p.status)
	}
	// p.status = StatusTakingAction

	// deduplicate actions
	availableMap := make(map[ActionType]Action)
	for _, action := range available {
		availableMap[action.Type] = action
	}

	// drain active and action channels
Drain:
	for {
		select {
		case <-p.activeChan:
		case <-p.actionChan:
		default:
			if len(p.activeChan) == 0 && len(p.actionChan) == 0 {
				break Drain
			}
		}
	}

	p.activeChan <- available

	// p.availableActions = availableActions
	// defer func() {
	// p.availableActions = nil
	// }()

	action := available[0]
	select {
	case <-ctx.Done():
		<-p.activeChan

	case action = <-p.actionChan:
	}

	// action concluded
	// if p.status != StatusTakingAction {
	// 	return nil, fmt.Errorf("player is not take action")
	// }
	p.status = action.Type.ToStatus()
	return &action, nil
}
