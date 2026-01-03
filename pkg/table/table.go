package table

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/yshngg/holdem/pkg/player"
	"github.com/yshngg/holdem/pkg/round"
)

const defaultMinBet = 2

type Table struct {
	sync.Once
	eventChan chan Event
	round     *round.Round

	waiting []*player.Player
	// auxiliary used to quick determine whether the player on table or not.

	minBet int
}

func New(opts ...Option) *Table {
	t := &Table{
		// event:     make(chan Event, 0),
		waiting: make([]*player.Player, 0),
		// auxiliary: make(map[string]struct{}, 0),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

type Option func(t *Table)

func WithMinBet(minBet int) Option {
	return func(t *Table) {
		t.minBet = minBet
	}
}

func WithPlayers(players []*player.Player) Option {
	return func(t *Table) {
		for _, p := range players {
			t.waiting = append(t.waiting, p)
			// t.auxiliary[p.ID().String()] = struct{}{}
		}
	}
}

func (t *Table) Join(p *player.Player) error {
	if t.round != nil && t.round.Status() == round.StatusReady {
		err := t.round.AddPlayer(p)
		if err != nil {
			return fmt.Errorf("join table, err: %v", err)
		}
	}
	// _, exists := t.auxiliary[p.ID().String()]
	exists := slices.ContainsFunc(t.waiting, func(pp *player.Player) bool {
		return p.ID() == pp.ID()
	})
	if exists {
		return fmt.Errorf("player %s (id: %s) have sat at the table", p.Name(), p.ID())
	}
	t.waiting = append(t.waiting, p)
	// t.auxiliary[p.ID().String()] = struct{}{}
	// t.event <- newEvent(EventPlayerJoin, p)
	return nil
}

func (t *Table) Leave(p *player.Player) error {
	if t.round != nil && t.round.Status() == round.StatusReady {
		err := t.round.RemovePlayer(p)
		if err != nil {
			return fmt.Errorf("join table, err: %v", err)
		}
	}
	_, exists := t.auxiliary[p.ID().String()]
	if !exists {
		return fmt.Errorf("player %s (id: %s) did not sit at the table", p.Name(), p.ID())
	}
	// t.event <- newEvent(EventPlayerLeave, p)
	t.players = append(t.players, p)
	t.auxiliary[p.ID().String()] = struct{}{}
	return nil
}

func (t *Table) Players() []*player.Player {
	return nil
}

func (t *Table) Start(ctx context.Context) error {
	go func(ctx context.Context, eventChan chan Event) {
		select {
		case event := <-eventChan:

		case <-ctx.Done():
			return
		}
	}(ctx, t.eventChan)

	for {
		round := round.New(t.Players())
		players := round.Players()
		round.Start(ctx)
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

func (t *Table) Destroy() {
	t.Once.Do(
		func() {
			close(t.eventChan)
		},
	)
}
