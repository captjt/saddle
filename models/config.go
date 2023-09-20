package models

type (
	// Config contains the configuration(s) model for the saddled service.
	Config struct {
		// Saddle contains the configuration(s) model(s) for the saddled service.
		Saddle struct {
			// CloudTrace contains the configuration(s) for the Google® Cloud Trace open-telemetry exporter.
			CloudTrace *CloudTrace `mapstructure:"cloud_trace" validate:"omitempty,excluded_with=Jaeger StdOut"`
			// Jaeger contains the configuration(s) for the Jaeger® open-telemetry exporter.
			Jaeger *Jaeger `mapstructure:"jaeger" validate:"omitempty,excluded_with=CloudTrace StdOut"`
			// StdOut contains the configuration(s) for the stdout open-telemetry exporter.
			StdOut *StdOut `mapstructure:"stdout" validate:"omitempty,excluded_with=CloudTrace Jaeger None"`
			// None contains the configuration(s) for no trace exporter.
			None *None `mapstructure:"none" validate:"omitempty,excluded_with=CloudTrace Jaeger StdOut"`
		} `mapstructure:"saddle"`
	}

	// CloudTrace contains the configuration(s) for the Google® Cloud Trace open-telemetry exporter.
	CloudTrace struct {
		// ProjectID is the Google Cloud Project identifier to export Cloud Tracing telemetry.
		ProjectID string `mapstructure:"project_id" validate:"required"`
		// SampleRate contains the percentage rate of total requests to collect and export.
		SampleRate float64 `mapstructure:"sample_rate" validate:"required,min=0,max=100"`
	}

	// Jaeger contains the configuration(s) for the Jaeger® open-telemetry exporter.
	Jaeger struct {
		// URI contains the host address of the Jaeger® service in which to export.
		URI string `mapstructure:"uri" validate:"required,uri"`
		// SampleRate contains the percentage rate of total requests to collect and export.
		SampleRate float64 `mapstructure:"sample_rate" validate:"required,min=0,max=100"`
	}

	StdOut struct {
		// SampleRate contains the percentage rate of total requests to collect and export.
		SampleRate float64 `mapstructure:"sample_rate" validate:"required,min=0,max=100"`
	}

	// None contains toggle for disabling configurations.
	None struct {
		Disabled bool `mapstructure:"disabled" validate:"required"`
	}
)
