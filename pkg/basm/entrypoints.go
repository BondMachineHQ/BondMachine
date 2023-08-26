package basm

import (
	"errors"
	"fmt"
	"strconv"
)

// section entry points detection, the pass detects the symbol used as entry point of the section and sign it as metadata.
func entryPoints(bi *BasmInstance) error {
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\tProcessing section: " + sectName))
			}

			checkEntry := false
			checkSymbol := ""
			checkLine := 0
			checkExistence := false
			checkPosition := -1

			prevSymbols := make([]string, 0)

			body := section.sectionBody
			for i, line := range body.Lines {

				varName := line.Operation.GetValue()

				if varName == "entry" {
					if checkEntry {
						return errors.New("multiple entry points specified")
					}
					checkEntry = true

					if len(line.Elements) != 1 {
						return errors.New("line has too many elements")
					}

					checkSymbol = line.Elements[0].GetValue()
					checkLine = i

					for _, symbol := range prevSymbols {
						if symbol == checkSymbol {
							if checkExistence {
								return errors.New("multiple entry points detected")
							} else {
								checkExistence = true
							}
						}
					}
				}

				if checkEntry {
					if symbol := line.GetMeta("symbol"); symbol != "" {
						if symbol == checkSymbol {
							if checkExistence {
								return errors.New("multiple entry points detected")
							} else {
								checkExistence = true
								checkPosition = i - 1
							}
						}
					}
				} else {
					if symbol := line.GetMeta("symbol"); symbol != "" {
						prevSymbols = append(prevSymbols, symbol)
					}
				}
			}

			if !checkEntry {
				return errors.New("entry point not specified")
			}

			if !checkExistence {
				return errors.New("entry point not detected")
			}

			// Removing the entry directive line
			copy(body.Lines[checkLine:], body.Lines[checkLine+1:])
			body.Lines[len(body.Lines)-1] = nil
			body.Lines = body.Lines[:len(body.Lines)-1]
			// Removing the line means that symbols are shifted, so we need to correct the symbol correction
			if body.GetMeta("symbcorrection") != "" {
				correction, _ := strconv.Atoi(body.GetMeta("symbcorrection"))
				correction--
				body.BasmMeta = body.BasmMeta.SetMeta("symbcorrection", fmt.Sprintf("%d", correction))
			} else {
				body.BasmMeta = body.BasmMeta.SetMeta("symbcorrection", "-1")
			}

			if checkPosition < 0 {
				for i, line := range body.Lines {
					if symbol := line.GetMeta("symbol"); symbol != "" {
						if symbol == checkSymbol {
							checkPosition = i
							break
						}
					}
				}
			}

			body.BasmMeta = body.SetMeta("entry", fmt.Sprintf("%d", checkPosition))

		}
	}
	return nil
}
