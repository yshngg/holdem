package player

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/watch"
)

const defaultChips = 100

type ErrNotEnoughChips struct {
	Have int
	Want int
}

func (e ErrNotEnoughChips) Error() string {
	return fmt.Sprintf("not enough chips, have: %d, want: %d", e.Have, e.Want)
}

type Player struct {
	name        string
	id          uuid.UUID
	holeCards   [2]*card.Card
	chips       int
	broadcaster watch.Broadcaster
	status      Status

	// for action handling
	activeChan       chan struct{}
	actionChan       chan Action
	availableActions []Action
}

func New(opts ...Option) *Player {
	id := uuid.New()
	p := &Player{
		id:         id,
		status:     StatusReady,
		activeChan: make(chan struct{}, 1),
		actionChan: make(chan Action, 1),
	}
	for _, opt := range opts {
		opt(p)
	}
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

func (p *Player) Watch(w watch.Interface) {
	defer w.Stop()
	for e := range w.ResultChan() {
		switch e.Type {
		case Check:
		case Fold:
		case Bet:
		case Call:
		case Raise:
		case AllIn:
		}
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

func (p *Player) TakeAction(ctx context.Context, action Action) error {
	switch action.Type {
	case ActionCheck:
	case ActionFold:

	}
}

func (p *Player) Active() chan struct{} {
	return p.activeChan
}

func (p *Player) WaitForAction(availableActions []Action) Action {
	p.status = StatusActive
	p.activeChan <- struct{}{}
	p.availableActions = availableActions
	action := <-p.actionChan
	p.status = action.Type.IntoStatus()
	return action
}
