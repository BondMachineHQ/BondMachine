package main

import (
	"github.com/BondMachineHQ/BondMachine/nnef2bm"
	"io/ioutil"
)

func main() {
	// Setup the input

	inputdata, _ := ioutil.ReadFile("small_net3.nnef")

	nnef2bm.NnefBuildBM(string(inputdata))
}
