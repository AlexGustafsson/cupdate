package httputil

import (
	"fmt"
	"net/http"
)

var _ error = (*Error)(nil)

type Error struct {
	Status     string
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
