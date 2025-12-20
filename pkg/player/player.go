package player

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yshngg/holdem/pkg/card"
)

const (
	MinBet = 2
)

type Player struct {
	name      string
	id        uuid.UUID
	holeCards [2]*card.Card
	chip      int
}

// TODO(@yshngg): Implement BestFivePockerHand method
func (p *Player) BestFivePockerHand(communityCards [5]*card.Card) [5]*card.Card {
	var bestFive [5]*card.Card

	return bestFive
}

func (p *Player) PostSmallBlind() error {
	bet := MinBet / 2
	if p.chip < bet {
		return fmt.Errorf("not enough chips, have: %d, want: %d", p.chip, bet)
	}
	p.chip -= bet
	return nil
}

func (p *Player) PostBigBlind() error {
	bet := MinBet
	if p.chip < bet {
		return fmt.Errorf("not enough chips, have: %d, want: %d", p.chip, bet)
	}
	p.chip -= bet
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

func (p *Player) Chip() int {
	return p.chip
}

func (p *Player) Bet(number int) {}
