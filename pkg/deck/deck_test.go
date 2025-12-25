package deck

import (
	"testing"

	"github.com/yshngg/holdem/pkg/rank"
	"github.com/yshngg/holdem/pkg/suit"
)

func TestLen(t *testing.T) {
	_deck := New()
	if _deck.Len() != 52 {
		t.Errorf("deck length: expected 52, got %d", _deck.Len())
	}
}

func TestPop(t *testing.T) {
	_deck := New()
	remainder := 52
	for s := suit.Clubs; s <= suit.Diamonds; s++ {
		_card := _deck.Pop()
		remainder--
		if remainder != _deck.Len() {
			t.Errorf("deck length: expected %d, got %d", remainder, _deck.Len())
		}
		if _card == nil {
			t.Errorf("pop card: expected not nil, got nil")
		}
		if _card.Rank() != rank.Ace || _card.Suit() != s {
			t.Errorf("card mismatch: expected %s of %s, got %s", rank.Ace, s, _card)
		}
	}
	for r := rank.Two; r <= rank.King; r++ {
		for s := suit.Clubs; s <= suit.Diamonds; s++ {
			// _card := _deck.cards[i*4+j]
			_card := _deck.Pop()
			remainder--
			if remainder != _deck.Len() {
				t.Errorf("deck length: expected %d, got %d", remainder, _deck.Len())
			}
			if _card == nil {
				t.Errorf("pop card: expected not nil, got nil")
			}
			if _card.Rank() != r || _card.Suit() != s {
				t.Errorf("card mismatch: expected %s of %s, got %s", r, s, _card)
			}
		}
	}
}

func TestSwap(t *testing.T) {
	_deck := New()
	first := _deck.cards[0]
	last := _deck.cards[_deck.Len()-1]
	_deck.Swap(0, _deck.Len()-1)
	if _deck.cards[0] != last {
		t.Errorf("swap cards: expected %s, got %s", last, _deck.cards[0])
	}
	if _deck.cards[_deck.Len()-1] != first {
		t.Errorf("swap cards: expected %s, got %s", first, _deck.cards[_deck.Len()-1])
	}
}
