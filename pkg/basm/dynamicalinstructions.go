package basm

import (
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

func dynamicalInstructions(bi *BasmInstance) error {

	// Loop over the sections
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			body := section.sectionBody

			for _, line := range body.Lines {
				eventualInstruction := line.Operation.GetValue()

				if created, err := procbuilder.EventuallyCreateInstruction(eventualInstruction); err != nil {
					return err
				} else {
					if created {
						op := procbuilder.Allopcodes[len(procbuilder.Allopcodes)-1]
						for _, line := range op.HLAssemblerMatch(nil) {
							if mt, err := bmline.Text2BasmLine(line); err == nil {
								bi.matchers = append(bi.matchers, mt)
								bi.matchersOps = append(bi.matchersOps, op)
							} else {
								bi.Warning(err)
							}
						}
					}
				}

				for j, matcher := range bi.dynMatchers {
					if bmline.MatchMatcher(matcher, line) {
						if bi.debug {
							fmt.Println(yellow("\t\t\t\tMatching " + matcher.String()))
						}
						dyn := bi.dynMatcherOps[j]
						for _, op := range dyn.HLAssemblerGeneratorList(nil, line) {
							eventualInstruction := op

							if created, err := procbuilder.EventuallyCreateInstruction(eventualInstruction); err != nil {
								return err
							} else {
								if created {
									op := procbuilder.Allopcodes[len(procbuilder.Allopcodes)-1]
									for _, line := range op.HLAssemblerMatch(nil) {
										if mt, err := bmline.Text2BasmLine(line); err == nil {
											bi.matchers = append(bi.matchers, mt)
											bi.matchersOps = append(bi.matchersOps, op)
										} else {
											bi.Warning(err)
										}
									}
								}
							}

						}
					}
				}

			}
		} else {
			if bi.debug {
				fmt.Println(yellow("\t\tSection type not handled: ") + sectName)
			}
		}
	}

	// Loop over the fragments
	for fragName, fragment := range bi.fragments {
		if bi.debug {
			fmt.Println(green("\t\tFragment: ") + fragName)
		}

		body := fragment.fragmentBody

		for _, line := range body.Lines {
			eventualInstruction := line.Operation.GetValue()

			if created, err := procbuilder.EventuallyCreateInstruction(eventualInstruction); err != nil {
				return err
			} else {
				if created {
					op := procbuilder.Allopcodes[len(procbuilder.Allopcodes)-1]
					for _, line := range op.HLAssemblerMatch(nil) {
						if mt, err := bmline.Text2BasmLine(line); err == nil {
							bi.matchers = append(bi.matchers, mt)
							bi.matchersOps = append(bi.matchersOps, op)
						} else {
							bi.Warning(err)
						}
					}
				}
			}

			for j, matcher := range bi.dynMatchers {
				if bmline.MatchMatcher(matcher, line) {
					if bi.debug {
						fmt.Println(yellow("\t\t\t\tMatching " + matcher.String()))
					}
					dyn := bi.dynMatcherOps[j]
					for _, op := range dyn.HLAssemblerGeneratorList(nil, line) {
						eventualInstruction := op

						if created, err := procbuilder.EventuallyCreateInstruction(eventualInstruction); err != nil {
							return err
						} else {
							if created {
								op := procbuilder.Allopcodes[len(procbuilder.Allopcodes)-1]
								for _, line := range op.HLAssemblerMatch(nil) {
									if mt, err := bmline.Text2BasmLine(line); err == nil {
										bi.matchers = append(bi.matchers, mt)
										bi.matchersOps = append(bi.matchersOps, op)
									} else {
										bi.Warning(err)
									}
								}
							}
						}

					}
				}
			}

		}
	}

	return nil
}
