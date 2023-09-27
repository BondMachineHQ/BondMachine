package basm

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

type choice []string

func mergeChoices(poss [][]string, level int) [][]string {
	if len(poss) == 0 {
		return [][]string{}
	}

	if level == len(poss)-1 {
		result := make([][]string, 0)
		for _, c := range poss[level] {
			result = append(result, []string{c})
		}
		return result
	} else {
		result := make([][]string, 0)
		for _, c := range poss[level] {
			for _, r := range mergeChoices(poss, level+1) {
				result = append(result, append(r, c))
			}
		}
		return result
	}
}

func matcherResolver(bi *BasmInstance) error {

	if bi.debug {
		fmt.Println(green("\tProcessing sections:"))
	}

	sectionsIncoming := make(map[string]*BasmSection)

	// Loop over the sections to find eventual alternatives
	for sectName, section := range bi.sections {

		sectionWithChoices := ""

		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			if section.sectionBody.GetMeta("template") == "true" {
				if bi.debug {
					fmt.Println(green("\t\t\tTemplated section, skipping"))
				}
				continue
			}

			// Setup the map of alternatives sections based on the choices
			sectAlts := make(map[string]choice)

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

				// Sort the matching slice to have always the same order
				slices.Sort(matching)
				// fmt.Println(matching)

				if !matched {
					return errors.New("no operator match")
				} else if len(matching) > 1 {
					// 	found := false
					choiceName := strings.Join(matching, "_")

					if _, ok := sectAlts[choiceName]; !ok {
						sectAlts[choiceName] = make([]string, len(matching))
						copy(sectAlts[choiceName], matching)
					}
				}
			}

			secAltsKeys := make([]string, 0)
			secChoices := make([][]string, 0)
			for k := range sectAlts {
				secAltsKeys = append(secAltsKeys, k)
				secChoices = append(secChoices, sectAlts[k])
			}

			allChoices := mergeChoices(secChoices, 0)

			if len(allChoices) == 0 {
				allChoices = make([][]string, 1)
			}

			for _, c := range allChoices {

				idx := 0
			idxLoop:
				// Identify the next available index
				for {
					if _, ok := bi.sections[sectName+"_"+strconv.Itoa(idx)]; !ok {
						if _, ok := sectionsIncoming[sectName+"_"+strconv.Itoa(idx)]; !ok {
							break idxLoop
						}
					}
					idx++
				}
				// Create the new section with the new name and the choice
				sectNameNew := sectName + "_" + strconv.Itoa(idx)

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

					// Sort the matching slice to have always the same order
					slices.Sort(matching)

					if len(matching) > 1 {
						choiceName := strings.Join(matching, "_")
						choiceId := 0
						for cid, cn := range secAltsKeys {
							if cn == choiceName {
								choiceId = cid
								break
							}
						}

						idx, _ := strconv.Atoi(c[choiceId])
						matchingOp = bi.matchersOps[idx]
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
				if sectionWithChoices == "" {
					sectionWithChoices = sectNameNew
				} else {
					sectionWithChoices = sectionWithChoices + ":" + sectNameNew
				}
			}

			section.sectionBody = new(bmline.BasmBody)
			section.sectionBody.BasmMeta = section.sectionBody.BasmMeta.SetMeta("alternatives", sectionWithChoices)
		}
	}

	for sn, s := range sectionsIncoming {
		bi.sections[sn] = s
	}

	return nil
}
