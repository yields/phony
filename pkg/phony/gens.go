package phony

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var smartdoubleArgs = map[string][]float64{}

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
	"smartdouble": func(g *Generator, args []string) (string, error) {
		var (
			desiredStdDev = 1000.0
			desiredMean   = 0.0
			num           = 0.0
			err           error
			floatArgs     = []float64{}
			ok bool
		)

		// Convert arguments to floats once and store them in private smartdoubleArgs
		if floatArgs, ok = smartdoubleArgs[strings.Join(args, ",")]; !ok {
			for i := 0; i < len(args); i++ {
				num, err = strconv.ParseFloat(args[i], 64)

				if err != nil {
					break
				}

				floatArgs = append(floatArgs, num)
			}

			smartdoubleArgs[strings.Join(args, ",")] = floatArgs
		}

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
}
