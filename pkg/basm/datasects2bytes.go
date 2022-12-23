package basm

import (
	"errors"
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

func dataSections2Bytes(bi *BasmInstance) error {
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomData {
			if bi.debug {
				fmt.Println(green("\tProcessing section: " + sectName))
			}

			varCheck := make(map[string]struct{})

			body := section.sectionBody
			for _, line := range body.Lines {
				varName := line.Operation.GetValue()

				if bi.debug {
					fmt.Println(green("\t\tvar " + varName))
				}

				if _, exists := varCheck[varName]; exists {
					return errors.New("Duplicate var " + varName)
				} else {
					varCheck[varName] = struct{}{}
				}

				if len(line.Elements) != 2 {
					return errors.New("Data elements expects 2 arguments")
				}

				dataOperator := line.Elements[0].GetValue()
				dataValue := line.Elements[1].GetValue()

				switch dataOperator {
				case "db":
					if decodedBytes, err := dbDataConverter(dataValue); err != nil {
						return err
					} else {
						newElements := make([]*bmline.BasmElement, len(decodedBytes))
						for i, obj := range decodedBytes {
							newArg := new(bmline.BasmElement)
							newArg.SetValue(obj)
							newElements[i] = newArg
						}
						line.Elements = newElements
					}
				case "equ":
				default:
					return errors.New("Unknown data operator " + dataOperator)
				}

			}
		}
	}
	return nil
}
