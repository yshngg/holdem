package session

import (
	"github.com/yshngg/holdem/pkg/game/round"
	"github.com/yshngg/holdem/pkg/player"
)

type Session struct {
	event  []Event
	active bool
	round  round.Round
}

func New() *Session {
	return &Session{
		event: make([]Event, 0),
	}
}

func (s *Session) Join(player *player.Player) error {
	s.event = append(s.event, newEvent(EventTypeJoin, player))
	return nil
}

func (s *Session) Leave(player *player.Player) error {
	s.event = append(s.event, newEvent(EventTypeLeave, player))
	return nil
}

func (s *Session) Start() error {
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
