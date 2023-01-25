package procbuilder

import (
	"strconv"
	"strings"
)

type Sharedmem struct{}

func (op Sharedmem) Shr_get_name() string {
	return "sharedmem"
}

func (op Sharedmem) Shortname() string {
	return "sh"
}

func (op Sharedmem) GetArchHeader(arch *Arch, shared_constraint string, seq int) string {
	shname := "sh" + strconv.Itoa(seq)
	return ", " + shname + "din, " + shname + "dout, " + shname + "addr, " + shname + "wren, " + shname + "en"
}

func (op Sharedmem) GetArchParams(arch *Arch, shared_constraint string, seq int) string {
	shname := "sh" + strconv.Itoa(seq)
	result := ""
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "din;\n"
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "addr;\n"
	result += "	output " + shname + "wren;\n"
	result += "	output " + shname + "en;\n"
	result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "dout;\n"
	return result
}

func (op Sharedmem) GetCPParams(arch *Arch, shared_constraint string, seq int) string {

	result := ""

	sharemem_num := 0
	if arch.Shared_constraints != "" {
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if soname == "sharedmem" {
				sharemem_num++
			}
		}
	}

	shname := "sh" + strconv.Itoa(seq)

	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "din;\n"
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "addr;\n"
	result += "	output " + shname + "wren;\n"
	result += "	output " + shname + "en;\n"
	result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + shname + "dout;\n"

	if seq == 0 {
		result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] sh_din_i[" + strconv.Itoa(sharemem_num-1) + ":0];\n"
		result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] sh_addr_i [" + strconv.Itoa(sharemem_num-1) + ":0];\n"
		result += "\treg [" + strconv.Itoa(sharemem_num-1) + ":0] sh_wren_i;\n"
		result += "\twire [" + strconv.Itoa(sharemem_num-1) + ":0] sh_en_i;\n"
		result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] sh_dout_i [" + strconv.Itoa(sharemem_num-1) + ":0];\n"
	}

	result += "\tassign " + shname + "din = sh_din_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + shname + "addr = sh_addr_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + shname + "wren = sh_wren_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + shname + "en = sh_en_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign sh_dout_i[" + strconv.Itoa(seq) + "] = " + shname + "dout;\n"
	return result
}
