package phony

import (
	"testing"
	"time"

	"strconv"

	"github.com/stretchr/testify/assert"
)

func TestGetReturnsNonEmptyString(t *testing.T) {
	g := NewGenerator()

	a, _ := g.Get("name", nil)
	assert.NotEqual(t, a, "")
}

func TestGetWithArgsBehavesAsGetByDefault(t *testing.T) {
	g := NewGenerator()

	a, _ := g.GetWithArgs("name", []string{}, nil)
	assert.NotEqual(t, a, "")
}

func TestSmartdouble(t *testing.T) {
	g := NewGenerator()

	a, err := g.GetWithArgs("smartdouble", []string{"100", "10000", "1000", "1000000"}, nil)
	assert.Nil(t, err)
	assert.NotEqual(t, a, "")

	num, err := strconv.ParseFloat(a, 64)
	assert.Nil(t, err)

	assert.True(t, num >= 1000.0, "Generated number is smaller than 1000")
	assert.True(t, num <= 1000000.0, "Generated number is larger than 1000000")
}

func TestSmartunixtimeDeviation(t *testing.T) {
	g := NewGenerator()

	a, err := g.GetWithArgs("smartunixtime", []string{"100"}, nil)
	assert.Nil(t, err)

	b, err := g.GetWithArgs("smartunixtime", []string{"1"}, nil)
	assert.Nil(t, err)

	aInt, err := strconv.ParseInt(a, 10, 64)
	assert.Nil(t, err)

	bInt, err := strconv.ParseInt(b, 10, 64)
	assert.Nil(t, err)

	assert.True(t, aInt > bInt, "The first date should be larger than the second")

	bIntPlus99Days := bInt + 99*int64(time.Hour)*24
	assert.True(t, aInt < bIntPlus99Days, "The first date should be less than 99 days ahead of the second date")
}

func TestSmartunixtimeScatter(t *testing.T) {
	g := NewGenerator()

	reference := time.Now().UnixNano()
	oneHourAgo := reference - int64(time.Minute)*10
	oneHourLater := reference + int64(time.Minute)*10

	rawTimestamp, err := g.GetWithArgs("smartunixtime", []string{"0", "100"}, nil)
	assert.Nil(t, err)

	timestamp, err := strconv.ParseInt(rawTimestamp, 10, 64)
	assert.Nil(t, err)

	assert.True(t, timestamp < oneHourAgo || timestamp > oneHourLater, "Time should be significantly different from now")

	hundredDaysAgo := reference - int64(time.Hour)*24*100
	hundredDaysLater := reference + int64(time.Hour)*24*100
	assert.True(t, timestamp > hundredDaysAgo && timestamp < hundredDaysLater, "Time should be within the given range")
}

func TestSmartdate(t *testing.T) {
	g := NewGenerator()

	format := "SqlDatetime"
	date, err := g.GetWithArgs("smartdate", []string{format}, nil)
	assert.Nil(t, err)

	today := time.Now().Format(supportedDateFormats[format])

	assert.Equal(t, date, today)
}

func TestEmpty(t *testing.T) {
	g := NewGenerator()

	a, _ := g.Get("foo", nil)
	assert.Equal(t, a, "")
}

func TestAll(t *testing.T) {
	g := NewGenerator()

	for _, p := range g.List() {
		a, _ := g.Get(p, nil)
		assert.NotEqual(t, a, "")
	}
}
