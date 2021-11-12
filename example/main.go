package main

import (
	"github.com/ch4rl1e5/ccdb/server"
	"github.com/ch4rl1e5/go-common/constants"
	"github.com/ch4rl1e5/go-common/logger"
	"go.uber.org/zap"
)

var ZapLogger *zap.Logger

func main() {
	logger.StartupLogger(logger.Config{Enviroment: constants.DevEnvironment})
	err := server.New("8080").Start()
	if err != nil {
		logger.ZapLogger.Fatal(server.ErrServerMessage, zap.Error(err))
	}
}
