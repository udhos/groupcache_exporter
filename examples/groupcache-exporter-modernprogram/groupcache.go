package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/modernprogram/groupcache/v2"
)

func startGroupcache() *groupcache.Group {

	ttl := time.Minute

	log.Printf("groupcache ttl: %v", ttl)

	//
	// create groupcache pool
	//

	workspace := groupcache.NewWorkspace()

	groupcachePort := ":5000"

	myURL := "http://127.0.0.1" + groupcachePort

	log.Printf("groupcache my URL: %s", myURL)

	pool := groupcache.NewHTTPPoolOptsWithWorkspace(workspace, myURL, &groupcache.HTTPPoolOptions{})

	//
	// start groupcache server
	//

	serverGroupCache := &http.Server{Addr: groupcachePort, Handler: pool}

	go func() {
		log.Printf("groupcache server: listening on %s", groupcachePort)
		err := serverGroupCache.ListenAndServe()
		log.Printf("groupcache server: exited: %v", err)
	}()

	pool.Set(myURL)

	//
	// create cache
	//

	const purgeExpired = true
	const groupcacheSizeBytes = 1_000_000

	// https://talks.golang.org/2013/oscon-dl.slide#46
	//
	// 64 MB max per-node memory usage

	options := groupcache.Options{
		Workspace:    workspace,
		Name:         "files",
		PurgeExpired: purgeExpired,
		CacheBytes:   groupcacheSizeBytes,
		Getter: groupcache.GetterFunc(
			func(_ /*ctx*/ context.Context, key string, dest groupcache.Sink) error {

				log.Printf("getter: loading: key:%s, ttl:%v", key, ttl)

				data, errFile := os.ReadFile(key)
				if errFile != nil {
					return errFile
				}

				expire := time.Now().Add(ttl)
				return dest.SetBytes(data, expire)
			}),
	}

	cache := groupcache.NewGroupWithWorkspace(options)

	return cache
}
