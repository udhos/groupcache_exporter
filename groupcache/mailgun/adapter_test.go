package mailgun

import (
	"context"
	"sync"
	"testing"

	"github.com/mailgun/groupcache/v2"
)

// go test -race -count 1 -run '^TestRaceMailgun$' ./...
func TestRaceMailgun(_ *testing.T) {

	g := groupcache.NewGroup("files", 1_000_000, groupcache.GetterFunc(
		func(_ /*ctx*/ context.Context, _ string, _ groupcache.Sink) error {
			return nil
		}),
	)

	eg := &exportGroup{group: g}

	var wg sync.WaitGroup

	const n = 10000

	for range n {
		go func() {
			eg.group.Stats.Gets.Add(1)
		}()
	}

	for range n {
		go func() {
			eg.Collect()
		}()
	}

	wg.Wait()
}
