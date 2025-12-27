package round

import (
	"context"
	"fmt"

	"github.com/yshngg/holdem/pkg/dealer"
	"github.com/yshngg/holdem/pkg/player"
	"github.com/yshngg/holdem/pkg/watch"
)

const (
	defaultMinBet = 2
	defaultButton = 0

	MinPlayerCount = 2
	MaxPlayerCount = 10
)

type Round struct {
	players []*player.Player
	dealer  *dealer.Dealer
	button  int
	minBet  int
	status  Status

	watcher     watch.Interface
	broadcaster watch.Broadcaster
}

func New(players []*player.Player, opts ...Option) *Round {
	r := &Round{
		players: players,
		button:  -1,
		minBet:  -1,
		status:  StatusReady,
	}
	for _, opt := range opts {
		opt(r)
	}

	if r.dealer == nil {
		r.dealer = dealer.New()
	}
	if r.button < 0 {
		r.button = defaultButton
	}
	if r.minBet < 0 {
		r.minBet = defaultMinBet
	}
	if r.broadcaster == nil {
		queueLength := len(players) * 2
		r.broadcaster = watch.NewBroadcaster(queueLength, queueLength)
	}
	if r.watcher == nil {
		watcher, err := r.broadcaster.Watch()
		// should not have an error
		if err != nil {
			panic(err)
		}
		r.watcher = watcher
	}
	return r
}

type Option func(*Round)

func WithButton(button int) Option {
	return func(r *Round) {
		r.button = button
	}
}

func WithMinBet(minBet int) Option {
	return func(r *Round) {
		r.minBet = minBet
	}
}

func WithDealer(dealer *dealer.Dealer) Option {
	return func(r *Round) {
		r.dealer = dealer
	}
}

func WithBroadcaster(broadcaster watch.Broadcaster) Option {
	return func(r *Round) {
		r.broadcaster = broadcaster
	}
}

func WithWatcher(watcher watch.Interface) Option {
	return func(r *Round) {
		r.watcher = watcher
	}
}

func (r *Round) playerStartWatching(ctx context.Context) error {
	for _, p := range r.players {
		if p == nil {
			continue
		}
		watcher, err := r.broadcaster.Watch()
		if err != nil {
			return fmt.Errorf("failed to watch broadcaster: %w", err)
		}
		watcher = watch.Filter(watcher, func(in watch.Event) (out watch.Event, keep bool) {
			// TODO(@yshngg): implement filter function
			return in, true
		})
		go func(p *player.Player) {
			select {
			case <-ctx.Done():
				return
			default:
				p.Watch(watcher)
			}
		}(p)
	}

	return nil
}

func (r *Round) playerCount() (count int) {
	for _, p := range r.players {
		if p != nil {
			count++
		}
	}
	return
}

func (r *Round) blindPositions() (int, int) {
	small := (r.button + 1) % len(r.players)
	big := 0
	for range len(r.players) {
		if r.players[small] == nil {
			small = (small + 1) % len(r.players)
			big = small
			continue
		}
		big = (big + 1) % len(r.players)
		if r.players[big] != nil {
			break
		}
	}
	return small, big
}

func (r *Round) betBlind() error {
	small, big := r.blindPositions()

	if err := r.players[small].Bet(r.minBet / 2); err != nil {
		return fmt.Errorf("post small blind: %w", err)
	}
	if err := r.players[big].Bet(r.minBet); err != nil {
		return fmt.Errorf("post big blind: %w", err)
	}
	return nil
}

func (r *Round) Start(ctx context.Context) error {
	if r.players[r.button] == nil {
		return fmt.Errorf("button position does not have player")
	}
	playerCount := r.playerCount()
	if playerCount < MinPlayerCount || playerCount > MaxPlayerCount {
		return fmt.Errorf("invalid player count: %d", playerCount)
	}

	// players start watching events
	err := r.playerStartWatching(ctx)
	if err != nil {
		return fmt.Errorf("start round, err: %w", err)
	}

	// ready to start the round
	r.status = StatusStart

	// dealer shuffle deck
	r.dealer.Shuffle()

	// compulsory bets
	if err := r.betBlind(); err != nil {
		return err
	}

	// pre-flop
	r.status = StatusPreFlop
	cards := r.dealer.DealHoleCards(playerCount)
	// for i, p := range r.effectivePlayers() {
	// p.SetHoleCards(cards[(r.button+i)%len(cards)])
	// }
	for i := range playerCount {
		p := r.players[i]
		if p == nil {
			continue
		}
		r.players[i].SetHoleCards(cards[i])
	}

	// flop
	r.status = StatusFlop

	// flop
	r.status = StatusFlop

	// turn
	r.status = StatusTurn

	// river
	r.status = StatusRiver

	// showdown
	r.status = StatusShowdown

	return nil
}

func (r *Round) End() error {
	r.status = StatusEnd
	return nil
}

func (r *Round) effectivePlayers() []*player.Player {
	players := make([]*player.Player, 0)
	for _, p := range r.players {
		if p != nil {
			players = append(players, p)
		}
	}
	return players
}
