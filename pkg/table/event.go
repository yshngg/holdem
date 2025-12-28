package table

import (
	"time"

	"github.com/yshngg/holdem/pkg/player"
	"github.com/yshngg/holdem/pkg/watch"
)

type EventType int

const (
	PlayerJoin  = watch.TablePlayerJoin
	PlayerLeave = watch.TablePlayerLeave
)

type EventObject struct {
	player    *player.Player
	timestamp time.Time
}
