package procbuilder

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Call struct {
	callName string
	s        int
	sn       string
	opType   uint8
}

func (op Call) Op_get_name() string {
	return op.callName
}

func (op Call) Op_get_desc() string {
	switch op.opType {
	case OP_CALLO:
		return "Call a rom subroutine via an hardware stack called " + op.sn + " with depth " + strconv.Itoa(op.s)
	case OP_CALLA:
		return "Call a ram subroutine via an hardware stack called " + op.sn + " with depth " + strconv.Itoa(op.s)
	case OP_RET:
		return "Return from a subroutine via an hardware stack called " + op.sn + " with depth " + strconv.Itoa(op.s)
	}
	return ""
}

func (op Call) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := ""
	switch op.opType {
	case OP_CALLO:
		result += op.callName + " [" + strconv.Itoa(opBits) + "(Location)]	// Call a rom subroutine via an hardware stack called " + op.sn + " with depth " + strconv.Itoa(op.s) + " [" + strconv.Itoa(opBits+int(arch.O)) + "]\n"
	case OP_CALLA:
		result += op.callName + " [" + strconv.Itoa(opBits) + "(Location)]	// Call a ram subroutine via an hardware stack called " + op.sn + " with depth " + strconv.Itoa(op.s) + " [" + strconv.Itoa(opBits+int(arch.L)) + "]\n"
	case OP_RET:
		result += op.callName + " [" + strconv.Itoa(opBits) + "]	// Return from a subroutine via an hardware stack called " + op.sn + " with depth " + strconv.Itoa(op.s) + " [" + strconv.Itoa(opBits) + "]\n"
	}
	return result
}

func (op Call) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	switch op.opType {
	case OP_CALLO:
		return opBits + int(arch.O) // The bits for the opcode + bits for a location
	case OP_CALLA:
		return opBits + int(arch.L) // The bits for the opcode + bits for a location
	case OP_RET:
		return opBits // The bits for the opcode
	}
	return 0
}

func (op Call) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	return ""
}

func (op Call) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (op Call) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op Call) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Call) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	locationBits := arch.O

	switch arch.Modes[0] {
	case "ha":
		locationBits = arch.O
	case "vn":
		locationBits = arch.L
	case "hy":
		if arch.O > arch.L {
			locationBits = arch.O
		} else {
			locationBits = arch.L
		}
	}

	result := ""
	result += "					J: begin\n"
	if locationBits == 1 {
		result += NextInstruction(conf, arch, 6, "current_instruction["+strconv.Itoa(rom_word-opbits-1)+"]")
		result += "						$display(\"J \", current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "]);\n"
	} else {
		result += NextInstruction(conf, arch, 6, "current_instruction["+strconv.Itoa(rom_word-opbits-1)+":"+strconv.Itoa(rom_word-opbits-int(locationBits))+"]")
		result += "						$display(\"J \", current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(locationBits)) + "]);\n"
	}
	result += "					end\n"
	return result
}

func (op Call) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Call) Assembler(arch *Arch, words []string) (string, error) {
	opBits := arch.Opcodes_bits()
	romWord := arch.Max_word()

	locationBits := arch.O

	switch op.opType {
	case OP_CALLO:
		locationBits = arch.O
		if len(words) != 1 {
			return "", Prerror{"Wrong arguments number"}
		}
	case OP_CALLA:
		locationBits = arch.L
		if len(words) != 1 {
			return "", Prerror{"Wrong arguments number"}
		}
	case OP_RET:
		locationBits = 0
		if len(words) != 0 {
			return "", Prerror{"Wrong arguments number"}
		}
	}

	result := ""
	if op.opType != OP_RET {
		if partial, err := Process_number(words[0]); err == nil {
			result += zeros_prefix(int(locationBits), partial)
		} else {
			return "", Prerror{err.Error()}
		}
	}
	for i := opBits + int(locationBits); i < romWord; i++ {
		result += "0"
	}
	return result, nil
}

func (op Call) Disassembler(arch *Arch, instr string) (string, error) {

	locationBits := arch.O

	switch op.opType {
	case OP_CALLO:
		locationBits = arch.O
	case OP_CALLA:
		locationBits = arch.L
	}
	result := ""
	if op.opType != OP_RET {
		value := get_id(instr[:locationBits])
		result += strconv.Itoa(value)
	}
	return result, nil
}

func (op Call) Simulate(vm *VM, instr string) error {
	value := get_id(instr[:vm.Mach.O])
	if value < len(vm.Mach.Slocs) {
		vm.Pc = uint64(value)
	} else {
		vm.Pc = vm.Pc + 1
	}
	return nil
}

func (op Call) Generate(arch *Arch) string {
	max_value := 1 << arch.O
	value := rand.Intn(max_value)
	return zeros_prefix(int(arch.O), get_binary(value))
}

func (op Call) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Call) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Call) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Call) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (op Call) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "j", I_NIL}
	result[0] = newnot
	return result, nil
}

func (op Call) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (op Call) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	switch op.opType {
	case OP_CALLO:
		result = append(result, op.callName+"::*--type=number")
		result = append(result, op.callName+"::*--type=symbol")
		result = append(result, op.callName+"::*--type=rom--romaddressing=symbol")
	case OP_CALLA:
		result = append(result, op.callName+"::*--type=number")
		result = append(result, op.callName+"::*--type=symbol")
		result = append(result, op.callName+"::*--type=ram--ramaddressing=symbol")
	case OP_RET:
		result = append(result, op.callName)
	}
	return result
}
func (op Call) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case op.callName:
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (op Call) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (op Call) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
