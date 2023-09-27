package basm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

const (
	FIRSTLINE = -2
	LASTLINE  = -3
	NULLLINE  = -1
)

func (bi *BasmInstance) fragmentResUsage(body *bmline.BasmBody, circular bool) error {
	//TODO finish this
	if bi.debug {
		// TODO Better formatting
		fmt.Println("fragmentResUsage", body, circular)
	}

	// Get all resources used by the fragment
	resUsed := body.GetMeta("resused")
	resIn := body.GetMeta("resin")
	resOut := body.GetMeta("resout")

	resInUse := strings.Split(resIn, ":")
	resOutUse := strings.Split(resOut, ":")
	_ = resInUse
	resStart := make(map[string]int)
	resEnd := make(map[string]int)

	for _, res := range strings.Split(resUsed, ":") {
		if stringInSlice(res, resInUse) {
			resStart[res] = FIRSTLINE
		} else {
			resStart[res] = NULLLINE
		}
		if stringInSlice(res, resOutUse) {
			resEnd[res] = LASTLINE
		} else {
			resEnd[res] = NULLLINE
		}
	}

	if circular {
	} else {
		for line := len(body.Lines) - 1; line >= 0; line-- {
			// get defined resources in this line and remove them from the resUsed map
			for _, cd := range strings.Split(body.Lines[line].GetMeta("inv"), ":") {
				if resEnd[cd] != NULLLINE {
					if resEnd[cd] == LASTLINE {
						for d := line; d < len(body.Lines); d++ {
							// Insert the metadata
							body.Lines[d].BasmMeta = body.Lines[d].AddMeta("inuse", cd)
						}
					} else {
						for d := line; d <= resEnd[cd]; d++ {
							// Insert the metadata
							body.Lines[d].BasmMeta = body.Lines[d].AddMeta("inuse", cd)
						}
					}
					resEnd[cd] = NULLLINE
				}
			}

			// get used resources in this line and add them to the resUsed map
			for _, cu := range strings.Split(body.Lines[line].GetMeta("use"), ":") {
				if resEnd[cu] == NULLLINE {
					resEnd[cu] = line
				}
			}
		}
	}

	for cd, line := range resStart {
		if line != NULLLINE && resEnd[cd] != NULLLINE {
			if resEnd[cd] == LASTLINE {
				for d := 0; d < len(body.Lines); d++ {
					// Insert the metadata
					body.Lines[d].BasmMeta = body.Lines[d].AddMeta("inuse", cd)
				}
			} else {
				for d := 0; d <= resEnd[cd]; d++ {
					// Insert the metadata
					body.Lines[d].BasmMeta = body.Lines[d].AddMeta("inuse", cd)
				}
			}
		}
	}

	fmt.Println(body)

	return nil
}

func fragmentAnalyzer(bi *BasmInstance) error {
	//TODO finish this

	if bi.debug {
		fmt.Println(green("\tProcessing fragments:"))
	}

	// Filter the matchers to select only the symbol based one
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
		if bmline.FilterMatcher(matcher, "symbol") {
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

		if fragment.fragmentBody.GetMeta("template") == "true" {
			if bi.debug {
				fmt.Println(green("\t\t\tFragment is templated, cannot be processed"))
			}
			continue
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
		// use the symbol.
		// This will only work with symbols, not with numbers. For this reason, it will execute after symboltagger and before symbolresolver.

		// Map all the symbols
		symbols := make(map[string]int)

		for i, line := range fBody.Lines {
			if symbol := line.GetMeta("symbol"); symbol != "" {
				if _, exists := symbols[symbol]; exists {
					return errors.New("symbol is specified multiple time: " + symbol)
				} else {
					symbols[symbol] = i
				}
			}
		}

		// Get where the branching instructions are

		for i, line := range fBody.Lines {

			for _, matcher := range filteredMatchers {
				if bmline.MatchMatcher(matcher, line) {
					// TODO Handling the operand
					for _, arg := range line.Elements {
						if j, ok := symbols[arg.GetValue()]; ok {
							fmt.Println("the branch is at line", i, "and the symbol is at line", j)

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

		// TODO Eventually: Union of the copies

		// TODO Temporary copy the main body
		fragment.fragmentBody = branchingBlocks[0]

	}
	// panic("fragmentAnalyzer not finished")
	return nil
}
