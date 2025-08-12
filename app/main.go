package app

import (
	"ka-cache/config"
	"ka-cache/logger"
	"ka-cache/server/grpc"
	"ka-cache/server/http"
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
	eServer := http.NewHttpServer(cnfg, loggr)
	err := eServer.Run()
	if err != nil {
		os.Exit(1)
	}
}

func startGrpcServer(ch chan string, cnfg *config.Config) {
	err := grpc.NewGrpcServer(cnfg).Run()
	if err != nil {
		os.Exit(1)
	}
}
