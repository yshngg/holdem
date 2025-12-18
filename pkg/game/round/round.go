package round

import (
	"github.com/yshngg/holdem/pkg/dealer"
	"github.com/yshngg/holdem/pkg/player"
)

type Round struct {
	players []*player.Player
	dealer  dealer.Dealer
	deck    []int
	button  int
}

func New() *Round {
	return &Round{}
}

func (r *Round) Prepare() error {
	return nil
}

func (r *Round) Start() error {
	return nil
}

func (r *Round) End() error {
	return nil
}
