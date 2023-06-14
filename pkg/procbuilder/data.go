package procbuilder

// The machine is an architecture provided with and execution code and an intial state
type Data struct {
	Vars []string
}

func (d *Data) String() string {
	result := ""
	for _, line := range d.Vars {
		result = result + line + "\n"
	}
	return result
}
