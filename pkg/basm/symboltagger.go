package basm

import (
	"errors"
	"fmt"
	"strings"
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

			// Map all the symbolList
			symbolList := make(map[string]struct{})

			body := section.sectionBody
			for i, line := range body.Lines {
				if symbols := line.GetMeta("symbol"); symbols != "" {
					for _, symbol := range strings.Split(symbols, ":") {
						if _, exists := symbolList[symbol]; exists {
							return errors.New("symbol is specified multiple time: " + symbol)
						} else {
							symbolList[symbol] = struct{}{}
							bi.symbols[symbolPrefix+symbol] = int64(i)
						}
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

		// Map all the symbolList
		symbolList := make(map[string]struct{})

		body := fragment.fragmentBody
		for i, line := range body.Lines {
			if symbols := line.GetMeta("symbol"); symbols != "" {
				for _, symbol := range strings.Split(symbols, ":") {
					if _, exists := symbolList[symbol]; exists {
						return errors.New("symbol is specified multiple time: " + symbol)
					} else {
						symbolList[symbol] = struct{}{}
						bi.symbols[symbolPrefix+symbol] = int64(i)
					}
				}
			}
		}

		for _, line := range body.Lines {

			for _, arg := range line.Elements {
				if _, ok := symbolList[arg.GetValue()]; ok {
					arg.BasmMeta = arg.SetMeta("type", "symbol")
				}
			}
		}
	}

	return nil
}
