package google

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/groupcache"
)

// go test -race -count 1 -run '^TestRaceGoogle$' ./...
func TestRaceGoogle(_ *testing.T) {

	g := groupcache.NewGroup("files", 1_000_000, groupcache.GetterFunc(
		func(_ /*ctx*/ context.Context, _ string, _ groupcache.Sink) error {
			return nil
		}),
	)

	eg := &exportGroup{group: g}

	var wg sync.WaitGroup

	const n = 10000

	wg.Add(2 * n)

	for range n {
		go func() {
			eg.group.Stats.Gets.Add(1)
			wg.Done()
		}()
	}

	for range n {
		go func() {
			eg.Collect()
			wg.Done()
		}()
	}

	wg.Wait()
}
