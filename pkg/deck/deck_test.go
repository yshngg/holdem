package deck

import (
	"testing"

	"github.com/yshngg/holdem/pkg/rank"
	"github.com/yshngg/holdem/pkg/suit"
)

func TestLen(t *testing.T) {
	_deck := New()
	if _deck.Len() != 52 {
		t.Errorf("deck length: %d, want 52", _deck.Len())
	}
}

func TestList(t *testing.T) {
	_deck := New()
	if len(_deck.List()) != 52 {
		t.Errorf("deck length: %d, want 52, ", len(_deck.List()))
	}
	firstCard := _deck.List()[0]
	if firstCard.Rank() != rank.Ace {
		t.Errorf("first card rank: %s, want %s", firstCard.Rank(), rank.Ace)
	}
	if firstCard.Suit() != suit.Clubs {
		t.Errorf("first card suit: %s, want %s", firstCard.Suit(), suit.Clubs)
	}

	lastCard := _deck.List()[_deck.Len()-1]
	if lastCard.Rank() != rank.King {
		t.Errorf("last card rank: %s, want %s", lastCard.Rank(), rank.King)
	}
	if lastCard.Suit() != suit.Diamonds {
		t.Errorf("last card suit: %s, want %s", lastCard.Suit(), suit.Diamonds)
	}
}

func TestPop(t *testing.T) {
	_deck := New()
	remainder := 52
	for s := suit.Clubs; s <= suit.Diamonds; s++ {
		_card := _deck.Pop()
		remainder--
		if remainder != _deck.Len() {
			t.Errorf("deck length: %d, want %d", _deck.Len(), remainder)
		}
		if _card == nil {
			t.Errorf("pop card: nil, want not nil")
		}
		if _card.Rank() != rank.Ace || _card.Suit() != s {
			t.Errorf("card mismatch: %s, want %s of %s", _card, rank.Ace, s)
		}
	}
	for r := rank.Two; r <= rank.King; r++ {
		for s := suit.Clubs; s <= suit.Diamonds; s++ {
			// _card := _deck.cards[i*4+j]
			_card := _deck.Pop()
			remainder--
			if remainder != _deck.Len() {
				t.Errorf("deck length: %d, want %d", _deck.Len(), remainder)
			}
			if _card == nil {
				t.Errorf("pop card: nil, want not nil")
			}
			if _card.Rank() != r || _card.Suit() != s {
				t.Errorf("card mismatch: %s, want %s of %s", _card, r, s)
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
		t.Errorf("swap cards: %s, want %s", _deck.cards[0], last)
	}
	if _deck.cards[_deck.Len()-1] != first {
		t.Errorf("swap cards: %s, want %s", _deck.cards[_deck.Len()-1], first)
	}
}
