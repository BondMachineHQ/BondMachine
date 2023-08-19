package basm

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bcof"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
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

	// Compute the register size padded to 8 bit multiples (the BCOF format requires bytes)
	rSizePad := 8 * ((rSize + 7) / 8)

	if bi.debug {
		fmt.Println("\t\t"+green("register size padded:"), rSizePad)
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

		var program *procbuilder.Program

		cpPayload := bcof.NewBCOFData(uint32(rSize))

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

					cpPayload.SetId(uint32(id))
					cpPayload.SetSignature("cp" + strconv.Itoa(id)) // TODO: temporary name for the signature

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

					if assembled, err := myArch.Assembler([]byte(prog)); err == nil {
						program = &assembled
					} else {
						return err
					}

				}
			}
		} else {

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

		for _, line := range program.Slocs {
			// Pad the line to the padded register size
			for len(line) < int(rSizePad) {
				line = "0" + line
			}
			// Read the binary number with bmnumbers
			// fmt.Println("0b<" + strconv.Itoa(int(rSizePad)) + ">" + line)
			num, _ := bmnumbers.ImportString("0b<" + strconv.Itoa(int(rSizePad)) + ">" + line)
			cpPayload.Payload = append(cpPayload.Payload, num.GetBytes()...)
		}

		bi.outBCOF.AddData(cpPayload)

		if bi.debug {
			fmt.Println(green("\t\tBCOF entry created dump: "))
			fmt.Println(green("\t\t----"))
			fmt.Println(cpPayload.Dump())
			fmt.Println(green("\t\t----"))
		}
	}

	return nil
}

func (bi *BasmInstance) GetBCOF() *bcof.BCOFEntry {
	return bi.outBCOF
}
