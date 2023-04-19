package procbuilder

type DynamicInstruction interface {
	GetName() string
	MatchName(string) bool
	CreateInstruction(string) (Opcode, error)
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
