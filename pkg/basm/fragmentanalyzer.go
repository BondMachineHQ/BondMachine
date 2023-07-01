package basm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

func (bi *BasmInstance) fragmentResUsage(body *bmline.BasmBody, circular bool) error {
	//TODO finish this
	fmt.Println("fragmentResUsage", body, circular)

	if circular {
	} else {

	}

	return nil
}

func fragmentAnalyzer(bi *BasmInstance) error {
	//TODO finish this

	if bi.debug {
		fmt.Println(green("\tProcessing fragments:"))
	}

	// Filter the matchers to select only the label based one
	filteredMatchers := make([]*bmline.BasmLine, 0)

	// The nop line
	nopLine := new(bmline.BasmLine)
	nop := new(bmline.BasmElement)
	nop.SetValue("nop")
	nopLine.Operation = nop
	nopLine.Elements = make([]*bmline.BasmElement, 0)

	if bi.debug {
		fmt.Println(green("\tFiltering matchers:"))
	}

	for _, matcher := range bi.matchers {
		if bmline.FilterMatcher(matcher, "label") {
			filteredMatchers = append(filteredMatchers, matcher)
			if bi.debug {
				fmt.Println(red("\t\tActive matcher:") + matcher.String())
			}
		} else {
			if bi.debug {
				fmt.Println(yellow("\t\tInactive matcher:") + matcher.String())
			}
		}
	}

	// Loop over the sections
	for fragName, fragment := range bi.fragments {
		if bi.debug {
			fmt.Println(green("\t\tFragment: ")+fragName, fragment)
		}

		fBody := fragment.fragmentBody

		resInS := fBody.GetMeta("resin")
		resIn := make([]string, 0)
		for _, res := range strings.Split(resInS, ":") {
			if res != "" {
				resIn = append(resIn, res)
			}
		}

		resOuts := fBody.GetMeta("resout")
		resOut := make([]string, 0)
		for _, res := range strings.Split(resOuts, ":") {
			if res != "" {
				resOut = append(resOut, res)
			}
		}

		resUsed := make(map[string]struct{})
		for _, line := range fBody.Lines {
			for _, elem := range line.Elements {
				ty := elem.GetMeta("type")
				switch ty {
				case "reg":
					resUsed[elem.GetValue()] = struct{}{}
				}
			}
		}

		// fmt.Println("resIn", resIn)
		// fmt.Println("resOut", resOut)
		// fmt.Println("resUsed", resUsed)

		resUseds := ""
		for res := range resUsed {
			resUseds += res + ":"
		}
		resUseds = strings.TrimSuffix(resUseds, ":")

		if resInS != "" {
			fBody.BasmMeta = fBody.SetMeta("resin", resInS)
		}
		if resOuts != "" {
			fBody.BasmMeta = fBody.SetMeta("resout", resOuts)
		}
		fBody.BasmMeta = fBody.SetMeta("resused", resUseds)

		// TODO rearrange resources in the order they are used

		for _, line := range fBody.Lines {

			if bi.debug {
				fmt.Println(green("\t\t\tLine: ") + line.String())
			}

			matched := false
			var matching procbuilder.Opcode

			for j, matcher := range bi.matchers {
				if bmline.MatchMatcher(matcher, line) {
					if bi.debug {
						fmt.Println(yellow("\t\t\t\tMatching " + matcher.String()))
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

			if meta, err := matching.HLAssemblerInstructionMetadata(nil, line); err != nil {
				return err
			} else {
				if meta != nil {
					for k, v := range meta.LoopMeta() {
						line.BasmMeta = line.SetMeta(k, v)
					}
				}
			}
		}

		branchingBlocks := make(map[int]*bmline.BasmBody, 1)
		circularBlocks := make(map[int]bool, 1)
		// Create a copy of the fragment body. The first copy will be the whole fragment body with index 0
		branchingBlocks[0] = fBody.Copy()
		circularBlocks[0] = false

		// Identify every branching instruction, the identifier will be the line number of the starting point. To identify the end point, we will
		// use the label.
		// This will only work with labels, not with numbers. For this reason, it will execute after labeltagger and before labelresolver.

		// Map all the labels
		labels := make(map[string]int)

		for i, line := range fBody.Lines {
			if label := line.GetMeta("label"); label != "" {
				if _, exists := labels[label]; exists {
					return errors.New("label is specified multiple time: " + label)
				} else {
					labels[label] = i
				}
			}
		}

		// Get where the branching instructions are

		for i, line := range fBody.Lines {

			for _, matcher := range filteredMatchers {
				if bmline.MatchMatcher(matcher, line) {
					// TODO Handling the operand
					for _, arg := range line.Elements {
						if j, ok := labels[arg.GetValue()]; ok {
							fmt.Println("the branch is at line", i, "and the label is at line", j)

							if i > j {
								branchingBlocks[i] = fBody.Copy()
								circularBlocks[i] = true
								// upper branch
								for k := 0; k < j; k++ {
									branchingBlocks[i].Lines[k] = nopLine.Copy()
								}
								for k := i + 1; k < len(fBody.Lines); k++ {
									branchingBlocks[i].Lines[k] = nopLine.Copy()
								}
							} else if i < j {
								branchingBlocks[i] = fBody.Copy()
								circularBlocks[i] = false
								// lower branch
								for k := i + 1; k < j; k++ {
									branchingBlocks[i].Lines[k] = nopLine.Copy()
								}
							} else {
								// Same line. Ignore it
							}

						}
					}

				}
			}
		}

		// Compute the usage of resources for each copy
		for i, body := range branchingBlocks {
			if err := bi.fragmentResUsage(body, circularBlocks[i]); err != nil {
				return err
			}
		}

		// Union of the copies
	}
	// panic("fragmentAnalyzer not finished")
	return nil
}
