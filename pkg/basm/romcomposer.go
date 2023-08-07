package basm

import (
	"errors"
	"fmt"
	"strconv"
)

func romComposer(bi *BasmInstance) error {

	if bi.debug {
		fmt.Println(green("\tConnecting CP code and ROM:"))
	}

	removalList := make(map[string]struct{})

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
		sectionData := bi.sections[data]

		sectionLength := len(sectionCode.sectionBody.Lines)

		if bi.debug {
			fmt.Println(green("\t\t\tCode section length: " + fmt.Sprintf("%d", sectionLength)))
		}

		cpNewSectionName := "coderom_" + code + "_" + data

		newSection := new(BasmSection)
		newSection.sectionName = cpNewSectionName
		newSection.sectionType = sectRomText
		newSection.sectionBody = sectionCode.sectionBody.Copy()

		// Collection locations
		locations := make(map[string]string)
		for _, vari := range sectionData.sectionBody.Lines {
			varName := vari.Operation.GetValue()
			offset := vari.Operation.GetMeta("offset")
			location, _ := strconv.Atoi(offset)
			location += sectionLength
			locations[varName] = strconv.Itoa(location)
		}

		// Searching for ROM variable to be translated into the ROM address

		body := newSection.sectionBody
		usefullSection := false

		for _, line := range body.Lines {
			for _, arg := range line.Elements {
				if arg.GetMeta("type") == "rom" && arg.GetMeta("romaddressing") == "variable" {
					removalList[code] = struct{}{}
					usefullSection = true
					romVariable := arg.GetMeta("variable")
					if romVariable == "" {
						return errors.New("ROM variable cannot be empty")
					}

					if loc, ok := locations[romVariable]; ok {
						arg.SetValue(loc)
						arg.BasmMeta = arg.BasmMeta.SetMeta("type", "number")
						arg.BasmMeta = arg.BasmMeta.SetMeta("numbertype", "unsigned")
						arg.BasmMeta.RmMeta("romaddressing")
						arg.BasmMeta.RmMeta("variable")
					} else {
						return errors.New("ROM variable " + romVariable + " not found")
					}

				}
			}
		}

		// TODO Temporary fix
		usefullSection = true

		bi.sections[newSection.sectionName] = newSection
		if usefullSection {
			cp.BasmMeta = cp.BasmMeta.SetMeta("romcode", newSection.sectionName)
		} else {
			removalList[newSection.sectionName] = struct{}{}
		}

	}

	// Removing processed sections
	for rem := range removalList {
		delete(bi.sections, rem)
	}

	return nil
}
