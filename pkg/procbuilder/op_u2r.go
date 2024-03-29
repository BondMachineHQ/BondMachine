package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The U2r opcode
type U2r struct{}

func (op U2r) getUartName(queueId int) string {
	return "u" + strconv.Itoa(queueId)
}

func (op U2r) Op_get_name() string {
	return "u2r"
}

func (op U2r) Op_get_desc() string {
	return "Get a value from a shared UART and put it in a register"
}

func (op U2r) Op_show_assembler(arch *Arch) string {
	uSo := Uart{}
	opBits := arch.Opcodes_bits()
	uartBits := arch.Shared_bits(uSo.Shr_get_name())
	result := "u2r [" + strconv.Itoa(int(arch.R)) + "(Reg)] [ " + strconv.Itoa(int(uartBits)) + " (Shared Uart)]	// " + op.Op_get_desc() + " [" + strconv.Itoa(opBits+int(arch.R)+uartBits) + "]\n"
	return result
}

func (op U2r) Op_get_instruction_len(arch *Arch) int {
	uSo := Uart{}
	opBits := arch.Opcodes_bits()
	queueBits := arch.Shared_bits(uSo.Shr_get_name())
	return opBits + int(arch.R) + int(queueBits) // The bits for the opcode + bits for a register + bits queues
}

func (op U2r) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pName string) string {
	uSo := Uart{}
	queueBits := arch.Shared_bits(uSo.Shr_get_name())
	queueNum := arch.Shared_num(uSo.Shr_get_name())

	result := ""
	if arch.OnlyOne(op.Op_get_name(), []string{"r2t", "t2r", "q2r", "r2q", "r2u", "u2r", "k2r"}) {
		result += "	reg stackqueueSM;\n"
	}
	if arch.OnlyOne(op.Op_get_name(), []string{"r2u", "u2r"}) {
		result += "	localparam "
		for i := 0; i < queueNum; i++ {
			result += strings.ToUpper(op.getUartName(i)) + "=" + strconv.Itoa(int(queueBits)) + "'d" + strconv.Itoa(i)
			if i < queueNum-1 {
				result += ",\n"
			} else {
				result += ";\n"
			}
		}
	}
	return result
}

func (Op U2r) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""
	return result
}

func (op U2r) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	uSo := Uart{}
	uartBits := arch.Shared_bits(uSo.Shr_get_name())
	uartNum := arch.Shared_num(uSo.Shr_get_name())
	rom_word := arch.Max_word()
	opBits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					U2R: begin\n"
	if uartNum > 0 {
		if arch.R == 1 {
			result += "						case (current_instruction[" + strconv.Itoa(rom_word-opBits-1) + "])\n"
		} else {
			result += "						case (current_instruction[" + strconv.Itoa(rom_word-opBits-1) + ":" + strconv.Itoa(rom_word-opBits-int(arch.R)) + "])\n"
		}
		for i := 0; i < reg_num; i++ {
			result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

			if uartBits == 1 {
				result += "							case (current_instruction[" + strconv.Itoa(rom_word-opBits-uartBits-1) + "])\n"
			} else {
				result += "							case (current_instruction[" + strconv.Itoa(rom_word-opBits-uartBits-1) + ":" + strconv.Itoa(rom_word-opBits-int(arch.R)-int(uartBits)) + "])\n"
			}

			for j := 0; j < uartNum; j++ {
				result += "							" + strings.ToUpper(op.getUartName(j)) + " : begin\n"
				result += "								if (" + strings.ToLower(op.getUartName((j))) + "receiverAck && " + strings.ToLower(op.getUartName(j)) + "receiverRead) begin\n"
				result += "								     " + strings.ToLower(op.getUartName(j)) + "receiverRead <= #1 1'b0;\n"
				result += "								     _" + strings.ToLower(Get_register_name(i)) + "[" + strconv.Itoa(int(arch.Rsize)-1) + ":0] <= #1 " + strings.ToLower(op.getUartName(j)) + "receiverData[" + strconv.Itoa(int(arch.Rsize)-1) + ":0];\n"
				result += NextInstruction(conf, arch, 8, "_pc + 1'b1")
				result += "								end\n"
				result += "								else begin\n"
				result += "								       " + strings.ToLower(op.getUartName(j)) + "receiverRead <= #1 1'b1;\n"
				result += "								end\n"
				result += "								$display(\"T2R " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(op.getUartName(j)) + "\");\n"
				result += "							end\n"

			}
			result += "							endcase\n"
			result += "						end\n"
		}
		result += "						endcase\n"
	} else {
		result += "						$display(\"NOP\");\n"
		result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
	}
	result += "					end\n"
	return result

}

func (op U2r) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	result := ""
	return result
}

func (op U2r) Assembler(arch *Arch, words []string) (string, error) {
	opBits := arch.Opcodes_bits()
	uSo := Uart{}
	uartNum := arch.Shared_num(uSo.Shr_get_name())
	uartBits := arch.Shared_bits(uSo.Shr_get_name())
	shortName := uSo.Shortname()
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

	if partial, err := Process_shared(shortName, words[1], uartNum); err == nil {
		result += zeros_prefix(uartBits, partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opBits + int(arch.R) + uartBits; i < romWord; i++ {
		result += "0"
	}

	return result, nil
}

func (op U2r) Disassembler(arch *Arch, instr string) (string, error) {
	uSo := Uart{}
	uartBits := arch.Shared_bits(uSo.Shr_get_name())
	shortname := uSo.Shortname()
	regId := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(regId)) + " "
	uId := get_id(instr[arch.R : int(arch.R)+uartBits])
	result += shortname + strconv.Itoa(uId)
	return result, nil
}

// The simulation does nothing
func (op U2r) Simulate(vm *VM, instr string) error {
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
func (op U2r) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op U2r) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op U2r) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op U2r) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op U2r) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	result := ""
	return result
}

func (Op U2r) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op U2r) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op U2r) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	// TODO Partial
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "u2r", I_NIL}
	result[0] = newnot
	return result, nil
}

func (Op U2r) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op U2r) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "fromuart::*--type=reg")
	result = append(result, "u2r::*--type=reg::*--type=somov--sotype=u")
	result = append(result, "mov::*--type=reg::*--type=somov--sotype=u")
	return result
}
func (Op U2r) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "u2r":
		regVal := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		soVal := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "sos", Value: soVal, Op: bmreqs.OpAdd})
		return line, nil
	case "fromuart":
		regVal := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		soVal := "u0" // Fromuart is always u0
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "sos", Value: soVal, Op: bmreqs.OpAdd})
		if regVal != "" && soVal != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("u2r")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regVal)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(soVal)
			newArg1.BasmMeta = newArg1.SetMeta("type", "somov")
			newArg1.BasmMeta = newArg1.SetMeta("sotype", "u")
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
			newOp.SetValue("u2r")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regVal)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(soVal)
			newArg1.BasmMeta = newArg1.SetMeta("type", "somov")
			newArg1.BasmMeta = newArg1.SetMeta("sotype", "u")
			newArgs[1] = newArg1
			newLine.Elements = newArgs
			return newLine, nil
		}
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op U2r) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op U2r) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
