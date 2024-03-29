package procbuilder

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Clr struct{}

func (op Clr) Op_get_name() string {
	return "clr"
}

func (op Clr) Op_get_desc() string {
	return "Clear register"
}

func (op Clr) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "clr [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to 0 [" + strconv.Itoa(opbits+int(arch.R)) + "]\n"
	return result
}

func (op Clr) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) // The bits for the opcode + bits for a register
}

func (op Clr) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	return ""
}

func (Op Clr) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}
func (Op Clr) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Clr) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Clr) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					CLR: begin\n"
	if arch.R == 1 {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "							_" + strings.ToLower(Get_register_name(i)) + " <= #1 'b0;\n"
		result += "							$display(\"CLR " + strings.ToUpper(Get_register_name(i)) + "\");\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
	result += "					end\n"
	return result
}

func (op Clr) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Clr) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()

	reg_num := 1 << arch.R

	if len(words) != 1 {
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

	for i := opbits + int(arch.R); i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op Clr) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	return result, nil
}

func (op Clr) Simulate(vm *VM, instr string) error {
	reg_bits := vm.Mach.R
	reg := get_id(instr[:reg_bits])
	switch vm.Mach.Rsize {
	case 8:
		vm.Registers[reg] = uint8(0)
	case 16:
		vm.Registers[reg] = uint16(0)
	case 32:
		vm.Registers[reg] = uint32(0)
	case 64:
		vm.Registers[reg] = uint64(0)
	default:
		return errors.New("go simulation only works on 8,16,32 or 64 bits registers")
	}
	vm.Pc = vm.Pc + 1
	return nil
}

func (op Clr) Generate(arch *Arch) string {
	reg_num := 1 << arch.R
	reg := rand.Intn(reg_num)
	return zeros_prefix(int(arch.R), get_binary(reg))
}

func (op Clr) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Clr) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Clr) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op Clr) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Clr) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq, types := Sequence_to_0(words[0])
	if len(seq) > 0 && types == O_REGISTER {

		result := make([]UsageNotify, 2)
		newnot0 := UsageNotify{C_OPCODE, "clr", I_NIL}
		result[0] = newnot0
		newnot1 := UsageNotify{C_REGSIZE, S_NIL, len(seq)}
		result[1] = newnot1

		return result, nil
	}

	return []UsageNotify{}, errors.New("Wrong register")

}

func (Op Clr) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}

func (Op Clr) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "clr::*--type=reg")
	return result
}

func (Op Clr) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "clr":
		regNeed := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Clr) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Clr) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	switch line.Operation.GetValue() {
	case "clr":
		regNeed := line.Elements[0].GetValue()
		if regNeed != "" {
			var meta *bmmeta.BasmMeta
			meta = meta.SetMeta("inv", regNeed)
			return meta, nil
		}
	}
	return nil, nil
}
