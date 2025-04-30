// Package google implements an adapter to extract metrics from google groupcache.
package google

import (
	"github.com/golang/groupcache"
	"github.com/udhos/groupcache_exporter"
)

// ListGroups issues an static list of exporter groups for groupcache groups.
// Other groupcache implementations are able to issue current list, but this is not.
func ListGroups(groups []*groupcache.Group) []groupcache_exporter.GroupStatistics {
	var exportGroups []groupcache_exporter.GroupStatistics
	for _, g := range groups {
		exportGroups = append(exportGroups, &exportGroup{group: g})
	}
	return exportGroups
}

// exportGroup implements interface GroupStatistics to extract metrics from groupcache group.
type exportGroup struct {
	group *groupcache.Group
}

// Collect requests metrics collection from implementation
func (g *exportGroup) Collect() groupcache_exporter.Stats {

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
func (g *exportGroup) Name() string {
	return g.group.Name()
}
