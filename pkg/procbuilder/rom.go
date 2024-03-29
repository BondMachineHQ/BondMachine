package procbuilder

import (
	"strconv"
)

// The Rom
type Rom struct {
	O uint8 // Number of ROM cells (the program storage)
}

func (rom *Rom) String() string {
	return ""
}

func (rom *Rom) Write_verilog(mach *Machine, rom_module_name string, flavor string) string {
	wordSize := mach.Max_word()

	result := ""

	// Module header
	result += "`timescale 1ns/1ps\n"
	result += "module " + rom_module_name + "(input [" + strconv.Itoa(int(rom.O)-1) + ":0] rom_bus, output [" + strconv.Itoa(wordSize-1) + ":0] rom_value);\n"

	result += "\treg [" + strconv.Itoa(wordSize-1) + ":0] _rom [0:" + strconv.Itoa((1<<rom.O)-1) + "];\n"
	result += "\tinitial\n"
	result += "\tbegin\n"

	i := 0
	for _, inst := range mach.Program.Slocs {
		result += "\t_rom[" + strconv.Itoa(i) + "] = " + strconv.Itoa(wordSize) + "'b" + inst + ";\n"
		i++
	}

	for _, inst := range mach.Data.Vars {
		result += "\t_rom[" + strconv.Itoa(i) + "] = " + strconv.Itoa(wordSize) + "'b" + inst + ";\n"
		i++
	}

	result += "\tend\n"
	result += "\tassign rom_value = _rom[rom_bus];\n"
	result += "endmodule\n"

	return result
}
