package hand

import (
	"testing"

	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/dealer"
	"github.com/yshngg/holdem/pkg/deck"
	"github.com/yshngg/holdem/pkg/rank"
	"github.com/yshngg/holdem/pkg/suit"
)

func TestValue(t *testing.T) {
	testCases := []struct {
		cards    []card.Card
		expected HandValue
		err      error
	}{
		{
			cards: []card.Card{
				card.New(rank.Three, suit.Hearts),
				card.New(rank.Four, suit.Spades),
				card.New(rank.Five, suit.Hearts),
				card.New(rank.Six, suit.Diamonds),
				card.New(rank.Seven, suit.Clubs),
			},
			expected: Straight,
			err:      nil,
		},
		{
			cards: []card.Card{
				card.New(rank.Ace, suit.Hearts),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Queen, suit.Hearts),
				card.New(rank.Jack, suit.Hearts),
				card.New(rank.Ten, suit.Hearts),
			},
			expected: RoyalFlush,
			err:      nil,
		},
		{
			cards: []card.Card{
				card.New(rank.Ace, suit.Hearts),
				card.New(rank.Two, suit.Hearts),
				card.New(rank.Three, suit.Hearts),
				card.New(rank.Four, suit.Hearts),
				card.New(rank.Five, suit.Hearts),
			},
			expected: StraightFlush,
			err:      nil,
		},
	}

	for _, tc := range testCases {
		value, err := Value(tc.cards)
		if err != tc.err {
			t.Errorf("Value(%v); err: got %v, want %v", tc.cards, err, tc.err)
		}
		if value != tc.expected {
			t.Errorf("Value(%v); hand value: got %v, want %v", tc.cards, value, tc.expected)
		}
	}
}

func BenchmarkValue(b *testing.B) {
	dck := deck.New()
	d := dealer.New(dck)
	for b.Loop() {
		d.Shuffle()
		cards := []card.Card{
			*d.Deal(),
			*d.Deal(),
			*d.Deal(),
			*d.Deal(),
			*d.Deal(),
		}
		_, err := Value(cards)
		if err != nil {
			b.Fatal(err)
		}
	}
}
