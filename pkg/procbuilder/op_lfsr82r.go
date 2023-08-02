package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The Lfsr82r opcode is both a basic instruction and a template for other instructions.
type Lfsr82r struct{}

func (op Lfsr82r) Op_get_name() string {
	return "lfsr82r"
}

func (op Lfsr82r) Op_get_desc() string {
	return "Read a pseudo casual number from a Lfsr8 SO"
}

func (op Lfsr82r) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	lfbits := arch.Shared_bits("lfsr8")
	result := "lfsr82r [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(lfbits) + "(Lfsr8)]       // Read a pseudo casual number from a Lfsr8 SO [" + strconv.Itoa(opbits+int(arch.R)+lfbits) + "]\n"
	return result
}

func (op Lfsr82r) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	lfbits := arch.Shared_bits("lfsr8")
	return opbits + int(arch.R) + int(lfbits) // The bits for the opcode + bits for a register + bits for the barrier id
}

func (op Lfsr82r) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	// TODO
	return ""
}

func (Op Lfsr82r) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (Op Lfsr82r) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (Op Lfsr82r) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Lfsr82r) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()
	lfsr8_num := arch.Shared_num("lfsr8")
	lfbits := arch.Shared_bits("lfsr8")
	soname := "lfsr8"

	reg_num := 1 << arch.R

	result := ""
	result += "					LFSR82R: begin\n"
	if arch.M > 0 {
		if arch.R == 1 {
			result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
		} else {
			result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
		}
		for i := 0; i < reg_num; i++ {
			result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

			if lfbits == 1 {
				result += "							case (current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
			} else {
				result += "							case (current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(lfbits)) + "])\n"
			}

			for j := 0; j < lfsr8_num; j++ {
				result += "							" + zeros_prefix(lfbits, get_binary(j)) + " : begin\n"
				result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 " + strings.ToLower(soname+strconv.Itoa(j)+"out") + ";\n"
				result += "								$display(\"LFSR8 " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(soname+strconv.Itoa(j)) + "\");\n"
				result += "							end\n"

			}
			result += "							endcase\n"
			result += "						end\n"
		}
		result += "						endcase\n"
	} else {
		result += "						$display(\"NOP\");\n"
	}
	result += "						_pc <= #1 _pc + 1'b1 ;\n"
	result += "					end\n"
	return result
}

func (op Lfsr82r) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Lfsr82r) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	lfsr8so := Lfsr8{}
	lfsr8num := arch.Shared_num(lfsr8so.Shr_get_name())
	lfsr8bits := arch.Shared_bits(lfsr8so.Shr_get_name())
	shortname := lfsr8so.Shortname()
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

	if partial, err := Process_shared(shortname, words[1], lfsr8num); err == nil {
		result += zeros_prefix(lfsr8bits, partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opbits + int(arch.R) + lfsr8bits; i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op Lfsr82r) Disassembler(arch *Arch, instr string) (string, error) {
	lfsr8so := Lfsr8{}
	lfsr8bits := arch.Shared_bits(lfsr8so.Shr_get_name())
	shortname := lfsr8so.Shortname()
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	lfsr8_id := get_id(instr[arch.R : int(arch.R)+lfsr8bits])
	result += shortname + strconv.Itoa(lfsr8_id)
	return result, nil
}

// The simulation does nothing
func (op Lfsr82r) Simulate(vm *VM, instr string) error {
	// TODO
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Lfsr82r) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Lfsr82r) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Lfsr82r) Required_modes() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Lfsr82r) Forbidden_modes() (bool, []string) {
	// TODO
	return false, []string{}
}

func (Op Lfsr82r) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Lfsr82r) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq0, types0 := Sequence_to_0(words[0])
	seq1, types1 := Sequence_to_0(words[1])

	if len(seq0) > 0 && types0 == O_REGISTER && len(seq1) > 0 && types1 == O_OUTPUT {

		result := make([]UsageNotify, 2+len(seq1))
		newnot0 := UsageNotify{C_OPCODE, "lfsr82r", I_NIL}
		result[0] = newnot0
		newnot1 := UsageNotify{C_REGSIZE, S_NIL, len(seq0)}
		result[1] = newnot1

		for i, _ := range seq1 {
			result[i+2] = UsageNotify{C_OUTPUT, S_NIL, i + 1}
		}

		return result, nil

	}

	return []UsageNotify{}, errors.New("Wrong parameters")

}

func (Op Lfsr82r) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Lfsr82r) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Lfsr82r) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Lfsr82r) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Lfsr82r) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
