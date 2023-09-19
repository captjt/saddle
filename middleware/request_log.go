package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	log "github.com/captjt/saddle/pkg/logger"
)

func RequestLog(logger *log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()
			logger.Info("request received",
				zap.String("path", c.Path()),
				zap.String("request_id", c.Get(CTXRequestID).(string)),
			)

			err := next(c)
			stopTime := time.Now()
			latency := stopTime.Sub(startTime).Milliseconds()
			logger.Info("request finished",
				zap.String("path", c.Path()),
				zap.String("request_id", c.Get(CTXRequestID).(string)),
				zap.Int64("latency", latency),
			)
			return err
		}
	}
}
