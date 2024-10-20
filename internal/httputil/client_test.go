package httputil

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientDoCachedHappyPath(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	var handledRequests atomic.Int32

	// Respond to requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handledRequests.Add(1)

		w.Header().Set("X-Foo", "bar")
		w.Write([]byte("bar"))
	}))
	defer server.Close()

	testCache := cachetest.NewCache(t)
	client := NewClient(testCache, 5*time.Second)

	// Perform a request with the response not yet cached
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	res, err := client.DoCached(req)
	require.NoError(t, err)

	// Server was hit
	assert.Equal(t, int32(1), handledRequests.Load(), "expect cache miss")

	// Headers are kept
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "bar", res.Header.Get("X-Foo"))

	// Body is kept
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)
	assert.Equal(t, []byte("bar"), body)

	// Expect the response to be cached
	_, err = testCache.Get(context.TODO(), client.CacheKey(req))
	require.NoError(t, err)

	// Perform the request again, expecting the server to not be hit
	req, err = http.NewRequestWithContext(context.TODO(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	res, err = client.DoCached(req)
	require.NoError(t, err)

	// Server was not hit
	assert.Equal(t, int32(1), handledRequests.Load(), "expect cache hit")

	// Headers are kept
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "bar", res.Header.Get("X-Foo"))

	// Body is kept
	body, err = io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)
	assert.Equal(t, []byte("bar"), body)
}

func TestClientDoCachedServerError(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	var handledRequests atomic.Int32

	// Respond to requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handledRequests.Add(1)

		w.Header().Set("X-Foo", "bar")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("bar"))
	}))
	defer server.Close()

	testCache := cachetest.NewCache(t)
	client := NewClient(testCache, 5*time.Second)

	// Perform a request
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	res, err := client.DoCached(req)
	require.NoError(t, err)

	// Server was hit
	assert.Equal(t, int32(1), handledRequests.Load(), "expect cache miss")

	// Headers are kept
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "bar", res.Header.Get("X-Foo"))

	// Body is kept
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)
	assert.Equal(t, []byte("bar"), body)

	// Expect the response to not be cached
	_, err = testCache.Get(context.TODO(), client.CacheKey(req))
	require.Equal(t, cache.ErrNotExist, err)

	// Perform the request again, expecting the server to be hit again
	req, err = http.NewRequestWithContext(context.TODO(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	res, err = client.DoCached(req)
	require.NoError(t, err)

	// Server was hit
	assert.Equal(t, int32(2), handledRequests.Load(), "expect cache miss")

	// Headers are kept
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "bar", res.Header.Get("X-Foo"))

	// Body is kept
	body, err = io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)
	assert.Equal(t, []byte("bar"), body)
}

func TestClientDoCachedOutdatedEntry(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	var handledRequests atomic.Int32

	// Respond to requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handledRequests.Add(1)

		w.Header().Set("X-Foo", "bar")
		w.Write([]byte("bar"))
	}))
	defer server.Close()

	cache := cachetest.NewCache(t)
	client := NewClient(cache, 1*time.Second)

	// Perform a request with the response not yet cached
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	res, err := client.DoCached(req)
	require.NoError(t, err)

	// Server was hit
	assert.Equal(t, int32(1), handledRequests.Load(), "expect cache miss")

	// Headers are kept
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "bar", res.Header.Get("X-Foo"))

	// Body is kept
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)
	assert.Equal(t, []byte("bar"), body)

	// Expect the response to be cached
	_, err = cache.Get(context.TODO(), client.CacheKey(req))
	require.NoError(t, err)

	<-time.After(1 * time.Second)

	// Perform the request again, expecting the server to be hit
	req, err = http.NewRequestWithContext(context.TODO(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	res, err = client.DoCached(req)
	require.NoError(t, err)

	// Server was hit
	assert.Equal(t, int32(2), handledRequests.Load(), "expect cache miss")

	// Headers are kept
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "bar", res.Header.Get("X-Foo"))

	// Body is kept
	body, err = io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)
	assert.Equal(t, []byte("bar"), body)
}
