package procbuilder

import (
	//"fmt"
	"strconv"
	//"strings"
)

func (arch *Arch) Write_verilog_testbench(arch_module_name string, processor_name string, rom_name string, flavor string) string {
	regsize := int(arch.Rsize)

	result := ""
	result = result + "module main_tb;\n"
	result = result + "\n"
	result = result + "	reg clk_n, clk_p, reset;\n"
	result = result + "\n"

	for i := 0; i < int(arch.N); i++ {
		result = result + "	reg [" + strconv.Itoa(regsize) + ":0] " + Get_input_name(i) + ";\n"
	}

	for i := 0; i < int(arch.M); i++ {
		result = result + "	wire [" + strconv.Itoa(regsize) + ":0] " + Get_output_name(i) + ";\n"
	}

	result = result + "	" + arch_module_name + " " + arch_module_name + "_impl(clkn, clk_p, reset"

	for i := 0; i < int(arch.N); i++ {
		result = result + ", " + Get_input_name(i)
	}

	for i := 0; i < int(arch.M); i++ {
		result = result + ", " + Get_output_name(i)
	}

	result = result + ");\n"

	result = result + "\n"
	result = result + "	integer tickN;\n"
	result = result + "	localparam TICK=5000;\n"
	result = result + "\n"
	result = result + "	always\n"
	result = result + "	begin\n"
	result = result + "		clk_p = 1;\n"
	result = result + "		clk_n = 0;\n"
	result = result + "		#(TICK/2);\n"
	result = result + "		clk_p = 0;\n"
	result = result + "		clk_n = 1;\n"
	result = result + "		#(TICK/2);\n"
	result = result + "\n"
	result = result + "		tickN = tickN + 1;\n"
	result = result + "		$display(\"--------------Tick %d---------------\", tickN);\n"
	result = result + "	end\n"
	result = result + "\n"
	result = result + "	initial\n"
	result = result + "	begin\n"
	result = result + "		tickN = 1;\n"
	result = result + "		reset = 1;\n"
	result = result + "		$display(\"--------------Tick %d---------------\", tickN);\n"
	result = result + "\n"
	result = result + "		#(3000 * TICK);\n"
	result = result + "\n"
	result = result + "		reset = 0;\n"
	result = result + "		#(10000 * TICK);\n"
	result = result + "		reset = 1;\n"
	result = result + "		#30000;\n"
	result = result + "		reset = 0;\n"
	result = result + "		#(10 * TICK);\n"
	result = result + "\n"
	result = result + "		//$finish;\n"
	result = result + "	end\n"
	result = result + "endmodule\n"
	return result
}

func (arch *Arch) Write_verilog_main(processor_module_name string, rom_module_name string, processor_name string, rom_name string, flavor string) string {
	//rom_word := arch.Max_word()

	result := ""
	result = result + "module main;\n"
	result = result + "endmodule\n"
	return result
}
