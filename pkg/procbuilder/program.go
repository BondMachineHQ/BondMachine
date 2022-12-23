package procbuilder

import ()

// The machine is an architecture provided with and execution code and an intial state
type Program struct {
	Slocs []string
}

func (prog *Program) String() string {
	result := ""
	for _, line := range prog.Slocs {
		result = result + line + "\n"
	}
	return result
}
