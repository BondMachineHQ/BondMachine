package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/BondMachineHQ/BondMachine/neuralbond"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var registerSize = flag.Int("register-size", 8, "Number of bits per register (n-bit)")

var saveBasm = flag.String("save-basm", "", "Create a basm file")

var netFile = flag.String("net-file", "", "JSON description of the net")

func init() {
	flag.Parse()
	if *saveBasm == "" {
		*saveBasm = "out.basm"
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
	} else {
		panic("No net file specified")
	}

	// fmt.Println(net.Weights)
	net.Normalize()

	// fmt.Println(net)

	if *saveBasm != "" {
		if basmFile, err := net.WriteBasm(); err == nil {
			ioutil.WriteFile(*saveBasm, []byte(basmFile), 0644)
		} else {
			panic(err)
		}
	}

}
