package handlers

import (
	"html/template"
	"tunnelogs-server/logger"
	"tunnelogs-server/utils"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func LobbyPageHandler(c echo.Context) error {
	log := logger.Log.Named("[LobbyPageHandler]")
	log.Debug("started")

	lobbyName := c.Param("name")
	if lobbyName == "" {
		log.Error("lobby name cannot be empty")
		return nil
	}

	template := template.Must(template.ParseFiles(
		utils.GetTemplateFilePath("logs.html"),
		utils.GetTemplateFilePath("base.html"),
	))

	err := template.ExecuteTemplate(c.Response(), "base", map[string]any{"Lobby": lobbyName})

	if err != nil {
		log.Error("failed to execute template", zap.Error(err))
		return err
	}

	log.Debug("finished")
	return nil
}

func HomePageHandler(c echo.Context) error {
	log := logger.Log.Named("[HomePageHandler]")
	log.Debug("started")

	template := template.Must(template.ParseFiles(
		utils.GetTemplateFilePath("home.html"),
		utils.GetTemplateFilePath("base.html"),
	))

	err := template.ExecuteTemplate(c.Response(), "base", nil)

	if err != nil {
		log.Error("failed to execute template", zap.Error(err))
		return err
	}

	log.Debug("finished")
	return nil
}
