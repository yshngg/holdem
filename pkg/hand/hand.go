package hand

type Head int

const (
	_             Head = iota
	HighCard           // Simple value of the card. Lowest: 2 – Highest: Ace (King in the example)
	Pair               // Two cards with the same value
	TwoPairs           // Two times two cards with the same value
	ThreeOfAKind       // Three cards with the same value
	Straight           // Sequence of 5 cards in increasing value (Ace can precede 2 or follow up King, but not both), not of the same suit
	Flush              // 5 cards of the same suit, not in sequential order
	FullHouse          // Combination of three of a kind and a pair
	FourOfAKind        // Four cards of the same value
	StraightFlush      // Straight of the same suit
	RoyalFlush         // Highest straight of the same suit
)
