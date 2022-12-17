package main

import (
	"flag"

	"github.com/BondMachineHQ/BondMachine/pkg/brvga/brvgasdl"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var constraint = flag.String("constraints", "", "Use a BM constraint file to set the VGA up")
var sockFile = flag.String("sock", "/tmp/brvga.sock", "Socket file to use")
var headerPath = flag.String("header", "", "Header file to use")
var fontsPath = flag.String("fonts", "", "Fonts file to use")

func init() {
	flag.Parse()
	if *constraint == "" {
		panic("Must specify a constraint string")
	}
	if *headerPath == "" {
		panic("Must specify a header file")
	}
	if *fontsPath == "" {
		panic("Must specify a fonts file")
	}
}

func main() {
	if mem, err := brvgasdl.NewBrvgaSdlUnixSock(*constraint, *sockFile, *headerPath, *fontsPath); err != nil {
		panic(err)
	} else {
		mem.Run()
		defer mem.Close()
	}
}
