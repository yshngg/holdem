package suit

type Suit int

const (
	_ Suit = iota
	Clubs
	Spades
	Hearts
	Diamonds
)

func (s Suit) String() string {
	switch s {
	case Clubs:
		return "Clubs"
	case Spades:
		return "Spades"
	case Hearts:
		return "Hearts"
	case Diamonds:
		return "Diamonds"
	default:
		return "Invalid"
	}
}
