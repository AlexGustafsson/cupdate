package imageworkflow

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/stretchr/testify/require"
)

func TestWorkflow(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	cache, err := cache.NewDiskCache("cache")
	require.NoError(t, err)

	httpClient := httputil.NewClient(cache, 5*time.Minute)

	reference, err := oci.ParseReference("mongo:4")
	require.NoError(t, err)

	data := &Data{
		ImageReference: reference,
		Image:          "",
		LatestVersion:  &reference,
		Tags:           make([]string, 0),
		Description:    nil,
		ReleaseNotes:   nil,
		Links:          make([]models.ImageLink, 0),
	}

	workflow := New(httpClient, data)

	require.NoError(t, workflow.Run(context.TODO()))

	json.NewEncoder(os.Stderr).Encode(data)
}
