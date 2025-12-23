package card

import (
	"fmt"

	"github.com/yshngg/holdem/pkg/rank"
	"github.com/yshngg/holdem/pkg/suit"
)

type Card struct {
	rank rank.Rank
	suit suit.Suit
}

func New(rank rank.Rank, suit suit.Suit) Card {
	return Card{
		rank: rank,
		suit: suit,
	}
}

func (c Card) Rank() rank.Rank {
	return c.rank
}

func (c Card) Suit() suit.Suit {
	return c.suit
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.rank.String(), c.suit.String())
}
