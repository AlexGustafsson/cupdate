package httputil

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var _ AuthHandler = (*MockAuthHandler)(nil)

type MockAuthHandler struct {
	mock.Mock
}

// HandleAuth implements AuthHandler.
func (m *MockAuthHandler) HandleAuth(r *http.Request) error {
	args := m.Called(r)
	return args.Error(0)
}

func TestAuthMuxPatterns(t *testing.T) {
	testCases := []struct {
		Pattern  string
		URL      string
		Expected bool
	}{
		// From: https://kubernetes.io/docs/concepts/containers/images/#config-json
		{
			Pattern:  "*.kubernetes.io",
			URL:      "https://abc.kubernetes.io",
			Expected: true,
		},
		{
			Pattern:  "*.kubernetes.io",
			URL:      "https://kubernetes.io",
			Expected: false,
		},
		{
			Pattern:  "*.*.kubernetes.io",
			URL:      "https://abc.kubernetes.io",
			Expected: false,
		},
		{
			Pattern:  "*.*.kubernetes.io",
			URL:      "https://abc.def.kubernetes.io",
			Expected: true,
		},
		{
			Pattern:  "prefix.*.io",
			URL:      "https://prefix.kubernetes.io",
			Expected: true,
		},
		{
			Pattern:  "*-good.kubernetes.io",
			URL:      "https://prefix-good.kubernetes.io",
			Expected: true,
		},
		{
			Pattern:  "my-registry.io/images",
			URL:      "https://my-registry.io/images",
			Expected: true,
		},
		{
			Pattern:  "my-registry.io/images",
			URL:      "https://my-registry.io/images/my-image",
			Expected: true,
		},
		{
			Pattern:  "my-registry.io/images",
			URL:      "https://my-registry.io/images/another-image",
			Expected: true,
		},
		{
			Pattern:  "*.my-registry.io/images",
			URL:      "https://sub.my-registry.io/images/my-image",
			Expected: true,
		},
		{
			Pattern:  "*.my-registry.io/images",
			URL:      "https://a.sub.my-registry.io/images/my-image",
			Expected: false,
		},
		{
			Pattern:  "*.my-registry.io/images",
			URL:      "https://a.b.sub.my-registry.io/images/my-image",
			Expected: false,
		},
		{
			Pattern:  "my-registry.io/images",
			URL:      "https://a.sub.my-registry.io/images/my-image",
			Expected: false,
		},
		{
			Pattern:  "my-registry.io/images",
			URL:      "https://a.b.sub.my-registry.io/images/my-image",
			Expected: false,
		},
		// HTTP / HTTPS
		{
			Pattern:  "https://example.com",
			URL:      "https://example.com/images",
			Expected: true,
		},
		{
			Pattern:  "https://example.com",
			URL:      "http://example.com/images",
			Expected: false,
		},
		{
			Pattern:  "example.com",
			URL:      "https://example.com/images",
			Expected: true,
		},
		{
			Pattern:  "example.com",
			URL:      "http://example.com/images",
			Expected: true,
		},
		// IP / port
		{
			Pattern:  "example.com:8080",
			URL:      "https://example.com:8080/alpine",
			Expected: true,
		},
		{
			Pattern:  "192.168.1.100:8080",
			URL:      "https://192.168.1.100:8080/alpine",
			Expected: true,
		},
		{
			Pattern:  "192.168.1.100",
			URL:      "https://192.168.1.100/alpine",
			Expected: true,
		},
		{
			Pattern:  "example.com",
			URL:      "https://example.com:8080/alpine",
			Expected: false,
		},
		{
			Pattern:  "192.168.1.100",
			URL:      "https://192.168.1.100:8080/alpine",
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s matches %s: %v", testCase.Pattern, testCase.URL, testCase.Expected), func(t *testing.T) {
			authMux := NewAuthMux()

			handler := &MockAuthHandler{}
			handler.On("HandleAuth", mock.Anything).Return(nil)

			if testCase.Expected {
				authMux.Handle(testCase.Pattern, handler)
			} else {
				// Register default handler to test that no handler was matched
				authMux.Handle("", handler)
			}

			req, err := http.NewRequest(http.MethodGet, testCase.URL, nil)
			require.NoError(t, err)

			err = authMux.HandleAuth(req)
			require.NoError(t, err)

			handler.AssertExpectations(t)
		})
	}
}

func TestAuthMuxModifiesRequest(t *testing.T) {
	authMux := NewAuthMux()

	handler := &MockAuthHandler{}
	handler.On("HandleAuth", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		args.Get(0).(*http.Request).Header.Set("Authorization", "Bearer <token>")
	})

	authMux.Handle("", handler)

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	err = authMux.HandleAuth(req)
	require.NoError(t, err)

	assert.Equal(t, "Bearer <token>", req.Header.Get("Authorization"))

	handler.AssertExpectations(t)
}

func TestAuthMuxSetsHeaders(t *testing.T) {
	authMux := NewAuthMux()

	authMux.SetHeader("Authorization", "Bearer <token>")

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	err = authMux.HandleAuth(req)
	require.NoError(t, err)

	assert.Equal(t, "Bearer <token>", req.Header.Get("Authorization"))
}
