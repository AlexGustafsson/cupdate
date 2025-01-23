package httputil

import "fmt"

var _ error = (*Error)(nil)

type Error struct {
	Status     string
	StatusCode int
}

func (e Error) Error() string {
	return fmt.Sprintf("http: server returned error code %d - %s", e.StatusCode, e.Status)
}
