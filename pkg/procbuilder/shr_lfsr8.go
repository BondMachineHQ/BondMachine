package procbuilder

import (
	"strconv"
	//	"strings"
)

type Lfsr8 struct{}

func (op Lfsr8) Shr_get_name() string {
	return "lfsr8"
}

func (op Lfsr8) Shortname() string {
	return "lfsr8"
}

func (op Lfsr8) GetArchHeader(arch *Arch, shared_constraint string, seq int) string {
	brname := "lfsr8" + strconv.Itoa(seq)
	return ", " + brname + "out"
}

func (op Lfsr8) GetArchParams(arch *Arch, shared_constraint string, seq int) string {

	brname := "lfsr8" + strconv.Itoa(seq)

	result := ""
	result += "	input [7:0] " + brname + "out;\n"

	return result
}

func (op Lfsr8) GetCPParams(arch *Arch, shared_constraint string, seq int) string {
	result := op.GetArchParams(arch, shared_constraint, seq)
	return result
}
