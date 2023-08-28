package basm

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

func symbolResolver(bi *BasmInstance) error {

	// Filter the matchers to select only the symbol based one
	filteredMatchers := make([]*bmline.BasmLine, 0)

	if bi.debug {
		fmt.Println(green("\tFiltering matchers:"))
	}

	for _, matcher := range bi.matchers {
		if bmline.FilterMatcher(matcher, "symbol") {
			filteredMatchers = append(filteredMatchers, matcher)
			if bi.debug {
				fmt.Println(red("\t\tActive matcher:") + matcher.String())
			}
		} else {
			if bi.debug {
				fmt.Println(yellow("\t\tInactive matcher:") + matcher.String())
			}
		}
	}

	if bi.debug {
		fmt.Println(green("\tProcessing sections:"))
	}

	// Loop over the sections TODO over the cps (the unused sections are supposed to be removed at this point)
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			body := section.sectionBody

			for _, line := range body.Lines {

				for _, arg := range line.Elements {
					if arg.GetMeta("type") == "symbol" {
						symbol := arg.GetValue()

						// Search the symbol in local symbols
						localSymbol := ""
						if section.sectionType == sectRomText {
							localSymbol = "rom." + sectName + "." + symbol
						} else {
							localSymbol = "ram." + sectName + "." + symbol
						}

						if loc, ok := bi.symbols[localSymbol]; ok {
							// Apply the correction if any
							// if body.GetMeta("symbcorrection") != "" {
							// 	correction, _ := strconv.Atoi(body.GetMeta("symbcorrection"))
							// 	loc += int64(correction)
							// }
							arg.SetValue(strconv.Itoa(int(loc)))
							arg.SetMeta("type", "number")
							continue
						}
						// TODO: Finish this
						return errors.New("symbol not found: " + symbol)
					}
				}
			}
		} else {
			if bi.debug {
				fmt.Println(yellow("\t\tSection type not handled: ") + sectName)
			}
		}
	}
	return nil
}
