package procbuilder

import (
	"strconv"
	"strings"
)

type Stack struct{}

func (op Stack) Shr_get_name() string {
	return "stack"
}

func (op Stack) Shortname() string {
	return "st"
}

func (op Stack) Get_header(arch *Arch, shared_constraint string, seq int) string {
	stackName := "st" + strconv.Itoa(seq)
	return ", " + stackName + "din, " + stackName + "dout, " + stackName + "addr, " + stackName + "wren, " + stackName + "en"
}

func (op Stack) Get_params(arch *Arch, shared_constraint string, seq int) string {
	stackName := "st" + strconv.Itoa(seq)
	result := ""
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + stackName + "din;\n"
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + stackName + "addr;\n"
	result += "	output " + stackName + "wren;\n"
	result += "	output " + stackName + "en;\n"
	result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + stackName + "dout;\n"
	return result
}

func (op Stack) Get_internal_params(arch *Arch, shared_constraint string, seq int) string {

	result := ""

	stackNum := 0
	if arch.Shared_constraints != "" {
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if soname == "stack" {
				stackNum++
			}
		}
	}

	stackName := "st" + strconv.Itoa(seq)

	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + stackName + "din;\n"
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + stackName + "addr;\n"
	result += "	output " + stackName + "wren;\n"
	result += "	output " + stackName + "en;\n"
	result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + stackName + "dout;\n"

	if seq == 0 {
		result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] q_din_i[" + strconv.Itoa(stackNum-1) + ":0];\n"
		result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] q_addr_i [" + strconv.Itoa(stackNum-1) + ":0];\n"
		result += "\treg [" + strconv.Itoa(stackNum-1) + ":0] q_wren_i;\n"
		result += "\twire [" + strconv.Itoa(stackNum-1) + ":0] q_en_i;\n"
		result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] q_dout_i [" + strconv.Itoa(stackNum-1) + ":0];\n"
	}

	result += "\tassign " + stackName + "din = q_din_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + stackName + "addr = q_addr_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + stackName + "wren = q_wren_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + stackName + "en = q_en_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign q_dout_i[" + strconv.Itoa(seq) + "] = " + stackName + "dout;\n"
	return result
}
