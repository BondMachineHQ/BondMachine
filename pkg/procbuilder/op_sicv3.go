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

const (
	SICV3IDLE = uint8(0) + iota
	SICV3WAIT
)

type Sicv3 struct{}

func (op Sicv3) Op_get_name() string {
	return "sicv3"
}

func (op Sicv3) Op_get_desc() string {
	return "Wait for an input change via valid and increments a register"
}

func (op Sicv3) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	inBits := arch.Inputs_bits()
	result := "sicv3 [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(inBits) + "(Input)]	// Wait for an input change via valid and increments a register [" + strconv.Itoa(opBits+int(arch.R)+inBits) + "]\n"
	return result
}

func (op Sicv3) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	inBits := arch.Inputs_bits()
	return opBits + int(arch.R) + int(inBits) // The bits for the opcode + bits for a register + bits for the input
}

func (op Sicv3) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""

	result += "\n\t// Signals for sicv3 instruction\n"
	result += "\treg sicv3_state;\n"
	result += "\tlocalparam SICV3IDLE = 1'b0;\n"
	result += "\tlocalparam SICV3WAIT = 1'b1;\n"
	result += "\n"
	result += "\tinitial begin\n"
	result += "\t\tsicv3_state = SICV3IDLE;\n"
	result += "\tend\n"

	if arch.OnlyOne(op.Op_get_name(), unique["inputrecv"]) {

		opBits := arch.Opcodes_bits()
		romWord := arch.Max_word()

		result += "\n"
		// Data received fro inputs
		for j := 0; j < int(arch.N); j++ {
			result += "\treg " + strings.ToLower(Get_input_name(j)) + "_recv;\n"
		}

		result += "\n"

		for j := 0; j < int(arch.N); j++ {

			objects := make([]string, 1)
			objects[0] = strings.ToLower(Get_input_name(j))

			// Process for data outputs data valid
			result += "\talways @(posedge clock_signal, posedge reset_signal)\n"
			result += "\tbegin\n"

			result += "\t\tif (reset_signal)\n"
			result += "\t\tbegin\n"
			result += "\t\t\t" + strings.ToLower(Get_input_name(j)) + "_recv <= #1 1'b0;\n"
			result += "\t\tend\n"
			result += "\t\telse\n"
			result += "\t\tbegin\n"

			if opBits == 1 {
				result += "\t\t\tcase(current_instruction[" + strconv.Itoa(romWord-1) + "])\n"
			} else {
				result += "\t\t\tcase(current_instruction[" + strconv.Itoa(romWord-1) + ":" + strconv.Itoa(romWord-opBits) + "])\n"
			}

			for _, currOp := range arch.Op {
				result += currOp.Op_instruction_verilog_extra_block(arch, flavor, uint8(4), "input_data_received", objects)
			}

			result += "\t\t\t\tdefault: begin\n"
			result += "\t\t\t\t\tif (!" + strings.ToLower(Get_input_name(j)) + "_valid)\n"
			result += "\t\t\t\t\tbegin\n"
			result += "\t\t\t\t\t\t" + strings.ToLower(Get_input_name(j)) + "_recv <= #1 1'b0;\n"
			result += "\t\t\t\t\tend\n"
			result += "\t\t\t\tend\n"

			result += "\t\t\tendcase\n"

			result += "\t\tend\n"

			result += "\tend\n"
		}
	}

	return result
}

func (op Sicv3) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()
	inpBits := arch.Inputs_bits()

	regNum := 1 << arch.R

	pref := strings.Repeat("\t", 6)

	result := ""
	result += "					SICV3: begin\n"
	if arch.N > 0 {
		if arch.R == 1 {
			result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + "])\n"
		} else {
			result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)) + "])\n"
		}
		for i := 0; i < regNum; i++ {
			result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

			if inpBits == 1 {
				result += "							case (current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-1) + "])\n"
			} else {
				result += "							case (current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)-int(inpBits)) + "])\n"
			}

			for j := 0; j < int(arch.N); j++ {
				result += "							" + strings.ToUpper(Get_input_name(j)) + " : begin\n"

				result += pref + "\t\tif (" + strings.ToLower(Get_input_name(j)) + "_valid) begin\n"
				result += pref + "\t\t\tif (sicv3_state == SICV3IDLE) begin\n"
				result += pref + "\t\t\t\t_" + strings.ToLower(Get_register_name(i)) + " <= #1 " + strconv.Itoa(int(arch.Rsize)) + "'d0;\n"
				result += pref + "\t\t\t\tsicv3_state <= SICV3WAIT;\n"
				result += pref + "\t\t\tend else begin\n"
				result += pref + "\t\t\t\tsicv3_state <= SICV3IDLE;\n"
				result += pref + "\t\t\tend\n"
				result += pref + NextInstruction(conf, arch, 3, "_pc + 1'b1")
				result += pref + "\t\t\t$display(\"SICV3 " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_input_name(j)) + "\");\n"
				result += pref + "\t\tend else begin\n"
				result += pref + "\t\t\tif (sicv3_state == SICV3WAIT) begin\n"
				result += pref + "\t\t\t\t_" + strings.ToLower(Get_register_name(i)) + " <= #1 _" + strings.ToLower(Get_register_name(i)) + " + 1;\n"
				result += pref + "\t\t\tend\n"
				result += pref + "\t\tend\n"
				result += "							end\n"

			}
			result += "							endcase\n"
			result += "						end\n"
		}
		result += "						endcase\n"
	} else {
		result += "						$display(\"NOP\");\n"
	}
	result += "					end\n"
	return result
}

func (op Sicv3) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Sicv3) Assembler(arch *Arch, words []string) (string, error) {
	// "reference": {"support_asm": "ok"}
	opBits := arch.Opcodes_bits()
	inBits := arch.Inputs_bits()
	romWord := arch.Max_word()

	regNum := 2
	regNum = regNum << (arch.R - 1)

	if len(words) != 2 {
		return "", errors.New("wrong arguments number")
	}

	result := ""
	for i := 0; i < regNum; i++ {
		if words[0] == strings.ToLower(Get_register_name(i)) {
			result += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if result == "" {
		return "", errors.New("Unknown register name " + words[0])
	}

	if partial, err := Process_input(words[1], int(arch.N)); err == nil {
		result += zeros_prefix(inBits, partial)
	} else {
		return "", err
	}

	for i := opBits + int(arch.R) + inBits; i < romWord; i++ {
		result += "0"
	}

	return result, nil
}

func (op Sicv3) Disassembler(arch *Arch, instr string) (string, error) {
	// "reference": {"support_disasm": "ok"}
	inBits := arch.Inputs_bits()
	regId := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(regId)) + " "
	inId := get_id(instr[arch.R : int(arch.R)+inBits])
	result += strings.ToLower(Get_input_name(inId))
	return result, nil
}

// Deferred instruction to wait for input to be received when other instructions are executed
func (vm *VM) waitRecvSicv3(inp int) bool {
	if !vm.InputsValid[inp] {
		vm.InputsRecv[inp] = false
		return true
	}
	return false
}

func (op Sicv3) Simulate(vm *VM, instr string) error {
	// "reference": {"support_gosim": "ok"}
	inBits := vm.Mach.Inputs_bits()
	regBits := vm.Mach.R
	reg := get_id(instr[:regBits])
	inp := get_id(instr[regBits : int(regBits)+inBits])

	var sicv3State uint8
	if state, ok := vm.Extra_states["sicv3_state"]; ok {
		sicv3State = state.(uint8)
	} else {
		vm.Extra_states["sicv3_state"] = SICV3IDLE
		sicv3State = SICV3IDLE
	}

	if vm.InputsValid[inp] {
		if sicv3State == SICV3IDLE {
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
			sicv3State = SICV3WAIT
			vm.Extra_states["sicv3_state"] = sicv3State
		} else {
			sicv3State = SICV3IDLE
			vm.Extra_states["sicv3_state"] = sicv3State
		}
		vm.InputsRecv[inp] = true
		vm.Pc = vm.Pc + 1
		// Spawn a deferred instruction to wait for the input to be received
		vm.AddDeferredInstruction("waitRecvSicv3"+strconv.Itoa(inp), func(vm *VM) bool {
			return vm.waitRecvSicv3(inp)
		})
	} else {
		vm.InputsRecv[inp] = false
		if sicv3State == SICV3WAIT {
			switch vm.Mach.Rsize {
			case 8:
				vm.Registers[reg] = vm.Registers[reg].(uint8) + 1
			case 16:
				vm.Registers[reg] = vm.Registers[reg].(uint16) + 1
			case 32:
				vm.Registers[reg] = vm.Registers[reg].(uint32) + 1
			case 64:
				vm.Registers[reg] = vm.Registers[reg].(uint64) + 1
			default:
				return errors.New("go simulation only works on 8,16,32 or 64 bits registers")
			}
		}
	}

	return nil
}

func (op Sicv3) Generate(arch *Arch) string {
	inpbits := arch.Inputs_bits()
	reg_num := 1 << arch.R
	reg := rand.Intn(reg_num)
	inp := rand.Intn(int(arch.N))
	return zeros_prefix(int(arch.R), get_binary(reg)) + zeros_prefix(inpbits, get_binary(inp))
}

func (op Sicv3) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Sicv3) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Sicv3) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Sicv3) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Sicv3) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Sicv3) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Sicv3) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Sicv3) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Sicv3) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	return []UsageNotify{}, errors.New("abstract Assembly not supported for this instruction")
}

func (Op Sicv3) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	opbits := arch.Opcodes_bits()
	inbits := arch.Inputs_bits()
	rom_word := arch.Max_word()

	result := ""

	pref := strings.Repeat("\t", int(level))

	switch blockname {
	case "input_data_received":
		result += pref + "SICV3: begin\n"
		if inbits == 1 {
			result += pref + "\tcase (current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
		} else {
			result += pref + "\tcase (current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(inbits)) + "])\n"
		}

		result += pref + "\t" + strings.ToUpper(objects[0]) + " : begin\n"
		result += pref + "\t\tif (" + strings.ToLower(objects[0]) + "_valid)\n"
		result += pref + "\t\tbegin\n"
		result += pref + "\t\t\t" + strings.ToLower(objects[0]) + "_recv <= #1 1'b1;\n"
		result += pref + "\t\tend else begin\n"
		result += pref + "\t\t\t" + strings.ToLower(objects[0]) + "_recv <= #1 1'b0;\n"
		result += pref + "\t\tend\n"
		result += pref + "\tend\n"

		result += pref + "\tdefault: begin\n"
		result += pref + "\t\tif (!" + strings.ToLower(objects[0]) + "_valid)\n"
		result += pref + "\t\tbegin\n"
		result += pref + "\t\t\t" + strings.ToLower(objects[0]) + "_recv <= #1 1'b0;\n"
		result += pref + "\t\tend\n"
		result += pref + "\tend\n"

		result += pref + "\tendcase\n"

		result += pref + "end\n"
	default:
		result = ""
	}
	return result
}
func (Op Sicv3) HLAssemblerMatch(arch *Arch) []string {
	// "reference": {"support_hlasm": "ok"}
	result := make([]string, 0)
	result = append(result, "sicv3::*--type=reg::*--type=input")
	return result
}
func (Op Sicv3) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "sicv3":
		regNeed := line.Elements[0].GetValue()
		inNeed := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "inputs", Value: inNeed, Op: bmreqs.OpAdd})
	}
	return line, nil
}
func (Op Sicv3) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Sicv3) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
