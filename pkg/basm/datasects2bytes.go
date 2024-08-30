package basm

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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

				re := regexp.MustCompile("^db$")
				if re.MatchString(dataOperator) {
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
					continue
				}

				re = regexp.MustCompile("^(?P<num>[0-9]+):db$")
				if re.MatchString(dataOperator) {
					if decodedBytes, err := dbDataConverter(dataValue); err != nil {
						return err
					} else {
						numS := re.ReplaceAllString(dataOperator, "${num}")
						num, _ := strconv.Atoi(numS)
						decoded := len(decodedBytes)
						newElements := make([]*bmline.BasmElement, num*decoded)
						for i, obj := range decodedBytes {
							for j := 0; j < num; j++ {
								newArg := new(bmline.BasmElement)
								newArg.SetValue(obj)
								newElements[j*decoded+i] = newArg
							}
						}
						line.Elements = newElements
						offset += num * len(decodedBytes)
					}
					continue
				}
				re = regexp.MustCompile("^dd$")
				if re.MatchString(dataOperator) {
					if decodedBytes, err := dbDataConverter(dataValue); err != nil {
						return err
					} else {
						newElements := make([]*bmline.BasmElement, 0)
						currHex := ""
						for i, obj := range decodedBytes {
							if i%4 == 0 {
								if i != 0 {
									newArg := new(bmline.BasmElement)
									newArg.SetValue("0x" + currHex)
									newElements = append(newElements, newArg)
									currHex = ""
								}
							}
							currHex += strings.TrimPrefix(obj, "0x")
						}
						newArg := new(bmline.BasmElement)
						newArg.SetValue("0x" + currHex)
						newElements = append(newElements, newArg)

						line.Elements = newElements
						offset += len(newElements)
					}
					continue
				}
				re = regexp.MustCompile("^sym$")
				if re.MatchString(dataOperator) {
					// The symbol is resolved in the next steps
					newElements := make([]*bmline.BasmElement, 1)
					newArg := new(bmline.BasmElement)
					newArg.SetValue(dataValue)
					newElements[0] = newArg
					line.Elements = newElements
					// the amount of bytes is not known yet, anyway, the symbol will be stored in a single cell so the offset is incremented by 1
					offset += 1
					continue
				}
				//	case "equ":
				// TODO: finish this and the other data operators
				return errors.New("Unknown data operator " + dataOperator)
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
