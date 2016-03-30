package phony

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringsToFloatParsesStringsAsFloats(t *testing.T) {
	a := []string{"10", "-100.1", "1e-3", "-1e3"}

	b := stringsToFloats(a)

	assert.Equal(t, []float64{10.0, -100.1, 0.001, -1000}, b)
}

func TestStringsToFloatStopsParsingOnFirstFailure(t *testing.T) {
	a := []string{"10", "abc", "1e-3", "-1e3"}

	b := stringsToFloats(a)

	assert.Equal(t, []float64{10.0}, b)
}
