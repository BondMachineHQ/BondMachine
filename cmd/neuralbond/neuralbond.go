package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/BondMachineHQ/BondMachine/neuralbond"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var registerSize = flag.Int("register-size", 8, "Number of bits per register (n-bit)")

var saveBondMachine = flag.String("save-bondmachine", "", "Create a BondMachine JSON file")

var netFile = flag.String("net-file", "", "JSON description of the net")

func init() {
	flag.Parse()
	if *saveBondMachine == "" {
		*saveBondMachine = "a.out.json"
	}
}

func main() {
	net := new(neuralbond.TrainedNet)

	// Load net from a JSON file
	if *netFile != "" {
		if netFileJSON, err := ioutil.ReadFile(*netFile); err == nil {
			if err := json.Unmarshal(netFileJSON, net); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	fmt.Println(net)
}
