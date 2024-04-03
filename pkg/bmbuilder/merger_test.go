package bmbuilder

import (
	"testing"

	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

func TestMerge(t *testing.T) {
	bm1 := new(bondmachine.Bondmachine)
	bm2 := new(bondmachine.Bondmachine)
	bld := BMBuilder{}
	bld.BMBuilderInit()
	_, err := bld.BMMerge(bm1, bm2, nil)
	if err != nil {
		t.Errorf("Error merging bondmachines: %v", err)
	}
}
