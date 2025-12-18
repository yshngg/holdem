package player

import (
	"github.com/google/uuid"
	"github.com/yshngg/holdem/pkg/card"
)

const (
	MinimumBet          = 2
	MaximumPalyerNumber = 10
)

type Player struct {
	name      string
	id        uuid.UUID
	holeCards [2]*card.Card
	token     int
}

// TODO(@yshngg): Implement BestFivePockerHand method
func (p *Player) BestFivePockerHand(communityCards [5]*card.Card) [5]*card.Card {
	var bestFive [5]*card.Card

	return bestFive
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

func (p *Player) Token() int {
	return p.token
}

func (p *Player) Bet(number int) {}
