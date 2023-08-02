package procbuilder

import (
	//"fmt"
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Wrd struct{}

func (op Wrd) Op_get_name() string {
	return "wrd"
}

func (op Wrd) Op_get_desc() string {
	return "Want read from a channel"
}

func (op Wrd) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	chanbits := arch.Shared_bits("channel")
	result := "wrd [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(chanbits) + "(Channel)]	// Want read from a channel to a register [" + strconv.Itoa(opbits+int(arch.R)+chanbits) + "]\n"
	return result
}

func (op Wrd) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	chanbits := arch.Shared_bits("channel")
	return opbits + int(arch.R) + int(chanbits) // The bits for the opcode + bits for a register + bits for the channel id
}

func (op Wrd) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""

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

	setflag := true
	for _, currop := range arch.Op {
		if currop.Op_get_name() == "wwr" {
			setflag = false
			break
		} else if currop.Op_get_name() == "wrd" {
			break
		}
	}
	if setflag {
		result += "\treg [" + strconv.Itoa(chanbits-1) + ":0] ch_num;\n"
		result += "\treg [" + strconv.Itoa(channel_num-1) + ":0] op_channel;\n"
		result += "\treg [" + strconv.Itoa(int(arch.R)-1) + ":0] reg_num;\n"
		result += "\treg [" + strconv.Itoa((chanbits<<1)-1) + ":0] count_seq_ch;\n"
	}

	result += "\treg [" + strconv.Itoa(channel_num-1) + ":0] wrd_strobe;\n"
	result += "\treg wrd_ch;\n"
	result += "\treg [" + strconv.Itoa(channel_num-1) + ":0] ack_wrd_d1;\n"

	return result
}

func (op Wrd) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	opbits := arch.Opcodes_bits()
	chso := Channel{}
	chanbits := arch.Shared_bits(chso.Shr_get_name())
	rom_word := arch.Max_word()

	result := ""
	result += "\t				WRD: begin\n"
	result += "\t					ch_num <= #1 current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-chanbits) + "];\n"
	result += "\t					reg_num <= #1 current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "];\n"
	result += "\t					wrd_ch <= #1 1'b1;\n"
	result += "\t					op_channel[ch_num] <= #1 1'b1;\n"
	result += "\t					if(ack_wrd_i[ch_num] == 1'b1) begin //ack of the chanel for the operation done\n"
	result += "\t                 		_pc <= #1  _pc + 1'b1;\n"
	result += "\t						count_seq_ch <= #1 count_seq_ch + 1;     //increment the sequence of the channel operation\n"
	result += "\t					end\n"
	result += "\t					$display(\"WRD\");\n"
	result += "\t				end\n"
	return result
}

func (op Wrd) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	result := ""

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

	result += "\n"
	result += "\talways @(posedge clock_signal)\n"
	result += "\tbegin\n"
	result += "\t	if(reset_signal | reset_flag_ch) begin //reset flag of all the cahnnels\n"
	result += "\t		ack_wrd_d1 <= #1 'b0;\n"
	result += "\t	end\n"
	result += "\t	else begin\n"
	result += "\t		ack_wrd_d1 <= #1 ack_wrd_i;\n"
	result += "\t	end\n"
	result += "\tend\n"
	result += "\t\n"
	result += "\t//define flag for write strobe for the channel\n"
	result += "\tgenvar idx_wrd;\n"
	result += "\tgenerate\n"
	result += "\tfor(idx_wrd = 0; idx_wrd < " + strconv.Itoa(channel_num) + "; idx_wrd = idx_wrd + 1) begin\n"
	result += "\t	always @(posedge clock_signal) begin\n"
	result += "\t   	if(reset_signal | reset_flag_ch) begin //reset flag of all the cahnnels\n"
	result += "\t			wrd_strobe[idx_wrd]   <= #1 'b0;\n"
	result += "\t		end\n"
	result += "\t		else if(wrd_ch & (ch_num == idx_wrd)) begin\n"
	result += "\t			wrd_strobe[idx_wrd] <= #1 1'b1;\n"
	result += "\t			end\n"
	result += "\t		end\n"
	result += "\t	end\n"
	result += "\tendgenerate\n"
	result += "\n"
	result += "\tassign ch_wrd_i = wrd_strobe;\n"

	return result
}

func (op Wrd) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	chso := Channel{}
	channum := arch.Shared_num(chso.Shr_get_name())
	chanbits := arch.Shared_bits(chso.Shr_get_name())
	shortname := chso.Shortname()
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

	if partial, err := Process_shared(shortname, words[1], channum); err == nil {
		result += zeros_prefix(chanbits, partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opbits + int(arch.R) + chanbits; i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op Wrd) Disassembler(arch *Arch, instr string) (string, error) {
	chso := Channel{}
	chanbits := arch.Shared_bits(chso.Shr_get_name())
	shortname := chso.Shortname()
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	ch_id := get_id(instr[arch.R : int(arch.R)+chanbits])
	result += shortname + strconv.Itoa(ch_id)
	return result, nil
}

func (op Wrd) Simulate(vm *VM, instr string) error {
	//TODO
	outbits := vm.Mach.Outputs_bits()
	reg_bits := vm.Mach.R
	reg := get_id(instr[:reg_bits])
	inp := get_id(instr[reg_bits : int(reg_bits)+outbits])
	vm.Outputs[inp] = vm.Registers[reg]
	vm.Pc = vm.Pc + 1
	return nil
}

func (op Wrd) Generate(arch *Arch) string {
	//TODO
	outbits := arch.Outputs_bits()
	reg_num := 1 << arch.R
	reg := rand.Intn(reg_num)
	inp := rand.Intn(int(arch.M))
	return zeros_prefix(int(arch.R), get_binary(reg)) + zeros_prefix(outbits, get_binary(inp))
}

func (op Wrd) Required_shared() (bool, []string) {
	return true, []string{"channel"}
}

func (op Wrd) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Wrd) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op Wrd) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""
	result += "\t\t\twrd_ch <= #1 1'b0;\n"

	setflag := true
	for _, currop := range arch.Op {
		if currop.Op_get_name() == "wwr" {
			setflag = false
			break
		} else if currop.Op_get_name() == "wrd" {
			break
		}
	}
	if setflag {
		result += "\t\t\tcount_seq_ch <= #1 'b0;\n"
		result += "\t\t\tch_num <= #1 'b0;\n"
		result += "\t\t\treg_num <= #1 'b0;\n"
		result += "\t\t\top_channel <= #1 'b0;\n"
	}

	return result
}

func (Op Wrd) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	result := ""
	result += "\t\t\t\twrd_ch <= #1 1'b0;\n"

	setflag := true
	for _, currop := range arch.Op {
		if currop.Op_get_name() == "wwr" {
			setflag = false
			break
		} else if currop.Op_get_name() == "wrd" {
			break
		}
	}
	if setflag {
		result += "\t\t\tch_num <= #1 'b0;\n"
		result += "\t\t\treg_num <= #1 'b0;\n"
	}

	return result
}

func (Op Wrd) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Wrd) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Wrd) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Wrd) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Wrd) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Wrd) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Wrd) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Wrd) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
