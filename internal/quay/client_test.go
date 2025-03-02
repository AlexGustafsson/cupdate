package quay

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httptest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVulnerabilities(t *testing.T) {
	testCases := []struct {
		Reference string
		Requests  []httptest.Request
		Expected  []Vulnerability
		Error     bool
	}{
		{
			Reference: "quay.io/openshift-release-dev/ocp-release@sha256:7708f832ae02919f2cdb2798fdbc64e17ce7a576d1e3baabdd78a000d2d62f40",
			Requests: []httptest.Request{
				{
					Expectations: httptest.RequestExpectations{
						URL: "https://quay.io/api/v1/repository/openshift-release-dev/ocp-release/manifest/sha256:7708f832ae02919f2cdb2798fdbc64e17ce7a576d1e3baabdd78a000d2d62f40/security?vulnerabilities=true",
					},
					Response: httptest.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       "./testdata/scan-happy-path.json",
					},
				},
			},
			Expected: []Vulnerability{
				{
					Name:        "CVE-2020-11023",
					Description: "A flaw was found in jQuery. HTML containing \\<option\\> elements from untrusted sources are passed, even after sanitizing, to one of jQuery's DOM manipulation methods, which may execute untrusted code. The highest threat from this vulnerability is to data confidentiality and integrity.",
					Links: []string{
						"https://access.redhat.com/security/cve/CVE-2020-11023",
						"https://bugzilla.redhat.com/show_bug.cgi?id=1850004",
						"https://www.cve.org/CVERecord?id=CVE-2020-11023",
						"https://nvd.nist.gov/vuln/detail/CVE-2020-11023",
						"https://blog.jquery.com/2020/04/10/jquery-3-5-0-released/",
						"https://www.cisa.gov/known-exploited-vulnerabilities-catalog",
						"https://security.access.redhat.com/data/csaf/v2/vex/2020/cve-2020-11023.json",
						"https://access.redhat.com/errata/RHSA-2025:1304",
					},
					FeatureName:    "libgomp",
					FeatureVersion: "11.4.1-3.el9",
					Layer:          "sha256:ae0badd537673e93bcbcf384ce6acda3cdfef75d43bd2f7bc766ef5ffba3e51a",
					Severity:       VulnerabilitySeverityMedium,
				},
			},
		},
		{
			Reference: "example.com/invalid-body/invalid-body@sha256:7708f832ae02919f2cdb2798fdbc64e17ce7a576d1e3baabdd78a000d2d62f40",
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
			Reference: "example.com/not-scanned/not-scanned@sha256:7708f832ae02919f2cdb2798fdbc64e17ce7a576d1e3baabdd78a000d2d62f40",
			Requests: []httptest.Request{
				{
					Response: httptest.Response{
						Body: `{"status": "not scanned"}`,
					},
				},
			},
			Expected: nil,
			Error:    false,
		},
		{
			Reference: "example.com/null-data/null-data@sha256:7708f832ae02919f2cdb2798fdbc64e17ce7a576d1e3baabdd78a000d2d62f40",
			Requests: []httptest.Request{
				{
					Response: httptest.Response{
						Body: `{"status": "scanned", "data": null}`,
					},
				},
			},
			Expected: nil,
			Error:    false,
		},
		{
			Reference: "example.com/non-existing/non-existing@sha256:7708f832ae02919f2cdb2798fdbc64e17ce7a576d1e3baabdd78a000d2d62f40",
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
			Reference: "example.com/server-error/server-error@sha256:7708f832ae02919f2cdb2798fdbc64e17ce7a576d1e3baabdd78a000d2d62f40",
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
			Reference: "example.com/no-digest/no-digest",
			Expected:  nil,
			Error:     true,
		},
		{
			Reference: "example.com/request-error/request-error@sha256:7708f832ae02919f2cdb2798fdbc64e17ce7a576d1e3baabdd78a000d2d62f40",
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.Reference, func(t *testing.T) {
			httpClient := &httptest.Client{
				T:        t,
				Requests: testCase.Requests,
			}

			ref, err := oci.ParseReference(testCase.Reference)
			require.NoError(t, err)

			client := &Client{Client: httpClient}

			actual, err := client.GetVulnerabilities(context.TODO(), ref)

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

	ref, err := oci.ParseReference("quay.io/openshift-release-dev/ocp-release@sha256:7708f832ae02919f2cdb2798fdbc64e17ce7a576d1e3baabdd78a000d2d62f40")
	require.NoError(t, err)

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	vulnerabilities, err := client.GetVulnerabilities(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Printf("%+v\n", vulnerabilities)
}
