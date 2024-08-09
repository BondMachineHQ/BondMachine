package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BondMachineHQ/BondMachine/pkg/bcof"
	"google.golang.org/protobuf/proto"
)

func init() {
	flag.Parse()
}

func main() {

	for _, bcofFile := range flag.Args() {
		// Read the BCOF file.
		bcofBytes, err := os.ReadFile(bcofFile)
		if err != nil {
		}
		bcof := new(bcof.BCOFEntry)
		if err := proto.Unmarshal(bcofBytes, bcof); err != nil {
		}

		fmt.Println(bcof)
	}
}
