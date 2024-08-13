package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// requestIDHeader contains the HTTP header in which to reference the request ID.
const requestIDHeader = "X-Request-ID"

// RequestID handles any referenced request ID attached to an incoming request; will construct and attach a new ID to
// the incoming request if it does not already exist.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(requestIDHeader) // Get the request ID from the header
		if requestID == "" {
			requestID = uuid.NewString()      // Generate a new UUID if the header is not present
			c.Set(requestIDHeader, requestID) // Set the new UUID in the response header
		}

		// Store the request ID in the local context, if needed elsewhere
		c.Locals("requestID", requestID)

		return c.Next() // Continue with the next middleware or handler in the chain
	}
}
