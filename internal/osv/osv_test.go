package osv

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizedSeverityCompare(t *testing.T) {
	expected := []NormalizedSeverity{
		NormalizedSeverityCritical,
		NormalizedSeverityHigh,
		NormalizedSeverityMedium,
		NormalizedSeverityLow,
		NormalizedSeverityUnspecified,
		NormalizedSeverity("unknown or unsupported"),
	}

	// Property-based testing
	for range 100 {
		actual := slices.Clone(expected)

		rand.Shuffle(len(actual), func(i, j int) {
			actual[i], actual[j] = actual[j], actual[i]
		})

		slices.SortFunc(actual, func(a NormalizedSeverity, b NormalizedSeverity) int {
			return a.Compare(b)
		})

		assert.Equal(t, expected, actual)
	}
}
