package hand

import (
	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/rank"
)

type HandValue int

const (
	Unknown       HandValue = iota
	HighCard                // Simple value of the card. Lowest: 2 – Highest: Ace (King in the example)
	Pair                    // Two cards with the same value
	TwoPairs                // Two times two cards with the same value
	ThreeOfAKind            // Three cards with the same value
	Straight                // Sequence of 5 cards in increasing value (Ace can precede 2 or follow up King, but not both), not of the same suit
	Flush                   // 5 cards of the same suit, not in sequential order
	FullHouse               // Combination of three of a kind and a pair
	FourOfAKind             // Four cards of the same value
	StraightFlush           // Straight of the same suit
	RoyalFlush              // Highest straight of the same suit
)

func (hv HandValue) String() string {
	switch hv {
	case HighCard:
		return "High Card"
	case Pair:
		return "Pair"
	case TwoPairs:
		return "Two Pairs"
	case ThreeOfAKind:
		return "Three of a Kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	case RoyalFlush:
		return "Royal Flush"
	default:
		return "Unknown"
	}
}

type ErrInvalidHandSize struct{}

func (e ErrInvalidHandSize) Error() string {
	return "invalid hand size"
}

type ErrUnknownHandValue struct{}

func (e ErrUnknownHandValue) Error() string {
	return "unknown hand value"
}

type ErrExistSameCards struct{}

func (e ErrExistSameCards) Error() string {
	return "exist same cards"
}

func existSameCards(cards []card.Card) bool {
	m := make(map[card.Card]int)
	for _, c := range cards {
		if _, ok := m[c]; ok {
			return true
		}
		m[c]++
	}
	return false
}

func Value(c []card.Card) (HandValue, error) {
	cards := make([]card.Card, len(c))
	copy(cards, c)
	if len(cards) != 5 {
		return Unknown, ErrInvalidHandSize{}
	}
	if existSameCards(cards) {
		return Unknown, ErrExistSameCards{}
	}

	m := make(map[rank.Rank]int)
	for _, c := range cards {
		m[c.Rank()]++
	}

	if len(m) < 2 {
		return Unknown, nil
	}

	product := 1
	for _, v := range m {
		product *= v
	}

	switch product {
	case 6: // 3 * 2 Full house
		return FullHouse, nil
	case 4: // 4 * 1, 2 * 2 * 1
		if len(m) == 2 { // 4 * 1 Four of a kind
			return FourOfAKind, nil
		}
		if len(m) == 3 { // 2 * 2 * 1 Two pairs
			return TwoPairs, nil
		}
	case 3: // 3 * 1 * 1 Three of a kind
		return ThreeOfAKind, nil
	case 2: // 2 * 1 * 1 * 1 Pair
		return Pair, nil
	case 1: // 1 * 1 * 1 * 1 * 1 High card / Straight / Flush / Straight flush / Royal flush
		const (
			flush    = 0b01 // 2
			straight = 0b10 // 1
		)
		flag := flush

		suit := cards[0].Suit()
		for _, c := range cards[1:] {
			if c.Suit() != suit {
				flag -= flush
				break
			}
		}

		minRank, maxRank := cards[0].Rank(), cards[0].Rank()
		for _, card := range cards[1:] {
			iRank := card.Rank()
			if iRank > maxRank {
				maxRank = iRank
			} else if iRank < minRank {
				minRank = iRank
			}
		}
		// or A 2 3 4 5
		_, okAce := m[rank.Ace]
		_, okTwo := m[rank.Two]
		_, okThree := m[rank.Three]
		_, okFour := m[rank.Four]
		_, okFive := m[rank.Five]
		if maxRank-minRank == 4 || okAce && okTwo && okThree && okFour && okFive {
			flag += straight
		}

		switch flag {
		case straight:
			return Straight, nil
		case flush:
			return Flush, nil
		case straight | flush:
			if maxRank == rank.Ace && minRank == rank.Ten {
				return RoyalFlush, nil
			}
			return StraightFlush, nil
		default:
			return HighCard, nil
		}
	}
	return Unknown, ErrUnknownHandValue{}
}
