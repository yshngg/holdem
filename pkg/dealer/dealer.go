package dealer

import (
	"math/rand"

	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/deck"
)

type Dealer struct {
	deck *deck.Deck
}

func New(deck *deck.Deck) *Dealer {
	return &Dealer{deck: deck}
}

func (d *Dealer) deal() *card.Card {
	return d.deck.Pop()
}

func (d *Dealer) Shuffle() {
	rand.Shuffle(d.deck.Len(), d.deck.Swap)
}

func (d *Dealer) DealHoleCards(playerCont int) [][2]*card.Card {
	holeCards := make([][2]*card.Card, playerCont)
	for i := range 2 {
		for j := range playerCont {
			card := d.deal()
			if card == nil {
				panic("deck is empty")
			}
			holeCards[i][j] = card
		}
	}
	return holeCards
}

func (d *Dealer) DealFlop() [3]*card.Card {
	var flop [3]*card.Card
	for i := range 3 {
		card := d.deal()
		if card == nil {
			panic("deck is empty")
		}
		flop[i] = card
	}
	return flop
}

func (d *Dealer) DealTurn() *card.Card {
	card := d.deal()
	if card == nil {
		panic("deck is empty")
	}
	return card
}

func (d *Dealer) DealRiver() *card.Card {
	card := d.deal()
	if card == nil {
		panic("deck is empty")
	}
	return card
}

func (d *Dealer) Deal() *card.Card {
	return d.deal()
}
