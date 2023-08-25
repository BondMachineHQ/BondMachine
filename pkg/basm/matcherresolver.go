package basm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

type choice struct {
	choiceName string
	choiceSel  string
}
type choices []choice

func matcherResolver(bi *BasmInstance) error {
	//TODO finish this

	if bi.debug {
		fmt.Println(green("\tProcessing sections:"))
	}

	// sectionsIncoming := make(map[string]*BasmSection)
	// sectionsOutgoing := make(map[string]struct{})

	// Loop over the sections to find eventual alternatives
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			// Setup the map of alternatives sections based on the choices
			sectAlts := make(map[string]choices)

			body := section.sectionBody

			for _, line := range body.Lines {

				if bi.debug {
					fmt.Println(green("\t\t\tLine: ") + line.String())
				}

				matched := false
				matching := make([]string, 0)

				for j, matcher := range bi.matchers {
					if bmline.MatchMatcher(matcher, line) {
						if bi.debug {
							fmt.Println(yellow("\t\t\t\tMatching " + matcher.String()))
						}
						matched = true
						matching = append(matching, strconv.Itoa(j))
					}
				}
				// fmt.Println(matching)
				if !matched {
					return errors.New("no operator match")
				} else if len(matching) > 1 {
					found := false
					choiceName := strings.Join(matching, "_")
				sLoop:
					for _, cs := range sectAlts {
						for _, c := range cs {
							if c.choiceName == choiceName {
								found = true
								break sLoop
							}
						}
					}
					// There is a choice to be made, creating the alternatives names and the choices for the sections
					if !found {
						sectRem := make(map[string]struct{})

						if len(sectAlts) == 0 {
							for _, c := range matching {
								idx := 0
							idxLoop1:
								// Identify the next available index
								for {
									if _, ok := bi.sections[sectName+"_"+strconv.Itoa(idx)]; !ok {
										if _, ok := sectAlts[sectName+"_"+strconv.Itoa(idx)]; !ok {
											break idxLoop1
										}
									}
									idx++
								}
								// Create the new section with the new name and the choice
								sectAlts[sectName+"_"+strconv.Itoa(idx)] = choices{choice{choiceName: choiceName, choiceSel: c}}
							}
						} else {
							for s, cs := range sectAlts {
								for _, c := range matching {
									idx := 0
								idxLoop2:
									// Identify the next available index
									for {
										if _, ok := bi.sections[sectName+"_"+strconv.Itoa(idx)]; !ok {
											if _, ok := sectAlts[sectName+"_"+strconv.Itoa(idx)]; !ok {
												break idxLoop2
											}
										}
										idx++
									}
									// Create the new section with the new name and the choice
									sectAlts[sectName+"_"+strconv.Itoa(idx)] = append(cs, choice{choiceName: choiceName, choiceSel: c})
								}
								sectRem[s] = struct{}{}
							}

							// Remove the old sections
							for s := range sectRem {
								delete(sectAlts, s)
							}
						}
					}
				}
			}
			// TODO Finire da qui
		}
	}
	// Loop over the sections
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			if section.sectionType == sectRomText {
				bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts", T: bmreqs.ObjectSet, Name: "sections", Value: sectName, Op: bmreqs.OpAdd})
			} else {
				bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:ramtexts", T: bmreqs.ObjectSet, Name: "sections", Value: sectName, Op: bmreqs.OpAdd})
			}

			body := section.sectionBody

			for i, line := range body.Lines {

				if bi.debug {
					fmt.Println(green("\t\t\tLine: ") + line.String())
				}

				matched := false
				var matching procbuilder.Opcode

				for j, matcher := range bi.matchers {
					if bmline.MatchMatcher(matcher, line) {
						if bi.debug {
							fmt.Println(yellow("\t\t\t\tMaching " + matcher.String()))
						}
						if matched {
							return errors.New("ambiguous, more than one operator match")
						}
						matched = true
						matching = bi.matchersOps[j]
					}
				}

				if !matched {
					return errors.New("no operator match")
				}

				opName := matching.Op_get_name()

				if section.sectionType == sectRomText {
					bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts/sections:" + sectName, T: bmreqs.ObjectSet, Name: "opcodes", Value: opName, Op: bmreqs.OpAdd})
				} else {
					bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:ramtexts/sections:" + sectName, T: bmreqs.ObjectSet, Name: "opcodes", Value: opName, Op: bmreqs.OpAdd})
				}

				// Normalize instruction
				if section.sectionType == sectRomText {
					if normalized, err := matching.HLAssemblerNormalize(nil, bi.rg, "/code:romtexts/sections:"+sectName, line); err != nil {
						return err
					} else {
						if bi.debug {
							fmt.Println(green("\t\t\t\tNormalized line: ") + normalized.String())
						}
						body.Lines[i] = normalized
					}
				} else {
					if normalized, err := matching.HLAssemblerNormalize(nil, bi.rg, "/code:ramtexts/sections:"+sectName, line); err != nil {
						return err
					} else {
						if bi.debug {
							fmt.Println(green("\t\t\t\tNormalized line: ") + normalized.String())
						}
						body.Lines[i] = normalized
					}
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
