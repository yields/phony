package phony

import "math/rand"
import "strconv"
import "time"
import "fmt"

// Default gens.
var gens = map[string]func(g *Generator) string{
	"name": func(g *Generator) string {
		a := g.Get("name.first")
		b := g.Get("name.last")
		return a + " " + b
	},
	"email": func(g *Generator) string {
		username := g.Get("username")
		host := g.Get("domain")
		return username + "@" + host
	},
	"domain": func(g *Generator) string {
		name := g.Get("domain.name")
		tld := g.Get("domain.tld")
		return name + "." + tld
	},
	"avatar": func(g *Generator) string {
		// http://uifaces.com/authorized
		user := g.Get("username")
		return "https://s3.amazonaws.com/uifaces/faces/twitter/" + user + "/128.jpg"
	},
	"unixtime": func(g *Generator) string {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	},
	"id": func(g *Generator) string {
		chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		ret := make([]rune, 10)

		for i := range ret {
			ret[i] = chars[rand.Intn(len(chars))]
		}

		return string(ret)
	},
	"ipv4": func(g *Generator) string {
		return fmt.Sprintf("%d.%d.%d.%d", 1+rand.Intn(253), rand.Intn(255), rand.Intn(255), 1+rand.Intn(253))
	},
	"ipv6": func(g *Generator) string {
		return fmt.Sprintf("2001:cafe:%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	},
	"mac.address": func(g *Generator) string {
		return fmt.Sprintf("%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	},
	"latitude": func(g *Generator) string {
		lattitude := (rand.Float64() * 180) - 90
		return strconv.FormatFloat(lattitude, 'f', 6, 64)
	},
	"longitude": func(g *Generator) string {
		longitude := (rand.Float64() * 360) - 180
		return strconv.FormatFloat(longitude, 'f', 6, 64)
	},
	"double": func(g *Generator) string {
		return strconv.FormatFloat(rand.NormFloat64()*1000, 'f', 4, 64)
	},
}
