package store

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInitializeNew(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	uri := "file://" + t.TempDir() + "/sqlite.db"

	err := Initialize(context.TODO(), uri)
	require.NoError(t, err)

	store, err := New(uri, true)
	require.NoError(t, err)

	err = store.Close()
	require.NoError(t, err)
}

func TestInitializeRandom(t *testing.T) {
	for range 100 {
		uri := "file://" + t.TempDir() + "/sqlite.db"

		// Try a few "false starts"
		for range 5 {
			ctx, cancel := context.WithCancel(context.TODO())
			go func() {
				<-time.After(time.Duration(rand.IntN(1 * 1e6)))
				cancel()
			}()

			err := Initialize(ctx, uri)
			if err == nil {
				break
			}
		}

		// Expect initialization to work regardless of previous failures
		err := Initialize(context.TODO(), uri)
		require.NoError(t, err)

		store, err := New(uri, true)
		require.NoError(t, err)

		err = store.Close()
		require.NoError(t, err)
	}
}
