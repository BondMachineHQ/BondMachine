package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/BondMachineHQ/BondMachine/pkg/simbox"
)

var debug = flag.Bool("d", false, "Debug")
var verbose = flag.Bool("v", false, "Verbose")

var sbFile = flag.String("simbox-file", "", "Filename of the simulation data file")
var bondmachineFile = flag.String("bondmachine-file", "", "Bondmachine in JSON format")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {
	flag.Parse()
}

func main() {
	sBox := new(simbox.Simbox)

	if *sbFile != "" {
		if _, err := os.Stat(*sbFile); err == nil {
			// Open the simbox file is exists
			if sbJSON, err := os.ReadFile(*sbFile); err == nil {
				if err := json.Unmarshal([]byte(sbJSON), sBox); err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}

		// Write the simbox file
		f, errI := os.Create(*sbFile)
		check(errI)
		defer f.Close()
		b, errJ := json.Marshal(sBox)
		check(errJ)
		_, errI = f.WriteString(string(b))
		check(errI)
	} else {
		panic("simbox-file is a mandatory option")
	}
}
