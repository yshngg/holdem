package player

import "testing"

func TestActionTypeIntoStatus(t *testing.T) {
	testCases := []struct {
		name       string
		actionType ActionType
		want       Status
	}{
		{"Invalid", ActionInvalid, StatusReady},
		{"Check", ActionCheck, StatusChecked},
		{"Fold", ActionFold, StatusFolded},
		{"Bet", ActionBet, StatusBetted},
		{"Call", ActionCall, StatusCalled},
		{"Raise", ActionRaise, StatusRaised},
		{"All-In", ActionAllIn, StatusAllIn},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status := tc.actionType.IntoStatus()
			if status != tc.want {
				t.Errorf("IntoStatus(%v) = %v, want %v", tc.actionType, status, tc.want)
			}
		})
	}
}

func TestActionTypeString(t *testing.T) {
	testCases := []struct {
		name       string
		actionType ActionType
		want       string
	}{
		{"Invalid", ActionInvalid, "Invalid"},
		{"Check", ActionCheck, "Check"},
		{"Fold", ActionFold, "Fold"},
		{"Bet", ActionBet, "Bet"},
		{"Call", ActionCall, "Call"},
		{"Raise", ActionRaise, "Raise"},
		{"All-In", ActionAllIn, "AllIn"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			str := tc.actionType.String()
			if str != tc.want {
				t.Errorf("ActionTypeString(%v) = %v, want %v", tc.actionType, str, tc.want)
			}
		})
	}
}
