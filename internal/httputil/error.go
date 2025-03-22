package httputil

import (
	"fmt"
	"net/http"
)

var _ error = (*Error)(nil)

// Error is a base HTTP error useful for HTTP clients.
type Error struct {
	// Status is the HTTP status of the response.
	Status string
	// StatusCode is the HTTP status code of the response.
	StatusCode int
	// Message is a message communicated through other means, such as through an
	// Www-Authenticate header error parameter.
	Message string
}

func (e Error) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("http: server returned error code %d - %s", e.StatusCode, e.Status)
	}

	return fmt.Sprintf("http: server returned error code %d - %s: %s", e.StatusCode, e.Status, e.Message)
}

// AssertStatusCode returns an error if the response does not match the given
// status code.
func AssertStatusCode(r *http.Response, statusCode int) error {
	if r.StatusCode != statusCode {
		return Error{
			Status:     r.Status,
			StatusCode: r.StatusCode,
		}
	}

	return nil
}
