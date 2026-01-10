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

	MinPlayerCount = 2
	MaxPlayerCount = 22
)

type playerCount struct {
	max, min, current int
}

type Round struct {
	number int

	// position indicates the position relationship among players in round.
	position []string

	// players only includes players in round.
	players map[string]*player.Player

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

// New init a new round, button must greater than or equal to zero.
func New(players []*player.Player, button int, opts ...Option) *Round {
	if button <= 0 {
		panic(errors.New("MUST button >= 0"))
	}
	counter := button
	position := make([]string, 0)
	playersMap := make(map[string]*player.Player)

	// TODO(@yshngg):
	for i, p := range players {
		if counter != 0 {
			counter--
		}
		if p == nil {
			continue
		}
		position = append(position, p.ID())
		playersMap[p.ID()] = p
	}
	r := &Round{
		position: position,
		players:  playersMap,
		button:   button,
		minBet:   -1,
		status:   StatusReady,
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
	if r.playerCount.max > MaxPlayerCount {
		r.playerCount.max = MaxPlayerCount
	}
	if r.playerCount.max < MinPlayerCount {
		r.playerCount.min = MinPlayerCount
	}
	return r
}

type Option func(*Round)

func WithNumber(number int) Option {
	return func(r *Round) {
		r.number = number
	}
}

// func WithButton(button int) Option {
// 	return func(r *Round) {
// 		r.button = button
// 	}
// }

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

// func WithPlayers(players ...*player.Player) Option {
// 	return func(r *Round) {
// 		r.players = players
// 	}
// }

// func WithPlayerCount(min, max, current int) Option {
// 	return func(r *Round) {
// 		r.playerCount = playerCount{
// 			max:     max,
// 			min:     min,
// 			current: current,
// 		}
// 	}
// }

type ErrMaxPlayerCountReached struct{}

func (e ErrMaxPlayerCountReached) Error() string {
	return "max player count reached"
}

type ErrPlayerAlreadyExists struct{}

func (e ErrPlayerAlreadyExists) Error() string {
	return "player already exists"
}

func (r Round) CountPlayer() int {
	count := 0
	for _, p := range r.players {
		if p != nil && p.Status() != player.StatusIdle {
			count++
		}
	}
	return count
}

// func (r Round) ExistsPlayer(id string) bool {
// 	return slices.ContainsFunc(r.players, func(p *player.Player) bool {
// 		return p.ID() == id
// 	})
// }

type ErrRoundAlreadyStarted struct{}

func (e ErrRoundAlreadyStarted) Error() string {
	return "round has already started"
}

// func (r *Round) AddPlayer(p *player.Player) error {
// 	if r.ExistsPlayer(p.ID()) {
// 		return ErrPlayerAlreadyExists{}
// 	}
// 	if r.status.After(StatusStarted) {
// 		return ErrRoundAlreadyStarted{}
// 	}
// 	if r.CountPlayer() >= r.playerCount.max {
// 		return ErrMaxPlayerCountReached{}
// 	}
// 	for i, pp := range r.players {
// 		if pp == nil {
// 			r.players[i] = p
// 			return nil
// 		}
// 	}
// 	r.players = append(r.players, p)
// 	return nil
// }

type ErrPlayerNotFound struct {
	id string
}

// func (r Round) FindPlayer(id string) (*player.Player, error) {
// 	for _, p := range r.players {
// 		if p.ID() == id {
// 			return p, nil
// 		}
// 	}
// 	return nil, ErrPlayerNotFound{id}
// }

func (e ErrPlayerNotFound) Error() string {
	return fmt.Sprintf("player (id: %s) not found", e.id)
}

func (r *Round) Watch() (watch.Interface, error) {
	return r.broadcaster.Watch()
}

// func (r *Round) RemovePlayer(ctx context.Context, id string) error {
// 	p, err := r.FindPlayer(id)
// 	if err != nil {
// 		return ErrPlayerNotFound{id}
// 	}

// 	if r.status.Before(StatusStarted) {
// 		r.players = slices.DeleteFunc(r.players, func(pp *player.Player) bool {
// 			return pp.ID() == id
// 		})
// 		return nil
// 	}

// 	p.StopWatch()

// 	if p.Status() == player.StatusWaitingToAct {
// 		select {
// 		case <-p.Active():
// 			err := p.Fold(ctx)
// 			if err != nil {
// 				return fmt.Errorf("player %s (id: %s) fold, err: %w", p.Name(), p.ID(), err)
// 			}
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		}
// 	}
// 	p.Close()
// 	return nil
// }

// func (r *Round) prepare(ctx context.Context) error {
// 	for _, p := range r.players {
// 		if p == nil {
// 			continue
// 		}
// 		watcher, err := r.broadcaster.Watch()
// 		if err != nil {
// 			return fmt.Errorf("failed to watch broadcaster: %w", err)
// 		}
// 		watcher = watch.Filter(watcher, func(in watch.Event) (out watch.Event, keep bool) {
// 			// TODO(@yshngg): filter showdown event
// 			// TODO(@yshngg): hiden burned card (dealer)
// 			// hole cards are only visable to player own them
// 			if in.Kind() == dealer.EventKind {
// 				dealerEvent := in.(dealer.Event)
// 				dealerEventObject := dealerEvent.Related().(dealer.EventObject)
// 				if dealerEvent.Action() != dealer.EventDealHoleCards || dealerEventObject.To == dealer.ToPlayer(p) {
// 					return in, true
// 				}
// 				// zero other player's hole cards
// 				cards := make([]*card.Card, len(dealerEventObject.Cards))
// 				dealer.NewEvent(dealerEvent.Action(), dealerEventObject.To, cards...)
// 			}
// 			return in, true
// 		})
// 		p.Apply(player.WithWatcher(watcher))
// 	}

// 	return nil
// }

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
	small, big, err := r.positionBlind()
	if err != nil {
		return fmt.Errorf("blind positions, err: %v", err)
	}

	if err := r.players[small].Bet(ctx, r.minBet/2); err != nil {
		return fmt.Errorf("post small blind: %w", err)
	}
	if err := r.players[big].Bet(ctx, r.minBet); err != nil {
		return fmt.Errorf("post big blind: %w", err)
	}

	r.pots.AddChips(r.players[small].ID(), r.minBet/2)
	r.pots.AddChips(r.players[big].ID(), r.minBet)

	smallBlindEvent := player.NewEvent(player.EventPostSmallBlind, player.EventObject{
		ID:  r.players[small].ID(),
		Bet: r.minBet / 2,
	})
	if err := r.broadcaster.Action(smallBlindEvent); err != nil {
		return fmt.Errorf("broadcast event: %s, err: %w", player.EventPostSmallBlind, err)
	}
	bigBlindEvent := player.NewEvent(player.EventPostSmallBlind, player.EventObject{
		ID:  r.players[big].ID(),
		Bet: r.minBet,
	})
	if err := r.broadcaster.Action(bigBlindEvent); err != nil {
		return fmt.Errorf("broadcast event: %s, err: %w", player.EventPostBigBlind, err)
	}
	return nil
}

func (r *Round) Status() StatusType {
	return r.status
}

// TODO(@yshngg): broadcast action taken by the player
func (r *Round) openBettingRound(ctx context.Context) (err error) {
	switch r.status {
	case StatusPreFlop, StatusFlop, StatusTurn, StatusRiver:
	default:
		return fmt.Errorf("invalid status: %s", r.status)
	}

	maxBet := 0
	if r.status == StatusPreFlop {
		maxBet = r.minBet
	}
	minRaise := r.minBet
	bets := make(map[string]int) // chips that have bet in current betting round
	small, big, err := r.positionBlind()
	if err != nil {
		return fmt.Errorf("position blind, err: %w", err)
	}
	bets[r.players[small].ID()] = r.minBet / 2
	bets[r.players[big].ID()] = r.minBet
	acted := make(map[string]bool) // map player id to whether they have acted in the current betting round
	for _, p := range r.players {
		acted[p.ID()] = false
	}
	start, err := r.positionFirstToAct()
	if err != nil {
		return fmt.Errorf("position first to act, err: %w", err)
	}

	// TODO(@yshngg): player don't need to take any action after the player allin
	next := func() bool {
		for _, ok := range acted {
			if !ok {
				return true
			}
		}
		for id, chips := range bets {
			p, err := r.FindPlayer(id)
			if err != nil {
				// TODO(@yshngg): log error, but don't return or panic
				continue
			}
			if p.Status() != player.StatusWaitingToAct {
				continue
			}
			if chips != maxBet {
				return true
			}
		}
		return false
	}

	i := 0
	for next() { // reopen the betting action
		p := r.players[(start+i)%len(r.players)]
		if p == nil {
			// TODO(@yshngg): log error, but don't return or panic
			// fmt.Errorf("player %s (id: %s) is not waiting to act (status: %s)", p.Name(), p.ID(), p.Status())
			continue
		}
		if p.Status() == player.StatusFolded || p.Status() == player.StatusAllIn {
			acted[p.ID()] = true
			continue
		}
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
					Chips: (maxBet - bets[p.ID()]),
				},
				{
					Type:  player.ActionRaise,
					Chips: (maxBet - bets[p.ID()]) + minRaise,
				},
			}
		}
		action, err := p.WaitForAction(ctx, availableActions)
		if err != nil {
			return fmt.Errorf("wait for action, err: %w", err)
		}

		switch action.Type {
		case player.ActionAllIn, player.ActionRaise, player.ActionBet, player.ActionCall:
			r.pots.AddChips(p.ID(), action.Chips)
			bets[p.ID()] += action.Chips
		}

		switch action.Type {
		case player.ActionAllIn:
			if bets[p.ID()] > maxBet {
				maxBet = bets[p.ID()]
			}
		case player.ActionRaise:
			minRaise = bets[p.ID()] - maxBet
			maxBet = bets[p.ID()]
		case player.ActionBet:
			maxBet = action.Chips
			minRaise = maxBet
		}
		acted[p.ID()] = true
		i++
	}

	return nil
}

type ErrBroadcast struct {
	Event watch.Event
}

func (e ErrBroadcast) Error() string {
	return fmt.Sprintf("broadcast event %s", e.Event)
}

func (r *Round) Start(ctx context.Context) error {
	if r.players[r.button] == nil {
		return fmt.Errorf("button position does not have player")
	}
	playerCount := r.CountPlayer()
	if playerCount < r.playerCount.min || playerCount > r.playerCount.max {
		return fmt.Errorf("invalid player count: %d", playerCount)
	}
	r.playerCount.current = playerCount

	// ready to start the round
	r.status = StatusStarted

	// prepare players
	err := r.prepare(ctx)
	if err != nil {
		return fmt.Errorf("start round, err: %w", err)
	}

	roundStartEvent := NewEvent(EventStart, r.players, nil)
	if err := r.broadcaster.Action(roundStartEvent); err != nil {
		return fmt.Errorf("broadcast event: %v, err: %w", roundStartEvent, err)
	}

	// dealer shuffle deck
	r.dealer.Shuffle()
	dealerShuffleEvent := dealer.NewEvent(dealer.EventShuffle, dealer.ToAll())
	if err := r.broadcaster.Action(dealerShuffleEvent); err != nil {
		return fmt.Errorf("broadcast event: %v, err: %w", dealerShuffleEvent, err)
	}

	// compulsory bets
	if err := r.betBlind(ctx); err != nil {
		return fmt.Errorf("bet blind, err: %w", err)
	}

	// pre-flop
	r.status = StatusPreFlop
	holeCards := r.dealer.DealHoleCards(playerCount)
	for i := range len(r.players) {
		p := r.players[(r.button+i+1)%len(r.players)]
		if p == nil {
			continue
		}
		p.SetHoleCards(holeCards[i])
		dealHoleCardsEvent := dealer.NewEvent(dealer.EventDealHoleCards, dealer.ToPlayer(p), holeCards[i][:]...)
		if err := r.broadcaster.Action(dealHoleCardsEvent); err != nil {
			return fmt.Errorf("broadcast event: %v, err: %w", dealHoleCardsEvent, err)
		}
	}
	// pre-flop betting round
	err = r.openBettingRound(ctx)
	if err != nil {
		return fmt.Errorf("open betting round: err: %w", err)
	}

	// flop
	r.status = StatusFlop
	burnCard := r.dealer.BurnCard()
	burnCardEvent := dealer.NewEvent(dealer.EventDealFlopCards, dealer.ToCommunity(), burnCard)
	if err := r.broadcaster.Action(burnCardEvent); err != nil {
		return fmt.Errorf("broadcast event: %v, err: %w", burnCardEvent, err)
	}
	flopCards := r.dealer.DealFlopCards()
	for _, card := range flopCards {
		r.communityCards = append(r.communityCards, card)
	}
	dealFlopCardsEvent := dealer.NewEvent(dealer.EventDealFlopCards, dealer.ToCommunity(), flopCards[:]...)
	if err := r.broadcaster.Action(dealFlopCardsEvent); err != nil {
		return fmt.Errorf("broadcast event: %v, err: %w", dealFlopCardsEvent, err)
	}

	// flop betting round
	err = r.openBettingRound(ctx)
	if err != nil {
		return fmt.Errorf("open betting round: err: %w", err)
	}

	// turn
	r.status = StatusTurn
	burnCard = r.dealer.BurnCard()
	burnCardEvent = dealer.NewEvent(dealer.EventDealFlopCards, dealer.ToCommunity(), burnCard)
	if err := r.broadcaster.Action(burnCardEvent); err != nil {
		return fmt.Errorf("broadcast event: %v, err: %w", burnCardEvent, err)
	}
	turnCard := r.dealer.DealTurnCard()
	r.communityCards = append(r.communityCards, turnCard)
	dealTurnCardEvent := dealer.NewEvent(dealer.EventDealTurnCard, dealer.ToCommunity(), turnCard)
	if err := r.broadcaster.Action(dealTurnCardEvent); err != nil {
		return fmt.Errorf("broadcast event: %v, err: %w", dealTurnCardEvent, err)
	}

	// turn betting round
	err = r.openBettingRound(ctx)
	if err != nil {
		return fmt.Errorf("open betting round: err: %w", err)
	}

	// river
	r.status = StatusRiver
	burnCard = r.dealer.BurnCard()
	burnCardEvent = dealer.NewEvent(dealer.EventDealFlopCards, dealer.ToCommunity(), burnCard)
	if err := r.broadcaster.Action(burnCardEvent); err != nil {
		return fmt.Errorf("broadcast event: %v, err: %w", burnCardEvent, err)
	}
	riverCard := r.dealer.DealTurnCard()
	r.communityCards = append(r.communityCards, riverCard)
	dealRiverCardEvent := dealer.NewEvent(dealer.EventDealRiverCard, dealer.ToCommunity(), riverCard)
	if err := r.broadcaster.Action(dealRiverCardEvent); err != nil {
		return fmt.Errorf("broadcast event: %v, err: %w", dealRiverCardEvent, err)
	}

	// turn betting round
	err = r.openBettingRound(ctx)
	if err != nil {
		return fmt.Errorf("open betting round: err: %w", err)
	}

	// showdown
	r.status = StatusShowdown
	// TODO(@yshngg): showdown???
	// cards, err := r.Showdown(ctx)

	for _, pot := range r.pots.Settle() {
		players := make([]*player.Player, len(r.players))
		for id := range pot.Contributors() {
			p := new(player.Player)
			if p, err = r.FindPlayer(id); err != nil {
				return fmt.Errorf("find player (id: %s), err: %w", id, err)
			}
			if p.Status() != player.StatusTakingAction && p.Status() != player.StatusAllIn {
				continue
			}
			players = append(players, p)
			// TODO(@yshngg): compare player's best five cards and win the pot
		}
	}

	roundShowdownEvent := NewEvent(EventShowdown, players, r.communityCards...)
	if err := r.broadcaster.Action(roundShowdownEvent); err != nil {
		return fmt.Errorf("broadcast event: %v, err: %w", roundShowdownEvent, err)
	}

	return nil
}

func (r *Round) End() error {
	r.status = StatusEnd
	for _, p := range r.players {
		p.StopWatch()
		p.Reset()
	}
	return nil
}

func (r *Round) Showdown(ctx context.Context) (map[string][2]*card.Card, error) {
	holeCards := make(map[string][2]*card.Card, 0)
	for _, p := range r.players {
		if p.Status() == player.StatusFolded {
			continue
		}
		holeCards[p.ID()] = p.HoleCards()
	}

	if len(holeCards) < 1 {
		return nil, fmt.Errorf("game has long been over")
	}
	if len(holeCards) > 1 {
		return holeCards, nil
	}
	var id string
	for id = range holeCards {
	}
	p, err := r.FindPlayer(id)
	if err != nil {
		return nil, fmt.Errorf("find player, err: %w", err)
	}
	action, err := p.WaitForAction(ctx, []player.Action{
		player.Action{Type: player.ActionHideHoleCards},
		player.Action{Type: player.ActionShowHoleCards},
	})
	if err != nil {
		return nil, fmt.Errorf("wait for action, err: %w", err)
	}
	if action.Type != player.ActionShowHoleCards {
		// zero hole cards
		cards := holeCards[id]
		for i := range len(cards) {
			cards[i] = nil
		}
		holeCards[id] = cards
	}
	return holeCards, nil
}

func (r *Round) Players() []*player.Player {
	players := make([]*player.Player, 0, len(r.players))
	for _, p := range r.players {
		if p != nil || p.Status() != player.StatusIdle {
			players = append(players, p)
		}
	}
	return players
}

type ErrFirstToActPlayerNotFound struct {
	button int
}

func (e ErrFirstToActPlayerNotFound) Error() string {
	return fmt.Sprintf("no first player found at button %d", e.button)
}

type ErrStatusNotSupported struct {
	status StatusType
}

func (e ErrStatusNotSupported) Error() string {
	return fmt.Sprintf("status not supported: %s", e.status)
}

func (r Round) positionBlind() (int, int, error) {
	playerCount := r.playerCount.current
	length := len(r.players)
	if r.button < 0 || length <= r.button || r.players[r.button] == nil {
		return -1, -1, ErrInvalidButton{button: r.button}
	}
	if playerCount < 2 {
		return -1, -1, ErrInvalidPlayerCount{count: playerCount}
	}
	if playerCount == 2 {
		small := r.button
		big := (small + 1) % length
		for r.players[big] == nil {
			big = (big + 1) % length
		}
		return small, big, nil
	}

	small := (r.button + 1) % length
	big := (small + 1) % length
	for range length {
		if r.players[small] == nil {
			small = (small + 1) % length
			big = (small + 1) % length
			continue
		}
		if r.players[big] != nil {
			break
		}
		big = (big + 1) % length
	}
	return small, big, nil
}

func (r *Round) positionFirstToAct() (int, error) {
	playerCount := r.playerCount.current
	if playerCount < 2 {
		return -1, ErrInvalidPlayerCount{count: playerCount}
	}

	small, big, err := r.positionBlind()
	if err != nil {
		return -1, fmt.Errorf("position blind, err: %w", err)
	}

	switch r.status {
	case StatusPreFlop:
		if playerCount == 2 {
			return small, nil
		}
		length := len(r.players)
		for i := range length - 2 {
			p := r.players[(big+i+1)%length]
			if p != nil {
				return i, nil
			}
		}
		return -1, ErrFirstToActPlayerNotFound{r.button}
	case StatusFlop, StatusTurn, StatusRiver:
		if playerCount == 2 {
			return big, nil
		}
		return small, nil
	default:
		return -1, ErrStatusNotSupported{status: r.status}
	}
}
