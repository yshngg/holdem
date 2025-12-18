package session

import "github.com/yshngg/holdem/pkg/player"

type Session struct {
	join   chan *player.Player
	leave  chan *player.Player
	active bool
}

func New() *Session {
	return &Session{
		join:  make(chan *player.Player, 10),
		leave: make(chan *player.Player, 10),
	}
}

func (s *Session) Join(player *player.Player) error {
	s.join <- player
	return nil
}

func (s *Session) Leave(player *player.Player) error {
	s.leave <- player
	return nil
}

func (s *Session) Start() error {

	go func() {
		for {
			select {
			case player := <-s.join:
				// TODO: handle player joining
			case player := <-s.leave:
				// TODO: handle player leaving
			}
		}
	}()
	return nil
}

func (s *Session) Destroy() error {
	return nil
}
