package models

type (
	// Error contains an outgoing error response payload.
	Error struct {
		// Code contains a unique code identifier for the referenced error.
		Code string `json:"code,omitempty"`
		// Message contains a user-friendly message pertaining to the details of the referenced error.
		Message string `json:"message,omitempty"`
	}

	// Errors contains a collection of outgoing error response payloads.
	Errors struct {
		// Errors contains a collection of Error.
		Errors []*Error `json:"errors"`
	}
)

// NewErrorResponse constructs a new instance of Errors.
func NewErrorResponse(err error, code ...string) *Errors {
	return new(Errors).AppendError(err, code...)
}

// AppendError appends a referenced error and error code to the associated Errors response.
func (e *Errors) AppendError(err error, code ...string) *Errors {
	var m, c string

	if err != nil {
		m = err.Error()
	}
	if len(code) > 0 {
		c = code[0]
	}

	if m != "" || c != "" {
		e.Errors = append(e.Errors, &Error{
			Code:    c,
			Message: m,
		})
	}
	return e
}
