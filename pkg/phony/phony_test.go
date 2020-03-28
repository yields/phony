package phony

import "github.com/bmizerany/assert"
import "testing"

func TestGet(t *testing.T) {
	a, _ := Get("name")
	b, _ := Get("name")
	assert.NotEqual(t, a, "")
	assert.NotEqual(t, b, "")
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
