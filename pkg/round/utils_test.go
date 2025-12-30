package round

import (
	"testing"

	"github.com/yshngg/holdem/pkg/player"
)

func TestBlindPositions(t *testing.T) {
	type want struct {
		small int
		big   int
		err   error
	}
	testCases := []struct {
		name    string
		players []*player.Player
		button  int
		want
	}{
		{
			name: "InvalidPlayerCount",
			players: []*player.Player{
				player.New(),
			},
			button: 0,
			want:   want{-1, -1, ErrInvalidPlayerCount{count: 1}},
		},
		{
			name: "InvalidButtonPosition",
			players: []*player.Player{
				player.New(),
				player.New(),
			},
			button: 2,
			want:   want{-1, -1, ErrInvalidButton{button: 2}},
		},
		{
			name: "TowPlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
			},
			button: 0,
			want:   want{0, 1, nil},
		},
		{
			name: "ThreePlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
				player.New(),
			},
			button: 0,
			want:   want{1, 2, nil},
		},
		{
			name: "FourPlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
				player.New(),
				player.New(),
			},
			button: 3,
			want:   want{0, 1, nil},
		},
		{
			name: "FivePlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
				player.New(),
				player.New(),
				player.New(),
			},
			button: 3,
			want:   want{4, 0, nil},
		},
		{
			name: "Absence",
			players: []*player.Player{
				nil,
				player.New(),
				player.New(),
			},
			button: 1,
			want:   want{1, 2, nil},
		},
		{
			name: "Absence",
			players: []*player.Player{
				nil,
				player.New(),
				nil,
			},
			button: 1,
			want:   want{-1, -1, ErrInvalidPlayerCount{1}},
		},
		{
			name: "Absence",
			players: []*player.Player{
				player.New(),
				nil,
				player.New(),
				nil,
				player.New(),
				nil,
			},
			button: 2,
			want:   want{4, 0, nil},
		},
		{
			name: "Absence",
			players: []*player.Player{
				nil,
				player.New(),
				nil,
				nil,
				player.New(),
				nil,
				nil,
				player.New(),
				nil,
			},
			button: 4,
			want:   want{7, 1, nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			small, big, err := positionBlind(tc.players, tc.button)
			if err != tc.want.err {
				t.Errorf("err: %v, want: %v", err, tc.want.err)
			}
			if small != tc.want.small {
				t.Errorf("small blind positions: %v, want: %v", small, tc.want.small)
			}
			if big != tc.want.big {
				t.Errorf("big blind positions: %v, want: %v", big, tc.want.big)
			}
		})
	}
}

func TestEffectivePlayerCount(t *testing.T) {
	testCases := []struct {
		name    string
		players []*player.Player
		want    int
	}{
		{
			name:    "Empty",
			players: []*player.Player{},
			want:    0,
		},
		{
			name: "OnePlayer",
			players: []*player.Player{
				player.New(),
			},
			want: 1,
		},
		{
			name: "TwoPlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
			},
			want: 2,
		},
		{
			name: "ThreePlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
				player.New(),
			},
			want: 3,
		},
		{
			name: "ThreePlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
				player.New(),
				nil,
				nil,
			},
			want: 3,
		},
		{
			name: "ThreePlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
				player.New(),
				nil,
				player.New(player.WithStatus(player.StatusFolded)),
			},
			want: 3,
		},
		{
			name: "ThreePlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
				player.New(),
				player.New(player.WithStatus(player.StatusFolded)),
				player.New(player.WithStatus(player.StatusFolded)),
			},
			want: 3,
		},
		{
			name: "ThreePlayers",
			players: []*player.Player{
				player.New(player.WithStatus(player.StatusChecked)),
				player.New(player.WithStatus(player.StatusFolded)),
				player.New(player.WithStatus(player.StatusBetted)),
				player.New(player.WithStatus(player.StatusCalled)),
				player.New(player.WithStatus(player.StatusRaised)),
				player.New(player.WithStatus(player.StatusAllIn)),
			},
			want: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			count := effectivePlayerCount(tc.players)
			if count != tc.want {
				t.Errorf("effective player count: %v, want: %v", count, tc.want)
			}
		})
	}
}

func TestRealPlayerCount(t *testing.T) {
	testCases := []struct {
		name    string
		players []*player.Player
		want    int
	}{
		{
			name:    "Empty",
			players: []*player.Player{},
			want:    0,
		},
		{
			name: "OnePlayer",
			players: []*player.Player{
				player.New(),
			},
			want: 1,
		},
		{
			name: "TwoPlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
			},
			want: 2,
		},
		{
			name: "ThreePlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
				player.New(),
			},
			want: 3,
		},
		{
			name: "ThreePlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
				player.New(),
				nil,
			},
			want: 3,
		},
		{
			name: "ThreePlayers",
			players: []*player.Player{
				player.New(),
				player.New(),
				player.New(),
				nil,
				nil,
			},
			want: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			count := realPlayerCount(tc.players)
			if count != tc.want {
				t.Errorf("real player count: %v, want: %v", count, tc.want)
			}
		})
	}
}
