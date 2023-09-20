package models

type (
	HelloWorldRequest struct {
		Message string `json:"message" validate:"required"`
	}
	Response struct {
		Message string `json:"message"`
	}
)
