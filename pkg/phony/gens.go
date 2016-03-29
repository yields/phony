package phony

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	smartdouble   = "smartdouble"
	smartunixtime = "smartunixtime"
	smartdate     = "smartdate"
)

var supportedDateFormats = map[string]string{
	// Go default formats
	"ANSIC":       "Mon Jan _2 15:04:05 2006",
	"UnixDate":    "Mon Jan _2 15:04:05 MST 2006",
	"RubyDate":    "Mon Jan 02 15:04:05 -0700 2006",
	"RFC822":      "02 Jan 06 15:04 MST",
	"RFC822Z":     "02 Jan 06 15:04 -0700",
	"RFC850":      "Monday, 02-Jan-06 15:04:05 MST",
	"RFC1123":     "Mon, 02 Jan 2006 15:04:05 MST",
	"RFC1123Z":    "Mon, 02 Jan 2006 15:04:05 -0700",
	"RFC3339":     "2006-01-02T15:04:05Z07:00",
	"RFC3339Nano": "2006-01-02T15:04:05.999999999Z07:00",
	"Kitchen":     "3:04PM",
	"Stamp":       "Jan _2 15:04:05",
	"StampMilli":  "Jan _2 15:04:05.000",
	"StampMicro":  "Jan _2 15:04:05.000000",
	"StampNano":   "Jan _2 15:04:05.000000000",
	// Additional formats
	"SqlDatetime": "2006-01-02 15:04:05",
	"SqlDate":     "2006-01-02",
	"SqlTime":     "15:04:05",
}

type Generator struct {
	sync.Mutex
	cachedArguments map[string]map[string]interface{}
	argumentGetters map[string]func(g *Generator, args []string) interface{}
	generators      map[string]func(g *Generator, args []string, r *rand.Rand) (string, error)
	dictionary      map[string][]string
}

func (g *Generator) getSmartDate(args []string, r *rand.Rand) time.Time {
	date := time.Now()

	floatArgs := g.argumentGetters[smartunixtime](g, args).([]float64)

	// Deviation
	if len(floatArgs) > 0 {
		addDays := int64(floatArgs[0] * float64(time.Hour) * 24.0)
		date = date.Add(time.Duration(addDays))
	}

	// Scatter
	if len(floatArgs) > 1 {
		scatterFactor := (r.Float64() - 0.5) * 2
		scatterDays := int64(floatArgs[1] * float64(time.Hour) * 24.0 * scatterFactor)

		date = date.Add(time.Duration(scatterDays))
	}

	return date
}

var gen *Generator

func NewGenerator() *Generator {
	if gen != nil {
		return gen
	}

	var generator = Generator{}

	generator.argumentGetters = map[string]func(g *Generator, args []string) interface{}{
		smartdouble: func(g *Generator, args []string) interface{} {
			key := strings.Join(args, ",")

			if found, ok := g.cachedArguments[smartdouble][key]; ok {
				return found
			}

			g.Lock()
			defer g.Unlock()
			g.cachedArguments[smartdouble][key] = stringsToFloats(args)

			return g.cachedArguments[smartdouble][key]
		},
		smartunixtime: func(g *Generator, args []string) interface{} {
			key := strings.Join(args, ",")

			if found, ok := g.cachedArguments[smartdouble][key]; ok {
				return found
			}

			g.Lock()
			defer g.Unlock()
			g.cachedArguments[smartunixtime][key] = stringsToFloats(args)

			return g.cachedArguments[smartunixtime][key]
		},
	}

	generator.cachedArguments = map[string]map[string]interface{}{
		smartdouble:   map[string]interface{}{},
		smartunixtime: map[string]interface{}{},
	}

	// Default gens.
	generator.generators = map[string]func(g *Generator, args []string, r *rand.Rand) (string, error){
		"name": func(g *Generator, args []string, r *rand.Rand) (string, error) {
			a, _ := g.Get("name.first", r)
			b, _ := g.Get("name.last", r)
			return a + " " + b, nil
		},
		"email": func(g *Generator, args []string, r *rand.Rand) (string, error) {
			username, _ := g.Get("username", r)
			host, _ := g.Get("domain", r)
			return fmt.Sprintf("%s%d@%s", username, r.Intn(253), host), nil
		},
		"domain": func(g *Generator, args []string, r *rand.Rand) (string, error) {
			name, _ := g.Get("domain.name", r)
			tld, _ := g.Get("domain.tld", r)
			return name + "." + tld, nil
		},
		"avatar": func(g *Generator, args []string, r *rand.Rand) (string, error) {
			// http://uifaces.com/authorized
			user, _ := g.Get("username", r)
			return "https://s3.amazonaws.com/uifaces/faces/twitter/" + user + "/128.jpg", nil
		},
		"unixtime": func(g *Generator, args []string, r *rand.Rand) (string, error) {
			return strconv.FormatInt(time.Now().UnixNano(), 10), nil
		},
		"id": func(g *Generator, args []string, r *rand.Rand) (string, error) {
			chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
			ret := make([]rune, 10)

			for i := range ret {
				ret[i] = chars[r.Intn(len(chars))]
			}

			return string(ret), nil
		},
		"ipv4": func(g *Generator, args []string, r *rand.Rand) (string, error) {
			return fmt.Sprintf("%d.%d.%d.%d", 1+r.Intn(253), r.Intn(255), r.Intn(255), 1+r.Intn(253)), nil
		},
		"ipv6": func(g *Generator, args []string, r *rand.Rand) (string, error) {
			return fmt.Sprintf("2001:cafe:%x:%x:%x:%x:%x:%x", r.Intn(255), r.Intn(255), r.Intn(255), r.Intn(255), r.Intn(255), r.Intn(255)), nil
		},
		"mac.address": func(g *Generator, args []string, r *rand.Rand) (string, error) {
			return fmt.Sprintf("%x:%x:%x:%x:%x:%x", r.Intn(255), r.Intn(255), r.Intn(255), r.Intn(255), r.Intn(255), r.Intn(255)), nil
		},
		"latitude": func(g *Generator, args []string, r *rand.Rand) (string, error) {
			lattitude := (r.Float64() * 180) - 90
			return strconv.FormatFloat(lattitude, 'f', 6, 64), nil
		},
		"longitude": func(g *Generator, args []string, r *rand.Rand) (string, error) {
			longitude := (r.Float64() * 360) - 180
			return strconv.FormatFloat(longitude, 'f', 6, 64), nil
		},
		"double": func(g *Generator, args []string, r *rand.Rand) (string, error) {
			return strconv.FormatFloat(r.NormFloat64()*1000, 'f', 4, 64), nil
		},
		// Smartdouble returns a random double
		// First argument will be interpreted as desired standard deviation
		// Second argument will be interpreted as desired mean
		// Third argument will be interpreted as minimum value expected
		// Fourth argument will be interpreted as maximum value expected
		//
		// More info:
		// https://golang.org/pkg/math/rand/#Rand.NormFloat64
		smartdouble: func(g *Generator, args []string, r *rand.Rand) (string, error) {
			var (
				desiredStdDev = 1000.0
				desiredMean   = 0.0
			)

			floatArgs := g.argumentGetters[smartdouble](g, args).([]float64)

			if len(floatArgs) > 0 {
				desiredStdDev = floatArgs[0]
			}

			if len(floatArgs) > 1 {
				desiredMean = floatArgs[1]
			}

			randNum := r.NormFloat64()*desiredStdDev + desiredMean

			if len(floatArgs) > 2 && randNum < floatArgs[2] {
				randNum = floatArgs[2]
			}

			if len(floatArgs) > 3 && randNum > floatArgs[3] {
				randNum = floatArgs[3]
			}

			return strconv.FormatFloat(randNum, 'f', 4, 64), nil
		},
		// Smartunixtime returns a random unix time based on the current time, args can be an array of strings.
		// Each string is expected to be float64 parsable. First argument will be interpreted as deviation days,
		// second argument will be interpreted as days for allowed range in days.
		//
		// Example:
		//   Today is: 2006-01-02T15:04:05Z07:00
		//   Args: ["10", "5"]
		//     Deviation is: 10 (days)
		//     Scatter is: 5 (days)
		//   Result: unix timestamp representing date between 2006-01-07T15:04:05Z07:00 and 2006-01-17T15:04:05Z07:00
		smartunixtime: func(g *Generator, args []string, r *rand.Rand) (string, error) {
			date := g.getSmartDate(args, r)

			return strconv.FormatInt(date.UnixNano(), 10), nil
		},
		// Smartdate returns a random date based on the current time, args can be an array of strings.
		// First argument will be interpreted as expected date format
		// Second argument will be interpreted as deviation days,
		// Third argument will be interpreted as days for allowed range in days.
		//
		// Example:
		//   Today is: 2006-01-02T15:04:05Z07:00
		//   Args: ["SQL_DATE", "10", "5"]
		//     Deviation is: 10 (days)
		//     Scatter is: 5 (days)
		//   Result: string date between 2006-01-07 and 2006-01-17
		smartdate: func(g *Generator, args []string, r *rand.Rand) (string, error) {
			var (
				format       = "RFC3339"
				unixtimeArgs = []string{}
				ok           bool
			)

			if len(args) > 0 {
				format = args[0]
				unixtimeArgs = args[1:]
			}

			if format, ok = supportedDateFormats[format]; !ok {
				return "", errors.New("Invalid date format.")
			}

			date := g.getSmartDate(unixtimeArgs, r)

			return date.Format(format), nil
		},
	}

	generator.dictionary = dict

	gen = &generator

	return gen
}
