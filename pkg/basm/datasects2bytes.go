package basm

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

func dataSections2Bytes(bi *BasmInstance) error {

	for sectName, section := range bi.sections {
		if section.sectionType == sectRomData || section.sectionType == sectRamData {
			if bi.debug {
				fmt.Println(green("\tProcessing section: " + sectName))
			}

			if section.sectionType == sectRomData {
				bi.rg.Requirement(bmreqs.ReqRequest{Node: "code:romdatas", T: bmreqs.ObjectSet, Name: "sections", Value: section.sectionName, Op: bmreqs.OpAdd})
			} else {
				bi.rg.Requirement(bmreqs.ReqRequest{Node: "code:ramdatas", T: bmreqs.ObjectSet, Name: "sections", Value: section.sectionName, Op: bmreqs.OpAdd})
			}

			varCheck := make(map[string]struct{})

			body := section.sectionBody
			offset := 0
			for _, line := range body.Lines {
				varName := line.Operation.GetValue()
				symbolName := ""

				if section.sectionType == sectRomData {
					symbolName = "romdata." + sectName + "." + varName
				} else {
					symbolName = "ramdata." + sectName + "." + varName
				}

				if bi.debug {
					fmt.Println(green("\t\tvar " + varName))
				}

				if _, exists := varCheck[varName]; exists {
					return errors.New("Duplicate var " + varName)
				} else {
					varCheck[varName] = struct{}{}
					bi.symbols[symbolName] = -1
				}

				if len(line.Elements) != 2 {
					return errors.New("data elements expects 2 arguments")
				}

				dataOperator := line.Elements[0].GetValue()
				dataValue := line.Elements[1].GetValue()
				line.Operation.BasmMeta = line.Operation.BasmMeta.SetMeta("offset", fmt.Sprintf("%d", offset))

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
						offset += len(decodedBytes)
					}
				case "sym":
					// The symbol is resolved in the next steps
					newElements := make([]*bmline.BasmElement, 1)
					newArg := new(bmline.BasmElement)
					newArg.SetValue(dataValue)
					newElements[0] = newArg
					line.Elements = newElements
					// the amount of bytes is not known yet, anyway, the symbol will be stored in a single cell so the offset is incremented by 1
					offset += 1
				case "equ":
					// TODO: finish this and the other data operators
				default:
					return errors.New("Unknown data operator " + dataOperator)
				}
			}

			if section.sectionType == sectRomData {
				bi.rg.Requirement(bmreqs.ReqRequest{Node: "code:romdatas/sections:" + section.sectionName, T: bmreqs.ObjectMax, Name: "datalength", Value: strconv.Itoa(offset), Op: bmreqs.OpAdd})
			} else {
				bi.rg.Requirement(bmreqs.ReqRequest{Node: "code:ramdatas/sections:" + section.sectionName, T: bmreqs.ObjectMax, Name: "datalength", Value: strconv.Itoa(offset), Op: bmreqs.OpAdd})
			}
		}
	}
	return nil
}
