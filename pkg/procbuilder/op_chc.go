package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The Chc opcode is both a basic instruction and a template for other instructions.
type Chc struct{}

func (op Chc) Op_get_name() string {
	return "chc"
}

func (op Chc) Op_get_desc() string {
	return "Channel operation check"
}

func (op Chc) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "chc [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Check the channels operation [" + strconv.Itoa(opbits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op Chc) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op Chc) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""

	//regsize := int(arch.Rsize)

	chso := Channel{}
	chanbits := arch.Shared_bits(chso.Shr_get_name())

	channel_num := 0
	if arch.Shared_constraints != "" {
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if soname == "channel" {
				channel_num++
			}
		}
	}

	result += "\treg [1:0] count_chc;\n"

	setflag := true
	for _, currop := range arch.Op {
		if currop.Op_get_name() == "chw" {
			setflag = false
			break
		} else if currop.Op_get_name() == "chc" {
			break
		}
	}
	if setflag {
		result += "\treg reset_flag_ch;\n"
		result += "\treg wrd_ch_ok;\n"
		result += "\treg [" + strconv.Itoa(chanbits-1) + ":0] count_ready;\n"
		result += "\treg find_ready, find_op;\n"
		result += "\treg [" + strconv.Itoa(channel_num-1) + ":0] ready_ch;\n"
		result += "\treg [" + strconv.Itoa(chanbits-1) + ":0] ch_num_ack;\n"
		result += "\treg [" + strconv.Itoa(int(arch.R)-1) + ":0] reg_num_ack;\n"
		result += "\treg [" + strconv.Itoa(channel_num-1) + ":0] finish_channel_d1, finish_channel_d2;\n"
		result += "\treg [" + strconv.Itoa(2+chanbits+int(arch.R)-1) + ":0] stat_op_ch [0:" + strconv.Itoa(2*channel_num-1) + "];\n"
		result += "\treg [1:0] chech_stat_op_w_r [0:" + strconv.Itoa(2*channel_num-1) + "];\n"
		result += "\treg [" + strconv.Itoa(int(arch.R)-1) + ":0] chech_stat_op_reg_num [0:" + strconv.Itoa(2*channel_num-1) + "];\n"
		result += "\treg [" + strconv.Itoa(chanbits-1) + ":0] check_stat_op_ch_num [0:" + strconv.Itoa(2*channel_num-1) + "];\n"
		result += "\treg [" + strconv.Itoa(2*channel_num-1) + ":0] stat_op_int;\n"
		result += "\treg dataToCh;\n"
	}

	return result
}

func (Op Chc) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""

	setflag := true
	for _, currop := range arch.Op {
		if currop.Op_get_name() == "chw" {
			setflag = false
			break
		} else if currop.Op_get_name() == "chc" {
			break
		}
	}
	if setflag {
		result += "\t\t\treset_flag_ch <= #1 1'b0;\n"
		result += "\t\t\tch_op_ready_i <= #1 1'b0;\n"
	}

	return result
}

func (Op Chc) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	result := ""

	setflag := true
	for _, currop := range arch.Op {
		if currop.Op_get_name() == "chw" {
			setflag = false
			break
		} else if currop.Op_get_name() == "chc" {
			break
		}
	}
	if setflag {
		result += "\t\t\treset_flag_ch <= #1 1'b0;\n"
		result += "\t\t\tch_op_ready_i <= #1 1'b0;\n"
	}

	return result
}

func (op Chc) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "\t				CHC: begin\n"
	result += "\t					wwr_ch <= #1   1'b0;\n"
	result += "\t					wrd_ch <= #1   1'b0;\n"
	result += "\t					ch_op_ready_i <= #1 op_channel;\n"
	result += "\t					if(finish_channel_i[ch_num_ack]) begin\n"
	result += "\t						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	for i := 0; i < reg_num; i++ {
		result += "\t							" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "\t								_" + strings.ToLower(Get_register_name(i)) + " <= #1 " + strconv.Itoa(int(arch.R)) + "'b001;\n" //CHECK
		result += "\t								$display(\"CHC " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "\t							end\n"
	}
	result += "\t						endcase\n"
	result += "\t						case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.R)) + "])\n"
	for i := 0; i < reg_num; i++ {
		result += "\t							" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "\t								_" + strings.ToLower(Get_register_name(i)) + " <= #1 stat_op_int;\n" //CHECK
		result += "\t								$display(\"CHC " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "\t							end\n"
	}
	result += "\t						endcase\n"
	result += "\t					end\n"
	result += "\t					if(wrd_ch_ok  & finish_channel_i[ch_num_ack]) begin\n"
	result += "\t						case (reg_num_ack)\n"
	for i := 0; i < reg_num; i++ {
		result += "\t						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "\t							_" + strings.ToLower(Get_register_name(i)) + " <= #1 ch2proc_i[ch_num_ack];\n" //CHECK
		result += "\t							$display(\"CHC " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "\t						end\n"
	}
	result += "\t						endcase\n"
	result += "\t					end\n"
	result += "\t					if(finish_channel_d1[ch_num_ack]) begin\n"
	result += "\t						count_seq_ch <= #1 'b0;\n"
	result += "\t						reset_flag_ch <= #1 'b1;\n"
	result += "\t					end\n"
	result += "\t					if(finish_channel_d2[ch_num_ack]) begin\n"
	result += "\t						reset_flag_ch <= #1 'b0;\n"
	result += "\t						_pc <= #1  _pc + 1'b1;\n"
	result += "\t					end\n"
	result += "\t					if((~(|ch_ready_i)) & (count_chc==2'b11)) begin\n"
	result += "\t						reset_flag_ch <= #1 'b1;\n"
	result += "\t						count_seq_ch <= #1 'b0;\n"
	result += "\t						_pc <= #1  _pc + 1'b1;\n"
	result += "\t					end\n"
	result += "\t				end\n"

	return result
}

func (op Chc) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	result := ""

	//regsize := int(arch.Rsize)
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()
	reg_num := 1 << arch.R

	chso := Channel{}
	chanbits := arch.Shared_bits(chso.Shr_get_name())

	channel_num := 0
	if arch.Shared_constraints != "" {
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if soname == "channel" {
				channel_num++
			}
		}
	}

	result += "\talways @(posedge clock_signal)\n"
	result += "\tbegin\n"
	result += "\t    if(reset_signal | reset_flag_ch)\n"
	result += "\t        count_chc <= #1 'b0;\n"
	result += "\t    else\n"
	result += "\t    	if(rom_value[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "]==CHC)\n"
	result += "\t        count_chc <= #1 count_chc + 1;\n"
	result += "\tend\n"

	setflag := true
	for _, currop := range arch.Op {
		if currop.Op_get_name() == "chw" {
			setflag = false
			break
		} else if currop.Op_get_name() == "chc" {
			break
		}
	}
	if setflag {

		result += "\treg [" + strconv.Itoa(chanbits+1) + ":0] nextbit;\n"

		result += "\talways @(posedge clock_signal)\n"
		result += "\tbegin\n"
		result += "\t    if(reset_signal) begin\n"
		result += "\t        finish_channel_d1 <= #1 'b0;\n"
		result += "\t        finish_channel_d2 <= #1 'b0;\n"
		result += "\t    end\n"
		result += "\t    else begin\n"
		result += "\t        finish_channel_d1 <= #1 finish_channel_i;\n"
		result += "\t        finish_channel_d2 <= #1 finish_channel_d1;\n"
		result += "\t    end\n"
		result += "\tend\n"

		result += "\t//Define the array of status channel sequences: 01 wwr - 10 wwr - 00 no operation\n"
		result += "\tinteger k;\n"
		result += "\talways @(posedge clock_signal) begin\n"
		result += "\t	if(reset_signal | reset_flag_ch) begin\n"
		result += "\t		for(k=0;k<" + strconv.Itoa(channel_num*2) + ";k=k+1)\n"
		result += "\t			stat_op_ch[k]  <= #1 'b0;\n"
		result += "\t	end\n"
		result += "\t	else if(ack_wwr_i) begin\n"
		result += "\t		stat_op_ch[count_seq_ch] <= #1 {2'b01, reg_num, ch_num};\n"
		result += "\t	end\n"
		result += "\t	else if(ack_wrd_i) begin\n"
		result += "\t		stat_op_ch[count_seq_ch] <= #1 {2'b10, reg_num, ch_num};\n"
		result += "\t	end\n"
		result += "\tend\n"

		result += "\talways @* begin\n"
		result += "\t	for(k=0;k<" + strconv.Itoa(channel_num*2) + ";k=k+1) begin\n"
		result += "\t   	chech_stat_op_w_r[k] <= stat_op_ch[k][" + strconv.Itoa(2+chanbits+int(arch.R)-1) + ":" + strconv.Itoa(2+chanbits+int(arch.R)-2) + "];\n"
		result += "\t    	chech_stat_op_reg_num[k] <= stat_op_ch[k][" + strconv.Itoa(chanbits+int(arch.R)-1) + ":" + strconv.Itoa(chanbits) + "];\n"
		result += "\t    	check_stat_op_ch_num[k] <= stat_op_ch[k][" + strconv.Itoa(chanbits-1) + ":0];\n"
		result += "\t	end\n"
		result += "\tend\n"

		result += "\tinteger idy, idyy;\n"
		result += "\tinteger idz, idzz;\n"
		result += "\talways @(posedge clock_signal) begin\n"
		result += "\t	if(reset_signal) begin\n"
		result += "\t		ch_num_ack <= #1 'b0;\n"
		result += "\t		reg_num_ack <= #1 'b0;\n"
		result += "\t		ack_ch_ready_i <= #1 'b0;\n"
		result += "\t		wrd_ch_ok <= #1 1'b0;\n"
		result += "\t		count_ready <= #1 'b0;\n"
		result += "\t		stat_op_int <= #1 'b0;\n"
		result += "\t		find_ready <= #1 'b0;\n"
		result += "\t		find_op <= #1 'b0;\n"
		result += "\t		dataToCh <= #1 1'b0;\n"
		//result += "\t		for(k=0;k<" + strconv.Itoa(channel_num) + ";k=k+1)\n"
		//result += "\t			proc2ch_i[k] <= #1 'b0;\n"
		result += "\t		end\n"
		result += "\t	else begin\n"
		result += "\t		if(reset_flag_ch) begin\n"
		//result += "\t			for(k=0;k<" + strconv.Itoa(channel_num) + ";k=k+1)\n"
		//result += "\t				proc2ch_i[k] <= #1 'b0;\n"
		result += "\t			wrd_ch_ok <= #1 1'b0;\n"
		result += "\t			ack_ch_ready_i <= #1 'b0;\n"
		result += "\t			ch_num_ack <= #1 'b0;\n"
		result += "\t			find_ready <= #1 'b0;\n"
		result += "\t			find_op <= #1 'b0;\n"
		result += "\t			dataToCh <= #1 1'b0;\n"
		result += "\t		end\n"
		result += "\t		else begin\n"
		result += "\t			if(ch_ready_i[count_ready]==1 & find_ready==1'b0) begin  //serach the value after the counter\n"
		result += "\t				for( idz = 0; idz < " + strconv.Itoa(channel_num*2) + "; idz = idz + 1) begin\n"
		result += "\t					if(check_stat_op_ch_num[idz] == count_ready & find_op==1'b0) begin\n" //CEHCK
		result += "\t						if(ch_w_r_ready_i[count_ready] == chech_stat_op_w_r[idz]) begin\n"
		result += "\t							if(ch_w_r_ready_i[count_ready] == 2'b01) begin //wwr\n"
		result += "\t								ack_ch_ready_i[count_ready] <= #1 1'b1;\n"
		result += "\t								ch_num_ack <= #1 count_ready;\n"
		result += "\t								reg_num_ack <= #1 chech_stat_op_reg_num[idz];\n"
		result += "\t								dataToCh <= #1 1'b1;\n"
		//result += "\t								case (chech_stat_op_reg_num[idz])\n"
		//for i := 0; i < reg_num; i++ {
		//	result += "\t									" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		//	result += "\t										proc2ch_i[count_ready] <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
		//	result += "\t									end\n"
		//}
		//result += "\t								endcase\n"
		result += "\t								stat_op_int <= #1 idz;\n"
		result += "\t							end\n"
		result += "\t							if (ch_w_r_ready_i[count_ready] == 2'b10) begin //wrd\n"
		result += "\t								ack_ch_ready_i[count_ready] <= #1 1'b1;\n"
		result += "\t								wrd_ch_ok <= #1 1'b1;\n"
		result += "\t								ch_num_ack <= #1 count_ready;\n"
		result += "\t								reg_num_ack <= #1 chech_stat_op_reg_num[idz];\n"
		result += "\t								stat_op_int <= #1 idz;\n"
		result += "\t							end\n"
		result += "\t							find_op <= #1 1'b1;\n"
		result += "\t						end\n"
		result += "\t					end\n"
		result += "\t				end\n"
		result += "\t				find_ready <= #1 1'b1;\n"
		result += "\t			end\n"
		result += "\t			if(ch_ready_i == 'b0) begin\n"
		result += "\t				wrd_ch_ok <= #1 'b0;\n"
		result += "\t				ch_num_ack <= #1 'b0;\n"
		result += "\t				reg_num_ack <= #1 'b0;\n"
		result += "\t				stat_op_int <= #1 'b0;\n"
		result += "\t				ack_ch_ready_i <= #1 'b0;\n"
		result += "\t				find_ready <= #1 'b0;\n"
		result += "\t				find_op <= #1 'b0;\n"
		result += "\t			end\n"
		result += "\t			if(find_ready==1'b0) begin\n"
		if channel_num == 1 {
			result += "\t			count_ready <= #1 1'b0;\n"
		} else {
			result += "\t			count_ready <= #1 nextbit[" + strconv.Itoa(chanbits) + ":0];\n"
		}
		result += "\t			end\n"
		result += "\t		end\n"
		result += "\t	end\n"
		result += "\tend\n"

		result += "\talways @(posedge clock_signal) begin\n"
		result += "\t	if(reset_signal) begin\n"
		result += "\t		for(k=0;k<" + strconv.Itoa(channel_num) + ";k=k+1)\n"
		result += "\t			proc2ch_i[k] <= #1 'b0;\n"
		result += "\t	end\n"
		result += "\t	else begin\n"
		result += "\t		if(reset_flag_ch) begin\n"
		result += "\t			for(k=0;k<" + strconv.Itoa(channel_num) + ";k=k+1)\n"
		result += "\t				proc2ch_i[k] <= #1 'b0;\n"
		result += "\t		end\n"
		result += "\t		else if(dataToCh) begin\n"
		result += "\t			case (reg_num_ack)\n"
		for i := 0; i < reg_num; i++ {
			result += "\t				" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
			result += "\t					proc2ch_i[ch_num_ack] <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
			result += "\t				end\n"
		}
		result += "\t		endcase\n"
		result += "\t	end\n"
		result += "\t	end\n"
		result += "\tend\n"

		//result += "\treg [" + strconv.Itoa(chanbits) + ":0] nextbit;\n"
		result += "\talways @(posedge clock_signal) begin\n"
		result += "\t	if(reset_signal)\n"
		result += "\t		nextbit <= #1 'b1;\n"
		result += "\t	else\n"
		result += "\t		nextbit <= #1 { nextbit[" + strconv.Itoa(chanbits) + ":0], nextbit[" + strconv.Itoa(chanbits+1) + "] ^ nextbit[" + strconv.Itoa(chanbits) + "] };\n"
		result += "\tend\n"

		result += "\t/*wire [7:0] seed = 4'b11010010;\n"
		result += "\twire load = (_pc == 'b0) ? 1'b1 : 1'b0;\n"
		result += "\treg [7:0] state_in;\n"
		result += "\treg [7:0] state_out;\n"

		result += "\talways @(posedge clock_signal) begin\n"
		result += "\tif(reset_signal)\n"
		result += "\t	state_out <= #1 'b0;\n"
		result += "\telse\n"
		result += "\t	state_out <= #1 state_in;\n"
		result += "\tend\n"

		result += "\talways @ (state_in or load or seed or state_out)\n"
		result += "\tbegin : MUX\n"
		result += "\tif (load == 1'b0) begin\n"
		result += "\t	state_in[7:1] = state_out[6:0];\n"
		result += "\t	state_in[0] = nextbit;\n"
		result += "\tend else\n"
		result += "\t	state_in = seed;\n"
		result += "\tend\n"
		result += "\tassign nextbit = state_out[6] ^ state_out[7];*/\n"

	}
	return result
}

func (op Chc) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()

	reg_num := 1 << arch.R

	if len(words) != 2 {
		return "", Prerror{"Wrong arguments number"}
	}

	result := ""
	for i := 0; i < reg_num; i++ {
		if words[0] == strings.ToLower(Get_register_name(i)) {
			result += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if result == "" {
		return "", Prerror{"Unknown register name " + words[0]}
	}

	partial := ""
	for i := 0; i < reg_num; i++ {
		if words[1] == strings.ToLower(Get_register_name(i)) {
			partial += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if partial == "" {
		return "", Prerror{"Unknown register name " + words[1]}
	}

	result += partial

	for i := opbits + 2*int(arch.R); i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op Chc) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

// The simulation does nothing
func (op Chc) Simulate(vm *VM, instr string) error {
	// TODO
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Chc) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Chc) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Chc) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Chc) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op Chc) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Chc) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Chc) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Chc) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Chc) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Chc) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Chc) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
