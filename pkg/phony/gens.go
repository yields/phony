package phony

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
}
