package basm

import (
	"fmt"
	"strings"
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

		// TODO rearrange resorces in the order they are used

	}
	// panic("fragmentAnalyzer not finished")
	return nil
}
