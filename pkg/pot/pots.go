package pots

import "slices"

type Pots interface {
	// Sum calculates the amount of chips in the pot.
	Sum() int

	// AddChips adds chips to the pot.
	AddChips(id string, amount int)

	// GetChips returns the amount of chips contributed by a player.
	ChipsBy(id string) int

	// Settle returns a list of pots after settling the contributions.
	Settle() []Pot
}

func New() Pots {
	return &pots{
		contributions: make(map[string]int),
	}
}

type pots struct {
	contributions map[string]int
}

func (p *pots) AddChips(id string, amount int) {
	p.contributions[id] += amount
}

func (p *pots) ChipsBy(id string) int {
	return p.contributions[id]
}

func (p pots) Sum() int {
	sum := 0
	for _, amount := range p.contributions {
		sum += amount
	}
	return sum
}

type chips struct {
	amount int
	by     string
}

func (p pots) Settle() (potList []Pot) {
	chipsList := []chips{}
	for id, amount := range p.contributions {
		chipsList = append(chipsList, chips{amount, id})
	}
	if len(chipsList) == 0 {
		return potList
	}

	slices.SortFunc(chipsList, func(a, b chips) int {
		return a.amount - b.amount
	})
	for i := range chipsList {
		tmp := chipsList[i]
		amount := tmp.amount
		if amount == 0 {
			continue
		}
		chipsList[i].amount -= amount // That is: chipsList[i].amount = 0
		contributors := map[string]struct{}{tmp.by: {}}
		for j := i + 1; j < len(chipsList); j++ {
			chipsList[j].amount -= amount
			contributors[chipsList[j].by] = struct{}{}
		}
		potList = append(potList, newPot(contributors, amount*len(contributors)))
	}
	return potList
}
