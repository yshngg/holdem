package player

import "testing"

func TestStatusString(t *testing.T) {
	testCases := []struct {
		name   string
		status StatusType
		want   string
	}{
		{"Idle", StatusIdle, "Idle"},
		{"Ready", StatusReady, "Ready"},
		{"Waiting", StatusWaitingToAct, "Waiting"},
		{"TakingAction", StatusTakingAction, "TakingAction"},
		{"Folded", StatusFolded, "Folded"},
		{"AllIn", StatusAllIn, "AllIn"},
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
