package pots

type Pot interface {
	Contributors() []string
	Chips() int
}

type pot struct {
	contributors []string
	chips        int
}

var _ Pot = pot{}

func newPot(ids []string, chips int) Pot {
	return pot{
		contributors: ids,
		chips:        chips,
	}
}

func (p pot) Contributors() []string {
	return p.contributors
}

func (p pot) Chips() int {
	return p.chips
}
