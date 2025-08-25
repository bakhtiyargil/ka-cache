package main

import (
	"ka-cache/bootstrap"
	"ka-cache/cache"
	"ka-cache/server/grpc"
	"ka-cache/server/http"
)

func main() {
	c := cache.NewLruCache(bootstrap.App.Config.Cache.Capacity)
	go startHttpServer(c)
	startGrpcServer(c)
}

func startHttpServer(cache cache.Cache) {
	h := http.NewCacheHandler(cache)
	hServer := http.NewHttpServer(bootstrap.App.Config, bootstrap.App.Logger, h)
	hServer.Start()
}

func startGrpcServer(cache cache.Cache) {
	gServer := grpc.NewGrpcServer(bootstrap.App.Config, bootstrap.App.Logger, cache)
	gServer.Start()
}
