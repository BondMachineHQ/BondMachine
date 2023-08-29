package basm

import (
	"errors"
	"fmt"
	"strconv"
)

func symbolResolver(bi *BasmInstance) error {

	if bi.debug {
		fmt.Println(green("\tProcessing sections:"))
	}

	// Loop over the sections TODO over the cps (the unused sections are supposed to be removed at this point)
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}
			if err := bi.resolveSymbols(section, ""); err != nil {
				return err
			}
		} else {
			if bi.debug {
				fmt.Println(yellow("\t\tSection type not handled: ") + sectName)
			}
		}
	}
	return nil
}

func (bi *BasmInstance) resolveSymbols(section *BasmSection, name string) error {
	// If name is empty, it means that the section name is used
	// If name is not empty, named composed on it are used according the convention
	// Using name is useful to handle the romcode/romdata sections in combined mode

	body := section.sectionBody

	for _, line := range body.Lines {
		for _, arg := range line.Elements {

			// Search the symbol in local symbols
			if arg.GetMeta("type") == "symbol" {
				symbol := arg.GetValue()
				localSymbol := ""

				if name == "" {
					if section.sectionType == sectRomText {
						localSymbol = "rom." + section.sectionName + "." + symbol
					} else {
						localSymbol = "ram." + section.sectionName + "." + symbol
					}
				} else {
					if section.sectionType == sectRomText {
						localSymbol = "rom.romcode" + name + "." + symbol
					} else {
						localSymbol = "ram.ramcode" + name + "." + symbol
					}
				}

				if loc, ok := bi.symbols[localSymbol]; ok {
					arg.SetValue(strconv.Itoa(int(loc)))
					arg.BasmMeta = arg.BasmMeta.SetMeta("type", "number")
					arg.BasmMeta = arg.BasmMeta.SetMeta("numbertype", "unsigned")
					continue
				}
			}

			// Search the symbol in rom
			if arg.GetMeta("type") == "rom" && arg.GetMeta("romaddressing") == "symbol" {
				symbol := arg.GetMeta("symbol")
				if symbol == "" {
					return errors.New("ROM symbol cannot be empty")
				}

				// In romdata
				romSymbol := ""
				if name == "" {
					romSymbol = "romdata." + section.sectionName + "." + symbol
				} else {
					romSymbol = "romdata.romdata" + name + "." + symbol
				}
				if loc, ok := bi.symbols[romSymbol]; ok {
					arg.SetValue(strconv.Itoa(int(loc)))
					arg.BasmMeta = arg.BasmMeta.SetMeta("type", "number")
					arg.BasmMeta = arg.BasmMeta.SetMeta("numbertype", "unsigned")
					arg.BasmMeta.RmMeta("romaddressing")
					arg.BasmMeta.RmMeta("symbol")
					continue
				}

				// In romcode
				romSymbol = ""
				if name == "" {
					romSymbol = "rom." + section.sectionName + "." + symbol
				} else {
					romSymbol = "romcode.romcode" + name + "." + symbol
				}
				if loc, ok := bi.symbols[romSymbol]; ok {
					arg.SetValue(strconv.Itoa(int(loc)))
					arg.BasmMeta = arg.BasmMeta.SetMeta("type", "number")
					arg.BasmMeta = arg.BasmMeta.SetMeta("numbertype", "unsigned")
					arg.BasmMeta.RmMeta("romaddressing")
					arg.BasmMeta.RmMeta("symbol")
					continue
				}

				// return errors.New("symbol not found: " + symbol)

			}
		}
	}
	return nil
}
