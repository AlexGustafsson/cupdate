package platform

import (
	"testing"

	"github.com/AlexGustafsson/cupdate/internal/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabelsIgnore(t *testing.T) {
	testCases := []struct {
		Labels   Labels
		Expected bool
	}{
		{
			Labels: Labels{
				"config.cupdate/ignore": "true",
			},
			Expected: true,
		},
		{
			Labels: Labels{
				"config.cupdate/ignore": "false",
			},
			Expected: false,
		},
		{
			Labels: Labels{
				"cupdate.config.ignore": "true",
			},
			Expected: true,
		},
		{
			Labels: Labels{
				"cupdate.config.ignore": "false",
			},
			Expected: false,
		},
		{
			Labels:   Labels{},
			Expected: false,
		},
		{
			Labels:   nil,
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		actual := testCase.Labels.Ignore()
		assert.Equal(t, testCase.Expected, actual)
	}
}

func TestLabelsPin(t *testing.T) {
	testCases := []struct {
		Labels   Labels
		Expected bool
	}{
		{
			Labels: Labels{
				"config.cupdate/pin": "true",
			},
			Expected: true,
		},
		{
			Labels: Labels{
				"config.cupdate/pin": "false",
			},
			Expected: false,
		},
		{
			Labels: Labels{
				"cupdate.config.pin": "true",
			},
			Expected: true,
		},
		{
			Labels: Labels{
				"cupdate.config.pin": "false",
			},
			Expected: false,
		},
		{
			Labels:   Labels{},
			Expected: false,
		},
		{
			Labels:   nil,
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		actual := testCase.Labels.Pin()
		assert.Equal(t, testCase.Expected, actual)
	}
}

func TestLabelStayOnCurrentMajor(t *testing.T) {
	testCases := []struct {
		Labels   Labels
		Expected bool
	}{
		{
			Labels: Labels{
				"config.cupdate/stay-on-current-major": "true",
			},
			Expected: true,
		},
		{
			Labels: Labels{
				"config.cupdate/stay-on-current-major": "false",
			},
			Expected: false,
		},
		{
			Labels: Labels{
				"cupdate.config.stay-on-current-major": "true",
			},
			Expected: true,
		},
		{
			Labels: Labels{
				"cupdate.config.stay-on-current-major": "false",
			},
			Expected: false,
		},
		{
			Labels:   Labels{},
			Expected: false,
		},
		{
			Labels:   nil,
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		actual := testCase.Labels.StayOnCurrentMajor()
		assert.Equal(t, testCase.Expected, actual)
	}
}

func TestLabelsStayBelow(t *testing.T) {
	testCases := []struct {
		Labels        Labels
		Expected      string
		ExpectedError bool
	}{
		{
			Labels: Labels{
				"config.cupdate/stay-below": "8.0.0",
			},
			Expected:      "8.0.0",
			ExpectedError: false,
		},
		{
			Labels: Labels{
				"config.cupdate/stay-below": "8.0.0-rc.1",
			},
			Expected:      "",
			ExpectedError: true,
		},
		{
			Labels: Labels{
				"cupdate.config.stay-below": "8.0.0",
			},
			Expected: "8.0.0",
		},
		{
			Labels: Labels{
				"cupdate.config.stay-below": "8.0.0-rc.1",
			},
			Expected:      "",
			ExpectedError: true,
		},
		{
			Labels:        Labels{},
			Expected:      "",
			ExpectedError: false,
		},
		{
			Labels:        nil,
			Expected:      "",
			ExpectedError: false,
		},
	}

	for _, testCase := range testCases {
		actual, err := testCase.Labels.StayBelow()
		if testCase.Expected == "" {
			assert.Nil(t, actual)
		} else {
			expected, err := semver.ParseVersion(testCase.Expected)
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		}
		if testCase.ExpectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
