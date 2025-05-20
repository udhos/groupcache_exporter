package groupcache_exporter

// CacheDelta represents deltas for cache stats.
type CacheDelta struct {
	Gets              int64
	Hits              int64
	PeerLoads         int64
	PeerErrors        int64
	Loads             int64
	LoadsDeduped      int64
	LocalLoads        int64
	LocalLoadsErrs    int64
	ServerRequests    int64
	CrosstalkRefusals int64
}

// GetCacheDelta calculates deltas for cache stats.
func GetCacheDelta(prev, curr GroupStats) CacheDelta {
	return CacheDelta{
		Gets:              curr.CounterGets - prev.CounterGets,
		Hits:              curr.CounterHits - prev.CounterHits,
		PeerLoads:         curr.CounterPeerLoads - prev.CounterPeerLoads,
		PeerErrors:        curr.CounterPeerErrors - prev.CounterPeerErrors,
		Loads:             curr.CounterLoads - prev.CounterLoads,
		LoadsDeduped:      curr.CounterLoadsDeduped - prev.CounterLoadsDeduped,
		LocalLoads:        curr.CounterLocalLoads - prev.CounterLocalLoads,
		LocalLoadsErrs:    curr.CounterLocalLoadsErrs - prev.CounterLocalLoadsErrs,
		ServerRequests:    curr.CounterServerRequests - prev.CounterServerRequests,
		CrosstalkRefusals: curr.CounterCrosstalkRefusals - prev.CounterCrosstalkRefusals,
	}
}

// CacheTypeDelta represents deltas for per-type cache stats.
type CacheTypeDelta struct {
	Gets                int64
	Hits                int64
	Evictions           int64
	EvictionsNonExpired int64
}

// GetCacheTypeDelta calculates deltas for per-type cache stats.
func GetCacheTypeDelta(prev, curr CacheTypeStats) CacheTypeDelta {
	return CacheTypeDelta{
		Gets:                curr.CounterCacheGets - prev.CounterCacheGets,
		Hits:                curr.CounterCacheHits - prev.CounterCacheHits,
		Evictions:           curr.CounterCacheEvictions - prev.CounterCacheEvictions,
		EvictionsNonExpired: curr.CounterCacheEvictionsNonExpired - prev.CounterCacheEvictionsNonExpired,
	}
}
