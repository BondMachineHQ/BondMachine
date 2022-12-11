package main

import "flag"

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

func init() {
	flag.Parse()
}

func main() {

}
