package docker

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPlatform(t *testing.T) {
	testCases := []struct {
		// URI is the URI to the docker API.
		//
		// To handle dynamic ports, the following strings are replaced:
		//
		//   - <http host> -> server http host:port
		//   - <https host> -> server https host:port
		//   - <socket path> -> UNIX domain socket path
		URI         string
		UseTLS      bool
		ExpectedErr bool
	}{
		{
			URI: "http://<http host>",
		},
		{
			URI:    "https://<https host>",
			UseTLS: true,
		},
		{
			URI: "https://<http host>",
			// HTTPS on HTTP port
			ExpectedErr: true,
		},
		{
			URI: "tcp://<http host>",
		},
		{
			URI:    "tcp://<https host>",
			UseTLS: true,
		},
		{
			URI: "tcp://<http host>",
			// TLS on TCP port
			UseTLS:      true,
			ExpectedErr: true,
		},
		{
			URI: "unix://<socket path>",
		},
		{
			URI: "unix:///tmp/enoent",
			// Supported scheme, socket not found
			ExpectedErr: true,
		},
		{
			URI: "ftp://<http host>",
			// Unsupported scheme
			ExpectedErr: true,
		},
		{
			URI: "",
			// Missing URI
			ExpectedErr: true,
		},
		{
			URI: "http://<http host>!<-+",
			// Supported scheme, invalid URL
			ExpectedErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.URI, func(t *testing.T) {
			// NOTE: The path becomes quite long when using t.TempDir(), which
			// includes the test case's name. If it becomes too long, at least on
			// macOS, the binding of the socket will fail. Keep the test names short
			socketPath := filepath.Join(t.TempDir(), "socket")
			listener, err := net.Listen("unix", socketPath)
			// NOTE: Debugging weird errors here, such as binding failed? Look at the
			// above note
			require.NoError(t, err)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"ApiVersion":"1.0","MinAPIVersion":"1.0"}`))
			})

			server := &http.Server{
				Handler: mux,
			}
			go func() {
				err := server.Serve(listener)
				require.ErrorIs(t, err, http.ErrServerClosed)
			}()
			defer server.Close()

			httpProxyServer := httptest.NewServer(server.Handler)
			defer httpProxyServer.Close()

			httpProxyURL, err := url.Parse(httpProxyServer.URL)
			require.NoError(t, err)

			httpsProxyServer := httptest.NewTLSServer(server.Handler)
			defer httpsProxyServer.Close()

			httpsProxyURL, err := url.Parse(httpsProxyServer.URL)
			require.NoError(t, err)

			uri := testCase.URI
			uri = strings.ReplaceAll(uri, "<http host>", httpProxyURL.Host)
			uri = strings.ReplaceAll(uri, "<https host>", httpsProxyURL.Host)
			uri = strings.ReplaceAll(uri, "<socket path>", socketPath)

			options := &Options{}
			if testCase.UseTLS {
				options.TLSClientConfig = httpsProxyServer.TLS.Clone()
				options.TLSClientConfig.InsecureSkipVerify = true
			}

			_, err = NewPlatform(context.TODO(), uri, options)
			if testCase.ExpectedErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
