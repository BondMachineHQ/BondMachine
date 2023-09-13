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

type Jria struct{}

func (op Jria) Op_get_name() string {
	return "jria"
}

func (op Jria) Op_get_desc() string {
	return "Register indirect Jump to a program location on RAM"
}

func (op Jria) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := "jria [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Jump to the location contained in reg [" + strconv.Itoa(opBits+int(arch.R)) + "]\n"
	return result
}

func (op Jria) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	return opBits + int(arch.R) // The bits for the opcode + bits for a register
}

func (op Jria) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pName string) string {
	return ""
}

func (Op Jria) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}
func (Op Jria) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Jria) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Jria) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					JRIA: begin\n"
	if arch.Modes[0] == "hy" {
		result += "						exec_mode <= #1 1'b1;\n"
	}
	result += "						vn_state <= FETCH;\n"

	if arch.R == 1 {
		result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + "])\n"
	} else {
		result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "							_pc <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
		result += "							$display(\"JRIA " + strings.ToUpper(Get_register_name(i)) + "\");\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "					end\n"
	return result
}

func (op Jria) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Jria) Assembler(arch *Arch, words []string) (string, error) {
	opBits := arch.Opcodes_bits()
	romWord := arch.Max_word()

	regNum := 1 << arch.R

	if len(words) != 1 {
		return "", Prerror{"Wrong arguments number"}
	}

	result := ""
	for i := 0; i < regNum; i++ {
		if words[0] == strings.ToLower(Get_register_name(i)) {
			result += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if result == "" {
		return "", Prerror{"Unknown register name " + words[0]}
	}

	for i := opBits + int(arch.R); i < romWord; i++ {
		result += "0"
	}

	return result, nil
}

func (op Jria) Disassembler(arch *Arch, instr string) (string, error) {
	regId := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(regId)) + " "
	return result, nil
}

func (op Jria) Simulate(vm *VM, instr string) error {
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

func (op Jria) Generate(arch *Arch) string {
	reg_num := 1 << arch.R
	reg := rand.Intn(reg_num)
	return zeros_prefix(int(arch.R), get_binary(reg))
}

func (op Jria) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Jria) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Jria) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op Jria) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Jria) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq, types := Sequence_to_0(words[0])
	if len(seq) > 0 && types == O_REGISTER {

		result := make([]UsageNotify, 2)
		newnot0 := UsageNotify{C_OPCODE, "jria", I_NIL}
		result[0] = newnot0
		newnot1 := UsageNotify{C_REGSIZE, S_NIL, len(seq)}
		result[1] = newnot1

		return result, nil
	}

	return []UsageNotify{}, errors.New("Wrong register")

}

func (Op Jria) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}

func (Op Jria) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "jria::*--type=reg")
	result = append(result, "jmp::*--type=ram--ramaddressing=register")
	return result
}

func (Op Jria) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "jmp":
		regNeed := line.Elements[0].GetMeta("ramregister")
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		line.Operation.SetValue("jria")
		line.Elements[0].SetValue(regNeed)
		line.Elements[0].RmMeta("ramregister")
		line.Elements[0].RmMeta("ramaddressing")
		line.Elements[0].BasmMeta = line.Elements[0].SetMeta("type", "reg")
		return line, nil
	case "jria":
		regNeed := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Jria) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Jria) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	switch line.Operation.GetValue() {
	case "jria":
		regNeed := line.Elements[0].GetValue()
		if regNeed != "" {
			var meta *bmmeta.BasmMeta
			meta = meta.SetMeta("inv", regNeed)
			return meta, nil
		}
	}
	return nil, nil
}
