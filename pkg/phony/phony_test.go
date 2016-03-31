package phony

import (
	"testing"
	"time"

	"strconv"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	a, _ := Get("name")
	b, _ := Get("name")
	assert.NotEqual(t, "", a)
	assert.NotEqual(t, "", b)
}

func TestGetWithArgsBehavesAsGetByDefault(t *testing.T) {
	a, _ := GetWithArgs("name", []string{})
	assert.NotEqual(t, "", a)
}

func TestSmartdouble(t *testing.T) {
	a, err := GetWithArgs("smartdouble", []string{"100", "10000", "1000", "1000000"})
	assert.Nil(t, err)
	assert.NotEqual(t, "", a)

	num, err := strconv.ParseFloat(a, 64)
	assert.Nil(t, err)

	assert.True(t, num >= 1000.0, "Generated number is smaller than 1000")
	assert.True(t, num <= 1000000.0, "Generated number is larger than 1000000")
}

func TestSmartunixtimeDeviation(t *testing.T) {
	a, err := GetWithArgs("smartunixtime", []string{"100"})
	assert.Nil(t, err)

	b, err := GetWithArgs("smartunixtime", []string{"1"})
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
	reference := time.Now().UnixNano()
	oneHourAgo := reference - int64(time.Minute)*10
	oneHourLater := reference + int64(time.Minute)*10

	rawTimestamp, err := GetWithArgs("smartunixtime", []string{"0", "100"})
	assert.Nil(t, err)

	timestamp, err := strconv.ParseInt(rawTimestamp, 10, 64)
	assert.Nil(t, err)

	assert.True(t, timestamp < oneHourAgo || timestamp > oneHourLater, "Time should be significantly different from now")

	hundredDaysAgo := reference - int64(time.Hour)*24*100
	hundredDaysLater := reference + int64(time.Hour)*24*100
	assert.True(t, timestamp > hundredDaysAgo && timestamp < hundredDaysLater, "Time should be within the given range")
}

func TestSmartdate(t *testing.T) {
	format := "SqlDatetime"
	date, err := GetWithArgs("smartdate", []string{format})
	assert.Nil(t, err)

	today := time.Now().Format(supportedDateFormats[format])

	assert.Equal(t, date, today)
}

func TestEmpty(t *testing.T) {
	a, _ := Get("foo")
	assert.Equal(t, "", a)
}

func TestAll(t *testing.T) {
	for _, p := range List() {
		a, _ := Get(p)
		assert.NotEqual(t, "", a)
	}
}
