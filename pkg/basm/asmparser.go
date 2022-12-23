package basm

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

func basmParser(bi *BasmInstance, s string, lineNo uint32) error {
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
			} else if operand == "%macro" {
				// Macro starting
				if bi.isWithinMacro != "" || bi.isWithinSection != "" || bi.isWithinFragment != "" {
					return errors.New(line + ", macro, section or fragment definition is not closed")
				}
				if argN != 3 {
					return errors.New(line + ", macro definition has the wrong number of parameters")
				}

				macroName := argS[1]
				macroArgs, err := strconv.Atoi(argS[2])
				if err != nil {
					return errors.New(line + ", macro definition has inconsistences")
				}

				if _, ok := bi.macros[macroName]; ok {
					return errors.New(line + ", macro already defined")
				}

				newMacro := new(BasmMacro)
				newMacro.macroName = macroName
				newMacro.macroArgs = macroArgs
				newMacro.macroBody = new(bmline.BasmBody)
				newMacro.macroBody.Lines = make([]*bmline.BasmLine, 0)

				bi.macros[macroName] = newMacro

				bi.isWithinMacro = macroName

				if bi.debug {
					fmt.Print(yellow(" --> Starting macro definition"))
				}
			} else if operand == "%endmacro" {
				// Macro ending
				if bi.isWithinMacro == "" {
					return errors.New(line + ", endmacro outside macro definition")
				}

				bi.isWithinMacro = ""

				if bi.debug {
					fmt.Print(yellow(" --> Ending macro definition"))
				}

			} else if operand == "%section" {
				if bi.isWithinMacro != "" || bi.isWithinSection != "" || bi.isWithinFragment != "" {
					return errors.New(line + ", macro, section or fragment definition is not closed")
				}
				// New section starting
				if argN < 3 {
					return errors.New(line + ", section definition has the wrong number of parameters")
				}

				sectionName := argS[1]
				sectionType := uint8(0)

				switch argS[2] {
				case ".romdata":
					sectionType = sectRomData
				case ".romtext":
					sectionType = setcRomText
				default:
					return errors.New(line + ", section definition has an unknown type " + argS[2])
				}

				newSection := new(BasmSection)
				newSection.sectionName = sectionName
				newSection.sectionType = sectionType
				newSection.sectionBody = new(bmline.BasmBody)

				if argN > 3 {
					pairs := getPairs(strings.Join(argS[3:], ","))

					if len(pairs) > 0 {
						for _, pair := range pairs {
							keyVal := strings.Split(pair, ":")
							if len(keyVal) < 2 {
								return errors.New("Wrong format for entry: " + pair)
							} else {
								newSection.sectionBody.BasmMeta = newSection.sectionBody.SetMeta(keyVal[0], strings.Join(keyVal[1:], ":"))
							}
						}
					}
				}

				bi.sections[sectionName] = newSection

				bi.isWithinSection = sectionName

				if bi.debug {
					fmt.Print(yellow(" --> Starting section " + sectionName + " definition"))
				}

			} else if operand == "%endsection" {
				// Section ending
				if bi.isWithinSection == "" {
					return errors.New(line + ", endsection outside section definition")
				}

				bi.isWithinSection = ""

				if bi.debug {
					fmt.Print(yellow(" --> Ending section definition"))
				}

			} else if operand == "%fragment" {
				if bi.isWithinMacro != "" || bi.isWithinSection != "" || bi.isWithinFragment != "" {
					return errors.New(line + ", macro, section, fragment or chunk definition is not closed")
				}
				// New fragment starting
				if argN < 2 {
					return errors.New(line + ", fragment definition has the wrong number of parameters")
				}

				fragmentName := argS[1]

				newFragment := new(BasmFragment)
				newFragment.fragmentName = fragmentName
				newFragment.fragmentBody = new(bmline.BasmBody)

				if argN > 2 {
					pairs := getPairs(strings.Join(argS[2:], ","))
					if len(pairs) > 0 {
						for _, pair := range pairs {
							keyVal := strings.Split(pair, ":")
							if len(keyVal) < 2 {
								return errors.New("Wrong format for entry: " + pair)
							} else {
								newFragment.fragmentBody.BasmMeta = newFragment.fragmentBody.SetMeta(keyVal[0], strings.Join(keyVal[1:], ":"))
							}
						}
					}
				}
				bi.fragments[fragmentName] = newFragment

				bi.isWithinFragment = fragmentName

				if bi.debug {
					fmt.Print(yellow(" --> Starting fragment definition"))
				}

			} else if operand == "%endfragment" {
				// Fragment ending
				if bi.isWithinFragment == "" {
					return errors.New(line + ", endfragment outside fragment definition")
				}

				bi.isWithinFragment = ""

				if bi.debug {
					fmt.Print(yellow(" --> Ending fragment definition"))
				}
			} else if operand == "%chunk" {
				if bi.isWithinMacro != "" || bi.isWithinSection != "" || bi.isWithinFragment != "" || bi.isWithinChunk != "" {
					return errors.New(line + ", macro, section, fragment or chunk definition is not closed")
				}
				// New ckunk starting
				if argN != 2 {
					return errors.New(line + ", chunk definition has the wrong number of parameters")
				}

				chunkName := argS[1]

				newChunk := new(BasmChunk)
				newChunk.chunkName = chunkName
				newChunk.chunkBody = new(bmline.BasmBody)

				bi.chunks[chunkName] = newChunk

				bi.isWithinChunk = chunkName

				if bi.debug {
					fmt.Print(yellow(" --> Starting chunk definition"))
				}

			} else if operand == "%endchunk" {
				// Chunk ending
				if bi.isWithinChunk == "" {
					return errors.New(line + ", endchunk outside chunk definition")
				}

				bi.isWithinFragment = ""

				if bi.debug {
					fmt.Print(yellow(" --> Ending chunk definition"))
				}

			} else {
				if bi.isWithinMacro != "" || bi.isWithinSection != "" || bi.isWithinFragment != "" || bi.isWithinChunk != "" {
					// Processing labels
					if strings.HasSuffix(operand, ":") {
						if bi.isLabelled != "" {
							return errors.New(line + ", Multiple label")
						}
						newlabel := strings.TrimSuffix(operand, ":")
						bi.isLabelled = newlabel
						bi.lineMeta = strings.Join(argS[1:], "")
						if bi.debug {
							fmt.Print(yellow(" --> Label set"))
						}
					} else {

						newElem := new(bmline.BasmElement)
						newElem.SetValue(operand)

						newLine := new(bmline.BasmLine)
						newLine.Operation = newElem

						if bi.isWithinSection != "" && bi.sections[bi.isWithinSection].sectionType == sectRomData {
							if argN > 1 {
								argS := strings.Split(strings.TrimSpace(strings.Join(argS[1:], " ")), " ")
								if len(argS) > 1 {
									dataType := argS[0]
									dataValue := strings.Join(argS[1:], " ")
									newArgs := make([]*bmline.BasmElement, 2)
									newArg1 := new(bmline.BasmElement)
									newArg1.SetValue(dataType)
									newArgs[0] = newArg1
									newArg2 := new(bmline.BasmElement)
									newArg2.SetValue(dataValue)
									newArgs[1] = newArg2
									newLine.Elements = newArgs
								}
							}
						} else {
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
						}
						if bi.isLabelled != "" || bi.lineMeta != "" {
							if bi.isLabelled != "" {
								newLine.BasmMeta = newLine.SetMeta("label", bi.isLabelled)
							}
							if bi.lineMeta != "" {
								if len(argS) > 1 {
									pairs := getPairs(bi.lineMeta)

									if len(pairs) > 0 {
										for _, pair := range pairs {
											keyVal := strings.Split(pair, ":")
											if len(keyVal) < 2 {
												return errors.New("Wrong format for entry: " + pair)
											} else {
												newLine.BasmMeta = newLine.SetMeta(keyVal[0], keyVal[1])
											}
										}
									}
								}
							}
							bi.isLabelled = ""
							bi.lineMeta = ""
						}

						if bi.isWithinMacro != "" {
							if bi.debug {
								fmt.Print(yellow(" --> Macro line:"))
								fmt.Print(newLine)
							}
							bi.macros[bi.isWithinMacro].macroBody.Lines = append(bi.macros[bi.isWithinMacro].macroBody.Lines, newLine)
						}
						if bi.isWithinSection != "" {
							if bi.debug {
								fmt.Print(yellow(" --> Section line:"))
								fmt.Print(newLine)
							}
							bi.sections[bi.isWithinSection].sectionBody.Lines = append(bi.sections[bi.isWithinSection].sectionBody.Lines, newLine)
						}

						if bi.isWithinFragment != "" {
							if bi.debug {
								fmt.Print(yellow(" --> Fragment line:"))
								fmt.Print(newLine)
							}
							bi.fragments[bi.isWithinFragment].fragmentBody.Lines = append(bi.fragments[bi.isWithinFragment].fragmentBody.Lines, newLine)
						}

						if bi.isWithinChunk != "" {
							if bi.debug {
								fmt.Print(yellow(" --> Chunk line:"))
								fmt.Print(newLine)
							}
							bi.chunks[bi.isWithinChunk].chunkBody.Lines = append(bi.chunks[bi.isWithinChunk].chunkBody.Lines, newLine)
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
