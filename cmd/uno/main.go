package main

import (
	"flag"

	"github.com/appetito/uno/gen"
)

func main() {
	var file string
	flag.StringVar(&file, "f", "", "YAML file path")
	flag.Parse()

	gen.Gen(file)
}