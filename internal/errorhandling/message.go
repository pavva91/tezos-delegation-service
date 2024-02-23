package errorhandling

// SimpleErrorMessage represents request response with a SimpleErrorMessage
type SimpleErrorMessage struct {
	Message string `json:"error"`
}

type ValidationErrorsMessage struct {
	Message []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}
