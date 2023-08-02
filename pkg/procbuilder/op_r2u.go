package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The R2u opcode
type R2u struct{}

func (op R2u) getUartName(uartId int) string {
	return "u" + strconv.Itoa(uartId)
}

func (op R2u) Op_get_name() string {
	return "r2u"
}

func (op R2u) Op_get_desc() string {
	return "Copy a register value to a shared UART"
}

func (op R2u) Op_show_assembler(arch *Arch) string {
	uSo := Uart{}
	opBits := arch.Opcodes_bits()
	uartBits := arch.Shared_bits(uSo.Shr_get_name())
	result := "r2u [" + strconv.Itoa(int(arch.R)) + "(Reg)] [ " + strconv.Itoa(int(uartBits)) + " (Shared UART)]	// " + op.Op_get_desc() + " [" + strconv.Itoa(opBits+int(arch.R)+uartBits) + "]\n"
	return result
}

func (op R2u) Op_get_instruction_len(arch *Arch) int {
	uSo := Uart{}
	opBits := arch.Opcodes_bits()
	uartBits := arch.Shared_bits(uSo.Shr_get_name())
	return opBits + int(arch.R) + int(uartBits) // The bits for the opcode + bits for a register + bits stacks
}

func (op R2u) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pName string) string {
	uSo := Uart{}
	uartBits := arch.Shared_bits(uSo.Shr_get_name())
	uartNum := arch.Shared_num(uSo.Shr_get_name())

	result := ""
	if arch.OnlyOne(op.Op_get_name(), []string{"r2t", "t2r", "q2r", "r2q", "r2u", "u2r", "k2r"}) {
		result += "	reg stackqueueSM;\n"
	}
	if arch.OnlyOne(op.Op_get_name(), []string{"r2u", "u2r"}) {
		result += "	localparam "
		for i := 0; i < uartNum; i++ {
			result += strings.ToUpper(op.getUartName(i)) + "=" + strconv.Itoa(int(uartBits)) + "'d" + strconv.Itoa(i)
			if i < uartNum-1 {
				result += ",\n"
			} else {
				result += ";\n"
			}
		}
	}
	return result
}

func (Op R2u) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""
	return result
}

func (op R2u) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	uSo := Uart{}
	uartBits := arch.Shared_bits(uSo.Shr_get_name())
	uartNum := arch.Shared_num(uSo.Shr_get_name())
	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()

	regNum := 1 << arch.R

	result := ""
	result += "					R2U: begin\n"
	if uartNum > 0 {
		if arch.R == 1 {
			result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + "])\n"
		} else {
			result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)) + "])\n"
		}
		for i := 0; i < regNum; i++ {
			result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

			if uartBits == 1 {
				result += "							case (current_instruction[" + strconv.Itoa(romWord-opBits-uartBits-1) + "])\n"
			} else {
				result += "							case (current_instruction[" + strconv.Itoa(romWord-opBits-uartBits-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)-int(uartBits)) + "])\n"
			}

			for j := 0; j < uartNum; j++ {
				result += "							" + strings.ToUpper(op.getUartName(j)) + " : begin\n"
				result += "								case (stackqueueSM)\n"
				result += "								   1'b0: begin\n"
				result += "								     if (!" + strings.ToLower(op.getUartName((j))) + "senderAck) begin\n"
				result += "								     " + strings.ToLower(op.getUartName(j)) + "senderData[" + strconv.Itoa(int(arch.Rsize)-1) + ":0] <= #1 _" + strings.ToLower(Get_register_name(i)) + "[" + strconv.Itoa(int(arch.Rsize)-1) + ":0];\n"
				result += "								     " + strings.ToLower(op.getUartName(j)) + "senderWrite <= #1 1'b1;\n"
				result += "								     stackqueueSM <= 1'b1;\n"
				result += "								     end\n"
				result += "								   end\n"
				result += "								   1'b1: begin\n"
				result += "								     if (" + strings.ToLower(op.getUartName((j))) + "senderAck) begin\n"
				result += "								       " + strings.ToLower(op.getUartName(j)) + "senderWrite <= #1 1'b0;\n"
				result += NextInstruction(conf, arch, 8, "_pc + 1'b1")
				result += "								       stackqueueSM <= 1'b0;\n"
				result += "								     end\n"
				result += "								   end\n"
				result += "								endcase\n"
				result += "								$display(\"R2U " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(op.getUartName(j)) + "\");\n"
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

func (op R2u) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	result := ""
	return result
}

func (op R2u) Assembler(arch *Arch, words []string) (string, error) {
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

func (op R2u) Disassembler(arch *Arch, instr string) (string, error) {
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
func (op R2u) Simulate(vm *VM, instr string) error {
	// TODO

	regBits := vm.Mach.R
	regPay := get_id(instr[:regBits])
	posS := instr[regBits : regBits+8]

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
func (op R2u) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op R2u) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op R2u) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op R2u) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op R2u) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	result := ""
	return result
}

func (Op R2u) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op R2u) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op R2u) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	// TODO Partial
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "r2u", I_NIL}
	result[0] = newnot
	return result, nil
}

func (Op R2u) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op R2u) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "touart::*--type=reg")
	result = append(result, "r2u::*--type=reg::*--type=somov--sotype=u")
	result = append(result, "mov::*--type=somov--sotype=u::*--type=reg")
	return result
}
func (Op R2u) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "r2u":
		regVal := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		soVal := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "sos", Value: soVal, Op: bmreqs.OpAdd})
		return line, nil
	case "touart":
		regVal := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		soVal := "u0" // Push implicitely uses the first stack
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "sos", Value: soVal, Op: bmreqs.OpAdd})
		if regVal != "" && soVal != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("r2u")
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
		regVal := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		soVal := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "sos", Value: soVal, Op: bmreqs.OpAdd})
		if regVal != "" && soVal != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("r2u")
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
func (Op R2u) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op R2u) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
