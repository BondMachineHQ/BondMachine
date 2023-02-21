package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The  opcode
type K2r struct{}


func (op K2r) getKbdName(queueId int) string {
	return "k" + strconv.Itoa(queueId)
}

func (op K2r) Op_get_name() string {
	return "k2r"
}

func (op K2r) Op_get_desc() string {
	return "Get a value from a keyboard queue and put it in a register"
}

func (op K2r) Op_show_assembler(arch *Arch) string {
	qSo := Kbd{}
	opBits := arch.Opcodes_bits()
	queueBits := arch.Shared_bits(qSo.Shr_get_name())
	result := "k2r [" + strconv.Itoa(int(arch.R)) + "(Reg)] [ " + strconv.Itoa(int(queueBits)) + " (Shared Kbd)]	// " + op.Op_get_desc() + " [" + strconv.Itoa(opBits+int(arch.R)+queueBits) + "]\n"
	return result
}

func (op K2r) Op_get_instruction_len(arch *Arch) int {
	qSo := Kbd{}
	opBits := arch.Opcodes_bits()
	queueBits := arch.Shared_bits(qSo.Shr_get_name())
	return opBits + int(arch.R) + int(queueBits) // The bits for the opcode + bits for a register + bits queues
}

func (op K2r) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pName string) string {
	qSo := Kbd{}
	queueBits := arch.Shared_bits(qSo.Shr_get_name())
	queueNum := arch.Shared_num(qSo.Shr_get_name())

	result := ""
	if arch.OnlyOne(op.Op_get_name(), []string{"r2t", "t2r", "k2r", "r2q"}) {
		result += "	reg stackqueueSM;\n"
	}
	if arch.OnlyOne(op.Op_get_name(), []string{"r2q", "k2r"}) {
		result += "	localparam "
		for i := 0; i < queueNum; i++ {
			result += strings.ToUpper(op.getKbdName(i)) + "=" + strconv.Itoa(int(queueBits)) + "'d" + strconv.Itoa(i)
			if i < queueNum-1 {
				result += ",\n"
			} else {
				result += ";\n"
			}
		}
	}
	return result
}

func (Op K2r) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""
	return result
}

func (op K2r) Op_instruction_verilog_state_machine(arch *Arch, flavor string) string {
	qSo := Kbd{}
	queueBits := arch.Shared_bits(qSo.Shr_get_name())
	queueNum := arch.Shared_num(qSo.Shr_get_name())
	rom_word := arch.Max_word()
	opBits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					K2R: begin\n"
	if queueNum > 0 {
		if arch.R == 1 {
			result += "						case (rom_value[" + strconv.Itoa(rom_word-opBits-1) + "])\n"
		} else {
			result += "						case (rom_value[" + strconv.Itoa(rom_word-opBits-1) + ":" + strconv.Itoa(rom_word-opBits-int(arch.R)) + "])\n"
		}
		for i := 0; i < reg_num; i++ {
			result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

			if queueBits == 1 {
				result += "							case (rom_value[" + strconv.Itoa(rom_word-opBits-queueBits-1) + "])\n"
			} else {
				result += "							case (rom_value[" + strconv.Itoa(rom_word-opBits-queueBits-1) + ":" + strconv.Itoa(rom_word-opBits-int(arch.R)-int(queueBits)) + "])\n"
			}

			for j := 0; j < queueNum; j++ {
				result += "							" + strings.ToUpper(op.getKbdName(j)) + " : begin\n"
				result += "								if (" + strings.ToLower(op.getKbdName((j))) + "receiverAck && " + strings.ToLower(op.getKbdName(j)) + "receiverRead) begin\n"
				result += "								     " + strings.ToLower(op.getKbdName(j)) + "receiverRead <= #1 1'b0;\n"
				result += "								     _" + strings.ToLower(Get_register_name(i)) + "[" + strconv.Itoa(int(arch.Rsize)-1) + ":0] <= #1 " + strings.ToLower(op.getKbdName(j)) + "receiverData[" + strconv.Itoa(int(arch.Rsize)-1) + ":0];\n"
				result += "								       _pc <= #1 _pc + 1'b1 ;\n"
				result += "								end\n"
				result += "								else begin\n"
				result += "								       " + strings.ToLower(op.getKbdName(j)) + "receiverRead <= #1 1'b1;\n"
				result += "								end\n"
				result += "								$display(\"T2R " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(op.getKbdName(j)) + "\");\n"
				result += "							end\n"

			}
			result += "							endcase\n"
			result += "						end\n"
		}
		result += "						endcase\n"
	} else {
		result += "						$display(\"NOP\");\n"
		result += "						_pc <= #1 _pc + 1'b1 ;\n"
	}
	result += "					end\n"
	return result

}

func (op K2r) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	result := ""
	return result
}

func (op K2r) Assembler(arch *Arch, words []string) (string, error) {
	opBits := arch.Opcodes_bits()
	qSo := Kbd{}
	queueNum := arch.Shared_num(qSo.Shr_get_name())
	queueBits := arch.Shared_bits(qSo.Shr_get_name())
	shortName := qSo.Shortname()
	romWord := arch.Max_word()

	regNum := 1 << arch.R

	if len(words) != 2 {
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

	if partial, err := Process_shared(shortName, words[1], queueNum); err == nil {
		result += zeros_prefix(queueBits, partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opBits + int(arch.R) + queueBits; i < romWord; i++ {
		result += "0"
	}

	return result, nil
}

func (op K2r) Disassembler(arch *Arch, instr string) (string, error) {
	kSo := Kbd{}
	queueBits := arch.Shared_bits(kSo.Shr_get_name())
	shortname := kSo.Shortname()
	regId := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(regId)) + " "
	stId := get_id(instr[arch.R : int(arch.R)+queueBits])
	result += shortname + strconv.Itoa(stId)
	return result, nil
}

// The simulation does nothing
func (op K2r) Simulate(vm *VM, instr string) error {
	// TODO

	reg_bits := vm.Mach.R
	regPay := get_id(instr[:reg_bits])
	posS := instr[reg_bits : reg_bits+8]

	pos := uint8(get_id(posS))
	payload := vm.Registers[regPay].(uint8)

	cmd := make([]byte, 0)

	cmd = append(cmd, byte(vm.CpID))
	cmd = append(cmd, byte(pos))
	cmd = append(cmd, byte(payload))

	vm.CmdChan <- cmd

	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op K2r) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op K2r) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op K2r) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op K2r) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op K2r) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	result := ""
	return result
}

func (Op K2r) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op K2r) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op K2r) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	// TODO Partial
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "k2r", I_NIL}
	result[0] = newnot
	return result, nil
}

func (Op K2r) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op K2r) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "kbd::*--type=reg")
	result = append(result, "k2r::*--type=reg::*--type=somov--sotype=k")
	result = append(result, "mov::*--type=reg::*--type=somov--sotype=k")
	return result
}
func (Op K2r) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "k2r":
		regVal := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		soVal := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "sos", Value: soVal, Op: bmreqs.OpAdd})
		return line, nil
	case "kbd":
		regVal := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		soVal := "q0" // Pull implicitly uses queue 0
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "sos", Value: soVal, Op: bmreqs.OpAdd})
		if regVal != "" && soVal != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("k2r")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regVal)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(soVal)
			newArg1.BasmMeta = newArg1.SetMeta("type", "somov")
			newArg1.BasmMeta = newArg1.SetMeta("sotype", "q")
			newArgs[1] = newArg1
			newLine.Elements = newArgs
			return newLine, nil
		}
	case "mov":
		regVal := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		soVal := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "sos", Value: soVal, Op: bmreqs.OpAdd})
		if regVal != "" && soVal != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("k2r")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regVal)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(soVal)
			newArg1.BasmMeta = newArg1.SetMeta("type", "somov")
			newArg1.BasmMeta = newArg1.SetMeta("sotype", "q")
			newArgs[1] = newArg1
			newLine.Elements = newArgs
			return newLine, nil
		}
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op K2r) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
