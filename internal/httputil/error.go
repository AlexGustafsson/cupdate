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
}

func (e Error) Error() string {
	return fmt.Sprintf("http: server returned error code %d - %s", e.StatusCode, e.Status)
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
