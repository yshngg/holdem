package round

import (
	"fmt"
)

type ErrFirstToActPlayerNotFound struct {
	button int
}

func (e ErrFirstToActPlayerNotFound) Error() string {
	return fmt.Sprintf("no first player found at button %d", e.button)
}

type ErrStatusNotSupported struct {
	status StatusType
}

func (e ErrStatusNotSupported) Error() string {
	return fmt.Sprintf("status not supported: %s", e.status)
}

func (r Round) positionBlind() (int, int, error) {
	playerCount := r.playerCount.current
	length := len(r.players)
	if r.button < 0 || length <= r.button || r.players[r.button] == nil {
		return -1, -1, ErrInvalidButton{button: r.button}
	}
	if playerCount < 2 {
		return -1, -1, ErrInvalidPlayerCount{count: playerCount}
	}
	if playerCount == 2 {
		small := r.button
		big := (small + 1) % length
		for r.players[big] == nil {
			big = (big + 1) % length
		}
		return small, big, nil
	}

	small := (r.button + 1) % length
	big := (small + 1) % length
	for range length {
		if r.players[small] == nil {
			small = (small + 1) % length
			big = (small + 1) % length
			continue
		}
		if r.players[big] != nil {
			break
		}
		big = (big + 1) % length
	}
	return small, big, nil
}

// func (r *Round) positionFirstToAct(players []*player.Player, button int) (int, error) {
// 	length := len(players)
// 	if button < 0 || length <= button || players[button] == nil {
// 		return -1, ErrInvalidButton{button: button}
// 	}
// 	if r.status == StatusPreFlop {
// 		counter := 0
// 		for i := range length {
// 			p := players[(button+i+1)%len(players)]
// 			if p == nil {
// 				continue
// 			}
// 			counter++
// 			if counter > 2 {
// 				return i, nil
// 			}
// 		}
// 		return -1, ErrFirstToActPlayerNotFound{button}
// 	}
// 	for i := range length {
// 		p := players[(button+i+1)%len(players)]
// 		if p != nil {
// 			return i, nil
// 		}
// 	}
// 	return -1, ErrFirstToActPlayerNotFound{button}
// }

func (r *Round) positionFirstToAct() (int, error) {
	playerCount := r.playerCount.current
	if playerCount < 2 {
		return -1, ErrInvalidPlayerCount{count: playerCount}
	}

	small, big, err := r.positionBlind()
	if err != nil {
		return -1, fmt.Errorf("position blind, err: %w", err)
	}

	switch r.status {
	case StatusPreFlop:
		if playerCount == 2 {
			return small, nil
		}
		length := len(r.players)
		for i := range length - 2 {
			p := r.players[(big+i+1)%length]
			if p != nil {
				return i, nil
			}
		}
		return -1, ErrFirstToActPlayerNotFound{r.button}
	case StatusFlop, StatusTurn, StatusRiver:
		if playerCount == 2 {
			return big, nil
		}
		return small, nil
	default:
		return -1, ErrStatusNotSupported{status: r.status}
	}
}
