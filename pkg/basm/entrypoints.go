package basm

import (
	"errors"
	"fmt"
)

// section entry points detection, the pass detects the label used as entry point of the section and sign it as metadata.
func entryPoints(bi *BasmInstance) error {
	for sectName, section := range bi.sections {
		if section.sectionType == setcRomText {
			if bi.debug {
				fmt.Println(green("\tProcessing section: " + sectName))
			}

			checkEntry := false
			checkLabel := ""
			checkLine := 0
			checkExistence := false
			checkPosition := -1

			prevLabels := make([]string, 0)

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

					checkLabel = line.Elements[0].GetValue()
					checkLine = i

					for _, label := range prevLabels {
						if label == checkLabel {
							if checkExistence {
								return errors.New("multiple entry points detected")
							} else {
								checkExistence = true
							}
						}
					}
				}

				if checkEntry {
					if label := line.GetMeta("label"); label != "" {
						if label == checkLabel {
							if checkExistence {
								return errors.New("multiple entry points detected")
							} else {
								checkExistence = true
								checkPosition = i - 1
							}
						}
					}
				} else {
					if label := line.GetMeta("label"); label != "" {
						prevLabels = append(prevLabels, label)
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

			if checkPosition < 0 {
				for i, line := range body.Lines {
					if label := line.GetMeta("label"); label != "" {
						if label == checkLabel {
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
