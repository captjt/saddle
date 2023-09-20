package middleware

import (
	"github.com/captjt/saddle/middleware"
	"github.com/labstack/echo/v4"

	"saddle/examples/webserver/modules/v1/models"
)

func APIData() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			switch c.Path() {
			case "/v1/hello-world":
				c.Set(middleware.CTXRequest, new(models.HelloWorldRequest))
			}
			return next(c)
		}
	}
}
