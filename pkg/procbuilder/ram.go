package procbuilder

import (
	"strconv"
)

// The Ram ,actually not used since the processor only have internal memory
type Ram struct {
	L uint8 // Number of n-bit memory banks
}

func (ram *Ram) Write_verilog(mach *Machine, ram_module_name string, flavor string) string {
	ram_depth := 1 << ram.L

	result := ""

	result += "`timescale 1ns/1ps\n"
	result += "module " + ram_module_name + "(clk, rst, din, dout, addr, wren, en);\n"
	result += "\n"
	result += "	//--------------Input Ports-----------------------\n"
	result += "	input clk;\n"
	result += "	input rst;\n"
	result += "	input [" + strconv.Itoa(int(ram.L)-1) + ":0] addr;\n"
	result += "	input [" + strconv.Itoa(int(mach.Rsize)-1) + ":0] din;\n"
	result += "	input wren;\n"
	result += "	input en;\n"
	result += "\n"
	result += "	//--------------Inout Ports-----------------------\n"
	result += "	output [" + strconv.Itoa(int(mach.Rsize)-1) + ":0] dout;\n"
	result += "\n"
	result += "	//--------------Reg-------------------------------\n"
	result += "	reg [" + strconv.Itoa(int(mach.Rsize)-1) + ":0] mem [0:" + strconv.Itoa(ram_depth-1) + "];\n"
	result += "\n"
	result += "	reg [" + strconv.Itoa(int(mach.Rsize)-1) + ":0] dout_i;\n"
	result += "\n"
	result += "	// Memory Write Block  \n"
	result += "	// Write Operation we = 1 \n"
	result += "	always @ (posedge clk) \n"
	result += "	begin : MEM_WRITE \n"
	result += "		integer k; \n"
	result += "		if (rst)\n"
	result += "		begin \n"
	result += "			for(k=0;k<" + strconv.Itoa(ram_depth) + ";k=k+1) \n"
	result += "				mem[k] <= #1 " + strconv.Itoa(int(mach.Rsize)) + "'b0; \n"
	result += "		end \n"
	result += "		else if (wren)\n"
	result += "			mem[addr] <= #1 din;\n"
	result += "	end \n"
	result += "\n"
	result += "	// Memory Read Block\n"
	result += "	// Read Operation when we = 0 and oe = 1 \n"
	result += "	always @ (posedge clk) \n"
	result += "	begin : MEM_READ \n"
	result += "		if (!wren)\n"
	result += "			dout_i <= #1 mem[addr];\n"
	result += "	end\n"
	result += "\n"
	result += "	assign dout = dout_i;\n"
	result += "\n"
	result += "endmodule \n"

	return result
}
