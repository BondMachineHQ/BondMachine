package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/BondMachineHQ/BondMachine/pkg/simbox"
)

type string_slice []string

func (i *string_slice) String() string {
	return fmt.Sprint(*i)
}

func (i *string_slice) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		*i = append(*i, dt)
	}
	return nil
}

var debug = flag.Bool("d", false, "Debug")
var verbose = flag.Bool("v", false, "Verbose")

var sbFile = flag.String("simbox-file", "", "Filename of the simulation data file")
var machineFile = flag.String("machine-file", "", "Machine in JSON format")
var bondmachineFile = flag.String("bondmachine-file", "", "Bondmachine in JSON format")

var verify = flag.Bool("verify", false, "Verify the simbox against a machine file or a bondmachine file")

var list = flag.Bool("list", false, "List rules")
var add = flag.String("add", "", "Add e rule")
var del = flag.Int("del", -1, "Remove a rule")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {
	rand.Seed(int64(time.Now().Unix()))
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
		if *verify {
			if *machineFile != "" {
				// TODO machine verify
			} else if *bondmachineFile != "" {
				// TODO bondmachine verify
			} else {
				panic("Missing machine or bondmachine file")
			}
		} else if *list {
			fmt.Print(sBox.Print())
		} else if *add != "" {
			err := sBox.Add(*add)
			check(err)
		} else if *del != -1 {
			err := sBox.Del(*del)
			check(err)
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
