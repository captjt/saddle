// Package middleware contains saddle-level related middleware which utilizes the Echo framework.
package middleware

const (
	// CTXRequest contains the key in which the request payload is attached and referenced to the request context.
	CTXRequest = "ctxRequest"
	// CTXRequestID contains the key in which the request id is attached and referenced to the request context.
	CTXRequestID = "ctxRequestID"
)
