package scorecard

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httptest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientGetScorecard(t *testing.T) {
	testCases := []struct {
		Repository string
		Requests   []httptest.Request
		Expected   *Scorecard
		Error      bool
	}{
		{
			Repository: "github.com/jacobalberty/unifi-docker",
			Requests: []httptest.Request{
				{
					Expectations: httptest.RequestExpectations{
						URL: "https://api.scorecard.dev/projects/github.com/jacobalberty/unifi-docker",
					},
					Response: httptest.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       "./testdata/scorecard-happy-path.json",
					},
				},
			},
			Expected: &Scorecard{
				Date:  "2025-02-24",
				Score: 5.3,
			},
		},
		{
			Repository: "github.com/non-existing/non-existing",
			Requests: []httptest.Request{
				{
					Expectations: httptest.RequestExpectations{
						URL: "https://api.scorecard.dev/projects/github.com/non-existing/non-existing",
					},
					Response: httptest.Response{
						Status:     "NOT FOUND",
						StatusCode: 404,
					},
				},
			},
			Expected: nil,
			Error:    false,
		},
		{
			Repository: "github.com/request-error/request-error",
			Requests: []httptest.Request{
				{
					Response: httptest.Response{
						Error: fmt.Errorf("request error"),
					},
				},
			},
			Expected: nil,
			Error:    true,
		},
		{
			Repository: "github.com/body-error/body-error",
			Requests: []httptest.Request{
				{
					Response: httptest.Response{
						Body: fmt.Errorf("body error"),
					},
				},
			},
			Expected: nil,
			Error:    true,
		},
		{
			Repository: "github.com/invalid-body/invalid-body",
			Requests: []httptest.Request{
				{
					Response: httptest.Response{
						Body: "{invalid json}",
					},
				},
			},
			Expected: nil,
			Error:    true,
		},
		{
			Repository: "github.com/server-error/server-error",
			Requests: []httptest.Request{
				{
					Response: httptest.Response{
						StatusCode: 500,
					},
				},
			},
			Expected: nil,
			Error:    true,
		},
		{
			// Invalid URL, meaning new HTTP requests will fail
			Repository: "github.com/_20_%+off_60000_/_20_%+off_60000_",
			Expected:   nil,
			Error:      true,
		},
		{
			Repository: "",
			Expected:   nil,
			Error:      true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Repository, func(t *testing.T) {
			httpClient := &httptest.Client{
				T:        t,
				Requests: testCase.Requests,
			}

			client := &Client{Client: httpClient}

			actual, err := client.GetScorecard(context.TODO(), testCase.Repository)

			if testCase.Error {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, testCase.Expected, actual)
		})
	}
}

func TestIntegrationClientGetScorecard(t *testing.T) {
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

func TestIntegrationClientGetScorecardNotFound(t *testing.T) {
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
