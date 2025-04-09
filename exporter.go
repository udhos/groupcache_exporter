// Package groupcache_exporter exports prometheus metrics for groupcache.
package groupcache_exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Exporter implements interface prometheus.Collector to extract metrics from groupcache.
type Exporter struct {
	groups []GroupStatistics

	groupGets                     *prometheus.Desc
	groupCacheHits                *prometheus.Desc
	groupGetFromPeersLatencyLower *prometheus.Desc
	groupPeerLoads                *prometheus.Desc
	groupPeerErrors               *prometheus.Desc
	groupLoads                    *prometheus.Desc
	groupLoadsDeduped             *prometheus.Desc
	groupLocalLoads               *prometheus.Desc
	groupLocalLoadErrs            *prometheus.Desc
	groupServerRequests           *prometheus.Desc
	groupCrosstalkRefusals        *prometheus.Desc

	cacheBytes               *prometheus.Desc
	cacheItems               *prometheus.Desc
	cacheGets                *prometheus.Desc
	cacheHits                *prometheus.Desc
	cacheEvictions           *prometheus.Desc
	cacheEvictionsNonExpired *prometheus.Desc
}

// GroupStatistics is a plugable interface to extract metrics from a groupcache implementation.
// GroupStatistics is used by Exporter to collect the group statistics.
// The user must provide a concrete implementation of this interface that knows how to
// extract group statistics from the actual groupcache implementation.
type GroupStatistics interface {
	// Collect requests metrics collection from implementation
	Collect() Stats

	// Name returns the group's name
	Name() string
}

// NewExporter creates Exporter.
// namespace is usually the empty string.
func NewExporter(namespace string, labels map[string]string, groups ...GroupStatistics) *Exporter {

	const subsystem = "groupcache"

	return &Exporter{
		groups: groups,

		groupGets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "gets_total"),
			"Count of cache gets (including from peers)",
			[]string{"group"},
			labels,
		),
		groupCacheHits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "hits_total"),
			"Count of cache hits (from either main or hot cache)",
			[]string{"group"},
			labels,
		),
		groupGetFromPeersLatencyLower: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "get_from_peers_latency_slowest_milliseconds"),
			"Represent slowest duration to request value from peers.",
			[]string{"group"},
			labels,
		),
		groupPeerLoads: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "peer_loads_total"),
			"Count of non-error loads or cache hits from peers",
			[]string{"group"},
			labels,
		),
		groupPeerErrors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "peer_errors_total"),
			"Count of errors from peers",
			[]string{"group"},
			labels,
		),
		groupLoads: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "loads_total"),
			"Count of (gets - hits)",
			[]string{"group"},
			labels,
		),
		groupLoadsDeduped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "loads_deduped_total"),
			"Count of loads after singleflight",
			[]string{"group"},
			labels,
		),
		groupLocalLoads: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "local_load_total"),
			"Count of loads from local cache",
			[]string{"group"},
			labels,
		),
		groupLocalLoadErrs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "local_load_errs_total"),
			"Count of loads from local cache that failed",
			[]string{"group"},
			labels,
		),
		groupServerRequests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "server_requests_total"),
			"Count of gets that came over the network from peers",
			[]string{"group"},
			labels,
		),
		groupCrosstalkRefusals: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "crosstalk_refusals_total"),
			"Count of refusals for additional crosstalks",
			[]string{"group"},
			labels,
		),

		cacheBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cache_bytes"),
			"Gauge of current bytes in use",
			[]string{"group", "type"},
			labels,
		),
		cacheItems: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cache_items"),
			"Gauge of current items in use",
			[]string{"group", "type"},
			labels,
		),
		cacheGets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cache_gets_total"),
			"Count of cache gets",
			[]string{"group", "type"},
			labels,
		),
		cacheHits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cache_hits_total"),
			"Count of cache hits",
			[]string{"group", "type"},
			labels,
		),
		cacheEvictions: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cache_evictions_total"),
			"Count of cache evictions",
			[]string{"group", "type"},
			labels,
		),
		cacheEvictionsNonExpired: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cache_evictions_nonexpired_total"),
			"Count of cache evictions for non-expired keys due to memory full.",
			[]string{"group", "type"},
			labels,
		),
	}
}

// Describe sends metrics descriptors.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.groupGets
	ch <- e.groupCacheHits
	ch <- e.groupGetFromPeersLatencyLower
	ch <- e.groupPeerLoads
	ch <- e.groupPeerErrors
	ch <- e.groupLoads
	ch <- e.groupLoadsDeduped
	ch <- e.groupLocalLoads
	ch <- e.groupLocalLoadErrs
	ch <- e.groupServerRequests
	ch <- e.groupCrosstalkRefusals

	ch <- e.cacheBytes
	ch <- e.cacheItems
	ch <- e.cacheGets
	ch <- e.cacheHits
	ch <- e.cacheEvictions
	ch <- e.cacheEvictionsNonExpired
}

// Collect is called by the Prometheus registry when collecting metrics.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	for _, group := range e.groups {
		e.collectFromGroup(ch, group)
	}
}

func (e *Exporter) collectFromGroup(ch chan<- prometheus.Metric, group GroupStatistics) {
	stats := group.Collect()
	groupName := group.Name()
	e.collectStats(ch, stats.Group, groupName)
	e.collectCacheStats(ch, stats.Main, groupName, "main")
	e.collectCacheStats(ch, stats.Hot, groupName, "hot")
}

func (e *Exporter) collectStats(ch chan<- prometheus.Metric, stats GroupStats, groupName string) {
	ch <- prometheus.MustNewConstMetric(e.groupGets, prometheus.CounterValue, float64(stats.CounterGets), groupName)
	ch <- prometheus.MustNewConstMetric(e.groupCacheHits, prometheus.CounterValue, float64(stats.CounterHits), groupName)
	ch <- prometheus.MustNewConstMetric(e.groupGetFromPeersLatencyLower, prometheus.GaugeValue, stats.GaugeGetFromPeersLatencyLower, groupName)
	ch <- prometheus.MustNewConstMetric(e.groupPeerLoads, prometheus.CounterValue, float64(stats.CounterPeerLoads), groupName)
	ch <- prometheus.MustNewConstMetric(e.groupPeerErrors, prometheus.CounterValue, float64(stats.CounterPeerErrors), groupName)
	ch <- prometheus.MustNewConstMetric(e.groupLoads, prometheus.CounterValue, float64(stats.CounterLoads), groupName)
	ch <- prometheus.MustNewConstMetric(e.groupLoadsDeduped, prometheus.CounterValue, float64(stats.CounterLoadsDeduped), groupName)
	ch <- prometheus.MustNewConstMetric(e.groupLocalLoads, prometheus.CounterValue, float64(stats.CounterLocalLoads), groupName)
	ch <- prometheus.MustNewConstMetric(e.groupLocalLoadErrs, prometheus.CounterValue, float64(stats.CounterLocalLoadsErrs), groupName)
	ch <- prometheus.MustNewConstMetric(e.groupServerRequests, prometheus.CounterValue, float64(stats.CounterServerRequests), groupName)
	ch <- prometheus.MustNewConstMetric(e.groupCrosstalkRefusals, prometheus.CounterValue, float64(stats.CounterCrosstalkRefusals), groupName)
}

func (e *Exporter) collectCacheStats(ch chan<- prometheus.Metric, stats CacheTypeStats, groupName, cacheType string) {
	ch <- prometheus.MustNewConstMetric(e.cacheItems, prometheus.GaugeValue, float64(stats.GaugeCacheItems), groupName, cacheType)
	ch <- prometheus.MustNewConstMetric(e.cacheBytes, prometheus.GaugeValue, float64(stats.GaugeCacheBytes), groupName, cacheType)
	ch <- prometheus.MustNewConstMetric(e.cacheGets, prometheus.CounterValue, float64(stats.CounterCacheGets), groupName, cacheType)
	ch <- prometheus.MustNewConstMetric(e.cacheHits, prometheus.CounterValue, float64(stats.CounterCacheHits), groupName, cacheType)
	ch <- prometheus.MustNewConstMetric(e.cacheEvictions, prometheus.CounterValue, float64(stats.CounterCacheEvictions), groupName, cacheType)
	ch <- prometheus.MustNewConstMetric(e.cacheEvictionsNonExpired, prometheus.CounterValue, float64(stats.CounterCacheEvictionsNonExpired), groupName, cacheType)
}
