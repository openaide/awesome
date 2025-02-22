package server

type MCPServerError struct {
	Message string
}

func (e *MCPServerError) Error() string {
	return e.Message
}

func newMCPServerError(message string) *MCPServerError {
	return &MCPServerError{Message: message}
}

var ErrorNotSupported = newMCPServerError("not supported")
