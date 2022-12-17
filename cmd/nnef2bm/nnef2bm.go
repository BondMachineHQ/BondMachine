package main

import (
	"io/ioutil"

	"github.com/BondMachineHQ/BondMachine/pkg/nnef2bm"
)

func main() {
	// Setup the input

	inputdata, _ := ioutil.ReadFile("small_net3.nnef")

	nnef2bm.NnefBuildBM(string(inputdata))
}
