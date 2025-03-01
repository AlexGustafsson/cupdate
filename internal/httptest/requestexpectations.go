package httptest

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// RequestExpectations describes the expectations on the request.
type RequestExpectations struct {
	URL    string
	Header http.Header
	// Body is the response's body. Defaults to being empty.
	// Available types, in order of precedence:
	//	- string: path to file
	//	- string: verbatim string
	//	- []byte
	//	- io.Reader
	Body any
}

// Assert asserts the expectations on the request.
func (e RequestExpectations) Assert(t *testing.T, r *http.Request) {
	// Assert URL
	if e.URL != "" {
		assert.Equal(t, e.URL, r.URL.String())
	}

	// Assert header
	for k, v := range e.Header {
		assert.EqualValues(t, v, r.Header[k])
	}

	if r.Body != nil {
		// Read body
		actualBody, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		// Assert body
		if e.Body != nil {
			expectedBody := resolveRequestBody(t, e.Body)

			assert.Equal(t, expectedBody, actualBody)
		}
	}
}

func resolveRequestBody(t *testing.T, body any) []byte {
	if body == nil {
		return nil
	}

	switch b := body.(type) {
	case string:
		data, err := os.ReadFile(b)
		if errors.Is(err, os.ErrNotExist) {
			data = []byte(b)
			err = nil
		}
		require.NoError(t, err)

		return data
	case []byte:
		return b
	}

	if b, ok := body.(io.Reader); ok {
		data, err := io.ReadAll(b)
		require.NoError(t, err)
		return data
	}

	t.Fatalf("Unsupported request body type: %t, %+v", body, body)
	return nil
}
