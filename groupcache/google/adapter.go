// Package google implements an adapter to extract metrics from google groupcache.
package google

import (
	"github.com/golang/groupcache"
	"github.com/udhos/groupcache_exporter"
)

// Group implements interface GroupStatistics to extract metrics from groupcache group.
type Group struct {
	group *groupcache.Group
}

// New creates a new Group.
func New(group *groupcache.Group) *Group {
	return &Group{group: group}
}

// Collect requests metrics collection from implementation
func (g *Group) Collect() groupcache_exporter.Stats {

	var result groupcache_exporter.Stats

	stats := g.group.Stats

	result.Group = groupcache_exporter.GroupStats{
		CounterGets:           stats.Gets.Get(),
		CounterHits:           stats.CacheHits.Get(),
		CounterPeerLoads:      stats.PeerLoads.Get(),
		CounterPeerErrors:     stats.PeerErrors.Get(),
		CounterLoads:          stats.Loads.Get(),
		CounterLoadsDeduped:   stats.LoadsDeduped.Get(),
		CounterLocalLoads:     stats.LocalLoads.Get(),
		CounterLocalLoadsErrs: stats.LocalLoadErrs.Get(),
		CounterServerRequests: stats.ServerRequests.Get(),
	}

	result.Main = getCacheStats(g.group.CacheStats(groupcache.MainCache))
	result.Hot = getCacheStats(g.group.CacheStats(groupcache.HotCache))

	return result
}

func getCacheStats(cacheStats groupcache.CacheStats) groupcache_exporter.CacheTypeStats {
	return groupcache_exporter.CacheTypeStats{
		GaugeCacheItems:       cacheStats.Items,
		GaugeCacheBytes:       cacheStats.Bytes,
		CounterCacheGets:      cacheStats.Gets,
		CounterCacheHits:      cacheStats.Hits,
		CounterCacheEvictions: cacheStats.Evictions,
	}
}

// Name returns the group's name
func (g *Group) Name() string {
	return g.group.Name()
}
