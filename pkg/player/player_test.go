package player

import (
	"testing"

	"github.com/yshngg/holdem/pkg/watch"
)

func TestWatch(t *testing.T) {
	broadcaster := watch.NewBroadcaster(0, 0)
	watcher, err := broadcaster.Watch()
	if err != nil {
		t.Fatalf("watch broadcaster, err: %v", err)
	}
	player := New(WithWatcher(watcher))
	events := []watch.Event{
		{Type: EventCheck, Object: nil},
		{Type: EventFold, Object: nil},
		{Type: EventBet, Object: nil},
		{Type: EventCall, Object: nil},
		{Type: EventRaise, Object: nil},
		{Type: EventAllIn, Object: nil},
	}

	errChan := make(chan error, 0)
	go func() {
		defer close(errChan)

		// player.watcher.Stop()
		defer player.Done()
		for _, event := range events {
			err := broadcaster.Action(event.Type, event.Type)
			errChan <- err
		}
	}()

	for err := range errChan {
		if err != nil {
			t.Errorf("broadcaster action, err: %v", err)
		}
	}

	for _, want := range events {
		got := <-player.Watch()
		if got.Type != want.Type {
			t.Errorf("event type: %v, want: %v", got.Type, want.Type)
		}
	}
}

func TestAction(t *testing.T) {
	player := New()
	errChan := make(chan error, 1)
	go func(p *Player) {
		err := player.Check(t.Context())
		errChan <- err
	}(player)
	action := player.WaitForAction(t.Context(), map[ActionType]Action{
		ActionCheck: {ActionCheck, 0},
		ActionBet:   {ActionBet, 2},
	})
	if err := <-errChan; err != nil {
		t.Fatalf("player check action error: %v", err)
	}
	if action.Type != ActionCheck {
		t.Fatalf("action type: %v, want: %v", action.Type, ActionCheck)
	}
}
