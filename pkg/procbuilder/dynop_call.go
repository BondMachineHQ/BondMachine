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
	result := ""

	// Size of the return stack
	dataSize := 0
	switch arch.Modes[0] {
	case "ha":
		dataSize = int(arch.O)
	case "vn":
		dataSize = int(arch.L)
	case "hy":
		// +1 because the stack memorizes the origin of the return (ROM or RAM)
		if arch.O > arch.L {
			dataSize = int(arch.O) + 1
		} else {
			dataSize = int(arch.L) + 1
		}
	}
	stackName := "restack" + arch.Tag + "_" + strconv.Itoa(op.s) + op.sn

	// Connect the return stack only once
	suffix := strconv.Itoa(op.s) + op.sn
	if arch.OnlyOne(op.Op_get_name(), []string{"callo" + suffix, "calla" + suffix, "ret" + suffix}) {
		result += "	reg [1:0] " + stackName + "SM;\n"
		result += "	localparam	CALL1 = 2'b00,\n"
		result += "			CALL2 = 2'b01,\n"
		result += "			CALL3 = 2'b10,\n"
		result += "			CALL4 = 2'b11;\n"
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
		result += "\t\t.clk(clk),\n"
		result += "\t\t.reset(reset),\n"
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
		result += "	" + stackName + "SM <= CALL1;\n"
		result += "end\n"

	}

	return result
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
	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()
	stackName := "restack" + arch.Tag + "_" + strconv.Itoa(op.s) + op.sn
	locationBits := arch.O

	dataSize := 0
	switch arch.Modes[0] {
	case "ha":
		locationBits = arch.O
		dataSize = int(arch.O)
	case "vn":
		locationBits = arch.L
		dataSize = int(arch.L)
	case "hy":
		if arch.O > arch.L {
			locationBits = arch.O
			dataSize = int(arch.O) + 1
		} else {
			locationBits = arch.L
			dataSize = int(arch.L) + 1
		}
	}

	result := ""
	if op.opType == OP_CALLA || op.opType == OP_CALLO {
		result += "					" + strings.ToUpper(op.callName) + ": begin\n"
		result += "						case (" + stackName + "SM)\n"
		result += "						CALL1: begin\n"
		result += "							if (!" + stackName + "senderAck) begin\n"
		if arch.Modes[0] == "hy" {
			result += "							     " + stackName + "senderData[" + strconv.Itoa(dataSize-1) + ":0] <= #1 { exec_mode, _pc + 1 };\n"
		} else if arch.Modes[0] == "ha" {
			result += "							     " + stackName + "senderData[" + strconv.Itoa(dataSize-1) + ":0] <= #1 _pc + 1;\n"
		}
		result += "							     " + stackName + "senderWrite <= #1 1'b1;\n"
		result += "							     " + stackName + "SM <= CALL2;\n"
		result += "							end\n"
		result += "						end\n"
		result += "						CALL2: begin\n"
		result += "							if (" + stackName + "senderAck) begin\n"
		result += "								" + stackName + "senderWrite <= #1 1'b0;\n"
		if arch.Modes[0] == "hy" {
			if op.opType == OP_CALLO {
				result += "								exec_mode <= #1 1'b0;\n"
			} else {
				result += "								exec_mode <= #1 1'b1;\n"
			}
		}
		if locationBits == 1 {
			result += "								_pc <= #1 current_instruction[" + strconv.Itoa(romWord-opBits-1) + "];\n"
			result += "								$display(\"" + strings.ToUpper(op.callName) + " \", current_instruction[" + strconv.Itoa(romWord-opBits-1) + "]);\n"
		} else {
			result += "								_pc <= #1 current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(locationBits)) + "];\n"
			result += "								$display(\"" + strings.ToUpper(op.callName) + " \", current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(locationBits)) + "]);\n"
		}
		result += "								" + stackName + "SM <= CALL1;\n"
		result += "							end\n"
		result += "						end\n"
		result += "						endcase\n"
		result += "					end\n"
	} else {
		result += "					" + strings.ToUpper(op.callName) + ": begin\n"
		result += "						if (" + stackName + "receiverAck && " + stackName + "receiverRead) begin\n"
		result += "							" + stackName + "receiverRead <= #1 1'b0;\n"
		result += "							_pc[" + strconv.Itoa(int(locationBits)-1) + ":0] <= #1 " + stackName + "receiverData[" + strconv.Itoa(int(locationBits)-1) + ":0];\n"
		if arch.Modes[0] == "hy" {
			result += "							exec_mode <= #1 " + stackName + "receiverData[" + strconv.Itoa(int(locationBits)) + "];\n"
		}
		if arch.Modes[0] == "hy" || arch.Modes[0] == "vn" {
			result += "							vn_state <= FETCH;\n"
		}
		result += "						end\n"
		result += "						else begin\n"
		result += "							" + stackName + "receiverRead <= #1 1'b1;\n"
		result += "						end\n"
		result += "					end\n"
	}

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
	suffix := strconv.Itoa(op.s) + op.sn
	if arch.OnlyOne(op.Op_get_name(), []string{"callo" + suffix, "calla" + suffix, "ret" + suffix}) {

		s := bmstack.CreateBasicStack()
		s.ModuleName = "restack" + arch.Tag + "_" + suffix
		switch arch.Modes[0] {
		case "ha":
			s.DataSize = int(arch.O)
		case "vn":
			s.DataSize = int(arch.L)
		case "hy":
			// +1 because the stack memorizes the origin of the return (ROM or RAM)
			if arch.O > arch.L {
				s.DataSize = int(arch.O) + 1
			} else {
				s.DataSize = int(arch.L) + 1
			}
		}
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
