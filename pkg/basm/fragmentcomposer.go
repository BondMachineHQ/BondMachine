package basm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

func fragmentComposer(bi *BasmInstance) error {

	if bi.debug {
		fmt.Println(green("\tProcessing fragments based processors:"))
	}

	// Loop over the cpdef searching for fragments
	for _, cp := range bi.cps {
		if cp.GetMeta("fragcollapse") != "" {
			// Get the fragments name
			fragments := cp.GetMeta("fragcollapse")
			fragments = strings.Trim(fragments, ":")
			if bi.debug {
				fmt.Println("\t\tFragments collapsing: ", fragments)
			}

			fragList := strings.Split(fragments, ":")

			cpNewSectionName := "coll_" + strings.Join(fragList, "_")

			fragResin := make([][]string, len(fragList))
			fragResout := make([][]string, len(fragList))
			fragResused := make([][]string, len(fragList))

			newInputs := make([][]string, len(fragList))
			newOutputs := make([][]string, len(fragList))
			newRegAsInputs := make([][]string, len(fragList))
			newRegAsOutputs := make([][]string, len(fragList))

			currNewInput := "i0"
			currNewOutput := "o0"
			currNewReg := "t0" // t is temporary reg

			// Loop over the fragments instances
			for i, fiName := range fragList {

				if _, fi, err := bi.GetFI(fiName); err != nil {
					return err
				} else {
					frag := fi.GetMeta("fragment")
					resInS := bi.fragments[frag].fragmentBody.GetMeta("resin")
					resIn := make([]string, 0)
					for _, res := range strings.Split(resInS, ":") {
						if res != "" {
							resIn = append(resIn, res)
						}
					}
					fragResin[i] = resIn

					resOutS := bi.fragments[frag].fragmentBody.GetMeta("resout")
					resOut := make([]string, 0)
					for _, res := range strings.Split(resOutS, ":") {
						if res != "" {
							resOut = append(resOut, res)
						}
					}
					fragResout[i] = resOut

					resUseds := bi.fragments[frag].fragmentBody.GetMeta("resused")
					resUsed := make([]string, 0)
					for _, res := range strings.Split(resUseds, ":") {
						if res != "" {
							resUsed = append(resUsed, res)
						}
					}

					fragResused[i] = resUsed

					newIns := make([]string, len(resIn))
					newOuts := make([]string, len(resOut))
					newRegIns := make([]string, len(resIn))
					newRegOuts := make([]string, len(resOut))

					for i, _ := range resOut {
						links, err := bi.GetLinks(FILINK, fi.GetValue(), strconv.Itoa(i), "output")
						if err != nil {
							return err
						}

						doneOut := false
						doneRegOut := false
						for _, link := range links {
							if endP, _, err := bi.GetEndpoints(FILINK, link); err != nil {
								return err
							} else {
								if stringInSlice(endP.name, fragList) {
									// The enpoint is internal
									newRegOuts[i] = currNewReg

									if !doneRegOut {
										doneRegOut = true
									}
								} else {
									// The endpoint is external
									newOuts[i] = currNewOutput

									newIO := new(bmline.BasmElement)
									newIO.SetValue(link)
									newIO.BasmMeta = newIO.BasmMeta.SetMeta("type", "io")
									bi.ios = append(bi.ios, newIO)

									newAttach := new(bmline.BasmElement)
									newAttach.SetValue(link)
									newAttach.BasmMeta = newAttach.BasmMeta.SetMeta("type", "output")

									newAttach.BasmMeta = newAttach.BasmMeta.SetMeta("cp", cp.GetValue())
									newAttach.BasmMeta = newAttach.BasmMeta.SetMeta("index", currNewOutput[1:])
									bi.ioAttach = append(bi.ioAttach, newAttach)

									if !doneOut {
										doneOut = true
									}
								}
							}
						}
						if doneOut {
							currNewOutput = nextRes(currNewOutput)
						}
						if doneRegOut {
							currNewReg = nextRes(currNewReg)
						}
					}

					for i, _ := range resIn {
						links, err := bi.GetLinks(FILINK, fi.GetValue(), strconv.Itoa(i), "input")
						if err != nil {
							return err
						}

						if len(links) > 1 {
							return fmt.Errorf("fragment %s has more than one input link", fiName)
						}

						if len(links) == 1 {
							if _, endP, err := bi.GetEndpoints(FILINK, links[0]); err != nil {
								return err
							} else {
								if stringInSlice(endP.name, fragList) {
									// The endpoint is internal
									newIns[i] = ""
									// The internal input endpoint will be computed later. For now set it to element in fragList-output port
									for k, n := range fragList {
										if n == endP.name {
											newRegIns[i] = strconv.Itoa(k) + "-" + endP.index
										}
									}
								} else {
									// The endpoint is external
									newIns[i] = currNewInput
									newRegIns[i] = ""

									newIO := new(bmline.BasmElement)
									newIO.SetValue(links[0])
									newIO.BasmMeta = newIO.BasmMeta.SetMeta("type", "io")
									bi.ios = append(bi.ios, newIO)

									newAttach := new(bmline.BasmElement)
									newAttach.SetValue(links[0])
									newAttach.BasmMeta = newAttach.BasmMeta.SetMeta("type", "input")
									newAttach.BasmMeta = newAttach.BasmMeta.SetMeta("cp", cp.GetValue())
									newAttach.BasmMeta = newAttach.BasmMeta.SetMeta("index", currNewInput[1:])
									bi.ioAttach = append(bi.ioAttach, newAttach)

									currNewInput = nextRes(currNewInput)
								}
							}

						}
					}

					// TODO Populate the newRegIns with the correct values

					newInputs[i] = newIns
					newOutputs[i] = newOuts
					newRegAsInputs[i] = newRegIns
					newRegAsOutputs[i] = newRegOuts

				}
			}

			// Resolve missing reg inputs
			for i, _ := range newRegAsInputs {
				for j, val := range newRegAsInputs[i] {
					dest := strings.Split(val, "-")
					if len(dest) == 2 {
						oi, _ := strconv.Atoi(dest[0])
						oj, _ := strconv.Atoi(dest[1])
						newRegAsInputs[i][j] = newRegAsOutputs[oi][oj]
					}
				}
			}

			newSection := new(BasmSection)
			newSection.sectionName = cpNewSectionName
			newSection.sectionType = sectRomText
			newSection.sectionBody = new(bmline.BasmBody)

			newSection.sectionBody.BasmMeta = newSection.sectionBody.BasmMeta.SetMeta("entry", "0")
			switch bi.global.BasmMeta.GetMeta("iomode") {
			case "sync":
				newSection.sectionBody.BasmMeta = newSection.sectionBody.BasmMeta.SetMeta("iomode", "sync")
			case "async":
				newSection.sectionBody.BasmMeta = newSection.sectionBody.BasmMeta.SetMeta("iomode", "async")
			default:
				if bi.global.BasmMeta.GetMeta("iomode") == "" {
					newSection.sectionBody.BasmMeta = newSection.sectionBody.BasmMeta.SetMeta("iomode", "async")
				} else {
					return fmt.Errorf("wrong iomode")
				}
			}
			newSection.sectionBody.Lines = make([]*bmline.BasmLine, 0)

			newLineE := new(bmline.BasmLine)
			newOperationE := new(bmline.BasmElement)
			newOperationE.SetValue("entry")
			newRegE := new(bmline.BasmElement)
			newRegE.SetValue("_start")
			newLineE.Operation = newOperationE
			newLineE.Elements = make([]*bmline.BasmElement, 1)
			newLineE.Elements[0] = newRegE
			newSection.sectionBody.Lines = append(newSection.sectionBody.Lines, newLineE)

			firstLine := true

			for i, fiName := range fragList {
				for j, inp := range newInputs[i] {
					if inp != "" {
						newLine := new(bmline.BasmLine)
						if firstLine {
							newLine.BasmMeta = newLine.BasmMeta.SetMeta("label", "_start")
							firstLine = false
						}
						newOperation := new(bmline.BasmElement)
						newOperation.SetValue("mov")
						newReg := new(bmline.BasmElement)
						newReg.SetValue(fragResin[i][j])
						newReg.BasmMeta = newReg.BasmMeta.SetMeta("type", "reg")
						newIn := new(bmline.BasmElement)
						newIn.SetValue(newInputs[i][j])
						newIn.BasmMeta = newIn.BasmMeta.SetMeta("type", "input")
						newLine.Operation = newOperation
						newLine.Elements = make([]*bmline.BasmElement, 2)
						newLine.Elements[0] = newReg
						newLine.Elements[1] = newIn
						newSection.sectionBody.Lines = append(newSection.sectionBody.Lines, newLine)
					}
				}

				for j, inp := range newRegAsInputs[i] {
					if inp != "" {
						newLine := new(bmline.BasmLine)
						if firstLine {
							newLine.BasmMeta = newLine.BasmMeta.SetMeta("label", "_start")
							firstLine = false
						}
						newOperation := new(bmline.BasmElement)
						newOperation.SetValue("mov")
						newReg := new(bmline.BasmElement)
						newReg.SetValue(fragResin[i][j])
						newReg.BasmMeta = newReg.BasmMeta.SetMeta("type", "reg")
						newIn := new(bmline.BasmElement)
						newIn.SetValue(newRegAsInputs[i][j])
						newIn.BasmMeta = newIn.BasmMeta.SetMeta("type", "reg")
						newLine.Operation = newOperation
						newLine.Elements = make([]*bmline.BasmElement, 2)
						newLine.Elements[0] = newReg
						newLine.Elements[1] = newIn
						newSection.sectionBody.Lines = append(newSection.sectionBody.Lines, newLine)
					}
				}

				_, fi, _ := bi.GetFI(fiName)
				frag := fi.GetMeta("fragment")
				fragment := bi.fragments[frag]

				// Create a copy of the fragment
				fCopy := new(BasmFragment)
				fCopy.fragmentName = fragment.fragmentName
				fCopy.fragmentBody = fragment.fragmentBody.Copy()

				// Change the labels including a prefix
				prefmeta := new(bmline.BasmElement)
				prefmeta.BasmMeta = prefmeta.SetMeta("label", "frag"+strconv.Itoa(i))
				prefvalue := new(bmline.BasmElement)
				prefvalue.SetValue("frag" + strconv.Itoa(i))
				prefvalue.BasmMeta = prefvalue.SetMeta("type", "label")
				fCopy.fragmentBody.PrefixMeta(prefmeta)
				fCopy.fragmentBody.PrefixValue(prefvalue)

				for _, line := range fCopy.fragmentBody.Lines {
					newLine := copyLine(line)
					newSection.sectionBody.Lines = append(newSection.sectionBody.Lines, newLine)
				}

				for j, out := range newOutputs[i] {
					if out != "" {
						newLine := new(bmline.BasmLine)
						if firstLine {
							newLine.BasmMeta = newLine.BasmMeta.SetMeta("label", "_start")
							firstLine = false
						}
						newOperation := new(bmline.BasmElement)
						newOperation.SetValue("mov")
						newReg := new(bmline.BasmElement)
						newReg.SetValue(fragResout[i][j])
						newReg.BasmMeta = newReg.BasmMeta.SetMeta("type", "reg")
						newOut := new(bmline.BasmElement)
						newOut.SetValue(newOutputs[i][j])
						newOut.BasmMeta = newOut.BasmMeta.SetMeta("type", "output")
						newLine.Operation = newOperation
						newLine.Elements = make([]*bmline.BasmElement, 2)
						newLine.Elements[0] = newOut
						newLine.Elements[1] = newReg
						newSection.sectionBody.Lines = append(newSection.sectionBody.Lines, newLine)
					}
				}

				for j, out := range newRegAsOutputs[i] {
					if out != "" {
						newLine := new(bmline.BasmLine)
						if firstLine {
							newLine.BasmMeta = newLine.BasmMeta.SetMeta("label", "_start")
							firstLine = false
						}
						newOperation := new(bmline.BasmElement)
						newOperation.SetValue("mov")
						newReg := new(bmline.BasmElement)
						newReg.SetValue(fragResout[i][j])
						newReg.BasmMeta = newReg.BasmMeta.SetMeta("type", "reg")
						newOut := new(bmline.BasmElement)
						newOut.SetValue(newRegAsOutputs[i][j])
						newOut.BasmMeta = newOut.BasmMeta.SetMeta("type", "reg")
						newLine.Operation = newOperation
						newLine.Elements = make([]*bmline.BasmElement, 2)
						newLine.Elements[0] = newOut
						newLine.Elements[1] = newReg
						newSection.sectionBody.Lines = append(newSection.sectionBody.Lines, newLine)
					}
				}
			}

			newLine := new(bmline.BasmLine)
			newOperation := new(bmline.BasmElement)
			newOperation.SetValue("j")
			newStart := new(bmline.BasmElement)
			newStart.SetValue("_start")
			newStart.BasmMeta = newStart.BasmMeta.SetMeta("type", "lineno")
			newLine.Operation = newOperation
			newLine.Elements = make([]*bmline.BasmElement, 1)
			newLine.Elements[0] = newStart
			newSection.sectionBody.Lines = append(newSection.sectionBody.Lines, newLine)

			// Resolve the temporary registers

			regTemplate := new(bmline.BasmElement)
			regTemplate.SetValue("r")
			regTemplate.BasmMeta = regTemplate.BasmMeta.SetMeta("type", "reg")

			for tr := 0; tr < 10000; tr++ {
				tempReg := new(bmline.BasmElement)
				tempReg.SetValue("t" + strconv.Itoa(tr))
				tempReg.BasmMeta = tempReg.BasmMeta.SetMeta("type", "reg")

				if newSection.sectionBody.CheckArg(tempReg) {
					newReg := newSection.sectionBody.NextResource(regTemplate)
					newSection.sectionBody.ReplaceArg(tempReg, newReg)
				} else {
					break
				}
			}

			bi.sections[newSection.sectionName] = newSection

			cp.BasmMeta = cp.BasmMeta.SetMeta("romcode", newSection.sectionName)
			cp.BasmMeta.RmMeta("fragcollapse")

			if bi.debug {
				fmt.Println("\t\tFragments resources in: ", fragResin)
				fmt.Println("\t\tFragments resources out: ", fragResout)
				fmt.Println("\t\tFragments resources used: ", fragResused)
				fmt.Println("\t\tFragments new inputs: ", newInputs)
				fmt.Println("\t\tFragments new outputs: ", newOutputs)
				fmt.Println("\t\tFragments new reg inputs: ", newRegAsInputs)
				fmt.Println("\t\tFragments new reg outputs: ", newRegAsOutputs)
				fmt.Println("---")
			}

		}
	}

	for _, fi := range bi.fiLinkAttach {
		if fi.GetMeta("fi") == "ext" {
			newAtt := new(bmline.BasmElement)
			newAtt.SetValue(fi.GetValue())
			newAtt.BasmMeta = newAtt.BasmMeta.SetMeta("type", fi.GetMeta("type"))
			newAtt.BasmMeta = newAtt.BasmMeta.SetMeta("index", fi.GetMeta("index"))
			newAtt.BasmMeta = newAtt.BasmMeta.SetMeta("cp", "bm")
			// fmt.Println("newAtt: ", newAtt)
			bi.ioAttach = append(bi.ioAttach, newAtt)
		}
	}

	// TODO resolve t regs
	return nil
}

func copyLine(line *bmline.BasmLine) *bmline.BasmLine {
	newLine := line.Copy()
	return newLine
}
