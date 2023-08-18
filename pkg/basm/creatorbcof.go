package basm

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bcof"
)

// Assembler2BCOF transform an assembled instance into a BCOF file
func (bi *BasmInstance) Assembler2BCOF() error {

	if bi.debug {
		fmt.Println("\t" + green("BCOF metadata"))
	}

	registerSize := bi.global.GetMeta("registersize")

	if registerSize == "" {
		return errors.New("register size not specified")
	}

	var rSize uint8
	if size, err := strconv.Atoi(registerSize); err == nil {
		if 0 < size && size < 256 {
			rSize = uint8(size)
		} else {
			return errors.New("wrong value for register size")
		}
	} else {
		return errors.New("register size not valid")
	}
	if bi.debug {
		fmt.Println("\t\t"+green("register size:"), rSize)
	}

	if bi.debug {
		fmt.Println("\t" + green("BCOF creation"))
	}

	bi.outBCOF = bcof.NewBCOF(uint32(rSize))
	bi.outBCOF.SetId(1)
	bi.outBCOF.SetSignature("bm")

	idDone := make(map[int]struct{})

	for _, cp := range bi.cps {
		if bi.debug {
			fmt.Println("\t\t" + green("CP: ") + yellow(cp.GetValue()))
		}

		ramCode := cp.GetMeta("ramcode")
		if bi.debug {
			if ramCode != "" {
				fmt.Println("\t\t - " + green("ram code: ") + yellow(ramCode))
			} else {
				fmt.Println("\t\t - " + green("ram code: ") + yellow("not specified"))
			}
		}

		ramData := cp.GetMeta("ramdata")
		if bi.debug {
			if ramData != "" {
				fmt.Println("\t\t - " + green("ram data: ") + yellow(ramData))
			} else {
				fmt.Println("\t\t - " + green("ram data: ") + yellow("not specified"))
			}
		}

		if ramCode == "" && ramData == "" {
			if bi.debug {
				fmt.Println("\t\t - " + purple("no ram code or ram data specified, skipping"))
			}
			continue
		}

		// To identify the CP id we have to check the following:
		// - if the BM came from this assembler run, the BMinfo file will contain the CP id and we have to use it
		// - if the BM came from the CLI, and the BMinfo is provided, we have to use the CP id from the BMinfo file
		// - if the BM came from the CLI, and the BMinfo is not provided, we will use the convention of the CP name in the form of "cpN" (where N is a number).

		// Search th ID in the BMinfo data
		if bi.BMinfo != nil {
			for id, name := range bi.BMinfo.CPNames {
				if name == cp.GetValue() {
					if _, ok := idDone[id]; ok {
						return errors.New("CP id already used")
					}
					idDone[id] = struct{}{}
					if bi.debug {
						fmt.Println("\t\t - " + green("cpId: ") + yellow(id))
					}

					prog := ""
					for _, line := range bi.sections[ramCode].sectionBody.Lines {
						prog += line.Operation.GetValue()
						for _, arg := range line.Elements {
							prog += " " + arg.GetValue()
						}
						prog += "\n"
					}
					// Get the machine
					myMachine := bi.result
					myArch := myMachine.Domains[myMachine.Processors[id]].Arch
					// TODO Finish this
					if prog, err := myArch.Assembler([]byte(prog)); err == nil {
						fmt.Println(prog)
					} else {
						return err
					}

					// Create the BCOF entry
					// bcofData := &bcof.BCOFData{
					// 	Id:        uint32(id),
					// 	Rsize:     uint32(rSize),
					// 	Signature: cp.GetValue(), // TODO: for now we use the CP name as signature
					// 	Payload:   []byte(ramCode),
					// }

					break
				}
			}
		}

		re := regexp.MustCompile("^cp(?P<cpId>[0-9]+)$")
		if re.MatchString(cp.GetValue()) {
			cpId := re.ReplaceAllString(cp.GetValue(), "${cpId}")
			if bi.debug {
				fmt.Println("\t\t - " + green("cpId: ") + yellow(cpId))
			}
			//TODO
		} else {
			//TODO
		}
	}
	//TODO
	return nil
}
