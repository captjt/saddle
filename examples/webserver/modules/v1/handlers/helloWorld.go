package handlers

import (
	"net/http"

	"github.com/captjt/saddle/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"saddle/examples/webserver/modules/v1/models"
)

func (h *Handlers) helloWorld() echo.HandlerFunc {
	return func(c echo.Context) error {
		in := c.Get(middleware.CTXRequest).(*models.HelloWorldRequest)

		h.logger.Info("retrieving application embedding",
			zap.String("message", in.Message),
		)

		out := models.Response{
			Message: "Hello there, saddle up friend!",
		}

		return c.JSON(http.StatusOK, out)
	}
}
