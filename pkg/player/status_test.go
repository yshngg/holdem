package player

import "testing"

func TestStatusString(t *testing.T) {
	testCases := []struct {
		name   string
		status Status
		want   string
	}{
		{"Invalid", StatusInvalid, "Invalid"},
		{"Check", StatusChecked, "Check"},
		{"Fold", StatusFolded, "Fold"},
		{"Bet", StatusBetted, "Bet"},
		{"Call", StatusCalled, "Call"},
		{"Raise", StatusRaised, "Raise"},
		{"All-In", StatusAllIn, "AllIn"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status := tc.status.String()
			if status != tc.want {
				t.Errorf("StatusString(%v) = %v, want %v", tc.status, status, tc.want)
			}
		})
	}
}
