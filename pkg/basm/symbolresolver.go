package basm

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func symbolResolver(bi *BasmInstance) error {

	if bi.debug {
		fmt.Println(green("\tProcessing sections:"))
	}

	// Loop over the sections TODO over the cps (the unused sections are supposed to be removed at this point)
	// This are the sections that are not combined
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

	// Loop over the cps to resolve the symbols of the combined sections
	for _, cp := range bi.cps {
		romCodeName := cp.GetMeta("romcode")
		if strings.HasPrefix(romCodeName, "romcode") {
			if section, ok := bi.sections[romCodeName]; ok {
				name := romCodeName[7:]
				if err := bi.resolveSymbols(section, name); err != nil {
					return err
				}
			} else {
				return errors.New("romcode section not found: " + romCodeName)
			}
		}
		ramCodeName := cp.GetMeta("ramcode")
		if strings.HasPrefix(ramCodeName, "ramcode") {
			if section, ok := bi.sections[ramCodeName]; ok {
				name := ramCodeName[7:]
				if err := bi.resolveSymbols(section, name); err != nil {
					return err
				}
			} else {
				return errors.New("ramcode section not found: " + ramCodeName)
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

			// Search the revsymbols
			if arg.GetMeta("type") == "revsymbol" {
				symbol := arg.GetValue()
				localSymbol := ""

				if name == "" {
					// To catch the revsymbol, a guess is made on the name
					guessedName := ""
					re := regexp.MustCompile("^romcode(?P<name>[0-9a-zA-Z_]+)$")
					if re.MatchString(section.sectionName) {
						guessedName = re.ReplaceAllString(section.sectionName, "${name}")
					}
					re = regexp.MustCompile("^ramcode(?P<name>[0-9a-zA-Z_]+)$")
					if re.MatchString(section.sectionName) {
						guessedName = re.ReplaceAllString(section.sectionName, "${name}")
					}
					if guessedName != "" {
						if section.sectionType == sectRomText {
							localSymbol = "ram.ramcode" + guessedName + "." + symbol
						} else {
							localSymbol = "rom.romcode" + guessedName + "." + symbol
						}
					}
				} else {
					if section.sectionType == sectRomText {
						localSymbol = "ram.ramcode" + name + "." + symbol
					} else {
						localSymbol = "rom.romcode" + name + "." + symbol
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
					romSymbol = "rom.romcode" + name + "." + symbol
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

			// Search the symbol in ram
			if arg.GetMeta("type") == "ram" && arg.GetMeta("ramaddressing") == "symbol" {
				symbol := arg.GetMeta("symbol")
				if symbol == "" {
					return errors.New("RAM symbol cannot be empty")
				}

				// In ramdata
				ramSymbol := ""
				if name == "" {
					ramSymbol = "ramdata." + section.sectionName + "." + symbol
				} else {
					ramSymbol = "ramdata.ramdata" + name + "." + symbol
				}
				if loc, ok := bi.symbols[ramSymbol]; ok {
					arg.SetValue(strconv.Itoa(int(loc)))
					arg.BasmMeta = arg.BasmMeta.SetMeta("type", "number")
					arg.BasmMeta = arg.BasmMeta.SetMeta("numbertype", "unsigned")
					arg.BasmMeta.RmMeta("ramaddressing")
					arg.BasmMeta.RmMeta("symbol")
					continue
				}

				// In ramcode
				ramSymbol = ""
				if name == "" {
					ramSymbol = "ram." + section.sectionName + "." + symbol
				} else {
					ramSymbol = "ram.ramcode" + name + "." + symbol
				}
				if loc, ok := bi.symbols[ramSymbol]; ok {
					arg.SetValue(strconv.Itoa(int(loc)))
					arg.BasmMeta = arg.BasmMeta.SetMeta("type", "number")
					arg.BasmMeta = arg.BasmMeta.SetMeta("numbertype", "unsigned")
					arg.BasmMeta.RmMeta("ramaddressing")
					arg.BasmMeta.RmMeta("symbol")
					continue
				}
			}
		}
	}
	return nil
}
