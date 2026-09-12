package forgejo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/stretchr/testify/require"
)

func TestClientGetREADME(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	c := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	readme, err := c.GetREADME(context.TODO(), "https://codeberg.org/forgejo/forgejo")
	require.NoError(t, err)

	fmt.Printf("%s\n", readme)
}
