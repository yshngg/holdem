package table

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/yshngg/holdem/pkg/player"
	"github.com/yshngg/holdem/pkg/round"
)

const (
	defaultMinBet   = 2
	defaultCapacity = 8

	defaultMinPlayerCountInRound = 2
	defaultMaxPlayerCountInRound = 22
)

type Table struct {
	round *round.Round

	waiting []*player.Player
	// capacity is cap of waiting slice
	capacity int
	// waitingMap used to quick determine whether the player on table or not.
	waitingMap map[string]struct{}

	minBet int

	minPlayerCountInRound, maxPlayerCountInRound int
}

func New(opts ...Option) *Table {
	t := &Table{
		waitingMap: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(t)
	}
	if t.minBet == 0 {
		t.minBet = defaultMinBet
	}
	if t.capacity < 0 {
		t.capacity = defaultCapacity
	}
	t.waiting = make([]*player.Player, 0, t.capacity)
	if t.maxPlayerCountInRound > defaultMaxPlayerCountInRound {
		t.maxPlayerCountInRound = defaultMaxPlayerCountInRound
	}
	if t.minPlayerCountInRound > defaultMinPlayerCountInRound {
		t.minPlayerCountInRound = defaultMinPlayerCountInRound
	}
	return t
}

type Option func(t *Table)

func WithMinBet(minBet int) Option {
	return func(t *Table) {
		t.minBet = minBet
	}
}

func WithWaitingCapacity(capacity int) Option {
	return func(t *Table) {
		t.capacity = capacity
	}
}

func WithMaxPlayerCountInRound(count int) Option {
	return func(t *Table) {
		t.maxPlayerCountInRound = count
	}
}

func WithMinPlayerCountInRound(count int) Option {
	return func(t *Table) {
		t.minPlayerCountInRound = count
	}
}

func WithPlayers(players []*player.Player) Option {
	return func(t *Table) {
		for _, p := range players {
			t.waiting = append(t.waiting, p)
			t.waitingMap[p.ID().String()] = struct{}{}
		}
	}
}

func (t *Table) Join(p *player.Player) error {
	if t.round != nil {
		err := t.round.AddPlayer(p)
		if err != nil && !errors.Is(err, round.ErrMaxPlayerCountReached{}) {
			return fmt.Errorf("join table, err: %v", err)
		}
	}
	// exists := slices.ContainsFunc(t.waiting, func(pp *player.Player) bool {
	// 	return p.ID() == pp.ID()
	// })
	_, exists := t.waitingMap[p.ID().String()]
	if exists {
		return fmt.Errorf("player %s (id: %s) have sat at the table", p.Name(), p.ID())
	}
	if len(t.waiting) >= t.capacity {
		return fmt.Errorf("have reached the capacity of waiting")
	}
	t.waiting = append(t.waiting, p)
	t.waitingMap[p.ID().String()] = struct{}{}
	return nil
}

func (t *Table) Leave(ctx context.Context, p *player.Player) error {
	if t.round != nil {
		err := t.round.RemovePlayer(ctx, p)
		if err != nil && !errors.Is(err, round.ErrPlayerNotFound{}) {
			return fmt.Errorf("leave table, err: %v", err)
		}
	}
	// exists := slices.ContainsFunc(t.waiting, func(pp *player.Player) bool {
	// 	return p.ID() == pp.ID()
	// })
	_, exists := t.waitingMap[p.ID().String()]
	if !exists {
		return fmt.Errorf("player %s (id: %s) did not sit at the table", p.Name(), p.ID())
	}
	t.waiting = slices.DeleteFunc(t.waiting, func(pp *player.Player) bool {
		return p.ID() == pp.ID()
	})
	delete(t.waitingMap, p.ID().String())
	return nil
}

func (t *Table) Start(ctx context.Context) error {
	players := make([]*player.Player, 0, len(t.waiting))
	button := 0
	for {
		for len(players) < t.maxPlayerCountInRound {
			if len(t.waiting) < 1 {
				return fmt.Errorf("not enough players")
			}
			p := t.waiting[0]
			players = append(players, p)
			if len(t.waiting) < 2 {
				t.waiting = make([]*player.Player, 0)
				break
			}
			t.waiting = t.waiting[1:]
		}
		if len(players) < t.minPlayerCountInRound {
			return fmt.Errorf("not enough players, round cannot start")
		}
		t.round = round.New(
			players,
			round.WithMinBet(t.minBet),
			round.WithButton(button%len(players)),
			round.WithPlayerCount(t.minPlayerCountInRound, t.maxPlayerCountInRound, len(players)),
		)

		err := t.round.Start(ctx)
		t.round.End()
		if err != nil {
			return fmt.Errorf("start round, err: %w", err)
		}
		time.Sleep(5 * time.Second)

		players = t.round.RealPlayers()
		button++
	}
}
