package basm

import "fmt"

func templateFinalizer(bi *BasmInstance) error {

	if err := bi.templateAutoMark(); err != nil {
		return err
	}
	// Loop over the sections to find eventual alternatives
	for sectName, section := range bi.sections {

		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			if section.sectionBody.GetMeta("template") != "true" {
				if bi.debug {
					fmt.Println(green("\t\t\tNot templated section, skipping"))
				}
				continue
			}

			// TODO: Implement this

			// body := section.sectionBody

			// for _, line := range body.Lines {

			// 	operation := line.Operation

			// 	for _, arg := range line.Elements {

			// 	}

			// }

		}
	}
	return nil
}
