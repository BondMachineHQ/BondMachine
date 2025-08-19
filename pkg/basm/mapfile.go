package basm

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

func (bi *BasmInstance) CreateMappingFile(filename string) error {
	ioMap := new(bondmachine.IOmap)
	ioMap.Assoc = make(map[string]string)

	for key, value := range bi.global.LoopMeta() {
		switch key {
		case "mapclk":
			ioMap.Assoc["clk"] = value
		case "mapreset":
			ioMap.Assoc["reset"] = value
		}
	}

	if len(ioMap.Assoc) != 2 {
		return fmt.Errorf("invalid setting for mapping file, both 'mapclk' and 'mapreset' must be set")
	}

	for _, ioatt := range bi.ioAttach {
		if ioatt.GetMeta("cp") == "bm" {
			fromS := ""
			toS := ""
			indexS := ""
			name := ""
			typeS := ""
			for key, value := range ioatt.LoopMeta() {
				switch key {
				case "mapfrom":
					fromS = value
				case "mapto":
					toS = value
				case "index":
					indexS = value
				case "mapname":
					name = value
				case "type":
					typeS = value
				}
			}
			if fromS == "" || toS == "" || indexS == "" || name == "" || typeS == "" {
				if bi.debug {
					fmt.Println("Skipping I/O attachment with missing mapping metadata")
				}
				continue
			}

			key := ""
			if typeS == "input" {
				key = "i" + indexS
			} else {
				key = "o" + indexS
			}

			value := ""

			if fromS != "" && toS != "" {
				value = "[" + toS + ":" + fromS + "]"
			}

			value += " " + name

			ioMap.Assoc[key] = value
		}
	}

	// Write the file
	mapBytes, err := json.Marshal(ioMap)
	if err != nil {
		return fmt.Errorf("failed to marshal I/O map: %v", err)
	}
	if err := os.WriteFile(filename, mapBytes, 0644); err != nil {
		return fmt.Errorf("failed to write I/O map file: %v", err)
	}

	return nil
}
