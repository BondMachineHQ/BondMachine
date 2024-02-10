package bmbuilder

import (
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/basm"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

func BasmGenerator(b *BMBuilder, e *bmline.BasmElement, l *bmline.BasmLine) (*bondmachine.Bondmachine, error) {

	if b.debug {
		fmt.Println(green("\t\t\tBasmGenerator - Start"))
		defer fmt.Println(green("\t\t\tBasmGenerator - End"))
	}

	return basm.BasmGenerator(e, l)

}
