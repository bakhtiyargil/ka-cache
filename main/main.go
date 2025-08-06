package main

import (
	"ka-cache/config"
	gs "ka-cache/grpc/server"
	hs "ka-cache/http/server"
	"ka-cache/pkg/logger"
	"os"
)

func main() {
	lConfig, _ := config.LoadConfig("./config/config-local")
	ch := make(chan string)

	go startGrpcServer(ch, lConfig)
	go startDefaultServer(ch, lConfig)

	<-ch
	<-ch
}

func startDefaultServer(ch chan string, cnfg *config.Config) {
	loggr := logger.NewCustomLogger(cnfg)
	loggr.InitLogger()
	eServer := hs.NewHttpServer(cnfg, loggr)
	err := eServer.Run()
	if err != nil {
		os.Exit(1)
	}
}

func startGrpcServer(ch chan string, cnfg *config.Config) {
	err := gs.NewServer().Run(cnfg)
	if err != nil {
		os.Exit(1)
	}
}
