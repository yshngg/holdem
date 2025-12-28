package dealer

import (
	"math/rand"

	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/deck"
)

type Dealer struct {
	deck *deck.Deck
}

// New creates a new dealer with the given options.
func New(opts ...Option) *Dealer {
	d := &Dealer{}
	for _, opt := range opts {
		opt(d)
	}
	if d.deck == nil {
		d.deck = deck.New()
	}
	return d
}

type Option func(*Dealer)

func WithShuffle() Option {
	return func(d *Dealer) {
		if d == nil || d.deck == nil {
			return
		}
		d.Shuffle()
	}
}

func WithDeck(_deck *deck.Deck) Option {
	return func(d *Dealer) {
		d.deck = _deck
	}
}

func (d *Dealer) Reset() {
	d.deck = deck.New()
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

func (d *Dealer) DealFlopCards() [3]*card.Card {
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

func (d *Dealer) DealTurnCard() *card.Card {
	card := d.deal()
	if card == nil {
		panic("deck is empty")
	}
	return card
}

func (d *Dealer) DealRiverCard() *card.Card {
	card := d.deal()
	if card == nil {
		panic("deck is empty")
	}
	return card
}

func (d *Dealer) BurnCard() *card.Card {
	card := d.deal()
	if card == nil {
		panic("deck is empty")
	}
	return card
}

func (d *Dealer) Deal() *card.Card {
	return d.deal()
}
