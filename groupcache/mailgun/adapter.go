// Package mailgun implements an adapter to extract metrics from mailgun groupcache.
package mailgun

import (
	"github.com/mailgun/groupcache/v2"
	"github.com/udhos/groupcache_exporter"
)

// ListGroups issues current list of exporter groups for groupcache groups.
func ListGroups() []groupcache_exporter.GroupStatistics {
	groups := groupcache.GetGroups()
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
		CounterGets:                   group.Stats.Gets.Get(),
		CounterHits:                   group.Stats.CacheHits.Get(),
		GaugeGetFromPeersLatencyLower: float64(group.Stats.GetFromPeersLatencyLower.Get()),
		CounterPeerLoads:              group.Stats.PeerLoads.Get(),
		CounterPeerErrors:             group.Stats.PeerErrors.Get(),
		CounterLoads:                  group.Stats.Loads.Get(),
		CounterLoadsDeduped:           group.Stats.LoadsDeduped.Get(),
		CounterLocalLoads:             group.Stats.LocalLoads.Get(),
		CounterLocalLoadsErrs:         group.Stats.LocalLoadErrs.Get(),
		CounterServerRequests:         group.Stats.ServerRequests.Get(),
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
