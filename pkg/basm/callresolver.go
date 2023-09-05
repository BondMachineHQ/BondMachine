package basm

import (
	"fmt"
	"regexp"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

// section entry points detection, the pass detects the symbol used as entry point of the section and sign it as metadata.
func callResolver(bi *BasmInstance) error {
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\tProcessing section: " + sectName))
			}

			symbolPrefix := ""
			mod := ""
			if section.sectionType == sectRomText {
				symbolPrefix = "rom."
				mod = "o"
			} else {
				symbolPrefix = "ram."
				mod = "a"
			}

			newCode := make([]*bmline.BasmLine, 0)

			for i, line := range section.sectionBody.Lines {
				// Calls are only one element
				if len(line.Elements) != 1 {
					continue
				}
				opName := line.Operation.GetValue()
				re := regexp.MustCompile("^call(?P<stack>[0-9]+[a-zA-Z_]+)$")
				if re.MatchString(opName) {
					stack := re.ReplaceAllString(opName, "${stack}")

					if bi.debug {
						fmt.Println(green("\t\tCall found: ") + line.String())
						fmt.Println(green("\t\t\tStack: ") + stack)
					}

					arg := line.Elements[0]

					// If the type is already set, skip
					if arg.GetMeta("type") != "" {
						continue
					}

					// Check if the arg is an already defined symbol. If so, set the type change the operation to the rom/ram version of the call and continue
					if _, ok := bi.symbols[symbolPrefix+sectName+"."+arg.GetValue()]; ok {
						section.sectionBody.Lines[i].Operation.SetValue("call" + mod + stack)
						arg.BasmMeta = arg.SetMeta("type", "symbol")
						continue
					}

					// Check if the arg is a fragment
					if _, ok := bi.fragments[arg.GetValue()]; ok {
						// Rewrite the call to a jump to the fragment
						section.sectionBody.Lines[i].Operation.SetValue("call" + mod + stack)
						arg.BasmMeta = arg.SetMeta("type", "symbol")
						// It has to be a new symbol, otherwise it would have been caught by the previous check
						// Append the fragment to the newcode
						for j, appLine := range bi.fragments[arg.GetValue()].fragmentBody.Lines {
							newLine := appLine.Copy()
							if j == 0 {
								newLine.BasmMeta = newLine.SetMeta("symbol", arg.GetValue())
								bi.symbols[symbolPrefix+sectName+"."+arg.GetValue()] = -1
							}
							newCode = append(newCode, newLine)
						}
						newLine := new(bmline.BasmLine)
						newOp := new(bmline.BasmElement)
						newOp.SetValue("ret" + stack)
						newLine.Operation = newOp
						elements := make([]*bmline.BasmElement, 0)
						newLine.Elements = elements
						newCode = append(newCode, newLine)
						continue
					} else {
						return fmt.Errorf("call to undefined symbol: %s", arg.GetValue())
					}
				}
			}
			// Append the new code to the section
			section.sectionBody.Lines = append(section.sectionBody.Lines, newCode...)
		}
	}

	return nil
}
