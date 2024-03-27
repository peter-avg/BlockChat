package utils

import "fmt"

type CustomError struct {
	Code    int    // Error code
	Message string // Error message
}

// Error returns the error message of CustomError.
func (e *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}
