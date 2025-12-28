package hand

import (
	"testing"

	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/dealer"
	"github.com/yshngg/holdem/pkg/rank"
	"github.com/yshngg/holdem/pkg/suit"
)

func TestValue(t *testing.T) {
	testCases := []struct {
		name  string
		cards []card.Card
		want  Hand
		err   error
	}{
		{
			name: "HighCard",
			cards: []card.Card{
				card.New(rank.Ten, suit.Clubs),
				card.New(rank.Four, suit.Hearts),
				card.New(rank.Seven, suit.Diamonds),
				card.New(rank.King, suit.Clubs),
				card.New(rank.Two, suit.Spades),
			},
			want: HighCard,
			err:  nil,
		},
		{
			name: "Pair",
			cards: []card.Card{
				card.New(rank.King, suit.Clubs),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Seven, suit.Diamonds),
				card.New(rank.Two, suit.Clubs),
				card.New(rank.Five, suit.Spades),
			},
			want: Pair,
			err:  nil,
		},
		{
			name: "TwoPairs",
			cards: []card.Card{
				card.New(rank.King, suit.Clubs),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Seven, suit.Diamonds),
				card.New(rank.Seven, suit.Clubs),
				card.New(rank.Five, suit.Spades),
			},
			want: TwoPairs,
			err:  nil,
		},
		{
			name: "ThreeOfAKind",
			cards: []card.Card{
				card.New(rank.King, suit.Clubs),
				card.New(rank.King, suit.Hearts),
				card.New(rank.King, suit.Diamonds),
				card.New(rank.Seven, suit.Clubs),
				card.New(rank.Five, suit.Spades),
			},
			want: ThreeOfAKind,
			err:  nil,
		},
		{
			name: "Straight",
			cards: []card.Card{
				card.New(rank.Three, suit.Clubs),
				card.New(rank.Four, suit.Hearts),
				card.New(rank.Five, suit.Diamonds),
				card.New(rank.Six, suit.Clubs),
				card.New(rank.Seven, suit.Spades),
			},
			want: Straight,
			err:  nil,
		},
		{
			name: "Flush",
			cards: []card.Card{
				card.New(rank.King, suit.Clubs),
				card.New(rank.Queen, suit.Clubs),
				card.New(rank.Nine, suit.Clubs),
				card.New(rank.Eight, suit.Clubs),
				card.New(rank.Two, suit.Clubs),
			},
			want: Flush,
			err:  nil,
		},
		{
			name: "FullHouse",
			cards: []card.Card{
				card.New(rank.King, suit.Clubs),
				card.New(rank.King, suit.Hearts),
				card.New(rank.King, suit.Diamonds),
				card.New(rank.Seven, suit.Clubs),
				card.New(rank.Seven, suit.Spades),
			},
			want: FullHouse,
			err:  nil,
		},
		{
			name: "FourOfAKind",
			cards: []card.Card{
				card.New(rank.Six, suit.Spades),
				card.New(rank.Six, suit.Diamonds),
				card.New(rank.Six, suit.Hearts),
				card.New(rank.Six, suit.Clubs),
				card.New(rank.King, suit.Spades),
			},
			want: FourOfAKind,
			err:  nil,
		},
		{
			name: "StraightFlush",
			cards: []card.Card{
				card.New(rank.Two, suit.Spades),
				card.New(rank.Three, suit.Spades),
				card.New(rank.Four, suit.Spades),
				card.New(rank.Five, suit.Spades),
				card.New(rank.Six, suit.Spades),
			},
			want: StraightFlush,
			err:  nil,
		},
		{
			name: "RoyalFlush",
			cards: []card.Card{
				card.New(rank.Ten, suit.Hearts),
				card.New(rank.Jack, suit.Hearts),
				card.New(rank.Queen, suit.Hearts),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Ace, suit.Hearts),
			},
			want: RoyalFlush,
			err:  nil,
		},
		{
			name: "MinimalStraightFlush",
			cards: []card.Card{
				card.New(rank.Ace, suit.Hearts),
				card.New(rank.Two, suit.Hearts),
				card.New(rank.Three, suit.Hearts),
				card.New(rank.Four, suit.Hearts),
				card.New(rank.Five, suit.Hearts),
			},
			want: StraightFlush,
			err:  nil,
		},
		{
			name: "FourHands",
			cards: []card.Card{
				card.New(rank.Jack, suit.Hearts),
				card.New(rank.Queen, suit.Hearts),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Ace, suit.Hearts),
			},
			want: Unknown,
			err:  ErrInvalidHandSize{},
		},
		{
			name: "SixHands",
			cards: []card.Card{
				card.New(rank.Nine, suit.Hearts),
				card.New(rank.Ten, suit.Hearts),
				card.New(rank.Jack, suit.Hearts),
				card.New(rank.Queen, suit.Hearts),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Ace, suit.Hearts),
			},
			want: Unknown,
			err:  ErrInvalidHandSize{},
		},
		{
			name: "ExistSameCards",
			cards: []card.Card{
				card.New(rank.Jack, suit.Hearts),
				card.New(rank.Queen, suit.Hearts),
				card.New(rank.King, suit.Hearts),
				card.New(rank.Ace, suit.Hearts),
				card.New(rank.Ace, suit.Hearts),
			},
			want: Unknown,
			err:  ErrExistSameCards{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := Value(tc.cards)
			if err != tc.err {
				t.Errorf("Value(%v).err = %v, want %v", tc.cards, err, tc.err)
			}
			if value != tc.want {
				t.Errorf("Value(%v).value = %v, want %v", tc.cards, value, tc.want)
			}
		})
	}
}

func BenchmarkValue(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		d := dealer.New()
		d.Shuffle()
		cards := []card.Card{
			*d.Deal(),
			*d.Deal(),
			*d.Deal(),
			*d.Deal(),
			*d.Deal(),
		}
		b.StartTimer()
		_, err := Value(cards)
		if err != nil {
			b.Errorf("Value(%v).err = %v", cards, err)
		}
	}
}
