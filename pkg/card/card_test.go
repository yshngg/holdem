package card

import (
	"testing"

	"github.com/yshngg/holdem/pkg/rank"
	"github.com/yshngg/holdem/pkg/suit"
)

func TestCard(t *testing.T) {
	testCases := []struct {
		name  string
		_rank rank.Rank
		_suit suit.Suit
		want  Card
	}{
		{
			name:  "Ace of Spades",
			_rank: rank.Ace,
			_suit: suit.Spades,
			want:  Card{rank.Ace, suit.Spades},
		},
		{
			name:  "Two of Hearts",
			_rank: rank.Two,
			_suit: suit.Hearts,
			want:  Card{rank.Two, suit.Hearts},
		},
		{
			name:  "Three of Diamonds",
			_rank: rank.Three,
			_suit: suit.Diamonds,
			want:  Card{rank.Three, suit.Diamonds},
		},
		{
			name:  "Four of Clubs",
			_rank: rank.Four,
			_suit: suit.Clubs,
			want:  Card{rank.Four, suit.Clubs},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Card{tc._rank, tc._suit}
			if got != tc.want {
				t.Errorf("Card(%v, %v) = %v, want %v", tc._rank, tc._suit, got, tc.want)
			}
		})
	}
}
