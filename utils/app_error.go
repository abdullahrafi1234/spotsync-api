package utils

// AppError is a custom error type that carries an HTTP status code along with it.
// This lets services/handlers say "this error should become a 404" or "409" etc.,
// without the handler needing a big if/else chain for every error.
type AppError struct {
	Code    int    // HTTP status code
	Message string // human-readable message
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError is a helper constructor.
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}