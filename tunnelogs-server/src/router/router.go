package router

import (
	"tunnelogs-server/src/router/handlers"

	"github.com/labstack/echo"
)

func SetupRouter(server *echo.Echo) {
	server.GET("/", handlers.PageHandler)
	server.GET("/connect/:name/client", handlers.ClientWSHandler)
	server.GET("/connect/:name/server", handlers.ServerWSHandler)
}
