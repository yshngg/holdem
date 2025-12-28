package table

import (
	"github.com/yshngg/holdem/pkg/game/round"
	"github.com/yshngg/holdem/pkg/player"
)

type Table struct {
	event  []Event
	active bool
	round  round.Round
	minBet int
}

func New() *Table {
	return &Table{
		event: make([]Event, 0),
	}
}

func (s *Table) Join(player *player.Player) error {
	s.event = append(s.event, newEvent(EventTypeJoin, player))
	return nil
}

func (s *Table) Leave(player *player.Player) error {
	s.event = append(s.event, newEvent(EventTypeLeave, player))
	return nil
}

func (s *Table) Start() error {
	round := round.New()
	for {
		players := round.Players()
		for _, event := range s.event {
			switch event.Type() {
			case EventTypeJoin:
				players = append(players, event.Player())
			case EventTypeLeave:
				// handle leave
			}
		}
	}

	return nil
}

func (s *Session) Destroy() error {
	return nil
}
