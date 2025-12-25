package deck

import (
	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/rank"
	"github.com/yshngg/holdem/pkg/suit"
)

type Deck struct {
	cards []card.Card
}

func New() *Deck {
	d := &Deck{}
	for s := suit.Clubs; s <= suit.Diamonds; s++ {
		c := card.New(rank.Ace, s)
		d.cards = append(d.cards, c)
	}
	for r := rank.Two; r <= rank.King; r++ {
		for s := suit.Clubs; s <= suit.Diamonds; s++ {
			c := card.New(r, s)
			d.cards = append(d.cards, c)
		}
	}
	return d
}

func (d Deck) Len() int {
	return len(d.cards)
}

func (d *Deck) Pop() *card.Card {
	remain := d.Len()
	if remain == 0 {
		return nil
	}
	card := d.cards[0]
	if remain > 1 {
		d.cards = d.cards[1:]
	} else {
		d.cards = nil
	}
	return &card
}

func (d *Deck) Swap(i, j int) {
	if i < 0 || i >= d.Len() || j < 0 || j >= d.Len() {
		return
	}
	d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
}
