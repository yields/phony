package phony

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	smartdouble   = "smartdouble"
	smartunixtime = "smartunixtime"
	smartdate     = "smartdate"
)

var cachedArgs = map[string]map[string]interface{}{
	smartdouble:   map[string]interface{}{},
	smartunixtime: map[string]interface{}{},
}

func stringsToFloats(args []string) []float64 {
	floatArgs := []float64{}

	for i := 0; i < len(args); i++ {
		num, err := strconv.ParseFloat(args[i], 64)

		if err != nil {
			break
		}

		floatArgs = append(floatArgs, num)
	}

	return floatArgs
}

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

var getGenArgs = map[string]func(args []string) interface{}{
	smartdouble: func(args []string) interface{} {
		key := strings.Join(args, ",")

		if found, ok := cachedArgs[smartdouble][key]; ok {
			return found
		}

		cachedArgs[smartdouble][key] = stringsToFloats(args)

		return cachedArgs[smartdouble][key]
	},
	smartunixtime: func(args []string) interface{} {
		key := strings.Join(args, ",")

		if found, ok := cachedArgs[smartdouble][key]; ok {
			return found
		}

		cachedArgs[smartunixtime][key] = stringsToFloats(args)

		return cachedArgs[smartunixtime][key]
	},
}

func getSmartDate(args []string) time.Time {
	date := time.Now()

	floatArgs := getGenArgs["smartunixtime"](args).([]float64)

	// Deviation
	if len(floatArgs) > 0 {
		addDays := int64(floatArgs[0] * float64(time.Hour) * 24.0)
		date = date.Add(time.Duration(addDays))
	}

	// Scatter
	if len(floatArgs) > 1 {
		scatterFactor := (rand.Float64() - 0.5) * 2
		scatterDays := int64(floatArgs[1] * float64(time.Hour) * 24.0 * scatterFactor)

		date = date.Add(time.Duration(scatterDays))
	}

	return date
}

// Default gens.
var gens = map[string]func(g *Generator, args []string) (string, error){
	"name": func(g *Generator, args []string) (string, error) {
		a, _ := g.Get("name.first")
		b, _ := g.Get("name.last")
		return a + " " + b, nil
	},
	"email": func(g *Generator, args []string) (string, error) {
		username, _ := g.Get("username")
		host, _ := g.Get("domain")
		return username + "@" + host, nil
	},
	"domain": func(g *Generator, args []string) (string, error) {
		name, _ := g.Get("domain.name")
		tld, _ := g.Get("domain.tld")
		return name + "." + tld, nil
	},
	"avatar": func(g *Generator, args []string) (string, error) {
		// http://uifaces.com/authorized
		user, _ := g.Get("username")
		return "https://s3.amazonaws.com/uifaces/faces/twitter/" + user + "/128.jpg", nil
	},
	"unixtime": func(g *Generator, args []string) (string, error) {
		return strconv.FormatInt(time.Now().UnixNano(), 10), nil
	},
	"id": func(g *Generator, args []string) (string, error) {
		chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		ret := make([]rune, 10)

		for i := range ret {
			ret[i] = chars[rand.Intn(len(chars))]
		}

		return string(ret), nil
	},
	"ipv4": func(g *Generator, args []string) (string, error) {
		return fmt.Sprintf("%d.%d.%d.%d", 1+rand.Intn(253), rand.Intn(255), rand.Intn(255), 1+rand.Intn(253)), nil
	},
	"ipv6": func(g *Generator, args []string) (string, error) {
		return fmt.Sprintf("2001:cafe:%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255)), nil
	},
	"mac.address": func(g *Generator, args []string) (string, error) {
		return fmt.Sprintf("%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255)), nil
	},
	"latitude": func(g *Generator, args []string) (string, error) {
		lattitude := (rand.Float64() * 180) - 90
		return strconv.FormatFloat(lattitude, 'f', 6, 64), nil
	},
	"longitude": func(g *Generator, args []string) (string, error) {
		longitude := (rand.Float64() * 360) - 180
		return strconv.FormatFloat(longitude, 'f', 6, 64), nil
	},
	"double": func(g *Generator, args []string) (string, error) {
		return strconv.FormatFloat(rand.NormFloat64()*1000, 'f', 4, 64), nil
	},
	// Smartdouble returns a random double
	// First argument will be interpreted as desired standard deviation
	// Second argument will be interpreted as desired mean
	// Third argument will be interpreted as minimum value expected
	// Fourth argument will be interpreted as maximum value expected
	//
	// More info:
	// https://golang.org/pkg/math/rand/#Rand.NormFloat64
	smartdouble: func(g *Generator, args []string) (string, error) {
		var (
			desiredStdDev = 1000.0
			desiredMean   = 0.0
		)

		floatArgs := getGenArgs["smartdouble"](args).([]float64)

		if len(floatArgs) > 0 {
			desiredStdDev = floatArgs[0]
		}

		if len(floatArgs) > 1 {
			desiredMean = floatArgs[1]
		}

		randNum := rand.NormFloat64()*desiredStdDev + desiredMean

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
	smartunixtime: func(g *Generator, args []string) (string, error) {
		date := getSmartDate(args)

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
	smartdate: func(g *Generator, args []string) (string, error) {
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

		date := getSmartDate(unixtimeArgs)

		return date.Format(format), nil
	},
}
