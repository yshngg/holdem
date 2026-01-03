package round

import (
	"context"
	"errors"
	"fmt"

	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/dealer"
	"github.com/yshngg/holdem/pkg/player"
	pots "github.com/yshngg/holdem/pkg/pot"
	"github.com/yshngg/holdem/pkg/watch"
)

const (
	defaultMinBet = 2
	defaultButton = 0

	minPlayerCount = 2
	maxPlayerCount = 22
)

type playerCount struct {
	max, min, current int
}

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
	pots pots.Pots

	// recorder captures all game events for replay, debugging, or auditing purposes.
	// It logs actions like bets, folds, and card deals.
	recorder watch.Recorder

	// broadcaster delivers real-time game events to all connected players.
	// Ensures players receive synchronized updates about round state changes.
	broadcaster watch.Broadcaster

	// status tracks the current stage of the poker round.
	// See the Status type for all possible values (pre-flop, flop, turn, river, etc.).
	status StatusType

	playerCount playerCount
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
	if r.playerCount.max > maxPlayerCount {
		r.playerCount.max = maxPlayerCount
	}
	if r.playerCount.max < minPlayerCount {
		r.playerCount.min = minPlayerCount
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

func WithPlayers(players ...*player.Player) Option {
	return func(r *Round) {
		r.players = players
	}
}

func WithPlayerCount(min, max, current int) Option {
	return func(r *Round) {
		r.playerCount = playerCount{
			max:     max,
			min:     min,
			current: current,
		}
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

type ErrMaxPlayerCountReached struct{}

func (e ErrMaxPlayerCountReached) Error() string {
	return "max player count reached"
}

type ErrPlayerAlreadyExists struct{}

func (e ErrPlayerAlreadyExists) Error() string {
	return "player already exists"
}

func (r *Round) AddPlayer(p *player.Player) error {
	if realPlayerCount(r.players) >= r.playerCount.max {
		return ErrMaxPlayerCountReached{}
	}
	if existsPlayer(r.players, p) {
		return ErrPlayerAlreadyExists{}
	}
	if r.status.After(StatusStarted) && p.Status() != player.StatusSpectating {
		return errors.New("round has started, player can only spectate")
	}
	for _, pp := range r.players {
		if pp == nil {
			pp = p
			return nil
		}
	}
	r.players = append(r.players, p)
	return nil
}

type ErrPlayerNotFound struct{}

func (e ErrPlayerNotFound) Error() string {
	return "player not found"
}

func (r *Round) RemovePlayer(ctx context.Context, p *player.Player) error {
	if r.status.After(StatusReady) && p.Status() != player.StatusSpectating {
		return errors.New("round has started")
	}

	for _, pp := range r.players {
		if pp.ID() == p.ID() {
			pp.Leave(ctx)
			pp = nil
			return nil
		}
	}

	return ErrPlayerNotFound{}
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

	r.pots.AddChips(r.players[small].ID().String(), r.minBet/2)
	r.pots.AddChips(r.players[big].ID().String(), r.minBet)

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

func (r *Round) openBettingRound(ctx context.Context) (err error) {
	maxBet := 0
	minRaise := r.minBet
	betChips := make(map[string]int) // chips that have bet in current betting round
	start := -1
	if r.Status() == StatusPreFlop {
		maxBet = r.minBet
		start, err = positionUTG(r.players, r.button)
		if err != nil {
			return fmt.Errorf("betting round, err: %w", err)
		}
	} else {
		maxBet = 0
		start, err = positionFirstToAct(r.players, r.button)
		if err != nil {
			return fmt.Errorf("betting round, err: %w", err)
		}
	}
	if start < 0 {
		return fmt.Errorf("can not find first player to act")
	}

	keep := func() bool {
		if effectivePlayerCount(r.players) != len(betChips) {
			return true
		}

		chipsList := make([]int, 0, len(betChips))
		allInList := make([]int, 0, len(betChips))
		for pid, chips := range betChips {
			if chips < 0 {
				continue
			}
			p, err := findPlayerByID(r.players, pid)
			if err != nil {
				// TODO(@yshngg): log error, but don't return or panic
			}
			if p != nil && p.Status() == player.StatusWaitingToAct {
				chipsList = append(chipsList, chips)
			}
			if p != nil && p.Status() == player.StatusAllIn {
				allInList = append(allInList, chips)
			}
		}

		// correct?
		if len(chipsList) < 1 {
			return true
		}

		first := chipsList[0]
		for _, chips := range chipsList[1:] {
			if chips != first {
				return true
			}
		}

		return false
	}

	// keep := func() bool {
	// 	highestBetPosition := 0
	// 	for i, p := range r.players {
	// 		if p == nil || p.Status() != player.StatusWaitingToAct {
	// 			continue
	// 		}
	// 		bet, acted := playerBetChips[p.ID().String()]
	// 		if !acted || bet < 0 {
	// 			continue
	// 		}
	// 		if bet > highestBetChips {
	// 			highestBetChips = bet
	// 			highestBetPosition = i
	// 		}
	// 	}

	// 	for i := range len(r.players) - 1 {
	// 		p := r.players[(highestBetPosition+i+1)%len(r.players)]
	// 		if p == nil || p.Status() != player.StatusWaitingToAct {
	// 			continue
	// 		}
	// 		bet, acted := playerBetChips[p.ID().String()]
	// 		if bet < 0 {
	// 			continue
	// 		}
	// 		if !acted || bet < highestBetChips {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }

	p := r.players[(start)%len(r.players)]

	if p == nil || p.Status() != player.StatusWaitingToAct {
		return fmt.Errorf("err: %w", err)
	}
	haveBet := betChips[p.ID().String()]
	availableActions := []player.Action{
		{
			Type: player.ActionFold,
		},
		{
			Type: player.ActionAllIn,
		},
	}
	if maxBet == 0 {
		availableActions = append(availableActions, []player.Action{
			{
				Type:  player.ActionBet,
				Chips: r.minBet,
			},
			{
				Type: player.ActionCheck,
			},
		}...)
	} else {
		availableActions = []player.Action{
			{
				Type:  player.ActionCall,
				Chips: minRaise,
			},
			{
				Type:  player.ActionRaise,
				Chips: minRaise + (maxBet - haveBet),
			},
		}
	}
	action := p.WaitForAction(ctx, availableActions)
	r.pots.AddChips(p.ID().String(), action.Chips)

	if action.Type == player.ActionRaise || action.Type == player.ActionAllIn {
		minRaise = action.Chips - maxBet
		maxBet = action.Chips + haveBet
	}
	if action.Type == player.ActionBet {
		maxBet = action.Chips
		minRaise = action.Chips
	}

	for keep() { // reopen the betting action
		for i := range len(r.players) - 1 {
			p := r.players[(start+i+1)%len(r.players)]
			if p == nil || p.Status() != player.StatusWaitingToAct {
				continue
			}
			// bet, acted := playerBetChips[p.ID().String()]
			// if acted && bet < 0 {
			// 	continue
			// }
			availableActions := []player.Action{
				{
					Type: player.ActionFold,
				},
				{
					Type: player.ActionAllIn,
				},
			}
			if maxBet == 0 {
				availableActions = append(availableActions, []player.Action{
					{
						Type:  player.ActionBet,
						Chips: r.minBet,
					},
					{
						Type: player.ActionCheck,
					},
				}...)
			} else {
				availableActions = []player.Action{
					{
						Type:  player.ActionCall,
						Chips: minRaise,
					},
					{
						Type: player.ActionRaise,
						// Chips: (highestBetChips - bet) * 2,
						Chips: minRaise + (maxBet - haveBet),
					},
				}
			}
			action := p.WaitForAction(ctx, availableActions)
			r.pots.AddChips(p.ID().String(), action.Chips)

			if action.Type == player.ActionRaise || action.Type == player.ActionAllIn {
				minRaise = action.Chips - maxBet
				start += i
				break
			}
			if action.Type == player.ActionBet {
				maxBet = action.Chips
				minRaise = action.Chips
			}
		}
	}

	return nil
}

func (r *Round) Start(ctx context.Context) error {
	if r.players[r.button] == nil {
		return fmt.Errorf("button position does not have player")
	}
	realPlayerCount := realPlayerCount(r.players)
	effectivePlayerCount := effectivePlayerCount(r.players)
	if effectivePlayerCount < r.playerCount.min || effectivePlayerCount > r.playerCount.max {
		return fmt.Errorf("invalid player count: %d", effectivePlayerCount)
	}

	// ready to start the round
	r.status = StatusStarted

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

func (r *Round) Showdown(show bool) map[string][2]*card.Card {
	holeCards := make(map[string][2]*card.Card, 0)
	for _, p := range r.players {
		if p.Status() == player.StatusFolded {
			continue
		}
		holeCards[p.ID().String()] = p.HoleCards()
	}
	if len(holeCards) == 1 && !show {
		return nil
	}
	return holeCards
}

func (r *Round) RealPlayers() []*player.Player {
	players := make([]*player.Player, 0)
	for _, p := range r.players {
		if p != nil || p.Status() != player.StatusLeft {
			players = append(players, p)
		}
	}
	return players
}
