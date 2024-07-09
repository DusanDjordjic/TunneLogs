package main

import (
	"github.com/labstack/echo"
	"go.uber.org/zap"

	"tunnelogs-server/logger"
	"tunnelogs-server/src/db"
	"tunnelogs-server/src/router"
	"tunnelogs-server/utils"
)

func main() {
	err := logger.Setup()
	if err != nil {
		println("failed to create logger", err)
	}

	err = db.Connect()
	if err != nil {
		logger.Log.Fatal("failed to connect to db", zap.Error(err))
	}

	server := echo.New()
	server.Static("/static", utils.GetStaticFolderPath())
	router.SetupRouter(server)

	if err := server.Start("127.0.0.1:8080"); err != nil {
		logger.Log.Fatal("failed to close server", zap.Error(err))
	}

}
