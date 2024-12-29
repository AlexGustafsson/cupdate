package vulndb

import (
	"context"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	err := Fetch(context.TODO(), httputil.NewClient(cachetest.NewCache(t), 24*time.Hour), "vulndb.sqlite")
	require.NoError(t, err)
}
