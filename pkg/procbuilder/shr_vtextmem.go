package procbuilder

import (
	"strconv"
	"strings"
)

type Vtextmem struct{}

func (op Vtextmem) Shr_get_name() string {
	return "vtextmem"
}

func (op Vtextmem) Shortname() string {
	return "vtm"
}

func (op Vtextmem) Get_header(arch *Arch, shared_constraint string, seq int) string {
	shname := "vtm" + strconv.Itoa(seq)
	return ", " + shname + "din, " + shname + "addr, " + shname + "wren, " + shname + "en"
}

func (op Vtextmem) Get_params(arch *Arch, shared_constraint string, seq int) string {
	shname := "vtm" + strconv.Itoa(seq)
	result := ""
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "din;\n"
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "addr;\n"
	result += "	output " + shname + "wren;\n"
	result += "	output " + shname + "en;\n"
	//	result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "dout;\n"
	return result
}

func (op Vtextmem) Get_internal_params(arch *Arch, shared_constraint string, seq int) string {

	result := ""

	sharemem_num := 0
	if arch.Shared_constraints != "" {
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if soname == "Vtextmem" {
				sharemem_num++
			}
		}
	}

	shname := "vtm" + strconv.Itoa(seq)

	result += "	output [7:0] " + shname + "din;\n"
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "addr;\n"
	result += "	output " + shname + "wren;\n"
	result += "	output " + shname + "en;\n"
	//	result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "dout;\n"

	if seq == 0 {
		result += "\treg [7:0] " + shname + "_din_i;\n"
		result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "_addr_i;\n"
		result += "\treg " + shname + "_wren_i;\n"
		result += "\twire " + shname + "_en_i;\n"
		//result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "_dout_i [" + strconv.Itoa(sharemem_num-1) + ":0];\n"
	}

	result += "\tassign " + shname + "din = " + shname + "_din_i;\n"
	result += "\tassign " + shname + "addr = " + shname + "_addr_i;\n"
	result += "\tassign " + shname + "wren = " + shname + "_wren_i;\n"
	result += "\tassign " + shname + "en = " + shname + "_en_i;\n"
	//result += "\tassign " + shname + "_dout_i[" + strconv.Itoa(seq) + "] = " + shname + "dout;\n"
	return result
}
