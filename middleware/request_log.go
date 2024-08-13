package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	log "github.com/captjt/saddle/pkg/logger"
)

// RequestLog creates a middleware that logs request start and end times, along with latency.
func RequestLog(logger *log.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// You might need a way to generate or retrieve a request ID
		requestID := c.Get("X-Request-ID") // Assuming Request ID is sent in headers
		if requestID == "" {
			requestID = generateRequestID() // Implement this function to generate a request ID if not provided
			c.Set("X-Request-ID", requestID)
		}

		logger.Info("request received",
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.String("request_id", requestID),
		)

		// Proceed with chain
		err := c.Next()

		latency := time.Since(startTime).Milliseconds()
		logger.Info("request finished",
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.String("request_id", requestID),
			zap.Int64("latency", latency),
		)

		return err
	}
}

// generateRequestID returns a new unique UUID string for each request.
func generateRequestID() string {
	return uuid.NewString()
}
