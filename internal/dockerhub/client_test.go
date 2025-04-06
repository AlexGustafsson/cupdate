package dockerhub

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httptest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/osv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientGetRepository(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}
	ref, err := oci.ParseReference("mongo")
	require.NoError(t, err)
	repository, err := client.GetRepository(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Println(repository.FullDescription)

	json.NewEncoder(os.Stdout).Encode(repository)
}

func TestGetVulnerabilities(t *testing.T) {
	testCases := []struct {
		Repository string
		Digest     string
		Requests   []httptest.Request
		Expected   []osv.Vulnerability
		Error      bool
	}{
		{
			Repository: "traefik",
			Digest:     "sha256:ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078",
			Requests: []httptest.Request{
				{
					Expectations: httptest.RequestExpectations{
						URL: "https://api.dso.docker.com/v1/graphql",
					},
					Response: httptest.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       "./testdata/scan-happy-path.json",
					},
				},
			},
			Expected: []osv.Vulnerability{
				{
					ID:               "CVE-2023-45918",
					DatabaseSpecific: map[string]any{"severity": "LOW"},
					References: []osv.Reference{
						{
							Type: "WEB",
							URL:  "https://scout.docker.com/v/CVE-2023-45918?s=ubuntu&n=ncurses&ns=ubuntu&t=deb&osn=ubuntu&osv=22.04&vr=%3E%3D0",
						},
					},
					Summary:   "...",
					Withdrawn: nil,
				},
				{
					ID:               "CVE-2023-50495",
					DatabaseSpecific: map[string]any{"severity": "LOW"},
					References: []osv.Reference{
						{
							Type: "WEB",
							URL:  "https://scout.docker.com/v/CVE-2023-50495?s=ubuntu&n=ncurses&ns=ubuntu&t=deb&osn=ubuntu&osv=22.04&vr=%3E%3D0",
						},
					},
					Severities: []osv.Severity{
						{
							Type:  "CVSS_V3",
							Score: "6.50",
						},
					},
					Summary: "...",
				},
			},
		},
		{
			Repository: "invalid-body",
			Digest:     "1234",
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
			Repository: "not-scanned",
			Digest:     "1234",
			Requests: []httptest.Request{
				{
					Response: httptest.Response{
						Body: `{"data":{"imagePackageForImageCoords":{"sbomState":"NOT_INDEXED"}}}`,
					},
				},
			},
			Expected: nil,
			Error:    false,
		},
		{
			Repository: "not-existing",
			Digest:     "1234",
			Requests: []httptest.Request{
				{
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
			Repository: "server-error",
			Digest:     "1234",
			Requests: []httptest.Request{
				{
					Response: httptest.Response{
						Status:     "INTERNAL SERVER ERROR",
						StatusCode: 500,
					},
				},
			},
			Expected: nil,
			Error:    true,
		},
		{
			Repository: "request-error",
			Digest:     "1234", Requests: []httptest.Request{
				{
					Response: httptest.Response{
						Error: fmt.Errorf("request error"),
					},
				},
			},
			Expected: nil,
			Error:    true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Repository+"@"+testCase.Digest, func(t *testing.T) {
			httpClient := &httptest.Client{
				T:        t,
				Requests: testCase.Requests,
			}

			client := &Client{Client: httpClient}

			actual, err := client.GetVulnerabilities(context.TODO(), testCase.Repository, testCase.Digest)

			// Set modified time to zero as to not check it. We currently don't have
			// any means of finding the correct time from Docker Hub, so it's always
			// the current time
			for i := range actual {
				actual[i].Modified = time.Time{}
			}

			if testCase.Error {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, testCase.Expected, actual)
		})
	}
}

func TestIntegrationGetVulnerabilities(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	report, err := client.GetVulnerabilities(context.TODO(), "traefik", "sha256:ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078")
	require.NoError(t, err)

	json.NewEncoder(os.Stdout).Encode(report)
}
