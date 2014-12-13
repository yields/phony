package phony

import "math/rand"

// Default generator.
var gen = New(&Dataset{
	gens: gens,
	dict: dict,
})

// Dataset.
type Dataset struct {
	gens map[string]func(g *Generator) string
	dict map[string][]string
}

// Generator structure.
type Generator struct {
	set *Dataset
}

// Initialize Generator with `dataset`.
func New(set *Dataset) *Generator {
	return &Generator{set}
}

// Get `path`.
func (g *Generator) Get(p string) string {
	gens := g.set.gens
	dict := g.set.dict

	for k, f := range gens {
		if k == p {
			return f(g)
		}
	}

	for k, list := range dict {
		if k == p {
			i := rand.Intn(len(list))
			return list[i]
		}
	}

	return ""
}

// List all paths.
func (g *Generator) List() []string {
	gens := g.set.gens
	dict := g.set.dict
	ret := make([]string, 0)

	for k, _ := range gens {
		ret = append(ret, k)
	}

	for k, _ := range dict {
		ret = append(ret, k)
	}

	return ret
}

// Get `path`.
func Get(path string) string {
	return gen.Get(path)
}

// List all available paths.
func List() []string {
	return gen.List()
}
