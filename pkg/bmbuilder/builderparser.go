package bmbuilder

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

// TODO Horrific and temporary code, a proper parser/lexer is desireable
func idiotParser(s string) ([]string, int) {
	// Convert tabs into spaces
	tab := regexp.MustCompile(`\t`)
	st := tab.ReplaceAllString(s, " ")
	// Strip away all duplicates whitspace characters and comments
	comment := regexp.MustCompile(`;.*`)
	//space := regexp.MustCompile(`\s+`)
	//stripped := strings.TrimSpace(space.ReplaceAllString(comment.ReplaceAllString(s, ""), " "))
	stripped := strings.TrimSpace(comment.ReplaceAllString(st, ""))

	// Splitting the line using spaces
	splitted := strings.Split(stripped, " ")

	return splitted, len(splitted)
}

func basmParser(bi *BMBuilder, s string, lineNo uint32) error {
	line := strconv.Itoa(int(lineNo))
	argS, argN := idiotParser(s)

	if bi.debug {
		fmt.Print("\t" + green(lineNo))
		fmt.Println(blue(": ") + s)

		fmt.Print("\t" + green(lineNo))
		for i, val := range argS {
			fmt.Println("\t", i, val)
		}
		fmt.Print("\t" + green(lineNo))
		fmt.Print(blue(": "))
	}

	// Parsing not empty lines
	if argN > 0 {
		// Identifying the main operand
		operand := argS[0]

		if operand != "" && operand != " " {

			if bi.debug {
				fmt.Print(argS)
			}

			if operand == "%meta" {
				if argN > 0 {
					if err := bi.metaProcessor(strings.Join(argS[1:], " ")); err != nil {
						return err
					}
				}

			} else if operand == "%block" {
				if bi.isWithinBlock != "" {
					return errors.New(line + ", macro, section or fragment definition is not closed")
				}
				// New section starting
				if argN < 3 {
					return errors.New(line + ", section definition has the wrong number of parameters")
				}

				blockName := argS[1]
				blockType := uint8(0)

				switch argS[2] {
				case ".sequential":
					blockType = blockSequential
				default:
					return errors.New(line + ", section definition has an unknown type " + argS[2])
				}

				newBlock := new(BMBuilderBlock)
				newBlock.blockName = blockName
				newBlock.blockType = blockType
				newBlock.blockBody = new(bmline.BasmBody)

				if argN > 3 {
					pairs := getPairs(strings.Join(argS[3:], ","))

					if len(pairs) > 0 {
						for _, pair := range pairs {
							keyVal := strings.Split(pair, ":")
							if len(keyVal) < 2 {
								return errors.New("Wrong format for entry (1) : " + pair)
							} else {
								newBlock.blockBody.BasmMeta = newBlock.blockBody.SetMeta(keyVal[0], strings.Join(keyVal[1:], ":"))
							}
						}
					}
				}

				bi.blocks[blockName] = newBlock

				bi.isWithinBlock = blockName

				if bi.debug {
					fmt.Print(yellow(" --> Starting section " + blockName + " definition"))
				}

			} else if operand == "%endblock" {
				// Section ending
				if bi.isWithinBlock == "" {
					return errors.New(line + ", endsection outside section definition")
				}

				bi.isWithinBlock = ""

				if bi.debug {
					fmt.Print(yellow(" --> Ending section definition"))
				}

			} else {
				if bi.isWithinBlock != "" {
					// Processing symbols
					if strings.HasSuffix(operand, ":") {
						newSymbol := strings.TrimSuffix(operand, ":")
						if newSymbol != "" {
							if bi.isSymbolled != "" {
								newSymbol += ":" + bi.isSymbolled
							}
							bi.isSymbolled = newSymbol
						}
						if argN > 1 {
							if bi.lineMeta != "" {
								bi.lineMeta += "," + strings.Join(argS[1:], "")
							} else {
								bi.lineMeta = strings.Join(argS[1:], "")
							}
						}
						if bi.debug {
							fmt.Print(yellow(" --> Symbol set"))
						}
					} else {

						newElem := new(bmline.BasmElement)
						newElem.SetValue(operand)

						newLine := new(bmline.BasmLine)
						newLine.Operation = newElem

						if argN > 1 {
							arguments := strings.Split(strings.Join(argS[1:], " "), ",")
							newArgs := make([]*bmline.BasmElement, len(arguments))
							for i, arg := range arguments {
								newArg := new(bmline.BasmElement)
								newArg.SetValue(strings.TrimSpace(arg))
								newArgs[i] = newArg
							}
							newLine.Elements = newArgs
						}

						if bi.isSymbolled != "" || bi.lineMeta != "" {
							if bi.isSymbolled != "" {
								newLine.BasmMeta = newLine.SetMeta("symbol", bi.isSymbolled)
							}
							if bi.lineMeta != "" {
								pairs := getPairs(bi.lineMeta)

								if len(pairs) > 0 {
									for _, pair := range pairs {
										keyVal := strings.Split(pair, ":")
										if len(keyVal) < 2 {
											return errors.New("Wrong format for entry (3): " + pair)
										} else {
											newLine.BasmMeta = newLine.SetMeta(keyVal[0], strings.Join(keyVal[1:], ":"))
										}
									}
								}
							}
							bi.isSymbolled = ""
							bi.lineMeta = ""
						}

						if bi.isWithinBlock != "" {
							if bi.debug {
								fmt.Print(yellow(" --> Section line:"))
								fmt.Print(newLine)
							}
							bi.blocks[bi.isWithinBlock].blockBody.Lines = append(bi.blocks[bi.isWithinBlock].blockBody.Lines, newLine)
						}
					}
				} else {
					return errors.New("Unknown directive on line " + line)
				}
			}

			if bi.debug {
				fmt.Println()
			}
		} else {
			if bi.debug {
				fmt.Println()
			}
		}
	}

	return nil
}
