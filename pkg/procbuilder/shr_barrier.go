package procbuilder

import (
	"strconv"
	//	"strings"
)

type Barrier struct{}

func (op Barrier) Shr_get_name() string {
	return "barrier"
}

func (op Barrier) Shortname() string {
	return "br"
}

func (op Barrier) GetArchHeader(arch *Arch, shared_constraint string, seq int) string {
	brname := "br" + strconv.Itoa(seq)
	return ", " + brname + "hit, " + brname + "ishitted, " + brname + "tout"
}

func (op Barrier) GetArchParams(arch *Arch, shared_constraint string, seq int) string {

	brname := "br" + strconv.Itoa(seq)

	result := ""
	result += "	output " + brname + "hit;\n"
	result += "	input " + brname + "ishitted;\n"
	result += "	input " + brname + "tout;\n"

	return result
}

func (op Barrier) GetCPParams(arch *Arch, shared_constraint string, seq int) string {
	result := op.GetArchParams(arch, shared_constraint, seq)
	return result
}
