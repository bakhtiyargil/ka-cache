package main

import (
	"ka-cache/bootstrap"
	"ka-cache/cache"
	"ka-cache/server/grpc"
	"ka-cache/server/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := cache.NewLruCache[string, string](bootstrap.App.Config.Cache.Capacity, bootstrap.App.Logger)
	go c.StartCleanup(bootstrap.App.Config.Cache.CleanupInterval * time.Second)
	startServers(c)
}

func startServers(cache cache.Cache[string, string]) {
	h := http.NewCacheHandler(cache)
	hServer := http.NewHttpServer(bootstrap.App.Config, bootstrap.App.Logger, h)
	hServer.Start()

	gServer := grpc.NewGrpcServer(bootstrap.App.Config, bootstrap.App.Logger, cache)
	gServer.Start()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stopChan:
		hServer.Stop()
		gServer.Stop()
	}
}
