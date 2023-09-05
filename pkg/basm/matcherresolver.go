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

	sectionsIncoming := make(map[string]*BasmSection)
	sectionsOutgoing := make(map[string]struct{})

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

			sectionWithChoices := ""
			if len(sectAlts) == 0 {
				sectAlts[sectName] = nil
			} else {
				alts := ""
				for s := range sectAlts {
					alts += s + ":"
				}
				sectionWithChoices = alts[:len(alts)-1]
			}

			// Ranging over the alternatives to create the reals sections
			for sectNameNew, cs := range sectAlts {

				sectionNew := new(BasmSection)
				sectionNew.sectionName = sectNameNew
				sectionNew.sectionType = section.sectionType
				sectionNew.sectionBody = section.sectionBody.Copy()

				if bi.debug {
					fmt.Println(green("\t\tNew section: ") + sectNameNew)
				}

				if sectionNew.sectionType == sectRomText {
					bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts", T: bmreqs.ObjectSet, Name: "sections", Value: sectNameNew, Op: bmreqs.OpAdd})
				} else {
					bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:ramtexts", T: bmreqs.ObjectSet, Name: "sections", Value: sectNameNew, Op: bmreqs.OpAdd})
				}

				body := sectionNew.sectionBody

				for i, line := range body.Lines {

					if bi.debug {
						fmt.Println(green("\t\t\tLine: ") + line.String())
					}

					matching := make([]string, 0)
					var matchingOp procbuilder.Opcode

					for j, matcher := range bi.matchers {
						if bmline.MatchMatcher(matcher, line) {
							if bi.debug {
								fmt.Println(yellow("\t\t\t\tMatching " + matcher.String()))
							}
							matching = append(matching, strconv.Itoa(j))
						}
					}

					if len(matching) > 1 {
						choiceName := strings.Join(matching, "_")
						for _, c := range cs {
							if c.choiceName == choiceName {
								idx, _ := strconv.Atoi(c.choiceSel)
								matchingOp = bi.matchersOps[idx]
								break
							}
						}
					} else {
						idx, _ := strconv.Atoi(matching[0])
						matchingOp = bi.matchersOps[idx]
					}

					opName := matchingOp.Op_get_name()

					if sectionNew.sectionType == sectRomText {
						bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts/sections:" + sectNameNew, T: bmreqs.ObjectSet, Name: "opcodes", Value: opName, Op: bmreqs.OpAdd})
					} else {
						bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:ramtexts/sections:" + sectNameNew, T: bmreqs.ObjectSet, Name: "opcodes", Value: opName, Op: bmreqs.OpAdd})
					}

					// Normalize instruction
					if sectionNew.sectionType == sectRomText {
						if normalized, err := matchingOp.HLAssemblerNormalize(nil, bi.rg, "/code:romtexts/sections:"+sectNameNew, line); err != nil {
							return err
						} else {
							if bi.debug {
								fmt.Println(green("\t\t\t\tNormalized line: ") + normalized.String())
							}
							for k, v := range body.Lines[i].LoopMeta() {
								normalized.BasmMeta = normalized.BasmMeta.SetMeta(k, v)
							}
							body.Lines[i] = normalized
						}
					} else {
						if normalized, err := matchingOp.HLAssemblerNormalize(nil, bi.rg, "/code:ramtexts/sections:"+sectNameNew, line); err != nil {
							return err
						} else {
							if bi.debug {
								fmt.Println(green("\t\t\t\tNormalized line: ") + normalized.String())
							}
							for k, v := range body.Lines[i].LoopMeta() {
								normalized.BasmMeta = normalized.BasmMeta.SetMeta(k, v)
							}
							body.Lines[i] = normalized
						}
					}
				}

				sectionsIncoming[sectNameNew] = sectionNew
			}

			if sectionWithChoices != "" {
				section.sectionBody = new(bmline.BasmBody)
				section.sectionBody.BasmMeta = section.sectionBody.BasmMeta.SetMeta("alternatives", sectionWithChoices)
			} else {
				sectionsOutgoing[sectName] = struct{}{}
			}
		}
	}

	for sn := range sectionsOutgoing {
		delete(bi.sections, sn)
	}

	for sn, s := range sectionsIncoming {
		bi.sections[sn] = s
	}

	return nil
}
