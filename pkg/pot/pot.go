package pots

type Pot interface {
	Contributors() map[string]struct{}
	Chips() int
}

type pot struct {
	contributors map[string]struct{}
	chips        int
}

var _ Pot = pot{}

func newPot(contributors map[string]struct{}, chips int) Pot {
	return pot{
		contributors: contributors,
		chips:        chips,
	}
}

func (p pot) Contributors() map[string]struct{} {
	return p.contributors
}

func (p pot) Chips() int {
	return p.chips
}
