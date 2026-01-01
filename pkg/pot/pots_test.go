package pots

import (
	"slices"
	"testing"

	"github.com/google/uuid"
)

func TestPotsSettle(t *testing.T) {
	id1 := uuid.New().String()
	id2 := uuid.New().String()
	id3 := uuid.New().String()

	testCases := []struct {
		name          string
		contributions []chips
		want          []Pot
	}{
		{
			name: "",
			contributions: []chips{
				{amount: 100, by: id1},
				{amount: 200, by: id2},
				{amount: 300, by: id3},
			},
			want: []Pot{
				pot{contributors: []string{id1, id2, id3}, chips: 300},
				pot{contributors: []string{id2, id3}, chips: 200},
				pot{contributors: []string{id3}, chips: 100},
			},
		},
		{
			name: "",
			contributions: []chips{
				{amount: 200, by: id1},
				{amount: 100, by: id2},
				{amount: 300, by: id3},
			},
			want: []Pot{
				pot{contributors: []string{id1, id2, id3}, chips: 300},
				pot{contributors: []string{id1, id3}, chips: 200},
				pot{contributors: []string{id3}, chips: 100},
			},
		},
		{
			name: "",
			contributions: []chips{
				{amount: 300, by: id1},
				{amount: 200, by: id2},
				{amount: 100, by: id3},
			},
			want: []Pot{
				pot{contributors: []string{id1, id2, id3}, chips: 300},
				pot{contributors: []string{id1, id2}, chips: 200},
				pot{contributors: []string{id1}, chips: 100},
			},
		},
	}

	for _, tc := range testCases {
		pots := New()
		for _, c := range tc.contributions {
			pots.AddChips(c.by, c.amount)
		}
		ps := pots.Settle()
		for i := range len(ps) {
			gotPot := ps[i]
			wantPot := tc.want[i]
			if gotPot.Chips() != wantPot.Chips() { // need to ensure the order of pots is correct
				t.Errorf("got %v, want %v", gotPot.Chips(), wantPot.Chips())
			}

			// don't need to ensure the order of Contributors in the pot
			if len(gotPot.Contributors()) != len(wantPot.Contributors()) {
				t.Errorf("got %v, want %v", gotPot.Contributors(), wantPot.Contributors())
			}
			for _, wc := range wantPot.Contributors() {
				if !slices.ContainsFunc(gotPot.Contributors(),
					func(gc string) bool {
						return gc == wc
					},
				) {
					t.Errorf("got %v, want %v", gotPot.Contributors(), wantPot.Contributors())
				}
			}
		}
	}
}
