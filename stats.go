package groupcache_exporter

// Stats holds full stats for the group.
type Stats struct {
	// Group holds general stats for the group.
	Group GroupStats

	// Main holds stats for the group main cache.
	Main CacheTypeStats

	// Hot holds stats for the group hot cache.
	Hot CacheTypeStats
}

// GroupStats holds general stats for the group.
type GroupStats struct {

	// CounterGets represents any Get request, including from peers
	CounterGets int64

	// CounterHits represents either cache was good
	CounterHits int64

	// GaugeGetFromPeersLatencyLower represents slowest duration to request value from peers
	GaugeGetFromPeersLatencyLower float64

	// CounterPeerLoads represents either remote load or remote cache hit (not an error)
	CounterPeerLoads int64

	// CounterPeerErrors represents a count of errors from peers
	CounterPeerErrors int64

	// CounterLoads represents (gets - cacheHits)
	CounterLoads int64

	// CounterLoadsDeduped represents after singleflight
	CounterLoadsDeduped int64

	// CounterLocalLoads represents total good local loads
	CounterLocalLoads int64

	// CounterLocalLoadsErrs represents total bad local loads
	CounterLocalLoadsErrs int64

	// CounterServerRequests represents gets that came over the network from peers
	CounterServerRequests int64

	// CounterCrosstalkRefusals represents refusals for additional crosstalks
	CounterCrosstalkRefusals int64
}

// CacheTypeStats holds stats for the main cache or the hot cache.
type CacheTypeStats struct {
	// GaugeCacheItems represents number of items in the main/hot cache
	GaugeCacheItems int64

	// GaugeCacheBytes represents number of bytes in the main/hot cache
	GaugeCacheBytes int64

	// CounterCacheGets represents number of get requests in the main/hot cache
	CounterCacheGets int64

	// CounterCacheHits represents number of hit in the main/hot cache
	CounterCacheHits int64

	// CounterCacheEvictions represents number of evictions in the main/hot cache
	CounterCacheEvictions int64

	// CounterCacheEvictionsNonExpired represents number of evictions for non-expired keys in the main/hot cache
	CounterCacheEvictionsNonExpired int64
}
