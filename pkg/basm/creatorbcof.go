package basm

import "fmt"

// Assembler2BCOF transform an assembled instance into a BCOF file
func (bi *BasmInstance) Assembler2BCOF() error {

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
		// if the name of the CP (cp.GetValue()) is in the form of "cpN" (where N is a number) then we have to compile against the N-th CP of the result BM
		// otherwise we have to compile against the CP with the same name that can be found into the bminfo file
	}
	//TODO
	return nil
}
