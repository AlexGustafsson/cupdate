package worker

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"path/filepath"
	"sync"
	"testing"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type removingRequester struct {
	store *store.Store
	once  sync.Once
}

func (r *removingRequester) Do(req *http.Request) (*http.Response, error) {
	var err error
	r.once.Do(func() {
		_, err = r.store.DeleteNonPresent(req.Context(), nil)
	})
	if err != nil {
		return nil, err
	}

	return nil, errors.New("stop workflow after removing image")
}

func (r *removingRequester) DoCached(req *http.Request) (*http.Response, error) {
	return r.Do(req)
}

func TestProcessRawImageRemovedDuringWorkflow(t *testing.T) {
	uri := "file:" + filepath.ToSlash(filepath.Join(t.TempDir(), "sqlite.db"))
	require.NoError(t, store.Initialize(t.Context(), uri))

	store, err := store.New(t.Context(), uri, false)
	require.NoError(t, err)
	defer store.Close()

	reference, err := oci.ParseReference("example.com/test/image:latest")
	require.NoError(t, err)

	_, err = store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: reference.String(),
	})
	require.NoError(t, err)
	exists, err := store.RawImageExists(t.Context(), reference.String())
	require.NoError(t, err)
	require.True(t, exists)

	var logs bytes.Buffer
	logger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() {
		slog.SetDefault(logger)
	})

	worker := New(&removingRequester{store: store}, store, httputil.NewAuthMux(), nil)
	require.NoError(t, worker.ProcessRawImage(t.Context(), reference))

	assert.Contains(t, logs.String(), "Image removed during processing")
	assert.NotContains(t, logs.String(), "FOREIGN KEY constraint failed")
	exists, err = store.RawImageExists(t.Context(), reference.String())
	require.NoError(t, err)
	assert.False(t, exists)
}
