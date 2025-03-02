package gitlab

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/stretchr/testify/require"
)

func TestGetRepositoryDescription(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	res, err := client.GetRepositoryDescription(context.TODO(), "arm-research/smarter/smarter-device-manager")
	require.NoError(t, err)

	fmt.Printf("%+v\n", res)
}

func TestGetRepositoryREADMEBlob(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	res, err := client.GetRepositoryREADMEBlob(context.TODO(), "/arm-research/smarter/smarter-device-manager/-/blob/master/README.md")
	require.NoError(t, err)

	if res != nil {
		fmt.Printf("%+s\n", res.Raw)
	}
}
