package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// requestIDHeader contains the http header in which to reference the request id.
const requestIDHeader = "x-request-id"

// RequestID handles any referenced request id attached to an incoming request; will construct and attached a new id to
//   the incoming request if it does not already exist.
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if rid := c.Request().Header.Get(requestIDHeader); rid != "" {
				c.Set(CTXRequestID, rid)
				return next(c)
			}
			c.Set(CTXRequestID, uuid.New().String())
			return next(c)
		}
	}
}
