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
				card.New(rank.Ten, suit.Clubs),
				card.New(rank.Four, suit.Hearts),
				card.New(rank.Seven, suit.Diamonds),
				card.New(rank.King, suit.Clubs),
				card.New(rank.Two, suit.Spades),
			},
			expected: HighCard,
			err:      nil,
		},
		{
			cards: []card.Card{
				card.New(rank.King, suit.Clubs),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Seven, suit.Diamonds),
				card.New(rank.Two, suit.Clubs),
				card.New(rank.Five, suit.Spades),
			},
			expected: Pair,
			err:      nil,
		},
		{
			cards: []card.Card{
				card.New(rank.King, suit.Clubs),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Seven, suit.Diamonds),
				card.New(rank.Seven, suit.Clubs),
				card.New(rank.Five, suit.Spades),
			},
			expected: TwoPairs,
			err:      nil,
		},
		{
			cards: []card.Card{
				card.New(rank.King, suit.Clubs),
				card.New(rank.King, suit.Hearts),
				card.New(rank.King, suit.Diamonds),
				card.New(rank.Seven, suit.Clubs),
				card.New(rank.Five, suit.Spades),
			},
			expected: ThreeOfAKind,
			err:      nil,
		},
		{
			cards: []card.Card{
				card.New(rank.Three, suit.Clubs),
				card.New(rank.Four, suit.Hearts),
				card.New(rank.Five, suit.Diamonds),
				card.New(rank.Six, suit.Clubs),
				card.New(rank.Seven, suit.Spades),
			},
			expected: Straight,
			err:      nil,
		},
		{
			cards: []card.Card{
				card.New(rank.King, suit.Clubs),
				card.New(rank.Queen, suit.Clubs),
				card.New(rank.Nine, suit.Clubs),
				card.New(rank.Eight, suit.Clubs),
				card.New(rank.Two, suit.Clubs),
			},
			expected: Flush,
			err:      nil,
		},
		{
			cards: []card.Card{
				card.New(rank.King, suit.Clubs),
				card.New(rank.King, suit.Hearts),
				card.New(rank.King, suit.Diamonds),
				card.New(rank.Seven, suit.Clubs),
				card.New(rank.Seven, suit.Spades),
			},
			expected: FullHouse,
			err:      nil,
		},
		{
			cards: []card.Card{
				card.New(rank.Six, suit.Spades),
				card.New(rank.Six, suit.Diamonds),
				card.New(rank.Six, suit.Hearts),
				card.New(rank.Six, suit.Clubs),
				card.New(rank.King, suit.Spades),
			},
			expected: FourOfAKind,
			err:      nil,
		},
		{
			cards: []card.Card{
				card.New(rank.Two, suit.Spades),
				card.New(rank.Three, suit.Spades),
				card.New(rank.Four, suit.Spades),
				card.New(rank.Five, suit.Spades),
				card.New(rank.Six, suit.Spades),
			},
			expected: StraightFlush,
			err:      nil,
		},
		{
			cards: []card.Card{
				card.New(rank.Two, suit.Spades),
				card.New(rank.Three, suit.Spades),
				card.New(rank.Four, suit.Spades),
				card.New(rank.Five, suit.Spades),
				card.New(rank.Six, suit.Spades),
			},
			expected: StraightFlush,
			err:      nil,
		},
		{
			cards: []card.Card{
				card.New(rank.Ten, suit.Hearts),
				card.New(rank.Jack, suit.Hearts),
				card.New(rank.Queen, suit.Hearts),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Ace, suit.Hearts),
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
		{
			cards: []card.Card{
				card.New(rank.Jack, suit.Hearts),
				card.New(rank.Queen, suit.Hearts),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Ace, suit.Hearts),
			},
			expected: Unknown,
			err:      ErrInvalidHandSize{},
		},
		{
			cards: []card.Card{
				card.New(rank.Nine, suit.Hearts),
				card.New(rank.Ten, suit.Hearts),
				card.New(rank.Jack, suit.Hearts),
				card.New(rank.Queen, suit.Hearts),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Ace, suit.Hearts),
			},
			expected: Unknown,
			err:      ErrInvalidHandSize{},
		},
		{
			cards: []card.Card{
				card.New(rank.Jack, suit.Hearts),
				card.New(rank.Queen, suit.Hearts),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Ace, suit.Hearts),
				card.New(rank.Ace, suit.Hearts),
			},
			expected: Unknown,
			err:      ErrExistSameCards{},
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
