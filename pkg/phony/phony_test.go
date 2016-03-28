package phony

import "testing"
import "github.com/stretchr/testify/assert"
import "strconv"

func TestGet(t *testing.T) {
	a, _ := Get("name")
	b, _ := Get("name")
	assert.NotEqual(t, a, "")
	assert.NotEqual(t, b, "")
}

func TestGetWithArgsBehavesAsGetByDefault(t *testing.T) {
	a, _ := GetWithArgs("name", []string{})
	assert.NotEqual(t, a, "")
}

func TestSmartdouble(t *testing.T) {
	a, err := GetWithArgs("smartdouble", []string{"100", "10000", "1000", "1000000"})
	assert.Nil(t, err)
	assert.NotEqual(t, a, "")

	num, err := strconv.ParseFloat(a, 64)
	assert.Nil(t, err)
	assert.True(t, num >= 1000.0, "Generated number is smaller than 1000")
	assert.True(t, num <= 1000000.0, "Generated number is larger than 1000000")
}

func TestEmpty(t *testing.T) {
	a, _ := Get("foo")
	assert.Equal(t, a, "")
}

func TestAll(t *testing.T) {
	for _, p := range List() {
		a, _ := Get(p)
		assert.NotEqual(t, a, "")
	}
}
