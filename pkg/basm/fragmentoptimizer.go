package basm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

const (
	OPT_INVALID_UNUSED = 0 + iota
)

func (bi *BasmInstance) listOptimizations() string {
	return "invalidunused\n"
}

func (bi *BasmInstance) ActivateOptimization(opt string) error {
	if bi.optimizations == nil {
		bi.optimizations = make(map[int]struct{})
	}

	if opt == "all" {
		bi.optimizations[OPT_INVALID_UNUSED] = struct{}{}
		return nil
	}

	switch opt {
	case "invalidunused":
		bi.optimizations[OPT_INVALID_UNUSED] = struct{}{}
	default:
		return errors.New("Unknown optimization: " + opt)
	}
	return nil
}

func fragmentOptimizer(bi *BasmInstance) error {
	//TODO finish this

	if bi.debug {
		fmt.Println(green("\tProcessing fragments:"))
	}

	for fragName, fragment := range bi.fragments {
		if bi.debug {
			fmt.Println(green("\t\tFragment: ")+fragName, fragment)
		}

		fBody := fragment.fragmentBody

		// Start the optimization processes

		if _, ok := bi.optimizations[OPT_INVALID_UNUSED]; ok {

			removeList := make([]int, 0)

			if bi.debug {
				fmt.Println(green("\t\t\tUnused wrote resources removal:"))
			}

		lineLoop:
			for i, line := range fBody.Lines {

				if bi.debug {
					fmt.Print(green("\t\t\t\tLine: ") + line.String())
				}

				inv := line.GetMeta("inv")
				if inv != "" {
					invSplit := strings.Split(inv, ":")
					inUse := line.GetMeta("inuse")
					if inUse != "" {
						inUseSplit := strings.Split(inUse, ":")

						for _, invRes := range invSplit {

							if stringInSlice(invRes, inUseSplit) {
								if bi.debug {
									fmt.Println(" - " + green("kept"))
									continue lineLoop
								}
							}
						}

						if bi.debug {
							fmt.Println(" - " + red("removed"))
						}
						removeList = append(removeList, i)
						continue
					} else {
						if bi.debug {
							fmt.Println(" - " + red("removed"))
						}
						removeList = append(removeList, i)
						continue
					}
				}
			}

			newBody := new(bmline.BasmBody)
			for k, v := range fBody.BasmMeta.LoopMeta() {
				newBody.BasmMeta = newBody.BasmMeta.SetMeta(k, v)
			}
			newBody.Lines = make([]*bmline.BasmLine, len(fBody.Lines)-len(removeList))
			j := 0
			for i, line := range fBody.Lines {
				if j == len(removeList) || i != removeList[j] {
					newBody.Lines[i-j] = line.Copy()
				} else {
					j++
				}
			}

			fragment.fragmentBody = newBody
		}
	}
	// panic("fragmentOptimizer not finished")
	return nil
}
