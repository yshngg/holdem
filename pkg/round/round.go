package round

import (
	"context"
	"errors"
	"fmt"

	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/dealer"
	"github.com/yshngg/holdem/pkg/player"
	"github.com/yshngg/holdem/pkg/pot"
	"github.com/yshngg/holdem/pkg/watch"
)

const (
	defaultMinBet = 2
	defaultButton = 0

	MinPlayerCount = 2
	MaxPlayerCount = 10
)

type Round struct {
	// players is a slice of players indexed by their position at the table.
	// The slice length corresponds to the maximum number of seats.
	// Must not be modified after the round starts.
	players []*player.Player

	// dealer handles card shuffling and dealing operations.
	// Responsible for dealing hole cards to players and community cards to the table.
	dealer *dealer.Dealer

	// button is the position of the dealer button in the round.
	// The player immediately to the left of the button posts the small blind.
	button int

	// minBet is the minimum bet required to call in the current round.
	// Typically equals the big blind amount (double the small blind).
	minBet int

	// communityCards are the shared cards visible to all players.
	// The length progresses through 0 (pre-flop), 3 (flop), 4 (turn), and 5 (river).
	communityCards []*card.Card

	// pots holds all active pots in the round, including the main pot and side pots.
	// pots[0] is always the main pot; subsequent elements are side pots (if any).
	pots []pot.Pot

	// recorder captures all game events for replay, debugging, or auditing purposes.
	// It logs actions like bets, folds, and card deals.
	recorder watch.Recorder

	// broadcaster delivers real-time game events to all connected players.
	// Ensures players receive synchronized updates about round state changes.
	broadcaster watch.Broadcaster

	// status tracks the current stage of the poker round.
	// See the Status type for all possible values (pre-flop, flop, turn, river, etc.).
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

func (r *Round) betBlind(ctx context.Context) error {
	small, big, err := positionBlind(r.players, r.button)
	if err != nil {
		return fmt.Errorf("blind positions, err: %v", err)
	}

	if err := r.players[small].Bet(ctx, r.minBet/2); err != nil {
		return fmt.Errorf("post small blind: %w", err)
	}
	if err := r.players[big].Bet(ctx, r.minBet); err != nil {
		return fmt.Errorf("post big blind: %w", err)
	}

	if err = r.broadcaster.Action(player.EventPostSmallBlind, player.EventObject{Player: r.players[small], Bet: r.minBet / 2}); err != nil {
		return fmt.Errorf("broadcaster action, err: %w", err)
	}
	if err = r.broadcaster.Action(player.EventPostBigBlind, player.EventObject{Player: r.players[big], Bet: r.minBet}); err != nil {
		return fmt.Errorf("broadcaster action, err: %w", err)
	}
	return nil
}

func (r *Round) Status() StatusType {
	return r.status
}

func (r *Round) bettingRound(ctx context.Context) (err error) {
	playerBetChips := make(map[string]int)

	keep := func() bool {
		highestBetChips := 0
		highestBetPosition := 0
		for i, p := range r.players {
			if p == nil || p.Status() != player.StatusWaiting {
				continue
			}
			bet, acted := playerBetChips[p.ID().String()]
			if !acted || bet < 0 {
				continue
			}
			if bet > highestBetChips {
				highestBetChips = bet
				highestBetPosition = i
			}
		}

		for i := range len(r.players) - 1 {
			p := r.players[(highestBetPosition+i+1)%len(r.players)]
			if p == nil || p.Status() != player.StatusWaiting {
				continue
			}
			bet, acted := playerBetChips[p.ID().String()]
			if bet < 0 {
				continue
			}
			if !acted {
				return true
			}
			if 0 < bet && bet < highestBetChips {
				return true
			}
		}
		return false
	}

	start := 0
	if r.Status() == StatusPreFlop {
		start, err = positionUTG(r.players, r.button)
		if err != nil {
			return fmt.Errorf("betting round, err: %w", err)
		}
	} else {
		start, err = positionFirstToAct(r.players, r.button)
		if err != nil {
			return fmt.Errorf("betting round, err: %w", err)
		}
	}

	for keep() { // reopen the betting action
		for i := range len(r.players) {
			p := r.players[(start+i)%len(r.players)]
			if p == nil || p.Status() != player.StatusWaiting {
				continue
			}
			bet, acted := playerBetChips[p.ID().String()]
			if bet < 0 {
				continue
			}
			start = i
		}
		// p = r.players[(utgi+1)%len(r.players)]
	}

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
	if err = r.broadcaster.Action(dealer.EventShuffle, dealer.EventObject{}); err != nil {
		return fmt.Errorf("broadcaster action, err: %w", err)
	}

	// compulsory bets
	if err := r.betBlind(ctx); err != nil {
		return fmt.Errorf("bet blind, err: %w", err)
	}

	// pre-flop
	r.status = StatusPreFlop
	holeCards := r.dealer.DealHoleCards(realPlayerCount)
	for i := range len(r.players) {
		p := r.players[(r.button+i+1)%len(r.players)]
		if p == nil || p.Status() != player.StatusReady {
			continue
		}
		p.SetHoleCards(holeCards[i])
		if err = r.broadcaster.Action(dealer.EventHoleCards, dealer.EventObject{HoleCards: holeCards[i]}); err != nil {
			return fmt.Errorf("broadcaster action, err: %w", err)
		}
	}

	// pre-flop betting round
	preflopBet := r.minBet
	for i := range len(r.players) {
		p := r.players[(r.button+i+3)%len(r.players)]
		if p == nil || p.Status() != player.StatusReady {
			continue
		}
		action := p.WaitForAction(ctx, []player.Action{
			{Type: player.ActionFold, Chips: 0},
			{Type: player.ActionCall, Chips: preflopBet},
			{Type: player.ActionRaise, Chips: preflopBet * 2},
			// {Type: player.ActionAllIn, Chips: 0},
		})
		switch action.Type {
		case player.ActionFold:
		case player.ActionCall:
		case player.ActionRaise:
			preflopBet = action.Chips
		case player.ActionAllIn:
			// TODO(@yshngg): Handle all-in action
		default:
			return fmt.Errorf("invalid action type")
		}
		if err = r.broadcaster.Action(action.Type.IntoEventType(), player.EventObject{Player: p, Bet: action.Chips}); err != nil {
			return fmt.Errorf("broadcaster action, err: %w", err)
		}
	}

	// flop
	r.status = StatusFlop
	r.dealer.BurnCard()
	flopCards := r.dealer.DealFlopCards()
	for _, card := range flopCards {
		r.communityCards = append(r.communityCards, card)
	}
	flopBet := 0
	for _, p := range r.players {
		if p == nil || p.Status() == player.StatusFolded {
			continue
		}
		r.broadcaster.Action(dealer.EventFlopCards, dealer.EventObject{FlopCards: flopCards})
	}

	// turn
	r.status = StatusTurn
	r.dealer.BurnCard()
	turnCard := r.dealer.DealTurnCard()
	r.communityCards = append(r.communityCards, turnCard)
	for _, p := range r.players {
		if p == nil || p.Status() == player.StatusFolded {
			continue
		}
		r.broadcaster.Action(dealer.EventTurnCard, dealer.EventObject{TurnCard: turnCard})
		action := p.WaitForAction(ctx, []player.Action{
			{Type: player.ActionFold, Chips: 0},
			{Type: player.ActionCall, Chips: preflopBet},
			{Type: player.ActionRaise, Chips: preflopBet * 2},
			// {Type: player.ActionAllIn, Chips: 0},
		})
		if err = r.broadcaster.Action(action.Type.IntoEventType(), player.EventObject{Player: p, Bet: action.Chips}); err != nil {
			return fmt.Errorf("broadcaster action, err: %w", err)
		}
	}

	// river
	r.status = StatusRiver
	r.dealer.BurnCard()
	riverCard := r.dealer.DealTurnCard()
	r.communityCards = append(r.communityCards, riverCard)
	for _, p := range r.players {
		if p == nil || p.Status() == player.StatusFolded {
			continue
		}
		if err = r.broadcaster.Action(dealer.EventRiverCard, dealer.EventObject{RiverCard: riverCard}); err != nil {
			return fmt.Errorf("broadcaster action, err: %w", err)
		}
		if p.Status() == player.StatusAllIn {
			continue
		}
		action := p.WaitForAction(ctx, []player.Action{
			{Type: player.ActionFold, Chips: 0},
			{Type: player.ActionCall, Chips: preflopBet},
			{Type: player.ActionRaise, Chips: preflopBet * 2},
			// {Type: player.ActionAllIn, Chips: 0},
		})
		if err = r.broadcaster.Action(action.Type.IntoEventType(), player.EventObject{Player: p, Bet: action.Chips}); err != nil {
			return fmt.Errorf("broadcaster action, err: %w", err)
		}
	}

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
