package procbuilder

// TODO: implement
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

type Expf struct{}

func (op Expf) Op_get_name() string {
	return "expf"
}

func (op Expf) Op_get_desc() string {
	return "Exponential of a register (float)"
}

func (op Expf) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "expf [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// xponential of a register (float) [" + strconv.Itoa(opbits+int(arch.R)) + "]\n"
	return result
}

func (op Expf) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) // The bits for the opcode + bits for a register
}

func (op Expf) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	return ""
}

func (Op Expf) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Expf) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Expf) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Expf) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					EXPF: begin\n"
	if arch.R == 1 {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "							_" + strings.ToLower(Get_register_name(i)) + " <= #1 exp(_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "							$display(\"EXPF " + strings.ToUpper(Get_register_name(i)) + "\");\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
	result += "					end\n"
	return result
}

func (op Expf) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Expf) Assembler(arch *Arch, words []string) (string, error) {
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

func (op Expf) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

func (op Expf) Simulate(vm *VM, instr string) error {
	// TODO: implement
	// reg_bits := vm.Mach.R
	// reg := get_id(instr[:reg_bits])
	// switch vm.Mach.Rsize {
	// case 8:
	// 	vm.Registers[reg] = vm.Registers[reg].(uint8) + uint8(1)
	// case 16:
	// 	vm.Registers[reg] = vm.Registers[reg].(uint16) + uint16(1)
	// case 32:
	// 	vm.Registers[reg] = vm.Registers[reg].(uint32) + uint32(1)
	// case 64:
	// 	vm.Registers[reg] = vm.Registers[reg].(uint64) + uint64(1)
	// default:
	// 	return errors.New("go simulation only works on 8,16,32 or 64 bits registers")
	// }
	vm.Pc = vm.Pc + 1
	return nil
}

func (op Expf) Generate(arch *Arch) string {
	reg_num := 1 << arch.R
	reg := rand.Intn(reg_num)
	return zeros_prefix(int(arch.R), get_binary(reg))
}

func (op Expf) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Expf) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Expf) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op Expf) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Expf) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq, types := Sequence_to_0(words[0])
	if len(seq) > 0 && types == O_REGISTER {

		result := make([]UsageNotify, 2)
		newnot0 := UsageNotify{C_OPCODE, "expf", I_NIL}
		result[0] = newnot0
		newnot1 := UsageNotify{C_REGSIZE, S_NIL, len(seq)}
		result[1] = newnot1

		return result, nil
	}

	return []UsageNotify{}, errors.New("wrong register")
}

func (Op Expf) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}

func (Op Expf) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "expf::*--type=reg")
	return result
}

func (Op Expf) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "expf":
		regNeed := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Expf) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Expf) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
