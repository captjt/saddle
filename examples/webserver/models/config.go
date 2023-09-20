package models

type (
	Config struct {
		V1 *V1 `mapstructure:"v1" validate:"required"`
	}

	V1 struct {
		// Specific V1 handler configurations can go in here.
		Test string `mapstructure:"test" validate:"required"`
	}
)
