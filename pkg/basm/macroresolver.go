package basm

import "fmt"

func macroResolver(bi *BasmInstance) error {

	// Loop over the sections
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			body := section.sectionBody

			for _, line := range body.Lines {
				op := line.Operation.GetValue()
				if _, ok := bi.macros[op]; ok {
					if bi.debug {
						fmt.Println(yellow("\t\t\tMacro: ") + op)
					}
				}
			}
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

		body := fragment.fragmentBody

		for _, line := range body.Lines {
			op := line.Operation.GetValue()
			if _, ok := bi.macros[op]; ok {
				if bi.debug {
					fmt.Println(yellow("\t\t\tMacro: ") + op)
				}
			}
		}
	}

	return nil

}
