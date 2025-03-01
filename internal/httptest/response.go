package httptest

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// Response describes a response to a request made to [Client].
type Response struct {
	// Status defaults to the status text of StatusCode.
	Status string
	// StatusCode defaults to 200.
	StatusCode int
	// Header defaults to an empty header.
	Header http.Header
	// Body is the response's body. Defaults to being empty.
	// Available types, in order of precedence:
	//	- string: path to file
	//	- string: verbatim string
	//	- []byte
	// 	- io.ReadCloser
	//	- io.Reader
	// 	- error: alias for ErrorReader{Error: err}
	Body any
	// Error is the response error, if any.
	Error error
}

// Response returns the HTTP response and error just like an HTTP client would,
// meaning the returned response and error can be returned from a "client.Do"
// function.
func (r Response) Response(t *testing.T, req *http.Request) (*http.Response, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	body := resolveResponseBody(t, r.Body)

	statusCode := http.StatusOK
	if r.StatusCode != 0 {
		statusCode = r.StatusCode
	}

	status := http.StatusText(statusCode)
	if r.Status != "" {
		status = r.Status
	}

	return &http.Response{
		Status:     status,
		StatusCode: statusCode,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     r.Header,
		Body:       body,
	}, nil
}

func resolveResponseBody(t *testing.T, body any) io.ReadCloser {
	if body == nil {
		return io.NopCloser(bytes.NewReader([]byte{}))
	}

	switch b := body.(type) {
	case string:
		data, err := os.ReadFile(b)
		if errors.Is(err, os.ErrNotExist) {
			data = []byte(b)
			err = nil
		}
		require.NoError(t, err)

		return io.NopCloser(bytes.NewReader(data))
	case []byte:
		return io.NopCloser(bytes.NewReader(b))
	}

	if b, ok := body.(io.ReadCloser); ok {
		return b
	}

	if b, ok := body.(io.Reader); ok {
		return io.NopCloser(b)
	}

	if b, ok := body.(error); ok {
		return io.NopCloser(ErrorReader{Error: b})
	}

	t.Fatalf("Unsupported response body type: %t, %+v", body, body)
	return nil
}
