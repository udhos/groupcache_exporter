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

	group := g.group

	result.Group = groupcache_exporter.GroupStats{
		CounterGets:           group.Stats.Gets.Get(),
		CounterHits:           group.Stats.CacheHits.Get(),
		CounterPeerLoads:      group.Stats.PeerLoads.Get(),
		CounterPeerErrors:     group.Stats.PeerErrors.Get(),
		CounterLoads:          group.Stats.Loads.Get(),
		CounterLoadsDeduped:   group.Stats.LoadsDeduped.Get(),
		CounterLocalLoads:     group.Stats.LocalLoads.Get(),
		CounterLocalLoadsErrs: group.Stats.LocalLoadErrs.Get(),
		CounterServerRequests: group.Stats.ServerRequests.Get(),
	}

	result.Main = getCacheStats(group.CacheStats(groupcache.MainCache))
	result.Hot = getCacheStats(group.CacheStats(groupcache.HotCache))

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
