package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/modernprogram/groupcache/v2"
)

func startGroupcache(workspace *groupcache.Workspace) []*groupcache.Group {

	ttl := time.Minute

	log.Printf("groupcache ttl: %v", ttl)

	//
	// create groupcache pool
	//

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

	var caches []*groupcache.Group

	names := []string{"files1", "files2"}

	for _, name := range names {

		options := groupcache.Options{
			Workspace:       workspace,
			Name:            name,
			PurgeExpired:    purgeExpired,
			CacheBytesLimit: groupcacheSizeBytes,
			Getter: groupcache.GetterFunc(
				func(_ /*ctx*/ context.Context, key string, dest groupcache.Sink, _ *groupcache.Info) error {

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

		caches = append(caches, cache)
	}

	return caches
}
