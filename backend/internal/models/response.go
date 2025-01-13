package models

// GenericResponse represents a generic JSON response
type GenericResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

