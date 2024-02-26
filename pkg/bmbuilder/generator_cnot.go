package bmbuilder

import (
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

func CnotGenerator(b *BMBuilder, e *bmline.BasmElement, l *bmline.BasmLine) (*bondmachine.Bondmachine, error) {

	if b.debug {
		fmt.Println(green("\t\t\tCNOT Generator - Start"))
		defer fmt.Println(green("\t\t\tCNOT Generator - End"))
	}

	return nil, nil
}
