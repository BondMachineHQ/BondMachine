package bmnumbers

// TODO finish this

type DynamicalType interface {
	GetName() string
	MatchName(string) bool
	CreateType(string, interface{}) (BMNumberType, error)
	// Save the type to a file
	// Load the type from a file
}

func EventuallyCreateType(name string, param interface{}) (bool, error) {
	for _, dyn := range AllDynamicalTypes {
		if dyn.MatchName(name) {
			for _, op := range AllTypes {
				if op.GetName() == name {
					return false, nil
				}
			}

			if newType, err := dyn.CreateType(name, param); err != nil {
				return false, err
			} else {
				AllTypes = append(AllTypes, newType)
				for k, v := range newType.importMatchers() {
					AllMatchers[k] = v
				}
				return true, nil
			}
		}
	}
	return false, nil
}
