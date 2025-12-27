package dealer

import (
	"testing"

	"github.com/yshngg/holdem/pkg/deck"
)

func TestDeal(t *testing.T) {
	_deck := deck.New()
	dealer := New(WithDeck(_deck), WithShuffle())
	dealer.Shuffle()
	_card := dealer.Deal()
	t.Logf("Card: %v", _card)
}
