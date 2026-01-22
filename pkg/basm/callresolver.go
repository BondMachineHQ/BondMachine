package basm

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"golang.org/x/exp/maps"
)

// pass to resolve call instructions to their appropriate rom/ram versions and handle fragment inlining. It will only
// target sections of type romtext and ramtext.
func callResolver(bi *BasmInstance) error {
	symbolTagger(bi)

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
						fragName := arg.GetValue()

						// Make a copy of the fragment body
						fBody := bi.fragments[fragName].fragmentBody.Copy()

						if isTemplate(fBody.Flat()) {
							params := make(map[string]string)
							maps.Copy(params, section.sectionBody.LoopMeta())
							maps.Copy(params, line.LoopMeta())
							if err := bodyTemplateResolver(fBody, params); err != nil {
								return err
							}
						}
						// Rewrite the call to a jump to the fragment
						section.sectionBody.Lines[i].Operation.SetValue("call" + mod + stack)
						arg.BasmMeta = arg.SetMeta("type", "symbol")

						// It has to be a new symbol, otherwise it would have been caught by the previous checks
						// Append the fragment to the newcode
						for j, appLine := range fBody.Lines {
							newLine := appLine.Copy()
							if j == 0 {
								if symbols := newLine.GetMeta("symbol"); symbols != "" {
									newSymbols := make([]string, len(strings.Split(symbols, ":"))+1)
									for k, s := range strings.Split(symbols, ":") {
										newSymbols[k] = fragName + s
									}
									newSymbols[len(newSymbols)-1] = fragName
									newLine.BasmMeta = newLine.SetMeta("symbol", strings.Join(newSymbols, ":"))
								} else {
									newLine.BasmMeta = newLine.SetMeta("symbol", fragName)
								}
								bi.symbols[symbolPrefix+sectName+"."+fragName] = -1
							} else {
								if symbols := newLine.GetMeta("symbol"); symbols != "" {
									newSymbols := make([]string, len(strings.Split(symbols, ":")))
									for k, s := range strings.Split(symbols, ":") {
										newSymbols[k] = fragName + s
									}
									newLine.BasmMeta = newLine.SetMeta("symbol", strings.Join(newSymbols, ":"))
								}
							}

							for _, el := range newLine.Elements {
								s := el.GetValue()
								// if s is a symbol local to the fragment rewrite it
								if _, ok := bi.symbols["frag."+fragName+"."+s]; ok {
									el.SetValue(fragName + s)
								}
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
