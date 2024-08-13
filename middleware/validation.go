package middleware

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/captjt/saddle/models"
)

func Validate(v *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the model from context
		model := c.Locals("model")
		if model == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No model provided"})
		}

		// Create a new instance of the model's type to bind the request body
		modelType := reflect.TypeOf(model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		newModel := reflect.New(modelType).Interface()

		if err := c.BodyParser(newModel); err != nil {
			log.Errorw("unable to bind request", "path", c.Path(), "error", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err))
		}

		if err := v.Struct(newModel); err != nil {
			if ute, ok := err.(validator.ValidationErrors); ok {
				errs := overrideErrors(ute)
				log.Warnw("request validation(s) failed", "path", c.Path(), "errors", errs)
				return c.Status(fiber.StatusBadRequest).JSON(errs)
			}
			log.Errorw("request validation failed", "path", c.Path(), "error", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err))
		}

		// If validation passes, store the model in context for further use in handlers
		c.Locals("request", newModel)
		return c.Next()
	}
}

// overrideErrors overrides the default validation errors with custom-defined and cleaner error messages.
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
