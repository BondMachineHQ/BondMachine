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
	rom_word := mach.Max_word()

	result := ""

	// Module header
	result += "`timescale 1ns/1ps\n"
	result += "module " + rom_module_name + "(input [" + strconv.Itoa(int(rom.O)-1) + ":0] rom_bus, output [" + strconv.Itoa(rom_word-1) + ":0] rom_value);\n"

	result += "\treg [" + strconv.Itoa(rom_word-1) + ":0] _rom [0:" + strconv.Itoa((1<<rom.O)-1) + "];\n"
	result += "\tinitial\n"
	result += "\tbegin\n"

	for i, inst := range mach.Program.Slocs {
		result += "\t_rom[" + strconv.Itoa(i) + "] = " + strconv.Itoa(rom_word) + "'b" + inst + ";\n"
	}

	result += "\tend\n"
	result += "\tassign rom_value = _rom[rom_bus];\n"
	result += "endmodule\n"

	return result
}
