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
	players []*player.Player // position: player (if any)
	dealer  *dealer.Dealer
	button  int
	minBet  int
	status  Status

	recorder    watch.Recorder
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
	if r.recorder == nil {
		watcher, err := r.broadcaster.Watch()
		// should not have an error
		if err != nil {
			panic(err)
		}
		r.recorder = watch.NewRecorder(watcher)
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

func WithRecorder(recorder watch.Recorder) Option {
	return func(r *Round) {
		r.recorder = recorder
	}
}

func (r *Round) prepare(ctx context.Context) error {
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
		p.Apply(player.WithWatcher(watcher))
	}

	return nil
}

func (r *Round) playerCount() (count int) {
	for _, p := range r.players {
		if p != nil && p.Status() != player.StatusFolded && p.Status() != player.StatusInvalid {
			count++
		}
	}
	return
}

type ErrInvalidPlayerCount struct {
	count int
}

func (e ErrInvalidPlayerCount) Error() string {
	return fmt.Sprintf("invalid player count: %d", e.count)
}

type ErrInvalidButton struct {
	button int
}

func (e ErrInvalidButton) Error() string {
	return fmt.Sprintf("invalid button: %d", e.button)
}

func blindPositions(players []*player.Player, button int) (int, int, error) {
	count := 0
	length := len(players)
	for _, p := range players {
		if p != nil {
			count++
		}
	}
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
		return button, big, nil
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

func (r *Round) betBlind(ctx context.Context) error {
	small, big, err := blindPositions(r.players, r.button)
	if err != nil {
		return fmt.Errorf("blind positions, err: %v", err)
	}

	if err := r.players[small].Bet(ctx, r.minBet/2); err != nil {
		return fmt.Errorf("post small blind: %w", err)
	}
	if err := r.players[big].Bet(ctx, r.minBet); err != nil {
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

	// prepare players
	err := r.prepare(ctx)
	if err != nil {
		return fmt.Errorf("start round, err: %w", err)
	}

	// ready to start the round
	r.status = StatusStart

	// dealer shuffle deck
	r.dealer.Shuffle()
	r.broadcaster.Action(dealer.EventShuffle, dealer.EventObject{})

	// compulsory bets
	// if err := r.betBlind(); err != nil {
	// 	return err
	// }

	// pre-flop
	r.status = StatusPreFlop
	cards := r.dealer.DealHoleCards(playerCount)
	r.broadcaster.Action(dealer.EventHoleCards, dealer.EventObject{})
	for _, p := range r.players {
		if p == nil {
			continue
		}
		p.WaitForAction(ctx, map[player.ActionType]player.Action{
			player.ActionCheck: player.Action{},
			player.ActionCall:  player.Action{},
			player.ActionBet:   player.Action{},
			player.ActionRaise: player.Action{},
		})
	}
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
