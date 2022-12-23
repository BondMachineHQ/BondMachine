package bondmachine

import (
	"strconv"
	"strings"
)

// The placeholder struct

type Channel struct{}

func (op Channel) Shr_get_name() string {
	return "channel"
}

func (op Channel) Shr_get_desc() string {
	return "Channel"
}

func (op Channel) Shortname() string {
	return "ch"
}

func (op Channel) GV_config(element uint8) string {
	result := ""
	switch element {
	case GVNODEINPROC:
		result += "style=filled fillcolor=yellow color=black"
	case GVNODE:
		result += "style=filled fillcolor=yellow color=black"
	case GVEDGE:
		result += "arrowhead=none"
	case GVCLUS:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey60"
	case GVCLUSINPROC:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey60"
	}
	return result
}

func (op Channel) Instantiate(s string) (Shared_instance, bool) {
	if strings.HasPrefix(s, "channel:") {
		result := new(Channel_instance)
		result.Shared_element = op
		return *result, true
	}
	return nil, false
}

// The instance struct

type Channel_instance struct {
	Shared_element
}

func (sm Channel_instance) String() string {
	return "channel:"
}

func (sm Channel_instance) Write_verilog(bmach *Bondmachine, so_index int, channel_name string, flavor string) string {

	result := ""

	subresult := ""

	num_processors := 0

	for _, solist := range bmach.Shared_links {
		for _, so_id := range solist {
			if so_id == so_index {
				subresult += ", p" + strconv.Itoa(num_processors) + "chin"
				subresult += ", p" + strconv.Itoa(num_processors) + "w2w"
				subresult += ", p" + strconv.Itoa(num_processors) + "w2r"
				subresult += ", p" + strconv.Itoa(num_processors) + "ack_ch_ready"
				subresult += ", p" + strconv.Itoa(num_processors) + "op_check_ready"
				subresult += ", p" + strconv.Itoa(num_processors) + "finish_channel"
				subresult += ", p" + strconv.Itoa(num_processors) + "chout"
				subresult += ", p" + strconv.Itoa(num_processors) + "ack_w2w"
				subresult += ", p" + strconv.Itoa(num_processors) + "ack_w2r"
				subresult += ", p" + strconv.Itoa(num_processors) + "ch_ready"
				subresult += ", p" + strconv.Itoa(num_processors) + "ch_w_r_ready"
				num_processors++
			}
		}
	}
	//result += "module Channel(clk, reset, ch2proc_0 (pchin), proc2ch_0(pchout), ch2proc_1, proc2ch_1, wr_strobe, ack_w2w, rd_strobe, ack_w2r, ack_ch, ready);\n"
	result += "\n"
	result += "`timescale 1ns/1ps\n"
	result += "module " + channel_name + "(clk, reset" + subresult + ");\n"
	result += "\n"
	result += "	//--------------Input Ports-----------------------\n"
	result += "	input clk;\n"
	result += "	input reset;\n"

	subresult_in := ""
	subresult_out := ""
	num_processors = 0

	for _, solist := range bmach.Shared_links {
		for _, so_id := range solist {
			if so_id == so_index {
				subresult_in += "	input [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(num_processors) + "chin;\n"
				subresult_in += "	input p" + strconv.Itoa(num_processors) + "w2w;\n"
				subresult_in += "	input p" + strconv.Itoa(num_processors) + "w2r;\n"
				subresult_in += "	input p" + strconv.Itoa(num_processors) + "ack_ch_ready;\n"
				subresult_in += "	input p" + strconv.Itoa(num_processors) + "op_check_ready;\n"
				subresult_out += "	output p" + strconv.Itoa(num_processors) + "finish_channel;\n"
				subresult_out += "	output [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(num_processors) + "chout;\n"
				subresult_out += "	output p" + strconv.Itoa(num_processors) + "ack_w2w;\n"
				subresult_out += "	output p" + strconv.Itoa(num_processors) + "ack_w2r;\n"
				subresult_out += "	output p" + strconv.Itoa(num_processors) + "ch_ready;\n"
				subresult_out += "	output [1:0] p" + strconv.Itoa(num_processors) + "ch_w_r_ready;\n"
				num_processors++
			}
		}
	}

	result += subresult_in
	result += "\n"
	result += "	//--------------Output Ports-----------------------\n"
	result += subresult_out
	result += "\n"

	result += "	//--------------Generic Parameter-------------------\n"
	result += "	parameter RAM_CH_DEPTH = 1 << " + strconv.Itoa(int(bmach.Rsize)) + ";\n"

	result += "	//--------------Reg declaration---------------------------------------------\n"
	result += "	reg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] ch2proc_i [0:" + strconv.Itoa((num_processors)-1) + "];\n"
	result += "\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] p_w2w_i_d1;\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] p_w2r_i_d1;                                    \n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] en_wrd_storbe;\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] en_wwr_storbe;\n"
	result += "	reg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] wr_pointer_w2w; 	//write pointer for the W2W request\n"
	result += "	reg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] wr_pointer_w2r; 	//write pointer for the W2W request\n"
	result += "	reg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] status_w2w_pointer [" + strconv.Itoa((num_processors)-1) + ":0];\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] status_w2w_reg;\n"
	result += "	reg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] status_w2r_pointer [" + strconv.Itoa((num_processors)-1) + ":0];\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] status_w2r_reg;\n"

	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] wwr_tag_tmp;\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] wrd_tag_tmp;\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] count_wwr;\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] count_wrd;\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] w2w_proc;\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] w2r_proc;\n"

	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] wwr_ready_ch;\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] wrd_ready_ch;\n"

	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] find_wwr;\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] find_wrd;\n"
	result += "\treg wwr_finish, wrd_finish;\n"
	result += "\twire wwr_finish_pulse, wrd_finish_pulse;\n"
	result += "\treg wwr_finish_d1, wrd_finish_d1;\n"

	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] mem_tag_w2w [0:RAM_CH_DEPTH-1];\n"
	result += "	reg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] rd_pointer_w2w;\n"

	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] mem_tag_w2r [0:RAM_CH_DEPTH-1];\n"
	result += "	reg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] rd_pointer_w2r;\n"

	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] finish_channel_wwr;\n"
	result += "	reg [" + strconv.Itoa((num_processors)-1) + ":0] finish_channel_wrd;\n"
	result += "\n"

	result += "	//--------------Wire declaration--------------------------------------------\n"
	result += "	wire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] proc2ch_i [0:" + strconv.Itoa(num_processors) + "-1];\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] p_w2w_i;\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] p_w2r_i;\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] ch_ready;\n"
	result += "	wire [1:0] ch_w_r_ready [" + strconv.Itoa((num_processors)-1) + ":0];\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] ack_w2w_i;\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] ack_w2r_i;\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] ack_ch_ready_i;\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] op_check_ready_i;\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] finish_channel_i;\n"
	result += "\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] reset_w2w_storbe;\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] reset_w2r_storbe;\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] search_wwr;\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] search_w2r;\n"
	result += "	wire dv_w2w;\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] tag_w2w;\n"
	result += "	wire dv_w2r;\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] tag_w2r;\n"
	result += "	wire valid_w2w_pulse;\n"
	result += "	wire valid_w2r_pulse;\n"
	result += "	wire w2r_reg_check;\n"
	result += "	wire w2r_pointer_check;\n"
	result += "	wire w2w_reg_check;\n"
	result += "	wire w2w_pointer_check;\n"

	result += "	//--------------Define the tag parameters of the processor----------------------------\n"
	result += "	wire [" + strconv.Itoa((num_processors)-1) + ":0] TAG_CH [" + strconv.Itoa((num_processors)-1) + ":0];\n"
	for proc_id := 0; proc_id < num_processors; proc_id++ {
		result += "\tlocalparam TAG_CH_" + strconv.Itoa(proc_id) + " = 'b" + strconv.Itoa(proc_id) + ";\n"
		result += "\tassign TAG_CH[" + strconv.Itoa(proc_id) + "] = TAG_CH_" + strconv.Itoa(proc_id) + ";\n"
	}

	result += "	//--------------Internal input signal assignment----------------------------\n"
	for proc_id := 0; proc_id < num_processors; proc_id++ {
		result += "	assign proc2ch_i[" + strconv.Itoa(proc_id) + "] = p" + strconv.Itoa(proc_id) + "chin;\n"
		result += "	assign p_w2w_i[" + strconv.Itoa(proc_id) + "] = p" + strconv.Itoa(proc_id) + "w2w;\n" //wr_strobe
		result += "	assign p_w2r_i[" + strconv.Itoa(proc_id) + "] = p" + strconv.Itoa(proc_id) + "w2r;\n" //rd_strobe
		result += "	assign ack_ch_ready_i[" + strconv.Itoa(proc_id) + "] = p" + strconv.Itoa(proc_id) + "ack_ch_ready;\n"
		result += "	assign op_check_ready_i[" + strconv.Itoa(proc_id) + "] = p" + strconv.Itoa(proc_id) + "op_check_ready;\n"

	}
	result += "\n"

	result += "	//--------------Internal output signal assignment----------------------------------\n"
	for proc_id := 0; proc_id < num_processors; proc_id++ {
		result += "	assign p" + strconv.Itoa(proc_id) + "finish_channel = finish_channel_i[" + strconv.Itoa(proc_id) + "];\n"
		result += "	assign p" + strconv.Itoa(proc_id) + "ch_ready = ch_ready[" + strconv.Itoa(proc_id) + "];\n"
		result += "	assign p" + strconv.Itoa(proc_id) + "ch_w_r_ready = ch_w_r_ready[" + strconv.Itoa(proc_id) + "];\n"
		result += "	assign p" + strconv.Itoa(proc_id) + "ack_w2w = ack_w2w_i[" + strconv.Itoa(proc_id) + "];\n"
		result += "	assign p" + strconv.Itoa(proc_id) + "ack_w2r = ack_w2r_i[" + strconv.Itoa(proc_id) + "];\n"
		result += "	assign p" + strconv.Itoa(proc_id) + "chout = ch2proc_i[" + strconv.Itoa(proc_id) + "];\n"
	}

	result += "\n"
	result += "	//---------------------Delay Register---------------------------------------------\n"
	result += "	always@(posedge clk)\n"
	result += "	begin\n"
	result += "		if(reset) begin\n"
	result += "			p_w2w_i_d1 <= #1 'b0;\n"
	result += "        	p_w2r_i_d1 <= #1 'b0;\n"
	result += "    	end\n"
	result += "    	else begin\n"
	result += "			p_w2w_i_d1 <= #1 p_w2w_i;\n"
	result += "        	p_w2r_i_d1 <= #1 p_w2r_i;\n"
	result += "     end\n"
	result += "	end\n"
	result += "\n"

	result += "	//-------------- Logic design ---------------------------------------------\n"
	result += "	//Define the TAG ID for each processor attached to the channel\n"
	result += "	/*genvar i_tag;\n"
	result += "	generate\n"
	result += "		for (i_tag=0; i_tag < " + strconv.Itoa(num_processors) + "; i_tag=i_tag+1) begin\n"
	result += "			always @(posedge clk) begin\n"
	result += "				if(reset)\n"
	result += "					TAG_CH[i_tag] <= #1 'b0;\n"
	result += "				else\n"
	result += "					TAG_CH[i_tag] <= #1 i_tag;\n"
	result += "			end\n"
	result += "		end\n"
	result += "	endgenerate*/\n"
	result += "\n"

	result += "	//insert the satus register for each processor the Write Operation\n"
	result += "	//the register is enabled when the write_strobe is asserted\n"
	result += "	//the register is reset when the write_storge is asserted\n"
	result += "	//the write pointer is inserted in the status register\n"
	result += "	genvar i_status;\n"
	result += "	generate\n"
	result += "		for (i_status=0; i_status < " + strconv.Itoa(num_processors) + "; i_status = i_status + 1) begin\n"
	result += "			always @(posedge clk or posedge reset) begin\n"
	result += "				if(reset) begin\n"
	result += "					status_w2w_pointer[i_status] <= #1 'b0;\n"
	result += "					status_w2w_reg[i_status] <= #1 1'b0;\n"
	result += "				end\n"
	result += "				else begin\n"
	result += "					if(reset_w2w_storbe[i_status]) begin\n"
	result += "						status_w2w_pointer[i_status] <= #1 'b0;\n"
	result += "						status_w2w_reg[i_status] <= #1 1'b0;\n"
	result += "					end\n"
	result += "					else if(en_wwr_storbe[i_status]) begin\n"
	result += "						status_w2w_pointer[i_status] <= #1 wr_pointer_w2w;\n"
	result += "						status_w2w_reg[i_status] <= #1 1'b1;\n"
	result += "					end\n"
	result += "				end\n"
	result += "			end\n"
	result += "		end\n"
	result += "	endgenerate\n"
	result += "\n"
	result += "	assign reset_w2w_storbe = (~p_w2w_i) & p_w2w_i_d1;\n"
	result += "\n"

	result += "	//insert the satus register for each processor the Read Operation\n"
	result += "	//the register is enabled when the read_strobe is asserted\n"
	result += "	//the register is reset when the read_strobe is asserted\n"
	result += "	//the write pointer for the read is inserted in the status register\n"
	result += "	genvar i_status_w2r;\n"
	result += "	generate\n"
	result += "		for (i_status_w2r=0; i_status_w2r < " + strconv.Itoa(num_processors) + "; i_status_w2r=i_status_w2r+1) begin\n"
	result += "			always @(posedge clk or posedge reset) begin\n"
	result += "				if(reset) begin\n"
	result += "					status_w2r_pointer[i_status_w2r] <= #1 'b0;\n"
	result += "					status_w2r_reg[i_status_w2r] <= #1 1'b0;\n"
	result += "				end\n"
	result += "				else begin\n"
	result += "					if(reset_w2r_storbe[i_status_w2r]) begin\n"
	result += "						status_w2r_pointer[i_status_w2r] <= #1 'b0;\n"
	result += "						status_w2r_reg[i_status_w2r] <= #1 1'b0;\n"
	result += "					end\n"
	result += "					else if(en_wrd_storbe[i_status_w2r]) begin\n"
	result += "						status_w2r_pointer[i_status_w2r] <= #1 wr_pointer_w2r;\n"
	result += "						status_w2r_reg[i_status_w2r] <= #1 1'b1;\n"
	result += "					end\n"
	result += "				end\n"
	result += "			end\n"
	result += "		end\n"
	result += "	endgenerate\n"
	result += "\n"
	result += "	assign reset_w2r_storbe = (~p_w2r_i) & p_w2r_i_d1;\n"
	result += "\n"

	result += "	//define the counter to randmly select the tag                                                   \n"
	result += "	//the down and up signal define the searching range to                                           \n"
	result += "	//select the tag                                                                                 \n"
	result += "	assign search_wwr =  p_w2w_i & (p_w2w_i ^ status_w2w_reg);                               \n"
	result += "\n"
	result += "\tinteger idy;\n"
	result += "\talways @(posedge clk or posedge reset) begin\n"
	result += "\t	if(reset) begin\n"
	result += "\t		wwr_tag_tmp <= #1 'b0;\n"
	result += "\t      	find_wwr <= #1 'b0;\n"
	result += "\t      	en_wwr_storbe <= #1 'b0;\n"
	result += "\t	end\n"
	result += "\t   else begin\n"
	result += "\t   	en_wwr_storbe <= #1 'b0;\n"
	result += "\t       for( idy = 0; idy < " + strconv.Itoa(num_processors) + "; idy = idy + 1) begin\n"
	result += "\t           if(reset_w2r_storbe[idy]) begin\n"
	result += "\t               find_wwr[idy] <= #1 1'b0;\n"
	result += "\t           end\n"
	result += "\t           else if(search_wwr[idy]==1 & find_wwr[idy]==1'b0) begin\n"
	result += "\t                wwr_tag_tmp <= #1 TAG_CH[idy];\n"
	result += "\t                en_wwr_storbe[idy] <= #1 1'b1;\n"
	result += "\t                find_wwr[idy] <= #1 'b1;\n"
	result += "\t           end\n"
	result += "\t       end\n"
	result += "\t   end\n"
	result += "\tend\n"
	/*result += "\tinteger idy;\n"
	result += "\talways @(posedge clk) begin\n"
	result += "\t	if(reset) begin\n"
	result += "\t    	wwr_tag_tmp <= #1 'b0;\n"
	result += "\t       count_wwr <= #1 'b0;\n"
	result += "\t       find_wwr <= #1 'b0;\n"
	result += "\t       en_wwr_storbe <= #1 'b0;\n"
	result += "\t	end\n"
	result += "\t   else begin\n"
	result += "\t   	find_wwr <= #1 'b0;\n"
	result += "\t       en_wwr_storbe <= #1 'b0;\n"
	result += "\t       if(|p_w2w_i==0) begin\n"
	result += "\t        	count_wwr = #1 count_wwr +1;\n"
	result += "\t       end\n"
	result += "\t       else begin\n"
	result += "\t       	for( idy = 0; idy < " + strconv.Itoa(num_processors) + "; idy = idy + 1) begin\n"
	result += "\t            	if(search_wwr[idy]==1 & idy >= count_wwr & find_wwr==1'b0) begin\n"
	result += "\t           	   	wwr_tag_tmp <= #1 TAG_CH[idy];\n"
	//result += "\t                  	en_wwr_storbe <= #1 'b0;\n"
	result += "\t                  	en_wwr_storbe[idy] <= #1 1'b1;\n"
	result += "\t                  	find_wwr <= #1 'b1;\n"
	result += "\t               end\n"
	result += "\t               if(search_wwr[idy]==1 & idy < count_wwr & find_wwr==1'b0) begin\n"
	result += "\t                	wwr_tag_tmp <= #1 TAG_CH[idy];\n"
	result += "\t                  	find_wwr <= #1 'b1;\n"
	//result += "\t                  	en_wwr_storbe <= #1 'b0;\n"
	result += "\t                  	en_wwr_storbe[idy] <= #1 1'b1;\n"
	result += "\t               end\n"
	result += "\t           end\n"
	result += "\t       end\n"
	result += "\t   end\n"
	result += "\tend\n"
	result += "\n"*/

	result += "	//signal valid to write in the memory                                                            \n"
	result += "	assign valid_w2w_pulse = |en_wwr_storbe;\n"
	result += "	assign ack_w2w_i = en_wwr_storbe;\n"
	result += "\n"

	result += "	//define the counter to randmly select the tag                                                   \n"
	result += "	//the down and up signal define the searching range to                                           \n"
	result += "	//select the tag                                                                                 \n"
	result += "	assign search_w2r =  p_w2r_i & (p_w2r_i ^ status_w2r_reg);                               \n"
	result += "                                                                                                 \n"
	result += "\tinteger idx;\n"
	result += "\talways @(posedge clk or posedge reset) begin\n"
	result += "\t	if(reset) begin\n"
	result += "\t		wrd_tag_tmp <= #1 'b0;\n"
	result += "\t      	find_wrd <= #1 'b0;\n"
	result += "\t      	en_wrd_storbe <= #1 'b0;\n"
	result += "\t	end\n"
	result += "\t   else begin\n"
	result += "\t   	en_wrd_storbe <= #1 'b0;\n"
	result += "\t       for( idx = 0; idx < " + strconv.Itoa(num_processors) + "; idx = idx + 1) begin\n"
	result += "\t           if(reset_w2r_storbe[idx]) begin\n"
	result += "\t               find_wrd[idx] <= #1 1'b0;\n"
	result += "\t           end\n"
	result += "\t           else if(search_w2r[idx]==1 & find_wrd[idx]==1'b0) begin\n"
	result += "\t                wrd_tag_tmp <= #1 TAG_CH[idx];\n"
	result += "\t                en_wrd_storbe[idx] <= #1 1'b1;\n"
	result += "\t                find_wrd[idx] <= #1 'b1;\n"
	result += "\t           end\n"
	result += "\t       end\n"
	result += "\t   end\n"
	result += "\tend\n"

	/*result += "\talways @(posedge clk) begin\n"
	result += "\t	if(reset) begin\n"
	result += "\t    	wrd_tag_tmp <= #1 'b0;\n"
	result += "\t       count_wrd <= #1 'b0;\n"
	result += "\t       find_wrd <= #1 'b0;\n"
	result += "\t       en_wrd_storbe <= #1 'b0;\n"
	result += "\t   end\n"
	result += "\t   else begin\n"
	result += "\t   	find_wrd <= #1 'b0;\n"
	result += "\t       en_wrd_storbe <= #1 'b0;\n"
	result += "\t       if(|p_w2r_i==0) begin\n"
	result += "\t        	count_wrd = #1 count_wrd +1;\n"
	result += "\t       end\n"
	result += "\t       else begin\n"
	result += "\t       	for( idx = 0; idx < " + strconv.Itoa(num_processors) + "; idx = idx + 1) begin\n"
	result += "\t              	if(search_w2r[idx]==1 & idx >= count_wrd & find_wrd==1'b0) begin\n"
	result += "\t                  	wrd_tag_tmp <= #1 TAG_CH[idx];\n"
	//result += "\t                  	en_wrd_storbe <= #1 'b0;\n"
	result += "\t                  	en_wrd_storbe[idx] <= #1 1'b1;\n"
	result += "\t                  	find_wrd <= #1 'b1;\n"
	result += "\t              	end\n"
	result += "\t              	if(search_w2r[idx]==1 & idx < count_wrd & find_wrd==1'b0) begin \n"
	result += "\t                  	wrd_tag_tmp <= #1 TAG_CH[idx];\n"
	result += "\t                  	find_wrd <= #1 'b1;\n"
	//result += "\t                  	en_wrd_storbe <= #1 'b0;\n"
	result += "\t                  	en_wrd_storbe[idx] <= #1 1'b1;\n"
	result += "\t              	end\n"
	result += "\t           end\n"
	result += "\t       end\n"
	result += "\t   end\n"
	result += "\tend\n"
	result += "\n"*/

	result += "	//signal valid to write in the memory                                                            \n"
	result += "	assign valid_w2r_pulse = |en_wrd_storbe;   //change as function of processor     \n"
	result += "	assign ack_w2r_i = en_wrd_storbe;\n"
	result += "\n"

	result += "	// Memory Circular Block to store the W2W                                                        \n"
	result += "	//this memeory store the tag of the incoming processor                                           \n"
	result += "	//for the reading oepration                                                                      \n"
	result += "	always @ (posedge clk)                                                                           \n"
	result += "	begin : TAG_MEM                                                                                  \n"
	result += " 	integer k;                                                                                     \n"
	result += "    	if (reset) begin                                                                             \n"
	result += "     	for(k=0;k<RAM_CH_DEPTH;k=k+1)                                                                \n"
	result += "     		mem_tag_w2w[k] <= #1 'b0;                                                                \n"
	result += "    	end                                                                                          \n"
	result += "    	else if (valid_w2w_pulse) begin                                                              \n"
	result += "     	mem_tag_w2w[wr_pointer_w2w] <= #1  wwr_tag_tmp;                                          \n"
	result += "    	end                                                                                          \n"
	result += "	end                                                                                              \n"
	result += "                                                                                                 \n"
	result += "	always @ (posedge clk)                                                                          \n"
	result += "	begin                                                                                            \n"
	result += " 	if(reset)                                                                                    \n"
	result += "     	wr_pointer_w2w <= #1 'b0;                                                                \n"
	result += "    	else if (valid_w2w_pulse)                                                                    \n"
	result += "     	wr_pointer_w2w <= #1 wr_pointer_w2w + 1;                                                 \n"
	result += "	end                                                                                              \n"
	result += "                                                                                                 \n"
	result += "	assign w2w_reg_check = status_w2w_reg[mem_tag_w2w[rd_pointer_w2w]]==1'b1 ? 1'b1 : 1'b0;            \n"
	result += "	assign w2w_pointer_check = status_w2w_pointer[mem_tag_w2w[rd_pointer_w2w]]==rd_pointer_w2w ? 1'b1 : 1'b0; 				\n"
	result += "\n"

	result += "\treg of_wwr;\n"
	result += "\talways @(rd_pointer_w2w, wr_pointer_w2w)\n"
	result += "\tbegin\n"
	result += "\t	if(rd_pointer_w2w <= wr_pointer_w2w)\n"
	result += "\t   	of_wwr <= 1'b0;\n"
	result += "\t   else if (rd_pointer_w2w > wr_pointer_w2w)\n"
	result += "\t   	of_wwr <= 1'b1;\n"
	result += "\tend\n"

	result += "\talways @ (posedge clk )\n"
	result += "\tbegin\n"
	result += "\t	if(reset)\n"
	result += "\t    	rd_pointer_w2w <= #1 'b0;\n"
	result += "\t	else begin\n"
	result += "\t		if((rd_pointer_w2w < wr_pointer_w2w) | of_wwr) begin\n"
	result += "\t			if(~w2w_reg_check)\n"
	result += "\t				rd_pointer_w2w <= #1 rd_pointer_w2w + 1;\n"
	result += "\t			else if(wwr_finish_pulse)\n"
	result += "\t				rd_pointer_w2w <= #1 rd_pointer_w2w + 1;\n"
	result += "\t		end\n"
	result += "\t	end\n"
	result += "\tend\n"
	result += "\n"

	result += "	// Memory Circular Block to store the W2R                                                                            \n"
	result += "	//this memeory store the tag of the incoming processor                                                               \n"
	result += "	//for the reading oepration                                                                                          \n"
	result += "	always @ (posedge clk)                                                                                               \n"
	result += "	begin : TAG_MEM_W2R                                                                                                  \n"
	result += " 	integer k;                                                                                                         \n"
	result += "    	if (reset) begin                                                                                                 \n"
	result += "     	for(k=0;k<RAM_CH_DEPTH;k=k+1)                                                                                    \n"
	result += "        		mem_tag_w2r[k] <= #1 8'b0;                                                                                   \n"
	result += "   	end                                                                                                              \n"
	result += "    	else if (valid_w2r_pulse) begin                                                                                  \n"
	result += "     	mem_tag_w2r[wr_pointer_w2r] <= #1  wrd_tag_tmp;                                                              \n"
	result += "    	end                                                                                                              \n"
	result += "	end                                                                                                                  \n"
	result += "                                                                                                                     \n"
	result += "	assign w2r_reg_check = status_w2r_reg[mem_tag_w2r[rd_pointer_w2r]]==1'b1 ? 1'b1 : 1'b0;                                \n"
	result += "	assign w2r_pointer_check = status_w2r_pointer[mem_tag_w2r[rd_pointer_w2r]]==rd_pointer_w2r ? 1'b1 : 1'b0;              \n"
	result += "\n"

	result += "\treg of_wrd;\n"
	result += "\talways @(rd_pointer_w2r, wr_pointer_w2r)\n"
	result += "\tbegin\n"
	result += "\t	if(rd_pointer_w2r <= wr_pointer_w2r)\n"
	result += "\t   	of_wrd <= 1'b0;\n"
	result += "\t   else if (rd_pointer_w2r > wr_pointer_w2r)\n"
	result += "\t   	of_wrd <= 1'b1;\n"
	result += "\tend\n"

	result += "\talways @ (posedge clk )                                                                                              \n"
	result += "\tbegin                                                                                                                \n"
	result += "\t	if(reset)                                                                                                        \n"
	result += "\t   	wr_pointer_w2r <= #1 'b0;                                                                                    \n"
	result += "\t   else if (valid_w2r_pulse)                                                                                        \n"
	result += "\t   	wr_pointer_w2r <= #1 wr_pointer_w2r + 1;                                                                     \n"
	result += "\tend                                                                                                                  \n"
	result += "\t                                                                                                                    \n"
	result += "\talways @ (posedge clk )                                                                                              \n"
	result += "\tbegin                                                                                                                \n"
	result += "\t	if(reset)                                                                                                        \n"
	result += "\t		rd_pointer_w2r <= #1 'b0;                                                                                    \n"
	result += "\t   else begin\n"
	result += "\t		if((rd_pointer_w2r < wr_pointer_w2r) | of_wrd) begin\n"
	result += "\t			if(~w2r_reg_check)\n"
	result += "\t				rd_pointer_w2r <= #1 rd_pointer_w2r + 1;\n"
	result += "\t			else if(wrd_finish_pulse)\n"
	result += "\t				rd_pointer_w2r <= #1 rd_pointer_w2r + 1;\n"
	result += "\t		end\n"
	result += "\t	end\n"
	result += "\tend\n"
	result += "\n"

	result += "	//data valid for generate the ack connection                                                                         \n"
	result += "	assign tag_w2w = mem_tag_w2w[rd_pointer_w2w];\n"
	result += "	assign tag_w2r = mem_tag_w2r[rd_pointer_w2r];\n"
	result += "	assign dv_w2w = w2w_reg_check & w2w_pointer_check; //status_tag_mem_w2w;//status_w2w_reg[mem_tag[rd_pointer_w2w]];     \n"
	result += "	assign dv_w2r = w2r_reg_check & w2r_pointer_check; 																	\n"

	result += "\t//logic to assert the ACK operation\n"
	result += "\tassign ch_ready = (wwr_ready_ch | wrd_ready_ch) & {" + strconv.Itoa(num_processors) + "{|wwr_ready_ch}} & {" + strconv.Itoa(num_processors) + "{|wrd_ready_ch}};\n"
	for proc_id := 0; proc_id < num_processors; proc_id++ {
		result += "\tassign ch_w_r_ready[" + strconv.Itoa(proc_id) + "] = {wrd_ready_ch[" + strconv.Itoa(proc_id) + "], wwr_ready_ch[" + strconv.Itoa(proc_id) + "]};\n"
	}
	result += "\n"
	result += "\talways @(posedge clk or posedge reset)\n"
	result += "\tbegin\n"
	result += "\t	if(reset) begin\n"
	result += "\t		wwr_ready_ch <= #1 'b0;\n"
	result += "\t		wrd_ready_ch <= #1 'b0;\n"
	result += "\t	end\n"
	result += "\t	else begin\n"
	result += "\t		if(dv_w2w & dv_w2r & op_check_ready_i[tag_w2w]  & tag_w2w!=tag_w2r)\n"
	result += "\t			wwr_ready_ch[tag_w2w] <= #1 1'b1;\n"
	result += "\t		else\n"
	result += "\t			wwr_ready_ch[tag_w2w] <= #1 1'b0;\n"
	result += "\t		if(dv_w2r & dv_w2w & op_check_ready_i[tag_w2r]  & tag_w2w!=tag_w2r)\n"
	result += "\t			wrd_ready_ch[tag_w2r] <= #1 1'b1;\n"
	result += "\t		else\n"
	result += "\t			wrd_ready_ch[tag_w2r] <= #1 1'b0;\n"
	result += "\t	end\n"
	result += "\tend\n"

	result += "//define the ready logic to pass the data\n"
	result += "\talways @(posedge clk or posedge reset)\n"
	result += "\tbegin\n"
	result += "\t	if(reset) begin\n"
	result += "\t		wwr_finish <= #1 1'b0;\n"
	result += "\t		finish_channel_wwr <= #1 'b0;\n"
	result += "\t		wrd_finish <= #1 1'b0;\n"
	result += "\t		finish_channel_wrd[tag_w2r] <= #1 'b0;\n"
	result += "\t	end\n"
	result += "\t	else begin\n"
	result += "\t		if(ack_ch_ready_i[tag_w2w]==1'b1) begin\n"
	result += "\t			wwr_finish <= #1 1'b1;\n"
	result += "\t			finish_channel_wwr[tag_w2w] <= #1 1'b1;\n"
	result += "\t		end\n"
	result += "\t		else begin\n"
	result += "\t			wwr_finish <= #1 1'b0;\n"
	result += "\t			finish_channel_wwr <= #1 'b0;\n"
	result += "\t		end\n"
	result += "\t		if(ack_ch_ready_i[tag_w2r]==1'b1) begin\n"
	result += "\t			wrd_finish <= #1 1'b1;\n"
	result += "\t			finish_channel_wrd[tag_w2r] <= #1 1'b1;\n"
	result += "\t		end\n"
	result += "\t		else begin\n"
	result += "\t			wrd_finish <= #1 1'b0;\n"
	result += "\t			finish_channel_wrd <= #1 'b0;\n"
	result += "\t		end\n"
	result += "\t	end\n"
	result += "\tend\n"
	result += "\t\n"
	//result += "\talways @(tag_w2r or tag_w2w)\n"
	//result += "\tbegin\n"
	//result += "\t	ch2proc_i[tag_w2r] <= proc2ch_i[tag_w2w];\n"
	//result += "\tend\n"

	result += "\t\n"
	result += "\tassign finish_channel_i = (finish_channel_wrd | finish_channel_wwr) & {" + strconv.Itoa(num_processors) + "{wrd_finish}}  & {" + strconv.Itoa(num_processors) + "{wwr_finish}};\n"
	result += "\n"
	result += "\t//process finish\n"
	result += "\talways @ (posedge clk)\n"
	result += "\tbegin\n"
	result += "\t	if(reset) begin\n"
	result += "\t		wrd_finish_d1 <= #1 'b0;\n"
	result += "\t		wwr_finish_d1 <= #1 'b0;\n"
	result += "\t	end\n"
	result += "\t	else begin\n"
	result += "\t    	wrd_finish_d1 <= #1 wrd_finish;\n"
	result += "\t		wwr_finish_d1 <= #1 wwr_finish;\n"
	result += "\t	end\n"
	result += "\tend\n"
	result += "\n"
	result += "\tassign wrd_finish_pulse = ~wrd_finish & wrd_finish_d1 & finish_channel_i;\n"
	result += "\tassign wwr_finish_pulse = ~wwr_finish & wwr_finish_d1 & finish_channel_i;\n"
	result += "\n"

	result += "\tinteger i_ch2proc;\n"
	result += "\talways @ (*) begin//(posedge clk or posedge reset) begin\n"
	result += "\t	for (i_ch2proc=0; i_ch2proc < 2; i_ch2proc=i_ch2proc+1) begin\n"
	result += "\t		ch2proc_i[i_ch2proc] <= 'b0;\n"
	result += "\t	if(tag_w2r==i_ch2proc)\n"
	result += "\t		ch2proc_i[i_ch2proc] <= proc2ch_i[tag_w2w];\n"
	result += "\t	end\n"
	result += "\tend\n"

	result += "\n"
	result += "	// synthesis translate_off\n"
	result += "	//define assertion to check the enaable_wr_strobe assertion                                                                                                                  \n"
	result += "	generate                                                                                                                                                                     \n"
	result += "		for (genvar j = 0; j < " + strconv.Itoa(num_processors) + "; j++)                                                                                                                               \n"
	result += "			begin : assert_array                                                                                                                                                    \n"
	result += "			//controllo che dopo il read strobe si abilita l'enable                                                                                                                 \n"
	result += "			rd_strobe_enable:assert property (@(posedge clk) disable iff(reset) $rose(en_wrd_storbe[j]) |-> ##1 status_w2r_reg[j]);                                              \n"
	result += "			wr_strobe_enable:assert property (@(posedge clk) disable iff(reset) $rose(en_wwr_storbe[j]) |-> ##1 status_w2w_reg[j]);                                              \n"
	result += "			//controllo che dopo il read strobe basso si abilita il reset                                                                                                           \n"
	result += "			reset_wr_strobe:assert property (@(posedge clk) disable iff(reset) $fell(p_w2w_i[j]) |-> ##0 reset_w2w_storbe[j]);                                                     \n"
	result += "			reset_rd_strobe:assert property (@(posedge clk) disable iff(reset) $fell(p_w2r_i[j]) |-> ##0 reset_w2r_storbe[j]);                                                     \n"
	result += "			//controllo l'abilitazione del data valid                                                                                                                               \n"
	result += "		                                                                                                                                                                        \n"
	result += "			//il data valid non si deve abilitare se viene tolto il read                                                                                                            \n"
	result += "		                                                                                                                                                                        \n"
	result += "			//devo ricevere un solo ready dal canale per read e write                                                                                                               \n"
	result += "                                                                                                                                                                             \n"
	result += "	end\n"
	result += "	endgenerate\n"
	result += "\n"
	result += "\t//devo mandare dal canale un solo wwr o wrd ai canali perchè il canale ne sceglie sono una coppia\n"
	result += "\twwr_ch:assert property (@(posedge clk) disable iff(reset | ~("
	for proc_id := 0; proc_id < num_processors; proc_id++ {
		result += " ack_ch_ready_i[" + strconv.Itoa(proc_id) + "] "
		if proc_id < num_processors-1 {
			result += "&"
		}
	}
	result += ")) $onehot(wwr_ready_ch));\n"
	result += "\tw2r_ch:assert property (@(posedge clk) disable iff(reset | ~("
	for proc_id := 0; proc_id < num_processors; proc_id++ {
		result += " ack_ch_ready_i[" + strconv.Itoa(proc_id) + "] "
		if proc_id < num_processors-1 {
			result += "&"
		}
	}
	result += "))$onehot(wrd_ready_ch));\n"
	result += "\t// synthesis translate_on                                                                                                                                                    \n"
	result += "\n"
	result += "endmodule \n"

	return result
}

func (sm Channel_instance) GetPerProcPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(proc_id) + soname + "chin;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "w2w;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "w2r;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "ack_ch_ready;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "op_check_ready;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "finish_channel;\n"
		result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(proc_id) + soname + "chout;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "ack_w2w;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "ack_w2r;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "ch_ready;\n"
		result += "\twire [1:0] p" + strconv.Itoa(proc_id) + soname + "ch_w_r_ready;\n"
		result += "\n"
	}
	return result
}

func (sm Channel_instance) GetPerProcPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += ", p" + strconv.Itoa(proc_id) + soname + "chin"
		result += ", p" + strconv.Itoa(proc_id) + soname + "w2w"
		result += ", p" + strconv.Itoa(proc_id) + soname + "w2r"
		result += ", p" + strconv.Itoa(proc_id) + soname + "ack_ch_ready"
		result += ", p" + strconv.Itoa(proc_id) + soname + "op_check_ready"
		result += ", p" + strconv.Itoa(proc_id) + soname + "finish_channel"
		result += ", p" + strconv.Itoa(proc_id) + soname + "chout"
		result += ", p" + strconv.Itoa(proc_id) + soname + "ack_w2w"
		result += ", p" + strconv.Itoa(proc_id) + soname + "ack_w2r"
		result += ", p" + strconv.Itoa(proc_id) + soname + "ch_ready"
		result += ", p" + strconv.Itoa(proc_id) + soname + "ch_w_r_ready"
	}
	return result
}

func (sm Channel_instance) GetExternalPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Channel_instance) GetExternalPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}
