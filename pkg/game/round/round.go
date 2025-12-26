package round

import (
	"fmt"

	"github.com/yshngg/holdem/pkg/dealer"
	"github.com/yshngg/holdem/pkg/deck"
	"github.com/yshngg/holdem/pkg/player"
	"github.com/yshngg/holdem/pkg/watch"
)

const (
	MinPlayerCount = 2
	MaxPlayerCount = 10
)

type Round struct {
	players []*player.Player
	dealer  *dealer.Dealer
	button  int

	watcher     watch.Interface
	broadcaster *watch.Broadcaster
}

func New(players []*player.Player, button int) *Round {
	d := dealer.New(deck.New())
	queueLength := len(players) * 2
	return &Round{
		players:     players,
		dealer:      d,
		button:      button,
		broadcaster: watch.NewBroadcaster(queueLength, queueLength),
	}
}

func (r *Round) Prepare() error {
	r.dealer.Shuffle()
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

	if err := r.players[small].PostSmallBlind(); err != nil {
		return fmt.Errorf("failed to post small blind: %w", err)
	}
	if err := r.players[big].PostBigBlind(); err != nil {
		return fmt.Errorf("failed to post big blind: %w", err)
	}
	return nil
}

func (r *Round) Start() error {
	if r.players[r.button] == nil {
		return fmt.Errorf("button position is empty")
	}
	playerCount := r.playerCount()
	if playerCount < MinPlayerCount || playerCount > MaxPlayerCount {
		return fmt.Errorf("invalid player count: %d", playerCount)
	}

	// compulsory bets
	if err := r.betBlind(); err != nil {
		return err
	}

	// pre-flop
	cards := r.dealer.DealHoleCards(playerCount)
	for i, p := range r.effectivePlayers() {
		// p.SetHoleCards(cards[(r.button+i)%len(cards)])
	}
	for i := range playerCount {
		p := r.players[i]
		if p == nil {
			continue
		}
		r.players[i].SetHoleCards(cards[i])
	}

	return nil
}

func (r *Round) End() error {
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
