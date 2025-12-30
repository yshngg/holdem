package round

import (
	"context"
	"errors"
	"fmt"

	"github.com/yshngg/holdem/pkg/card"
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
	// players represents the relationship between positions
	// and their corresponding players in the round.
	// Must not be modified after round start.
	players        []*player.Player
	dealer         *dealer.Dealer
	button         int
	minBet         int
	communityCards []*card.Card

	// recorder is used to record events during the round.
	// Such as: logging events, etc.
	recorder watch.Recorder

	// broadcaster is used to broadcast events to all players.
	broadcaster watch.Broadcaster

	// status indicates the current stage of the round.
	// Such as: pre-flop, flop, turn, river, showdown, etc.
	status Status
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

func existsPlayer(players []*player.Player, player *player.Player) bool {
	for _, p := range players {
		if p.ID() == player.ID() {
			return true
		}
	}
	return false
}

func (r *Round) AddPlayer(player *player.Player) error {
	if r.status.After(StatusReady) {
		return errors.New("round has started")
	}
	if existsPlayer(r.players, player) {
		return errors.New("player already exists")
	}
	for _, p := range r.players {
		if p == nil {
			p = player
			return nil
		}
	}
	r.players = append(r.players, player)
	return nil
}

func (r *Round) RemovePlayer(player *player.Player) error {
	if r.status.After(StatusReady) {
		return errors.New("round has started")
	}
	for _, p := range r.players {
		if p.ID() == player.ID() {
			p = nil
			return nil
		}
	}
	return errors.New("player does not exist")
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

func realPlayerCount(players []*player.Player) int {
	count := 0
	for _, p := range players {
		if p != nil {
			count++
		}
	}
	return count
}

func effectivePlayerCount(players []*player.Player) int {
	count := 0
	for _, p := range players {
		if p != nil && p.Status() != player.StatusFolded && p.Status() != player.StatusInvalid {
			count++
		}
	}
	return count
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

func positionBlind(players []*player.Player, button int) (int, int, error) {
	count := 0
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
	small, big, err := positionBlind(r.players, r.button)
	if err != nil {
		return fmt.Errorf("blind positions, err: %v", err)
	}

	if err := r.players[small].Bet(ctx, r.minBet/2); err != nil {
		return fmt.Errorf("post small blind: %w", err)
	}
	r.broadcaster.Action(player.EventPostSmallBlind, player.EventObject{Player: r.players[small], Bet: r.minBet / 2})
	if err := r.players[big].Bet(ctx, r.minBet); err != nil {
		return fmt.Errorf("post big blind: %w", err)
	}
	r.broadcaster.Action(player.EventPostBigBlind, player.EventObject{Player: r.players[big], Bet: r.minBet})
	return nil
}

func (r *Round) Start(ctx context.Context) error {
	if r.players[r.button] == nil {
		return fmt.Errorf("button position does not have player")
	}
	realPlayerCount := realPlayerCount(r.players)
	effectivePlayerCount := effectivePlayerCount(r.players)
	if effectivePlayerCount < MinPlayerCount || effectivePlayerCount > MaxPlayerCount {
		return fmt.Errorf("invalid player count: %d", effectivePlayerCount)
	}

	// ready to start the round
	r.status = StatusStart

	// prepare players
	err := r.prepare(ctx)
	if err != nil {
		return fmt.Errorf("start round, err: %w", err)
	}

	// dealer shuffle deck
	r.dealer.Shuffle()
	r.broadcaster.Action(dealer.EventShuffle, dealer.EventObject{})

	// compulsory bets
	if err := r.betBlind(ctx); err != nil {
		return fmt.Errorf("bet blind, err: %w", err)
	}

	// pre-flop
	r.status = StatusPreFlop
	r.dealer.BurnCard()
	holeCards := r.dealer.DealHoleCards(realPlayerCount)
	for i := range len(r.players) {
		p := r.players[(r.button+i+1)%len(r.players)]
		if p == nil {
			continue
		}
		p.SetHoleCards(holeCards[i])
		r.broadcaster.Action(dealer.EventHoleCards, dealer.EventObject{HoleCards: holeCards[i]})
	}

	// TODO(@yshngg): Implement pre-flop betting logic
	for _, p := range r.players {
		if p == nil || p.Status() == player.StatusInvalid {
			continue
		}
		p.WaitForAction(ctx, map[player.ActionType]player.Action{
			player.ActionCall:  {},
			player.ActionRaise: {},
		})
	}

	// flop
	r.status = StatusFlop
	r.dealer.BurnCard()
	flopCards := r.dealer.DealFlopCards()
	for _, card := range flopCards {
		r.communityCards = append(r.communityCards, card)
	}
	for _, p := range r.players {
		if p == nil || p.Status() == player.StatusInvalid || p.Status() == player.StatusFolded {
			continue
		}
		r.broadcaster.Action(dealer.EventFlopCards, dealer.EventObject{FlopCards: flopCards})
	}

	// turn
	r.status = StatusTurn
	r.dealer.BurnCard()
	turnCard := r.dealer.DealTurnCard()
	for _, p := range r.players {
		if p == nil || p.Status() == player.StatusInvalid {
			continue
		}
		r.broadcaster.Action(dealer.EventTurnCard, dealer.EventObject{TurnCard: turnCard})
		p.WaitForAction(ctx, map[player.ActionType]player.Action{
			player.ActionCall:  {},
			player.ActionRaise: {},
		})
	}

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
