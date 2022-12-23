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

func (op Barrier) Get_header(arch *Arch, shared_constraint string, seq int) string {
	brname := "br" + strconv.Itoa(seq)
	return ", " + brname + "hit, " + brname + "ishitted, " + brname + "tout"
}

func (op Barrier) Get_params(arch *Arch, shared_constraint string, seq int) string {

	brname := "br" + strconv.Itoa(seq)

	result := ""
	result += "	output " + brname + "hit;\n"
	result += "	input " + brname + "ishitted;\n"
	result += "	input " + brname + "tout;\n"

	return result
}

func (op Barrier) Get_internal_params(arch *Arch, shared_constraint string, seq int) string {
	result := op.Get_params(arch, shared_constraint, seq)
	return result
}
