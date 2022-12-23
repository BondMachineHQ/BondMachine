package basm

import (
	"errors"
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

func matcherResolver(bi *BasmInstance) error {
	//TODO finish this

	if bi.debug {
		fmt.Println(green("\tProcessing sections:"))
	}

	// Loop over the sections
	for sectName, section := range bi.sections {
		if section.sectionType == setcRomText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts", T: bmreqs.ObjectSet, Name: "sections", Value: sectName, Op: bmreqs.OpAdd})

			body := section.sectionBody

			for i, line := range body.Lines {

				if bi.debug {
					fmt.Println(green("\t\t\tLine: ") + line.String())
				}

				matched := false
				var maching procbuilder.Opcode

				for j, matcher := range bi.matchers {
					if bmline.MatchMatcher(matcher, line) {
						if bi.debug {
							fmt.Println(yellow("\t\t\t\tMaching " + matcher.String()))
						}
						if matched {
							return errors.New("Ambiguous, more than one operator match")
						}
						matched = true
						maching = bi.matchersOps[j]
					}
				}

				if !matched {
					return errors.New("no operator match")
				}

				opname := maching.Op_get_name()

				bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts/sections:" + sectName, T: bmreqs.ObjectSet, Name: "opcodes", Value: opname, Op: bmreqs.OpAdd})

				// Normalize instruction
				if normalized, err := maching.HLAssemblerNormalize(nil, bi.rg, "/code:romtexts/sections:"+sectName, line); err != nil {
					return err
				} else {
					if bi.debug {
						fmt.Println(green("\t\t\t\tNormalized line: ") + normalized.String())
					}
					body.Lines[i] = normalized
				}
			}
		} else {
			if bi.debug {
				fmt.Println(yellow("\t\tSection type not handled: ") + sectName)
			}
		}
	}
	return nil
}
