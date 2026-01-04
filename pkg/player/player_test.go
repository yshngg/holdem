package player

import (
	"testing"

	"github.com/yshngg/holdem/pkg/watch"
	"golang.org/x/sync/errgroup"
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

	g := new(errgroup.Group)
	g.Go(func() error {
		for _, event := range events {
			err := broadcaster.Action(event.Type, event.Type)
			if err != nil {
				return err
			}
		}
		return nil
	})

	for _, want := range events {
		got := <-player.Watch()
		if got.Type != want.Type {
			t.Errorf("event type: %v, want: %v", got.Type, want.Type)
		}
	}

	if err = g.Wait(); err != nil {
		t.Fatalf("broadcaster action, err: %v", err)
	}
}

func TestAction(t *testing.T) {
	player := New()
	g := errgroup.Group{}
	g.Go(func() error {
		err := player.Check(t.Context())
		return err
	})
	action := player.WaitForAction(t.Context(), []Action{
		{ActionCheck, 0},
		{ActionBet, 2},
	})

	if err := g.Wait(); err != nil {
		t.Fatalf("player check action error: %v", err)
	}
	if action.Type != ActionCheck {
		t.Fatalf("action type: %v, want: %v", action.Type, ActionCheck)
	}
}
