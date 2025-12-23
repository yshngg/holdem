package dealer

import (
	"testing"

	"github.com/yshngg/holdem/pkg/deck"
)

func TestDeal(t *testing.T) {
	_deck := deck.New()
	_dealer := New(_deck)
	_dealer.Shuffle()
	_card := _dealer.Deal()
	t.Logf("Card: %v", _card)
}
