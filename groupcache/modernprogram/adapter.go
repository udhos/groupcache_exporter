// Package modernprogram implements an adapter to extract metrics from modernprogram groupcache.
package modernprogram

import (
	"github.com/modernprogram/groupcache/v2"
	"github.com/udhos/groupcache_exporter"
)

// ListGroups issues current list of exporter groups for groupcache groups.
func ListGroups(ws *groupcache.Workspace) []groupcache_exporter.GroupStatistics {
	groups := groupcache.GetGroups(ws)
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
		CounterGets:                   stats.Gets.Get(),
		CounterHits:                   stats.CacheHits.Get(),
		GaugeGetFromPeersLatencyLower: float64(stats.GetFromPeersLatencyLower.Get()),
		CounterPeerLoads:              stats.PeerLoads.Get(),
		CounterPeerErrors:             stats.PeerErrors.Get(),
		CounterLoads:                  stats.Loads.Get(),
		CounterLoadsDeduped:           stats.LoadsDeduped.Get(),
		CounterLocalLoads:             stats.LocalLoads.Get(),
		CounterLocalLoadsErrs:         stats.LocalLoadErrs.Get(),
		CounterServerRequests:         stats.ServerRequests.Get(),
		CounterCrosstalkRefusals:      stats.CrosstalkRefusals.Get(),
	}

	result.Main = getCacheStats(g.group.CacheStats(groupcache.MainCache))
	result.Hot = getCacheStats(g.group.CacheStats(groupcache.HotCache))

	return result
}

func getCacheStats(cacheStats groupcache.CacheStats) groupcache_exporter.CacheTypeStats {
	return groupcache_exporter.CacheTypeStats{
		GaugeCacheItems:                 cacheStats.Items,
		GaugeCacheBytes:                 cacheStats.Bytes,
		CounterCacheGets:                cacheStats.Gets,
		CounterCacheHits:                cacheStats.Hits,
		CounterCacheEvictions:           cacheStats.Evictions,
		CounterCacheEvictionsNonExpired: cacheStats.EvictionsNonExpiredOnMemFull,
	}
}

// Name returns the group's name
func (g *exportGroup) Name() string {
	return g.group.Name()
}
