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

	// Loop over the sections
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			// Map all the symbols
			symbols := make(map[string]int)

			body := section.sectionBody
			for i, line := range body.Lines {
				if symbol := line.GetMeta("symbol"); symbol != "" {
					if _, exists := symbols[symbol]; exists {
						return errors.New("symbol is specified multiple time: " + symbol)
					} else {
						symbols[symbol] = i
					}
				}
			}

			for _, line := range body.Lines {

				for _, arg := range line.Elements {
					if _, ok := symbols[arg.GetValue()]; ok {
						arg.BasmMeta = arg.SetMeta("type", "symbol")
					}
				}

				for _, matcher := range filteredMatchers {
					if bmline.MatchMatcher(matcher, line) {
						// TODO Handling the operand
						for _, arg := range line.Elements {
							if lno, ok := symbols[arg.GetValue()]; ok {
								arg.SetValue(strconv.Itoa(lno))
								arg.BasmMeta = arg.SetMeta("type", "number")
							}
						}

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
