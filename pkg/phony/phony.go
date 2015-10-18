package phony

import "math/rand"

// Default generator.
var gen = New(&Dataset{
	gens: gens,
	dict: dict,
})

// Dataset.
type Dataset struct {
	gens map[string]func(g *Generator, args []string) (string, error)
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
func (g *Generator) Get(p string) (string, error) {
	return g.GetWithArgs(p, nil)
}

// Get `path`.
func (g *Generator) GetWithArgs(p string, args []string) (string, error) {
	gens := g.set.gens
	dict := g.set.dict

	for k, f := range gens {
		if k == p {
			return f(g, args)
		}
	}

	for k, list := range dict {
		if k == p {
			i := rand.Intn(len(list))
			return list[i], nil
		}
	}

	return "", nil
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
func Get(path string) (string, error) {
	return gen.Get(path)
}

// Get `path`.
func GetWithArgs(path string, args []string) (string, error) {
	return gen.GetWithArgs(path, args)
}

// List all available paths.
func List() []string {
	return gen.List()
}
