package handlers

import (
	"html/template"
	"tunnelogs-server/logger"
	"tunnelogs-server/utils"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func PageHandler(c echo.Context) error {
	log := logger.Log.Named("[PageHandler]")
	log.Debug("started")

	template := template.Must(template.ParseFiles(
		utils.GetTemplateFilePath("logs.html"),
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
