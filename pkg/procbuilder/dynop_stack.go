package procbuilder

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
	"github.com/BondMachineHQ/BondMachine/pkg/bmstack"
)

type DynOpStack struct {
	callName string
	s        int
	sn       string
	opType   uint8
}

func (op DynOpStack) Op_get_name() string {
	return op.callName
}

func (op DynOpStack) Op_get_desc() string {
	switch op.opType {
	case OP_PUSH:
		return "Push a register value into a hardware stack called " + op.sn + " with depth " + strconv.Itoa(op.s)
	case OP_PULL:
		return "Pull a register value from a hardware stack called " + op.sn + " with depth " + strconv.Itoa(op.s)
	}
	return ""
}

func (op DynOpStack) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := ""
	switch op.opType {
	case OP_PUSH:
		result += op.callName + " [" + strconv.Itoa(opBits) + "(Reg)]	// Push a register value into a hardware stack called " + op.sn + " with depth " + strconv.Itoa(op.s) + " [" + strconv.Itoa(opBits+int(arch.R)) + "]\n"
	case OP_PULL:
		result += op.callName + " [" + strconv.Itoa(opBits) + "(Reg)]	// Pull a register value from a hardware stack called " + op.sn + " with depth " + strconv.Itoa(op.s) + " [" + strconv.Itoa(opBits+int(arch.R)) + "]\n"
	}
	return result
}

func (op DynOpStack) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	switch op.opType {
	case OP_PUSH:
		return opBits + int(arch.R) // The bits for the opcode + bits for a register
	case OP_PULL:
		return opBits + int(arch.R) // The bits for the opcode + bits for a register
	}
	return 0
}

func (op DynOpStack) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""

	dataSize := int(arch.Rsize)
	stackName := "regstack" + arch.Tag + "_" + strconv.Itoa(op.s) + op.sn

	// Connect the return stack only once
	suffix := strconv.Itoa(op.s) + op.sn
	if arch.OnlyOne(op.Op_get_name(), []string{"pull" + suffix, "push" + suffix}) {
		result += "	reg [1:0] " + stackName + "SM;\n"
		result += "	localparam	REGST1 = 2'b00,\n"
		result += "			REGST2 = 2'b01,\n"
		result += "			REGST3 = 2'b10,\n"
		result += "			REGST4 = 2'b11;\n"
		result += "\n"
		result += "\treg [" + strconv.Itoa(dataSize-1) + ":0] " + stackName + "senderData;\n"
		result += "\treg " + stackName + "senderWrite;\n"
		result += "\twire " + stackName + "senderAck;\n"
		result += "\n"
		result += "\twire [" + strconv.Itoa(dataSize-1) + ":0] " + stackName + "receiverData;\n"
		result += "\treg " + stackName + "receiverRead;\n"
		result += "\twire " + stackName + "receiverAck;\n"
		result += "\n"
		result += "\twire " + stackName + "empty;\n"
		result += "\twire " + stackName + "full;\n"
		result += "\n"
		result += "\t" + stackName + " " + stackName + "_inst (\n"
		result += "\t\t.clk(clock_signal),\n"
		result += "\t\t.reset(reset_signal),\n"
		result += "\t\t.senderData(" + stackName + "senderData),\n"
		result += "\t\t.senderWrite(" + stackName + "senderWrite),\n"
		result += "\t\t.senderAck(" + stackName + "senderAck),\n"
		result += "\t\t.receiverData(" + stackName + "receiverData),\n"
		result += "\t\t.receiverRead(" + stackName + "receiverRead),\n"
		result += "\t\t.receiverAck(" + stackName + "receiverAck),\n"
		result += "\t\t.empty(" + stackName + "empty),\n"
		result += "\t\t.full(" + stackName + "full)\n"
		result += "\t);\n"
		result += "\n"
		result += "initial begin\n"
		result += "	" + stackName + "SM <= REGST1;\n"
		result += "end\n"

	}

	return result
}

func (op DynOpStack) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (op DynOpStack) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op DynOpStack) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op DynOpStack) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()
	stackName := "regstack" + arch.Tag + "_" + strconv.Itoa(op.s) + op.sn
	dataSize := int(arch.Rsize)
	locationBits := 4
	regNum := 1 << arch.R

	result := ""
	if op.opType == OP_PUSH {
		result += "					" + strings.ToUpper(op.callName) + ": begin\n"
		if arch.R == 1 {
			result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + "])\n"
		} else {
			result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)) + "])\n"
		}
		for i := 0; i < regNum; i++ {
			result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
			result += "							case (" + stackName + "SM)\n"
			result += "							REGST1: begin\n"
			result += "								if (!" + stackName + "senderAck) begin\n"
			result += "								     " + stackName + "senderData[" + strconv.Itoa(dataSize-1) + ":0] <= #1 _" + strings.ToLower(Get_register_name(i)) + "[" + strconv.Itoa(dataSize-1) + ":0];\n"
			result += "								     " + stackName + "senderWrite <= #1 1'b1;\n"
			result += "								     " + stackName + "SM <= REGST2;\n"
			result += "								end\n"
			result += "							end\n"
			result += "							REGST2: begin\n"
			result += "								if (" + stackName + "senderAck) begin\n"
			result += "									" + stackName + "senderWrite <= #1 1'b0;\n"
			result += NextInstruction(conf, arch, 9, "_pc + 1'b1")
			result += "									$display(\"" + strings.ToUpper(op.callName) + " \", current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(locationBits)) + "]);\n"
			result += "									" + stackName + "SM <= REGST1;\n"
			result += "								end\n"
			result += "							end\n"
			result += "							endcase\n"
			result += "						end\n"
		}
		result += "						endcase\n"
		result += "					end\n"
	} else if op.opType == OP_PULL {

		result += "					" + strings.ToUpper(op.callName) + ": begin\n"
		if arch.R == 1 {
			result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + "])\n"
		} else {
			result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)) + "])\n"
		}
		for i := 0; i < regNum; i++ {
			result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
			result += "							if (" + stackName + "receiverAck && " + stackName + "receiverRead) begin\n"
			result += "								" + stackName + "receiverRead <= #1 1'b0;\n"
			result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 " + stackName + "receiverData[" + strconv.Itoa(int(arch.Rsize)-1) + ":0];\n"
			result += NextInstruction(conf, arch, 8, "_pc + 1'b1")
			result += "							end\n"
			result += "							else begin\n"
			result += "								" + stackName + "receiverRead <= #1 1'b1;\n"
			result += "							end\n"
			result += "						end\n"
		}
		result += "						endcase\n"
		result += "					end\n"
	}

	return result
}

func (op DynOpStack) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op DynOpStack) Assembler(arch *Arch, words []string) (string, error) {
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

func (op DynOpStack) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	return result, nil
}

func (op DynOpStack) Simulate(vm *VM, instr string) error {
	value := get_id(instr[:vm.Mach.O])
	if value < len(vm.Mach.Slocs) {
		vm.Pc = uint64(value)
	} else {
		vm.Pc = vm.Pc + 1
	}
	return nil
}

func (op DynOpStack) Generate(arch *Arch) string {
	max_value := 1 << arch.O
	value := rand.Intn(max_value)
	return zeros_prefix(int(arch.O), get_binary(value))
}

func (op DynOpStack) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op DynOpStack) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op DynOpStack) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op DynOpStack) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	suffix := strconv.Itoa(op.s) + op.sn
	if arch.OnlyOne(op.Op_get_name(), []string{"pull" + suffix, "push" + suffix}) {

		s := bmstack.CreateBasicStack()
		s.ModuleName = "regstack" + arch.Tag + "_" + suffix
		s.DataSize = int(arch.Rsize)
		s.Depth = op.s
		s.MemType = "LIFO"
		s.Senders = []string{"sender"}
		s.Receivers = []string{"receiver"}

		r, _ := s.WriteHDL()

		result := r
		return []string{s.ModuleName}, []string{result}
	}
	return []string{}, []string{}
}

func (op DynOpStack) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "j", I_NIL}
	result[0] = newnot
	return result, nil
}

func (op DynOpStack) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (op DynOpStack) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	switch op.opType {
	case OP_CALLO:
		result = append(result, op.callName+"::*--type=reg")
	case OP_CALLA:
		result = append(result, op.callName+"::*--type=reg")
	}
	return result
}
func (op DynOpStack) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case op.callName:
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (op DynOpStack) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (op DynOpStack) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
