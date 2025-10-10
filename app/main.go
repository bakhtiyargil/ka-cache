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
	cfg := bootstrap.App.Config

	h := http.NewCacheHandler(cache)
	hServer := http.NewHttpServer(cfg, bootstrap.App.Logger, h)
	if cfg.Server.Default.EnableSecure {
		hServer.StartSecure()
	} else {
		hServer.Start()
	}

	gServer := grpc.NewGrpcServer(cfg, bootstrap.App.Logger, cache)
	if cfg.Server.Grpc.EnableSecure {
		gServer.StartSecure()
	} else {
		gServer.Start()
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-stopChan:
		if cfg.Server.Grpc.EnableSecure {
			gServer.StopSecure()
		} else {
			gServer.Stop()
		}

		if cfg.Server.Default.EnableSecure {
			hServer.StopSecure()
		} else {
			hServer.Stop()
		}
	}
}
