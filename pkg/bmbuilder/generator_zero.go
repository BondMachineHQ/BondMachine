package bmbuilder

import (
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

func ZeroGenerator(b *BMBuilder, e *bmline.BasmElement, l *bmline.BasmLine) (*bondmachine.Bondmachine, error) {

	if b.debug {
		fmt.Println(green("\t\t\tZero Generator - Start"))
		defer fmt.Println(green("\t\t\tZero Generator - End"))
	}

	return nil, nil
}
