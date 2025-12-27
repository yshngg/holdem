package round

import (
	"fmt"
	"testing"

	"github.com/yshngg/holdem/pkg/player"
)

func TestRound(t *testing.T) {
	playerCount := 5
	playerChips := 100
	players := make([]*player.Player, 0, playerCount)
	for i := range 5 {
		p := player.New(player.WithName(fmt.Sprintf("player-%d", i)), player.WithChips(playerChips))
		players = append(players, p)
	}
	r := New(players)
	r.Start(t.Context())
	r.End()
}
