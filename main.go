package main

import "github.com/yields/phony/pkg/phony"
import "github.com/tj/docopt"
import "math/rand"
import "io/ioutil"
import "strconv"
import "strings"
import "regexp"
import "time"
import "fmt"
import "os"

var usage = `
  Usage: phony
    [--tick d]
    [--max n]
    [--list]

    phony -h | --help
    phony -v | --version
    
  Options:
    --list          list all available generators
    --max n         generate data up to n [default: -1]
    --tick d        generate data every d [default: 10ms]
    -v, --version   show version information
    -h, --help      show help information

`

func main() {
	args, err := docopt.Parse(usage, nil, true, "0.0.1", false)
	check(err)

	if args["--list"].(bool) {
		all := phony.List()
		println()
		for _, name := range all {
			fmt.Printf("  %s\n", name)
		}
		println()
		os.Exit(0)
	}

	rand.Seed(time.Now().UnixNano())

	d := parseDuration(args["--tick"].(string))
	max := parseInt(args["--max"].(string))
	tmpl := readAll(os.Stdin)
	tick := time.Tick(d)
	f := compile(string(tmpl))
	it := 0

	for _ = range tick {
		fmt.Fprintf(os.Stdout, "%s", f())
		if it++; -1 != max && it == max {
			return
		}
	}
}

func compile(tmpl string) func() string {
	expr, err := regexp.Compile(`({{ *(([a-z]+(\.[a-z]+)?)+) *}})`)
	check(err)
	return func() string {
		return expr.ReplaceAllStringFunc(tmpl, func(s string) string {
			path := strings.Trim(s[2:len(s)-2], " ")
			return phony.Get(path)
		})
	}
}

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "phony: %s\n", err.Error())
		os.Exit(1)
	}
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	check(err)
	return d
}

func readAll(r *os.File) string {
	b, err := ioutil.ReadAll(r)
	check(err)
	return string(b)
}
