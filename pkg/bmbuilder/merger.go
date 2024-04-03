package bmbuilder

import (
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

const (
	FinalBMInput = uint8(1) << iota
	FinalBMOutput
	FirstBMInput
	FirstBMOutput
	SecondBMInput
	SecondBMOutput
)

type BMEndPoint struct {
	BType uint8
	BNum  int
}

type BMLink struct {
	E1 BMEndPoint
	E2 BMEndPoint
}

type BMConnections struct {
	Links []BMLink
}

func (bld *BMBuilder) BMMerge(bm1 *bondmachine.Bondmachine, bm2 *bondmachine.Bondmachine, l *BMConnections) (*bondmachine.Bondmachine, error) {
	// Merge two bondmachines
	// bm1 and bm2 are the bondmachines to merge
	// Returns the merged bondmachine

	return nil, nil
}
