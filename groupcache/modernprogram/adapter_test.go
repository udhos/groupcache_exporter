package modernprogram

import (
	"context"
	"sync"
	"testing"

	"github.com/modernprogram/groupcache/v2"
)

// go test -race -count 1 -run '^TestRaceModernProgram$' ./...
func TestRaceModernProgram(_ *testing.T) {

	options := groupcache.Options{
		Workspace: groupcache.DefaultWorkspace,
		Name:      "group1",
		Getter: groupcache.GetterFunc(
			func(_ /*ctx*/ context.Context, _ string, _ groupcache.Sink, _ *groupcache.Info) error {
				return nil
			}),
	}

	g := groupcache.NewGroupWithWorkspace(options)

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
