package procbuilder

import (
	"strconv"
	"strings"
)

type Queue struct{}

func (op Queue) Shr_get_name() string {
	return "queue"
}

func (op Queue) Shortname() string {
	return "q"
}

func (op Queue) GetArchHeader(arch *Arch, shared_constraint string, seq int) string {
	qName := "q" + strconv.Itoa(seq)
	return ", " + qName + "din, " + qName + "dout, " + qName + "addr, " + qName + "wren, " + qName + "en"
}

func (op Queue) GetArchParams(arch *Arch, shared_constraint string, seq int) string {
	qName := "q" + strconv.Itoa(seq)
	result := ""
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + qName + "din;\n"
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + qName + "addr;\n"
	result += "	output " + qName + "wren;\n"
	result += "	output " + qName + "en;\n"
	result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + qName + "dout;\n"
	return result
}

func (op Queue) GetCPParams(arch *Arch, shared_constraint string, seq int) string {

	result := ""

	qNum := 0
	if arch.Shared_constraints != "" {
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if soname == "queue" {
				qNum++
			}
		}
	}

	qName := "q" + strconv.Itoa(seq)

	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + qName + "din;\n"
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + qName + "addr;\n"
	result += "	output " + qName + "wren;\n"
	result += "	output " + qName + "en;\n"
	result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + qName + "dout;\n"

	if seq == 0 {
		result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] q_din_i[" + strconv.Itoa(qNum-1) + ":0];\n"
		result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] q_addr_i [" + strconv.Itoa(qNum-1) + ":0];\n"
		result += "\treg [" + strconv.Itoa(qNum-1) + ":0] q_wren_i;\n"
		result += "\twire [" + strconv.Itoa(qNum-1) + ":0] q_en_i;\n"
		result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] q_dout_i [" + strconv.Itoa(qNum-1) + ":0];\n"
	}

	result += "\tassign " + qName + "din = q_din_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + qName + "addr = q_addr_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + qName + "wren = q_wren_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + qName + "en = q_en_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign q_dout_i[" + strconv.Itoa(seq) + "] = " + qName + "dout;\n"
	return result
}
