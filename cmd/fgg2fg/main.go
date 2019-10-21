package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/rhu1/fgg/fg"
	"github.com/rhu1/fgg/fgg"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "not enough arguments (expected FGG file path)")
		os.Exit(1)
	}
	b, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	fggAdaptor := new(fgg.FGGAdaptor)
	fggProg := fggAdaptor.Parse(false, string(b))

	obliterate(fggProg.(fgg.FGGProgram))
}

func obliterate(prog fgg.FGGProgram) fg.FGProgram {
	return fg.FGProgram{}
}
