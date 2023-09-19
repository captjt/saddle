package models

import (
	"runtime/debug"
)

type (
	// StatusResponse contains the response of a status request.
	StatusResponse struct {
		// Version contains the version of the service.
		Version string `json:"version,omitempty"`
		// CompiledAt contains the datetime stamp representing when the service was built.
		CompiledAt string `json:"compiled_at,omitempty"`
		// ExecutedAt contains the datetime stamp representing when the service was executed.
		ExecutedAt string `json:"executed_at"`
		// Uptime contains the different of time between now and ExecutedAt.
		Uptime string `json:"uptime"`
		// BuildInfo contains details about the build information of the compiled the service executable.
		BuildInfo *debug.BuildInfo `json:"build_info"`
	}
)
