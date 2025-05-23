// Package main implements the example.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/mailgun/groupcache/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/udhos/groupcache_exporter"
	"github.com/udhos/groupcache_exporter/groupcache/mailgun"
)

func main() {

	var debug bool
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.Parse()

	appName := filepath.Base(os.Args[0])

	cache := startGroupcache()

	//
	// expose prometheus metrics
	//
	{
		metricsRoute := "/metrics"
		metricsPort := ":3000"

		log.Printf("starting metrics server at: %s %s", metricsPort, metricsRoute)

		labels := map[string]string{
			"app": appName,
		}
		namespace := ""
		options := groupcache_exporter.Options{
			Namespace:  namespace,
			Labels:     labels,
			Debug:      debug,
			ListGroups: func() []groupcache_exporter.GroupStatistics { return mailgun.ListGroups() },
		}
		collector := groupcache_exporter.NewExporter(options)

		prometheus.MustRegister(collector)

		go func() {
			http.Handle(metricsRoute, promhttp.Handler())
			log.Fatal(http.ListenAndServe(metricsPort, nil))
		}()
	}

	//
	// query cache periodically
	//

	const interval = 5 * time.Second

	for {
		var dst []byte
		cache.Get(context.TODO(), "/etc/passwd", groupcache.AllocatingByteSliceSink(&dst))
		log.Printf("cache answer: %d bytes, sleeping %v", len(dst), interval)
		time.Sleep(interval)
	}

}
