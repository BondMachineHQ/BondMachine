package bondmachine

import (
	"strconv"
	"strings"
)

// The placeholder struct

type Lfsr8 struct{}

func (op Lfsr8) Shr_get_name() string {
	return "lfsr8"
}

func (op Lfsr8) Shr_get_desc() string {
	return "Lfsr8"
}

func (op Lfsr8) Shortname() string {
	return "lfsr8"
}

func (op Lfsr8) GV_config(element uint8) string {
	result := ""
	switch element {
	case GVNODEINPROC:
		result += "style=filled fillcolor=mediumorchid1 color=black"
	case GVNODE:
		result += "style=filled fillcolor=mediumorchid1 color=black"
	case GVEDGE:
		result += "arrowhead=none"
	case GVCLUS:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey65"
	case GVCLUSINPROC:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey65"
	}
	return result
}

func (op Lfsr8) Instantiate(s string) (Shared_instance, bool) {
	if strings.HasPrefix(s, "lfsr8:") {
		if len(s) > 6 {
			if seed, ok := strconv.Atoi(s[6:]); ok == nil {
				result := new(Lfsr8_instance)
				result.Shared_element = op
				result.Seed = uint8(seed)
				return *result, true
			}
		}
	}
	return nil, false
}

// The instance struct

type Lfsr8_instance struct {
	Shared_element
	Seed uint8
}

func (sm Lfsr8_instance) String() string {
	return "lfsr8:" + strconv.Itoa(int(sm.Seed))
}

func (sm Lfsr8_instance) Write_verilog(bmach *Bondmachine, so_index int, lfsr8_name string, flavor string) string {

	result := ""

	result += "`timescale 1ns/1ps\n"
	result += "module " + lfsr8_name + "(clk, reset, lfsr8out);\n"
	result += "\n"
	result += "	//--------------Input Ports-----------------------\n"
	result += "	input clk;\n"
	result += "	input reset;\n"

	result += "\n"
	result += "	//--------------Output Ports-----------------------\n"
	result += "\n"

	result += "	output reg [7:0] lfsr8out;\n"
	result += "\n"

	result += "	wire feedback;\n"
	result += "\n"

	result += "	initial begin\n"
	result += "		lfsr8out = 8'd" + strconv.Itoa(int(sm.Seed)) + ";\n"
	result += "	end\n"
	result += "\n"

	result += "	assign feedback=(lfsr8out[7]^(lfsr8out[5]^(lfsr8out[4] ^ lfsr8out[3])));\n"
	result += "\n"

	result += "	always @ (posedge clk) begin\n"
	result += "		if (reset) begin\n"
	result += "			lfsr8out <= 8'd" + strconv.Itoa(int(sm.Seed)) + ";\n"
	result += "		end\n"
	result += "		else\n"
	result += "		begin\n"
	result += "			lfsr8out <= {lfsr8out[6:0],feedback};\n"
	result += "		end\n"
	result += "	end\n"
	result += "\n"

	result += "endmodule\n"
	result += "\n"

	return result
}

func (sm Lfsr8_instance) GetPerProcPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += "\twire [7:0] p" + strconv.Itoa(proc_id) + soname + "out;\n"
		result += "\n"
	}
	return result
}

func (sm Lfsr8_instance) GetPerProcPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += ", p" + strconv.Itoa(proc_id) + soname + "out"
	}
	return result
}

func (sm Lfsr8_instance) GetExternalPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Lfsr8_instance) GetExternalPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}
