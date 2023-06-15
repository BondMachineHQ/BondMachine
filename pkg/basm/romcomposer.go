package basm

import (
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

func romComposer(bi *BasmInstance) error {

	if bi.debug {
		fmt.Println(green("\tConnecting CP code and ROM:"))
	}

	// Loop over the CP code and data and connect them to the ROM
	for _, cp := range bi.cps {
		code := cp.GetMeta("romcode")
		data := cp.GetMeta("romdata")
		if code == "" || data == "" {
			continue
		}

		if bi.debug {
			fmt.Println(green("\t\tConnecting code " + code + " and data " + data))
		}

		sectionCode := bi.sections[code]

		cpNewSectionName := "coderom_" + code + "_" + data

		newSection := new(BasmSection)
		newSection.sectionName = cpNewSectionName
		newSection.sectionType = setcRomText
		newSection.sectionBody = new(bmline.BasmBody)

		newSection.sectionBody.Lines = make([]*bmline.BasmLine, 0)

		for _, line := range sectionCode.sectionBody.Lines {
			newLine := new(bmline.BasmLine)
			newLine.Operation = new(bmline.BasmElement)
			newLine.Operation.SetValue(line.Operation.GetValue())
			// TODO Finish this

			newSection.sectionBody.Lines = append(newSection.sectionBody.Lines, newLine)
		}

		bi.sections[newSection.sectionName] = newSection

	}

	return nil
}
