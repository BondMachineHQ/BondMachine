package bondmachine

import (
	//"fmt"

	"strconv"
	"strings"
)

// The placeholder struct

type Vtextmem struct{}

func (op Vtextmem) Shr_get_name() string {
	return "vtextmem"
}

func (op Vtextmem) Shr_get_desc() string {
	return "Vtextmem"
}

func (op Vtextmem) Shortname() string {
	return "vtm"
}

func (op Vtextmem) GV_config(element uint8) string {
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

// Instantiate is Textual Video RAM  creation from a SO string
func (op Vtextmem) Instantiate(s string) (Shared_instance, bool) {
	if strings.HasPrefix(s, "vtextmem:") {
		components := strings.Split(s, ":")
		componentsN := len(components)
		if (componentsN-1)%5 == 0 {
			boxes := make([]GraphBox, 0)

			for i := 1; i < componentsN; i = i + 5 {
				newBox := GraphBox{}
				if newCP, err := strconv.Atoi(components[i]); err == nil {
					newBox.CP = newCP
				} else {
					return nil, false
				}
				if newLeft, err := strconv.Atoi(components[i+1]); err == nil {
					newBox.Left = newLeft
				} else {
					return nil, false
				}
				if newTop, err := strconv.Atoi(components[i+2]); err == nil {
					newBox.Top = newTop
				} else {
					return nil, false
				}
				if newWidth, err := strconv.Atoi(components[i+3]); err == nil {
					newBox.Width = newWidth
				} else {
					return nil, false
				}
				if newHeight, err := strconv.Atoi(components[i+4]); err == nil {
					newBox.Height = newHeight
				} else {
					return nil, false
				}

				boxes = append(boxes, newBox)
			}

			result := new(Vtextmem_instance)
			result.Shared_element = op
			result.Boxes = boxes
			return *result, true
		}
	}
	return nil, false
}

// The instance struct

type GraphBox struct {
	CP     int
	Left   int
	Top    int
	Width  int
	Height int
}

type Vtextmem_instance struct {
	Shared_element
	Boxes []GraphBox
}

func (sm Vtextmem_instance) String() string {
	result := "vtextmem"
	for _, box := range sm.Boxes {
		result += ":" + strconv.Itoa(box.CP) + ":" + strconv.Itoa(box.Left) + ":" + strconv.Itoa(box.Top) + ":" + strconv.Itoa(box.Width) + ":" + strconv.Itoa(box.Height)
	}
	return result
}

func (sm Vtextmem_instance) Write_verilog(bmach *Bondmachine, so_index int, vtextmem_name string, flavor string) string {

	//ram_depth := 1 << uint8(8)

	result := ""
	subresult := ""
	num_processors := 0

	for _, solist := range bmach.Shared_links {
		for _, so_id := range solist {
			if so_id == so_index {
				// From/To CP
				subresult += ", p" + strconv.Itoa(num_processors) + "din"
				subresult += ", p" + strconv.Itoa(num_processors) + "addrfromcp"
				subresult += ", p" + strconv.Itoa(num_processors) + "wren"
				subresult += ", p" + strconv.Itoa(num_processors) + "en"
				// From/To External
				subresult += ", p" + strconv.Itoa(num_processors) + "dout"
				subresult += ", p" + strconv.Itoa(num_processors) + "addrfromext"
				num_processors++
			}
		}
	}

	// Parametric Video RAM, it will be instanced once for every CP exporting to Video outputs
	result += `
module cptextvideoram #(parameter ADDR_WIDTH=8, DATA_WIDTH=8, DEPTH=256) (
    input wire clk,
    input wire [ADDR_WIDTH-1:0] addr, 
    input wire wen,
    input wire [DATA_WIDTH-1:0] i_data,
    output reg [DATA_WIDTH-1:0] o_data 
    );

    reg [DATA_WIDTH-1:0] videoram [0:DEPTH-1];

    integer i;    
    initial begin
        for (i=0 ; i<DEPTH ; i=i+1) begin
            videoram[i] = 0;  
        end
    end

    always @ (posedge clk)
    begin
        if(wen) begin
            videoram[addr] <= i_data;
        end
        else begin
            o_data <= videoram[addr];
        end     
    end
endmodule
`

	result += "`timescale 1ns/1ps\n"
	result += "module " + vtextmem_name + "(clk, rst" + subresult + ");\n"
	result += "\n"
	result += "	input clk;\n"
	result += "	input rst;\n"
	result += "\n"

	num_processors = 0

	for _, solist := range bmach.Shared_links {
		for _, so_id := range solist {
			if so_id == so_index {
				depth := sm.Boxes[num_processors].Width * sm.Boxes[num_processors].Height
				depthBits := Needed_bits(depth)
				// From/To CP
				result = result + "	input [7:0] p" + strconv.Itoa(num_processors) + "din;\n"
				result = result + "	input [" + strconv.Itoa(depthBits-1) + ":0] p" + strconv.Itoa(num_processors) + "addrfromcp;\n"
				result = result + "	input p" + strconv.Itoa(num_processors) + "wren;\n"
				result = result + "	input p" + strconv.Itoa(num_processors) + "en;\n"
				// From/To external
				result = result + "	output [7:0] p" + strconv.Itoa(num_processors) + "dout;\n"
				result = result + "	input [" + strconv.Itoa(depthBits-1) + ":0] p" + strconv.Itoa(num_processors) + "addrfromext;\n"
				result += "\n"
				num_processors++
			}
		}
	}

	num_processors = 0
	for _, solist := range bmach.Shared_links {
		for _, so_id := range solist {
			if so_id == so_index {
				p := strconv.Itoa(num_processors)
				topS := strconv.Itoa(sm.Boxes[num_processors].Top)
				leftS := strconv.Itoa(sm.Boxes[num_processors].Left)
				widthS := strconv.Itoa(sm.Boxes[num_processors].Width)
				heightS := strconv.Itoa(sm.Boxes[num_processors].Height)
				depthS := strconv.Itoa(sm.Boxes[num_processors].Width * sm.Boxes[num_processors].Height)
				depth := sm.Boxes[num_processors].Width * sm.Boxes[num_processors].Height
				depthBitsS := strconv.Itoa(Needed_bits(depth))
				result += `
	localparam P` + p + `_VIDEORAM_LEFT = ` + leftS + `;
	localparam P` + p + `_VIDEORAM_TOP = ` + topS + `;
	localparam P` + p + `_VIDEORAM_ROWS = ` + heightS + `;
	localparam P` + p + `_VIDEORAM_COLS = ` + widthS + `;
	localparam P` + p + `_VIDEORAM_DEPTH = ` + depthS + `; 
	localparam P` + p + `_VIDEORAM_A_WIDTH = ` + depthBitsS + `;
	localparam P` + p + `_VIDEORAM_D_WIDTH = 8;

	reg [P` + p + `_VIDEORAM_A_WIDTH-1:0] p` + p + `address;
	wire [P` + p + `_VIDEORAM_D_WIDTH-1:0] p` + p + `dataout;
	reg [7:0] p` + p + `data;
	reg p` + p + `wantwrite;

	cptextvideoram #(
		.ADDR_WIDTH(P` + p + `_VIDEORAM_A_WIDTH), 
		.DATA_WIDTH(P` + p + `_VIDEORAM_D_WIDTH), 
		.DEPTH(P` + p + `_VIDEORAM_DEPTH))
		p` + p + `videoram (
		.addr(p` + p + `address[P` + p + `_VIDEORAM_A_WIDTH-1:0]), 
		.o_data(p` + p + `dataout[P` + p + `_VIDEORAM_D_WIDTH-1:0]),
		.clk(clk),
		.wen(p` + p + `wantwrite),
		.i_data(p` + p + `data[P` + p + `_VIDEORAM_D_WIDTH-1:0])
	);

	always @ (posedge clk)
	begin

		if (p` + p + `wren) begin
			p` + p + `address[P` + p + `_VIDEORAM_A_WIDTH-1:0] <= p` + p + `addrfromcp[P` + p + `_VIDEORAM_A_WIDTH-1:0];
			p` + p + `wantwrite <= 1;
			p` + p + `data[P` + p + `_VIDEORAM_D_WIDTH-1:0] <= p` + p + `din[P` + p + `_VIDEORAM_D_WIDTH-1:0];
		end
		else
		begin
			p` + p + `address[P` + p + `_VIDEORAM_A_WIDTH-1:0] <= p` + p + `addrfromext[P` + p + `_VIDEORAM_A_WIDTH-1:0];
			p` + p + `wantwrite <= 0;
		end
	end

	assign p` + p + `dout[P` + p + `_VIDEORAM_D_WIDTH-1:0] = p` + p + `dataout[P` + p + `_VIDEORAM_D_WIDTH-1:0];

`
				num_processors++
			}
		}
	}

	result += "\n"
	result += "endmodule \n"

	return result

}

func (sm Vtextmem_instance) GetPerProcPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += ", p" + strconv.Itoa(proc_id) + soname + "din"
		result += ", p" + strconv.Itoa(proc_id) + soname + "addrfromcp"
		result += ", p" + strconv.Itoa(proc_id) + soname + "wren"
		result += ", p" + strconv.Itoa(proc_id) + soname + "en"
	}
	return result
}

func (sm Vtextmem_instance) GetPerProcPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(proc_id) + soname + "din;\n"
		result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(proc_id) + soname + "addrfromcp;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "wren;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "en;\n"
		result += "\n"
	}
	return result
}

func (sm Vtextmem_instance) GetExternalPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""

	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += ", p" + strconv.Itoa(proc_id) + soname + "dout"
		result += ", p" + strconv.Itoa(proc_id) + soname + "addrfromext"
	}

	return result
}
func (sm Vtextmem_instance) GetExternalPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""

	if soname, ok := bmach.Get_so_name(so_id); ok {
		depth := sm.Boxes[proc_id].Width * sm.Boxes[proc_id].Height
		depthBits := Needed_bits(depth)
		//result = result + "	input [7:0] p" + strconv.Itoa(proc_id) + soname + "din;\n"
		result = result + "	output [7:0] p" + strconv.Itoa(proc_id) + soname + "dout;\n"
		//result = result + "	input [" + strconv.Itoa(depthBits-1) + ":0] p" + strconv.Itoa(proc_id) + soname + "addrfromcp;\n"
		result = result + "	input [" + strconv.Itoa(depthBits-1) + ":0] p" + strconv.Itoa(proc_id) + soname + "addrfromext;\n"
		//result = result + "	input p" + strconv.Itoa(proc_id) + soname + "wren;\n"
		//result = result + "	input p" + strconv.Itoa(proc_id) + soname + "en;\n"
	}

	return result
}
