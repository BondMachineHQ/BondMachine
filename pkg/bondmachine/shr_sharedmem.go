package bondmachine

import (
	//"fmt"
	"math"
	"strconv"
	"strings"
)

// The placeholder struct

type Sharedmem struct{}

func (op Sharedmem) Shr_get_name() string {
	return "sharedmem"
}

func (op Sharedmem) Shr_get_desc() string {
	return "Sharedmem"
}

func (op Sharedmem) Shortname() string {
	return "sh"
}

func (op Sharedmem) GV_config(element uint8) string {
	result := ""
	switch element {
	case GVNODEINPROC:
		result += "style=filled fillcolor=cyan color=black"
	case GVNODE:
		result += "style=filled fillcolor=cyan color=black"
	case GVEDGE:
		result += "arrowhead=none"
	case GVCLUS:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey70"
	case GVCLUSINPROC:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey70"
	}
	return result
}

func (op Sharedmem) Instantiate(s string) (Shared_instance, bool) {
	if strings.HasPrefix(s, "sharedmem:") {
		if len(s) > 10 {
			if depth, ok := strconv.Atoi(s[10:]); ok == nil {
				result := new(Sharedmem_instance)
				result.Shared_element = op
				result.Depth = depth
				return *result, true
			}
		}
	}
	return nil, false
}

// The instance struct

type Sharedmem_instance struct {
	Shared_element
	Depth int
}

func (sm Sharedmem_instance) String() string {
	return "sharedmem:" + strconv.Itoa(sm.Depth)
}

func (sm Sharedmem_instance) Write_verilog(bmach *Bondmachine, so_index int, sharedmem_name string, flavor string) string {

	ram_depth := 1 << uint8(sm.Depth)

	result := ""

	subresult := ""

	num_processors := 0

	for _, solist := range bmach.Shared_links {
		for _, so_id := range solist {
			if so_id == so_index {
				subresult += ", p" + strconv.Itoa(num_processors) + "din"
				subresult += ", p" + strconv.Itoa(num_processors) + "dout"
				subresult += ", p" + strconv.Itoa(num_processors) + "addr"
				subresult += ", p" + strconv.Itoa(num_processors) + "wren"
				subresult += ", p" + strconv.Itoa(num_processors) + "en"
				num_processors++
			}
		}
	}

	result += "`timescale 1ns/1ps\n"
	result += "module " + sharedmem_name + "(clk, rst" + subresult + ");\n"
	result += "\n"
	result += "	//--------------Input Ports-----------------------\n"
	result += "	input clk;\n"
	result += "	input rst;\n"

	num_processors = 0

	for _, solist := range bmach.Shared_links {
		for _, so_id := range solist {
			if so_id == so_index {
				result = result + "	input [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(num_processors) + "din;\n"
				result = result + "	input [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(num_processors) + "addr;\n"
				result = result + "	input p" + strconv.Itoa(num_processors) + "wren;\n"
				result = result + "	input p" + strconv.Itoa(num_processors) + "en;\n"
				result += "	//--------------Output Ports-----------------------\n"
				result = result + "	output [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(num_processors) + "dout;\n"
				num_processors++
			}
		}
	}

	result += "\n"

	result += "	//--------------Reg-------------------------------\n"
	result += "	reg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] mem [0:" + strconv.Itoa(ram_depth-1) + "];\n"
	result += "	reg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] dout_i;\n"
	result += "	reg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] addr_i;\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] counter_wren;\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] enable_wren;\n"
	result += "	reg find_wren_tag, find_wren_tag_up;\n"
	result += "\n"

	result += "	//--------------Wire-------------------------------\n"
	result += "	wire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] din_i;\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] wren_i;\n"
	result += "	wire wren;\n"
	result += "\n"

	result += "	//--------------logic Design ------------------------\n"
	result += "	//Control the operation to select the rigth processor\n"
	for proc_id := 0; proc_id < num_processors; proc_id++ {
		result += "	assign wren_i[" + strconv.Itoa(proc_id) + "] = p" + strconv.Itoa(proc_id) + "wren;\n"
	}

	result += "	assign wren = |wren_i;\n"
	result += "                                                                                             \n"
	result += "	integer idx;                                                                                 \n"
	result += "	always @(posedge clk) begin                                                                  \n"
	result += "		if(rst) begin                                                                            \n"
	result += "        	counter_wren <= #1 'b0;                                                              \n"
	result += "        	find_wren_tag <= #1 'b0;                                                             \n"
	result += "        	find_wren_tag_up <= #1 'b0;                                                          \n"
	result += "        	enable_wren <= #1 'b0;                                                               \n"
	result += "    	end                                                                                      \n"
	result += "    	else begin                                                                               \n"
	result += "        	find_wren_tag <= #1 'b0;                                                             \n"
	result += "        	find_wren_tag_up <= #1 'b0;                                                           \n"
	result += "        	enable_wren <= #1 'b0;                                                               \n"
	result += "        	if(|wren_i==0) begin                                                                \n"
	result += "            	counter_wren = #1 counter_wren +1;                                               \n"
	result += "        	end                                                                                  \n"
	result += "        	else begin                                                                           \n"
	result += "         	for( idx = 0; idx < " + strconv.Itoa(num_processors) + "; idx = idx + 1) begin                          \n"
	result += "                	if(wren_i[idx]==1 & idx <= counter_wren & find_wren_tag==1'b0) begin        \n"
	result += "                    	enable_wren <= #1 'b0;                                                   \n"
	result += "                    	enable_wren[idx] <= #1 1'b1;                                             \n"
	result += "                    	find_wren_tag <= #1 'b1;                                                 \n"
	result += "                	end                                                                          \n"
	result += "                	if(wren_i[idx]==1 & idx > counter_wren & find_wren_tag_up==1'b0) begin      \n"
	result += "                    	enable_wren <= #1 'b0;                                                   \n"
	result += "                    	enable_wren[idx] <= #1 1'b1;                                             \n"
	result += "                    	find_wren_tag_up <= #1 'b1;                                              \n"
	result += "                	end                                                                          \n"
	result += "            	end                                                                              \n"
	result += "        	end                                                                                  \n"
	result += "    	end                                                                                      \n"
	result += "	end                                                                                          \n"
	result += "\n"
	result += "	//Caswe definition for the wr_en and address\n"
	result += "	always @(enable_wren"
	for proc_id := 0; proc_id < num_processors; proc_id++ {
		result += ", p" + strconv.Itoa(proc_id) + "addr"
	}
	result += ")\n"
	result += "		case (enable_wren)\n"
	for proc_id := 0; proc_id < num_processors; proc_id++ {
		result += "			'd"
		result += strconv.Itoa(int(math.Pow(2, float64(proc_id)))) + ": addr_i = p" + strconv.Itoa(proc_id) + "addr;\n"
	}
	result += "		endcase                                                                                     \n"

	result += "	// Memory Write Block  \n"
	result += "	// Write Operation we = 1 \n"
	result += "	always @ (posedge clk) \n"
	result += "	begin : MEM_WRITE \n"
	result += "		integer k; \n"
	result += "		if (rst)\n"
	result += "		begin \n"
	result += "			for(k=0;k<" + strconv.Itoa(ram_depth) + ";k=k+1) \n"
	result += "				mem[k] <= #1 " + strconv.Itoa(int(bmach.Rsize)) + "'b0; \n"
	result += "		end \n"
	result += "		else if (wren)\n"
	result += "			mem[addr_i] <= #1 din_i;\n"
	result += "	end \n"
	result += "\n"
	result += "	// Memory Read Block\n"
	result += "	// Read Operation when we = 0 and oe = 1 \n"
	result += "	always @ (posedge clk) \n"
	result += "	begin : MEM_READ \n"
	result += "		if (!wren)\n"
	result += "			dout_i <= #1 mem[addr_i];\n"
	result += "	end\n"
	result += "\n"

	for proc_id, solist := range bmach.Shared_links {
		for _, so_id := range solist {
			if so_id == so_index {
				result += "	assign p" + strconv.Itoa(proc_id) + "dout = dout_i;\n"
			}
		}
	}
	result += "\n"
	result += "endmodule \n"

	return result

}

func (sm Sharedmem_instance) GetPerProcPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(proc_id) + soname + "din;\n"
		result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(proc_id) + soname + "dout;\n"
		result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(proc_id) + soname + "addr;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "wren;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "en;\n"
		result += "\n"
	}
	return result
}

func (sm Sharedmem_instance) GetPerProcPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += ", p" + strconv.Itoa(proc_id) + soname + "din"
		result += ", p" + strconv.Itoa(proc_id) + soname + "dout"
		result += ", p" + strconv.Itoa(proc_id) + soname + "addr"
		result += ", p" + strconv.Itoa(proc_id) + soname + "wren"
		result += ", p" + strconv.Itoa(proc_id) + soname + "en"
	}
	return result
}

func (sm Sharedmem_instance) GetExternalPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Sharedmem_instance) GetExternalPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}
