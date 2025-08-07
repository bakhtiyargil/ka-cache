package app

import (
	"ka-cache/config"
	"ka-cache/logger"
	gs "ka-cache/server/grpc"
	hs "ka-cache/server/http/server"
	"os"
)

func main() {
	lConfig := config.LoadConfig("./config/config-local")
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
