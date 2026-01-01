package round

import "github.com/yshngg/holdem/pkg/player"

func findPlayerByID(players []*player.Player, id string) *player.Player {
	for _, p := range players {
		if p != nil && p.ID().String() == id {
			return p
		}
	}
	return nil
}

func positionBlind(players []*player.Player, button int) (int, int, error) {
	count := realPlayerCount(players)
	length := len(players)
	if MaxPlayerCount < count || count < MinPlayerCount {
		return -1, -1, ErrInvalidPlayerCount{count: count}
	}
	if button < 0 || length <= button || players[button] == nil {
		return -1, -1, ErrInvalidButton{button: button}
	}
	if count == 2 {
		small := button
		big := (small + 1) % length
		for players[big] == nil {
			big = (big + 1) % length
		}
		return small, big, nil
	}

	small := (button + 1) % length
	big := (small + 1) % length
	for range length {
		if players[small] == nil {
			small = (small + 1) % length
			big = (small + 1) % length
			continue
		}
		if players[big] != nil {
			break
		}
		big = (big + 1) % length
	}
	return small, big, nil
}

// positionUTG positions the UTG (Under The Gun) player
func positionUTG(players []*player.Player, button int) (int, error) {
	count := realPlayerCount(players)
	length := len(players)
	if MaxPlayerCount < count || count < MinPlayerCount {
		return -1, ErrInvalidPlayerCount{count: count}
	}
	if button < 0 || length <= button || players[button] == nil {
		return -1, ErrInvalidButton{button: button}
	}
	counter := 0
	for i := range length {
		p := players[(button+i+1)%len(players)]
		if p == nil {
			continue
		}
		counter++
		if counter > 2 {
			return i, nil
		}
	}
	return button, nil
}

func positionFirstToAct(players []*player.Player, button int) (int, error) {
	count := realPlayerCount(players)
	length := len(players)
	if MaxPlayerCount < count || count < MinPlayerCount {
		return -1, ErrInvalidPlayerCount{count: count}
	}
	if button < 0 || length <= button || players[button] == nil {
		return -1, ErrInvalidButton{button: button}
	}
	for i := range length {
		p := players[(button+i+1)%len(players)]
		if p != nil {
			return i, nil
		}
	}
	return -1, ErrInvalidPlayerCount{count: count}
}

func realPlayerCount(players []*player.Player) int {
	count := 0
	for _, p := range players {
		if p != nil {
			count++
		}
	}
	return count
}

func realPlayers(players []*player.Player) []*player.Player {
	ps := make([]*player.Player, 0)
	for _, p := range players {
		if p != nil {
			ps = append(ps, p)
		}
	}
	return ps
}

func effectivePlayerCount(players []*player.Player) int {
	count := 0
	for _, p := range players {
		if p != nil && (p.Status() == player.StatusWaitingToAct || p.Status() == player.StatusTakingAction || p.Status() != player.StatusAllIn) {
			count++
		}
	}
	return count
}

func effectivePlayers(players []*player.Player) []*player.Player {
	ps := make([]*player.Player, 0)
	for _, p := range players {
		if p != nil && p.Status() != player.StatusFolded {
			ps = append(ps, p)
		}
	}
	return ps
}
