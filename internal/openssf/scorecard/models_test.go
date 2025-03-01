package scorecard

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScorecardTime(t *testing.T) {
	testCases := []struct {
		Date     string
		Expected time.Time
		Error    bool
	}{
		{
			Date:     "2025-02-27",
			Expected: time.Date(2025, 02, 27, 0, 0, 0, 0, time.UTC),
			Error:    false,
		},
		{
			Date:     "2025-02-27T17:41:30Z",
			Expected: time.Date(2025, 02, 27, 17, 41, 30, 0, time.UTC),
			Error:    false,
		},
		{
			Date:     "2025-02-27 10:00:00",
			Expected: time.Time{},
			Error:    true,
		},
		{
			Date:     "",
			Expected: time.Time{},
			Error:    true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Date, func(t *testing.T) {
			actual, err := (&Scorecard{
				Date: testCase.Date,
			}).Time()

			fmt.Println(testCase.Expected, actual)
			assert.Equal(t, 0, testCase.Expected.Compare(actual))
			if testCase.Error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
