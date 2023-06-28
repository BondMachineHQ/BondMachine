package basm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

func fragmentAnalyzer(bi *BasmInstance) error {
	//TODO finish this

	if bi.debug {
		fmt.Println(green("\tProcessing fragments:"))
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

	}
	// panic("fragmentAnalyzer not finished")
	return nil
}
