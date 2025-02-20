package scorecard

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/stretchr/testify/require"
)

func TestClientGetScorecard(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	scorecard, err := client.GetScorecard(context.TODO(), "github.com/home-assistant/core")
	require.NoError(t, err)

	fmt.Println(scorecard.Score)
	if scorecard != nil {
		time, err := scorecard.Time()
		require.NoError(t, err)

		fmt.Println(time.String())
	}
}

func TestClientGetScorecardNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	scorecard, err := client.GetScorecard(context.TODO(), "github.com/non-existent/no-existent")
	require.NoError(t, err)
	require.Nil(t, scorecard)
}
