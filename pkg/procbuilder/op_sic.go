package procbuilder

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Sic struct{}

func (op Sic) Op_get_name() string {
	return "sic"
}

func (op Sic) Op_get_desc() string {
	return "Wait for an input change and increments a register"
}

func (op Sic) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	inpbits := arch.Inputs_bits()
	result := "sic [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(inpbits) + "(Input)]	// Wait for an input change an increment the register [" + strconv.Itoa(opbits+int(arch.R)+inpbits) + "]\n"
	return result
}

func (op Sic) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	inpbits := arch.Inputs_bits()
	return opbits + int(arch.R) + int(inpbits) // The bits for the opcode + bits for a register + bits for the input
}

func (op Sic) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := "\n"
	result += "\t//Internal Regs for SIC\n"
	result += "\treg sic_state;\n"
	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] sic_reg;\n"
	result += "\n"

	result += "\n"
	// Data received fro inputs
	for j := 0; j < int(arch.N); j++ {
		result += "\treg " + strings.ToLower(Get_input_name(j)) + "_recv;\n"
	}

	result += "\n"
	// TODO Finire a allineare con i2r

	return result
}

func (op Sic) Op_instruction_verilog_state_machine(arch *Arch, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()
	inpbits := arch.Inputs_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					SIC: begin\n"
	if arch.N > 0 {
		if arch.R == 1 {
			result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
		} else {
			result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
		}
		for i := 0; i < reg_num; i++ {
			result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

			if inpbits == 1 {
				result += "							case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
			} else {
				result += "							case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(inpbits)) + "])\n"
			}

			for j := 0; j < int(arch.N); j++ {
				result += "							" + strings.ToUpper(Get_input_name(j)) + " : begin\n"
				result += "								if (sic_state == 1'b1)\n"
				result += "								begin\n"
				result += "									if (sic_reg == " + strings.ToLower(Get_input_name(j)) + ")\n"
				result += "									begin\n"
				result += "										_" + strings.ToLower(Get_register_name(i)) + " <= #1 _" + strings.ToLower(Get_register_name(i)) + " + 1'b1;\n"
				result += "									end\n"
				result += "									else\n"
				result += "									begin\n"
				result += "										sic_state <= #1 1'b0;\n"
				result += "										_pc <= #1 _pc + 1'b1 ;\n"
				result += "									end\n"
				result += "								end\n"
				result += "								else\n"
				result += "								begin\n"
				result += "									sic_state <= #1 1'b1;\n"
				result += "									_" + strings.ToLower(Get_register_name(i)) + " <= #1 0;\n"
				result += "									sic_reg <= #1 " + strings.ToLower(Get_input_name(j)) + ";\n"
				result += "								end\n"
				result += "							$display(\"SIC " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_input_name(j)) + "\");\n"
				result += "							end\n"

			}
			result += "							endcase\n"
			result += "						end\n"
		}
		result += "						endcase\n"
	} else {
		result += "						$display(\"NOP\");\n"
	}
	result += "					end\n"
	return result
}

func (op Sic) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Sic) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	inpbits := arch.Inputs_bits()
	rom_word := arch.Max_word()

	reg_num := 2
	reg_num = reg_num << (arch.R - 1)

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

	if partial, err := Process_input(words[1], int(arch.N)); err == nil {
		result += zeros_prefix(inpbits, partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opbits + int(arch.R) + inpbits; i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op Sic) Disassembler(arch *Arch, instr string) (string, error) {
	inpbits := arch.Inputs_bits()
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	inp_id := get_id(instr[arch.R : int(arch.R)+inpbits])
	result += strings.ToLower(Get_input_name(inp_id))
	return result, nil
}

func (op Sic) Simulate(vm *VM, instr string) error {
	inpbits := vm.Mach.Inputs_bits()
	reg_bits := vm.Mach.R
	reg := get_id(instr[:reg_bits])
	inp := get_id(instr[reg_bits : int(reg_bits)+inpbits])
	if sic_state, ok := vm.Extra_states["sic_state"]; ok {
		if sic_state == true {
			if vm.Extra_states["sic_reg"] == vm.Inputs[inp] {
				switch vm.Mach.Rsize {
				case 8:
					vm.Registers[reg] = vm.Registers[reg].(uint8) + 1
				case 16:
					vm.Registers[reg] = vm.Registers[reg].(uint16) + 1
				default:
					// TODO Fix
				}
			} else {
				vm.Extra_states["sic_state"] = false
				vm.Pc = vm.Pc + 1
			}
		} else {
			vm.Extra_states["sic_state"] = true
			vm.Extra_states["sic_reg"] = vm.Inputs[inp]
			switch vm.Mach.Rsize {
			case 8:
				vm.Registers[reg] = uint8(0)
			case 16:
				vm.Registers[reg] = uint16(0)
			default:
				// TODO Fix
			}
		}
	} else {
		vm.Extra_states["sic_state"] = true
		vm.Extra_states["sic_reg"] = vm.Inputs[inp]
		switch vm.Mach.Rsize {
		case 8:
			vm.Registers[reg] = uint8(0)
		case 16:
			vm.Registers[reg] = uint16(0)
		default:
			// TODO Fix
		}
	}
	return nil
}

func (op Sic) Generate(arch *Arch) string {
	inpbits := arch.Inputs_bits()
	reg_num := 1 << arch.R
	reg := rand.Intn(reg_num)
	inp := rand.Intn(int(arch.N))
	return zeros_prefix(int(arch.R), get_binary(reg)) + zeros_prefix(inpbits, get_binary(inp))
}

func (op Sic) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Sic) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Sic) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Sic) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Sic) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Sic) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Sic) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Sic) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Sic) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq0, types0 := Sequence_to_0(words[0])
	seq1, types1 := Sequence_to_0(words[1])

	if len(seq0) > 0 && types0 == O_REGISTER && len(seq1) > 0 && types1 == O_INPUT {

		result := make([]UsageNotify, 2+len(seq1))
		newnot0 := UsageNotify{C_OPCODE, "sic", I_NIL}
		result[0] = newnot0
		newnot1 := UsageNotify{C_REGSIZE, S_NIL, len(seq0)}
		result[1] = newnot1

		for i, _ := range seq1 {
			result[i+2] = UsageNotify{C_INPUT, S_NIL, i + 1}
		}

		return result, nil

	}

	return []UsageNotify{}, errors.New("Wrong parameters")
}

func (Op Sic) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Sic) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Sic) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Sic) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
