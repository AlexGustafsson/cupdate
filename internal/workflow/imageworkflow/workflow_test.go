package imageworkflow

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/stretchr/testify/require"
)

func TestWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	cache, err := cache.NewDiskCache(filepath.Join(t.TempDir(), "cache.boltdb"))
	require.NoError(t, err)

	httpClient := httputil.NewClient(cache, 5*time.Minute)

	reference, err := oci.ParseReference("mongo:4")
	require.NoError(t, err)

	data := &Data{
		ImageReference:  reference,
		Image:           "",
		LatestReference: &reference,
		Tags:            make([]string, 0),
		Links:           make([]models.ImageLink, 0),
		RegistryAuth:    httputil.NewAuthMux(),
	}

	workflow := New(httpClient, data)
	_, err = workflow.Run(context.TODO())

	require.NoError(t, err)

	json.NewEncoder(os.Stderr).Encode(data)
}
