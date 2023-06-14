package procbuilder

import (
	"fmt"
	"strconv"
	"strings"
)

func (mach *Machine) Specs() string {
	// TODO - implement this
	result := "    ROM width/Word size: " + strconv.Itoa(int(mach.Max_word())) + "\n"
	depth := uint64(1) << uint64(mach.O)
	result += "    ROM depth: " + strconv.Itoa(int(depth)) + "\n"
	result += "    RAM width: " + strconv.Itoa(int(mach.Rsize)) + "\n"
	depth = uint64(1) << uint64(mach.L)
	result += "    RAM depth: " + strconv.Itoa(int(depth)) + "\n"
	regs := uint64(1) << uint64(mach.R)
	result += "    Registers: " + strconv.Itoa(int(regs)) + "\n"
	result += "    Inputs: " + strconv.Itoa(int(mach.N)) + "\n"
	result += "    Outputs: " + strconv.Itoa(int(mach.M)) + "\n"
	ops := ""
	for i := 0; i < len(mach.Conproc.Op); i++ {
		ops += fmt.Sprint(mach.Conproc.Op[i].Op_get_name())
		if i < len(mach.Conproc.Op)-1 {
			ops += ","
		}
	}
	result += "    ISA: " + ops + "\n"
	modes := ""
	for i := 0; i < len(mach.Modes); i++ {
		modes += fmt.Sprint(mach.Modes[i])
		if i < len(mach.Modes)-1 {
			modes += ","
		}
	}
	result += "    Modes: " + modes + "\n"
	result += "    ROM Code:\n"
	code, _ := mach.Disassembler()
	for i, line := range strings.Split(code, "\n") {
		if line != "" {
			bins := mach.Program.Slocs[i]
			result += "      " + fmt.Sprintf("%03d", i) + " - " + bins + " - " + line + "\n"
		}
	}
	result += "    ROM Data:\n"
	for i, bins := range mach.Data.Vars {
		result += "      " + fmt.Sprintf("%03d", i) + " - " + bins + "\n"
	}

	return result

}
