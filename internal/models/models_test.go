package models

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeverityCompare(t *testing.T) {
	expected := []Severity{
		SeverityCritical,
		SeverityHigh,
		SeverityMedium,
		SeverityLow,
		SeverityUnspecified,
		Severity("unknown or unsupported"),
	}

	// Property-based testing
	for range 100 {
		actual := slices.Clone(expected)

		rand.Shuffle(len(actual), func(i, j int) {
			actual[i], actual[j] = actual[j], actual[i]
		})

		slices.SortFunc(actual, func(a Severity, b Severity) int {
			return a.Compare(b)
		})

		assert.Equal(t, expected, actual)
	}
}
