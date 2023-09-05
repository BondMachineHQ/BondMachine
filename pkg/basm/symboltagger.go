package basm

import (
	"errors"
	"fmt"
)

func symbolTagger(bi *BasmInstance) error {

	if bi.debug {
		fmt.Println(green("\tProcessing sections:"))
	}

	// Loop over the sections
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			symbolPrefix := ""
			if section.sectionType == sectRomText {
				symbolPrefix = "rom." + sectName + "."
			} else {
				symbolPrefix = "ram." + sectName + "."
			}

			// Map all the symbols
			symbols := make(map[string]struct{})

			body := section.sectionBody
			for i, line := range body.Lines {
				if symbol := line.GetMeta("symbol"); symbol != "" {
					if _, exists := symbols[symbol]; exists {
						return errors.New("symbol is specified multiple time: " + symbol)
					} else {
						symbols[symbol] = struct{}{}
						bi.symbols[symbolPrefix+symbol] = int64(i)
					}
				}
			}

			// for _, line := range body.Lines {

			// 	for _, arg := range line.Elements {
			// 		if _, ok := symbols[arg.GetValue()]; ok {
			// 			arg.BasmMeta = arg.SetMeta("type", "symbol")
			// 		}
			// 	}
			// }
		} else {
			if bi.debug {
				fmt.Println(yellow("\t\tSection type not handled: ") + sectName)
			}
		}
	}

	// Loop over the fragments
	for fragName, fragment := range bi.fragments {
		if bi.debug {
			fmt.Println(green("\t\tFragment: ") + fragName)
		}

		symbolPrefix := "frag." + fragName + "."

		// Map all the symbols
		symbols := make(map[string]struct{})

		body := fragment.fragmentBody
		for i, line := range body.Lines {
			if symbol := line.GetMeta("symbol"); symbol != "" {
				if _, exists := symbols[symbol]; exists {
					return errors.New("symbol is specified multiple time: " + symbol)
				} else {
					symbols[symbol] = struct{}{}
					bi.symbols[symbolPrefix+symbol] = int64(i)
				}
			}
		}

		for _, line := range body.Lines {

			for _, arg := range line.Elements {
				if _, ok := symbols[arg.GetValue()]; ok {
					arg.BasmMeta = arg.SetMeta("type", "symbol")
				}
			}
		}
	}

	return nil
}
