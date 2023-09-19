package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/captjt/saddle/models"
	log "github.com/captjt/saddle/pkg/logger"
)

// Validate is a middleware that handles the validation of a payload associated to a request; defined using the
//
//	'validate' tag in the model structure.
//	https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme
func Validate(logger *log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if r := c.Get(CTXRequest); r != nil {
				if err := c.Bind(r); err != nil {
					logger.Error("unable to bind request",
						zap.String("path", c.Path()),
						zap.String("request_id", c.Get(CTXRequestID).(string)),
						zap.Error(err),
					)
					return echo.NewHTTPError(http.StatusBadRequest, models.NewErrorResponse(err))
				}
				if err := c.Validate(r); err != nil {
					if ute, ok := err.(validator.ValidationErrors); ok {
						errs := overrideErrors(ute)
						logger.Warn("request validation(s) failed",
							zap.String("path", c.Path()),
							zap.String("request_id", c.Get(CTXRequestID).(string)),
							zap.Any("errors", errs),
						)
						return echo.NewHTTPError(http.StatusBadRequest, errs)
					}
					logger.Error("request validation failed",
						zap.String("path", c.Path()),
						zap.String("request_id", c.Get(CTXRequestID).(string)),
						zap.Error(err),
					)
					return echo.NewHTTPError(http.StatusBadRequest, models.NewErrorResponse(err))
				}
				c.Set(CTXRequest, r)
			}
			return next(c)
		}
	}
}

// overrideErrors overrides the default validation errors with custom-defined and cleaner error messages.  Validator
//
//	does have the capability for translations; however, it is overkill for what we need here.
//	https://pkg.go.dev/github.com/go-playground/universal-translator
//	https://pkg.go.dev/github.com/go-playground/universal-translator#Translator
func overrideErrors(errs validator.ValidationErrors) *models.Errors {
	ne := models.NewErrorResponse(nil)

	for _, err := range errs {
		switch err.Tag() {
		case "required":
			ne.AppendError(fmt.Errorf("missing required value for parameter | field: %s", err.StructField()))
		case "json":
			ne.AppendError(fmt.Errorf("invalid json value for parameter | field: %s", err.StructField()))
		default:
			ne.AppendError(fmt.Errorf("invalid value for parameter | field: %s", err.StructField()))
		}
	}
	return ne
}
