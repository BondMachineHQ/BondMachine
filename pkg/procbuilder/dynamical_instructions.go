package procbuilder

import (
	"github.com/BondMachineHQ/BondMachine/pkg/bmconfig"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

type DynamicInstruction interface {
	GetName() string
	MatchName(string) bool
	CreateInstruction(string) (Opcode, error)
	HLAssemblerGeneratorMatch(*bmconfig.BmConfig) []string
	HLAssemblerGeneratorList(*bmconfig.BmConfig, *bmline.BasmLine) []string
}

func EventuallyCreateInstruction(name string) (bool, error) {
	for _, dyn := range AllDynamicalInstructions {
		if dyn.MatchName(name) {
			for _, op := range Allopcodes {
				if op.Op_get_name() == name {
					return false, nil
				}
			}

			if newOp, err := dyn.CreateInstruction(name); err != nil {
				return false, err
			} else {
				Allopcodes = append(Allopcodes, newOp)
				return true, nil
			}
		}
	}
	return false, nil
}
