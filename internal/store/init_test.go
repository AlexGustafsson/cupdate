package store

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitializeNew(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	uri := "file://" + t.TempDir() + "/sqlite.db"

	err := Initialize(context.TODO(), uri)
	require.NoError(t, err)

	store, err := New(uri, true)
	require.NoError(t, err)
	store.Close()
}
