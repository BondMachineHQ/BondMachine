package procbuilder

import "strconv"

func (mach *Machine) Specs() string {
	// TODO - implement this
	result := "    ROM width/Word size: " + strconv.Itoa(int(mach.Max_word())) + "\n"
	depth := uint64(1) << uint64(mach.O)
	result += "    ROM depth: " + strconv.Itoa(int(depth)) + "\n"
	depth = uint64(1) << uint64(mach.L)
	result += "    RAM depth: " + strconv.Itoa(int(depth)) + "\n"
	return result
}
