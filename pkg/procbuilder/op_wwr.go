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

type Wwr struct{}

func (op Wwr) Op_get_name() string {
	return "wwr"
}

func (op Wwr) Op_get_desc() string {
	return "Want write to a channel"
}

func (op Wwr) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	chanBits := arch.Shared_bits("channel")
	result := "wwr [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(chanBits) + "(Channel)]	// Want write from a register to a channel  [" + strconv.Itoa(opBits+int(arch.R)+chanBits) + "]\n"
	return result
}

func (op Wwr) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	chanBits := arch.Shared_bits("channel")
	return opBits + int(arch.R) + int(chanBits) // The bits for the opcode + bits for a register + bits for the channel id
}

func (op Wwr) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {

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

	result := ""

	setflag := true
	for _, currop := range arch.Op {
		if currop.Op_get_name() == "wrd" {
			setflag = false
			break
		} else if currop.Op_get_name() == "wwr" {
			break
		}
	}
	if setflag {
		result += "\treg [" + strconv.Itoa(chanbits-1) + ":0] ch_num;\n"
		result += "\treg [" + strconv.Itoa(channel_num-1) + ":0] op_channel;\n"
		result += "\treg [" + strconv.Itoa(int(arch.R)-1) + ":0] reg_num;\n"
		result += "\treg [" + strconv.Itoa((chanbits<<1)-1) + ":0] count_seq_ch;\n"
	}

	result += "\treg [" + strconv.Itoa(channel_num-1) + ":0] wwr_strobe;\n"
	result += "\treg wwr_ch;\n"
	result += "\treg [" + strconv.Itoa(channel_num-1) + ":0] ack_wwr_d1;\n"

	return result
}

func (Op Wwr) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""
	result += "\t\t\twwr_ch <= #1 1'b0;\n"

	setflag := true
	for _, currop := range arch.Op {
		if currop.Op_get_name() == "wrd" {
			setflag = false
			break
		} else if currop.Op_get_name() == "wwr" {
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

func (Op Wwr) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Wwr) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	result := ""
	result += "\t\t\t\twwr_ch <= #1 1'b0;\n"

	setflag := true
	for _, currop := range arch.Op {
		if currop.Op_get_name() == "wrd" {
			setflag = false
			break
		} else if currop.Op_get_name() == "wwr" {
			break
		}
	}
	if setflag {
		result += "\t\t\tch_num <= #1 'b0;\n"
		result += "\t\t\treg_num <= #1 'b0;\n"
	}

	return result
}

func (op Wwr) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	opbits := arch.Opcodes_bits()
	chso := Channel{}
	chanbits := arch.Shared_bits(chso.Shr_get_name())
	rom_word := arch.Max_word()

	result := ""
	result += "\t				WWR: begin\n"
	result += "\t					ch_num <= #1 current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-chanbits) + "];\n"
	result += "\t					reg_num <= #1 current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "];\n"
	result += "\t					wwr_ch <= #1 1'b1;\n"
	result += "\t					op_channel[ch_num] <= #1 1'b1;\n"
	result += "\t					if(ack_wwr_i[ch_num] == 1'b1) begin //ack of the chanel for the operation done\n"
	result += NextInstruction(conf, arch, 5, "_pc + 1'b1")
	result += "\t						count_seq_ch <= #1 count_seq_ch + 1;     //increment the sequence of the channel operation\n"
	result += "\t					end\n"
	result += "\t					$display(\"WWR\");\n"
	result += "\t				end\n"
	return result
}

func (op Wwr) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
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
	result += "\t	if(reset_signal | reset_flag_ch) begin //reset flag of all the channels\n"
	result += "\t		ack_wwr_d1 <= #1 'b0;\n"
	result += "\t	end\n"
	result += "\t	else begin\n"
	result += "\t		ack_wwr_d1 <= #1 ack_wwr_i;\n"
	result += "\t	end\n"
	result += "\tend\n"
	result += "\t\n"
	result += "\t//define flag for write strobe for the channel\n"
	result += "\tgenvar idx_wwr;\n"
	result += "\tgenerate\n"
	result += "\tfor(idx_wwr = 0; idx_wwr < " + strconv.Itoa(channel_num) + "; idx_wwr = idx_wwr + 1) begin\n"
	result += "\t	always @(posedge clock_signal) begin\n"
	result += "\t   	if(reset_signal | reset_flag_ch) begin //reset flag of all the channels\n"
	result += "\t			wwr_strobe[idx_wwr]   <= #1 'b0;\n"
	result += "\t		end\n"
	result += "\t		else if(wwr_ch & (ch_num == idx_wwr)) begin\n"
	result += "\t			wwr_strobe[idx_wwr] <= #1 1'b1;\n"
	result += "\t			end\n"
	result += "\t		end\n"
	result += "\t	end\n"
	result += "\tendgenerate\n"
	result += "\n"
	result += "\tassign ch_wwr_i = wwr_strobe;\n"

	return result
}

func (op Wwr) Assembler(arch *Arch, words []string) (string, error) {
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

func (op Wwr) Disassembler(arch *Arch, instr string) (string, error) {
	chso := Channel{}
	chanbits := arch.Shared_bits(chso.Shr_get_name())
	shortname := chso.Shortname()
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	ch_id := get_id(instr[arch.R : int(arch.R)+chanbits])
	result += shortname + strconv.Itoa(ch_id)
	return result, nil
}

func (op Wwr) Simulate(vm *VM, instr string) error {
	//TODO
	outbits := vm.Mach.Outputs_bits()
	reg_bits := vm.Mach.R
	reg := get_id(instr[:reg_bits])
	inp := get_id(instr[reg_bits : int(reg_bits)+outbits])
	vm.Outputs[inp] = vm.Registers[reg]
	vm.Pc = vm.Pc + 1
	return nil
}

func (op Wwr) Generate(arch *Arch) string {
	//TODO
	outbits := arch.Outputs_bits()
	reg_num := 1 << arch.R
	reg := rand.Intn(reg_num)
	inp := rand.Intn(int(arch.M))
	return zeros_prefix(int(arch.R), get_binary(reg)) + zeros_prefix(outbits, get_binary(inp))
}

func (op Wwr) Required_shared() (bool, []string) {
	return true, []string{"channel"}
}

func (op Wwr) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Wwr) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op Wwr) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Wwr) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Wwr) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Wwr) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Wwr) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Wwr) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Wwr) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
