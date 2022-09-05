package main

import (
	"flag"
	"fmt"

	"github.com/BondMachineHQ/BondMachine/brvga"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var constraint = flag.String("constraints", "", "Use a BM constraint file to set the VGA up")
var sockFile = flag.String("sock", "/tmp/brvga.sock", "Socket file to use")

func init() {
	flag.Parse()
	if *constraint == "" {
		panic("Must specify a constraint string")
	}
}

func main() {
	if mem, err := brvga.NewBrvgaTextMemory(*constraint); err != nil {
		panic(err)
	} else {
		fmt.Print(mem.Dump())
	}
}
