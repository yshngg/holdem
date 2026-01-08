package table

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/yshngg/holdem/pkg/player"
	"github.com/yshngg/holdem/pkg/round"
	"github.com/yshngg/holdem/pkg/watch"
	"k8s.io/klog/v2"
)

const (
	defaultMinBet        = 2
	defaultCapacity      = 8 // [2: 22]
	defaultActionTimeout = 5 * time.Second

	MinPlayerCount = 2
	MaxPlayerCount = 22
)

type Table struct {
	round *round.Round

	// left indicates the players who left the table.
	// Need to remove from the Round after the round is finished.
	left map[string]struct{}

	// capacity is seating capacity of the table
	// equal to length of positions
	capacity int

	// position indicates the position relationship among players in round.
	position []*string

	// waiting indicates the players who are waiting to join the table.
	// waiting is a FIFO queue.
	waiting []string

	// players is all players in round, waiting queue and left.
	players map[string]*player.Player

	// minBet is minimum bet on the table.
	minBet int

	// player need at least `threshold` chips to join.
	// threshold must greater than minBet.
	// if `threshold <= 0`, the value will be `minBet * 4`.
	threshold int

	// actionTimeout indicates how long can player take actions
	actionTimeout time.Duration

	watcher watch.Interface

	broadcaster watch.Broadcaster
}

func New(opts ...Option) *Table {
	t := &Table{}
	for _, opt := range opts {
		opt(t)
	}
	if t.minBet == 0 {
		t.minBet = defaultMinBet
	}
	if t.capacity < MinPlayerCount || t.capacity > MaxPlayerCount {
		t.capacity = defaultCapacity
	}
	if t.threshold <= 0 {
		// if threshold is invalid, it's value will be four times of minBet
		t.threshold = t.minBet * 4
	}
	if t.actionTimeout <= 0 {
		t.actionTimeout = defaultActionTimeout
	}
	t.players = make(map[string]*player.Player, t.capacity)
	t.position = make([]*string, 0, t.capacity)
	t.waiting = make([]string, 0, t.capacity)
	t.left = make(map[string]struct{}, 0)
	queueLength := t.capacity * 2
	t.broadcaster = watch.NewBroadcaster(queueLength, queueLength)
	watcher, err := t.broadcaster.Watch()
	if err != nil {
		panic(err)
	}
	t.watcher = watcher
	return t
}

type Option func(t *Table)

func WithMinBet(minBet int) Option {
	return func(t *Table) {
		t.minBet = minBet
	}
}

func WithCapacity(capacity int) Option {
	return func(t *Table) {
		t.capacity = capacity
	}
}

func WithChipsThreshold(threshold int) Option {
	return func(t *Table) {
		t.threshold = threshold
	}
}

func WithActionTimeout(timeout time.Duration) Option {
	return func(t *Table) {
		t.actionTimeout = timeout
	}
}

func (t *Table) PlayerCount() int {
	return len(t.players) - len(t.left)
}

func (t *Table) Join(name, id string, chips int) (*player.Player, error) {
	// exists := slices.ContainsFunc(t.waiting, func(pp *player.Player) bool {
	// 	return p.ID() == pp.ID()
	// })

	_, exists := t.players[id]
	// _, exists := t.waitingMap[p.ID()]
	if exists {
		return nil, ErrPlayerNotFound{id: id}
	}
	if t.PlayerCount() >= t.capacity {
		return nil, fmt.Errorf("have reached the capacity of table")
	}

	watcher, err := t.broadcaster.Watch()
	if err != nil {
		return nil, fmt.Errorf("watch broadcaster, err: %w", err)
	}
	p := player.New(
		player.WithName(name),
		player.WithID(id),
		player.WithChips(chips),
		player.WithWatcher(watcher),
		player.WithActionTimeout(t.actionTimeout),
	)

	t.players[p.ID()] = p
	t.waiting = append(t.waiting, p.ID())
	t.sitDown(p.ID())
	return p, nil
}

func (t *Table) sitDown(id string) {
	for i, idPtr := range t.position {
		if idPtr != nil {
			continue
		}
		t.position[i] = &id
	}
}

type ErrPlayerNotFound struct {
	id string
}

func (e ErrPlayerNotFound) Error() string {
	return fmt.Sprintf("player (id: %s) did not sit at the table", e.id)
}

func (t *Table) Leave(ctx context.Context, id string) error {
	p, exists := t.players[id]
	if !exists {
		return ErrPlayerNotFound{id: id}
	}
	p.StopWatch()
	waiting := slices.ContainsFunc(t.waiting, func(idd string) bool {
		return idd == id
	})
	if !waiting {
		t.left[id] = struct{}{}
		return nil
	}
	t.waiting = slices.DeleteFunc(t.waiting, func(wid string) bool {
		return wid == id
	})
	return nil
}

func (t *Table) Start(ctx context.Context) error {
	t.logEvents(ctx)
	players := make([]*player.Player, 0, len(t.waiting))
	roundNumber := 0
	button := 0

	for {
		readyPosition := make([]*string, len(t.position))
		copy(readyPosition, t.position)
		readyPosition = slices.DeleteFunc(readyPosition, func(id *string) bool {
			if id == nil {
				return false
			}
			p := t.players[*id]
			if p.Status() == player.StatusReady {
				return false
			}
			return true
		})
		readyPlayers := make([]*player.Player, len(readyPosition))
		readyPlayerCount := 0
		for i, id := range readyPosition {
			if id == nil {
				continue
			}
			readyPlayers[i] = t.players[*id]
			readyPlayerCount++
		}

		if MinPlayerCount > readyPlayerCount || readyPlayerCount > MaxPlayerCount {
			break
		}

		for i := range t.capacity {
			idPtr := t.position[(button+i)%t.capacity]
			if idPtr == nil || t.players[*idPtr].Status() != player.StatusReady {
				continue
			}
			button = i
		}

		t.round = round.New(
			readyPlayers,
			round.WithNumber(roundNumber),
			round.WithMinBet(t.minBet),
			round.WithButton(button%len(players)),
			round.WithBroadcaster(t.broadcaster),
		)

		err := t.round.Start(ctx)
		t.round.End()
		if err != nil {
			return fmt.Errorf("start round, err: %w", err)
		}
		time.Sleep(5 * time.Second)

		button++
		roundNumber++
		t.clean()
	}
	return nil
}

func (t *Table) clean() {
	func() {
		for id := range t.left {
			p := t.players[id]
			p.Gone()
			delete(t.players, id)
			defer delete(t.left, id)
		}
	}()
}

func (t *Table) logEvents(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			return
		case event := <-t.watcher.Watch():
			klog.V(3).Info(event)
		}
	}()
}
