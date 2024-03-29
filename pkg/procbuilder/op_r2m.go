package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type R2m struct{}

func (op R2m) Op_get_name() string {
	return "r2m"
}

func (op R2m) Op_get_desc() string {
	return "Copy a register value to the ram"
}

func (op R2m) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "r2m [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.L)) + "(RAM address)]	// " + op.Op_get_desc() + " [" + strconv.Itoa(opbits+int(arch.R)+int(arch.L)) + "]\n"
	return result
}

func (op R2m) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.L) // The bits for the opcode + bits for a register + bits for RAM address
}

func (op R2m) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	if arch.OnlyOne(op.Op_get_name(), []string{"r2m", "r2mri"}) {
		result += "\treg [" + strconv.Itoa(int(arch.L)-1) + ":0] addr_ram_to_mem;\n"
		result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] ram_din_i;\n"
		result += "\treg wr_int_ram;\n"
	}
	return result
}

func (Op R2m) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""
	//result += "\t\t\twr_int_ram <= #1 1'b0;\n"
	//result += "\t\t\taddr_ram_to_mem <= #1 'b0;\n"
	//result += "\t\t\tram_din_i <= #1 'b0;\n"
	return result
}

func (op R2m) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R
	result := ""
	result += "					R2M: begin\n"
	result += "						if (wr_int_ram == 0) begin\n"
	result += "							wr_int_ram <= #1 1'b1;\n"
	if arch.R == 1 {
		result += "							case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "							case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "							" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "								addr_ram_to_mem <= current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.L)) + "];\n"
		result += "								ram_din_i <= _" + strings.ToLower(Get_register_name(i)) + ";\n"
		result += "								$display(\"R2M " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "							end\n"
	}
	result += "							endcase\n"
	result += "						end\n"
	result += "						else begin\n"
	result += "							wr_int_ram <= #1 1'b0;\n"
	result += NextInstruction(nil, arch, 7, "_pc + 1'b1")
	result += "						end\n"
	result += "					end\n"

	return result
}

func (op R2m) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	result := ""

	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	// The ram is enabled if any of the opcodes is active
	if arch.OnlyOne(op.Op_get_name(), []string{"r2mri", "r2m", "m2r", "m2rri"}) {
		result += "\tassign ram_en = 1'b1;\n"
	}

	// The ram write signals
	if arch.OnlyOne(op.Op_get_name(), []string{"r2mri", "r2m"}) {
		result += "\tassign ram_din = ram_din_i;\n"
		result += "\tassign ram_wren = wr_int_ram;\n"
	}

	ramAddr := ""
	if arch.HasAny([]string{"r2mri", "r2m"}) {
		ramAddr += "addr_ram_to_mem"
	}

	if arch.HasOp("m2r") {
		ramAddr = " (current_instruction[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "]==M2R) ? addr_ram_m2r : " + ramAddr
	}

	if arch.HasOp("m2rri") {
		ramAddr = " (current_instruction[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "]==M2RRI) ? addr_ram_m2rri: " + ramAddr
	}

	if arch.Modes[0] == "hy" || arch.Modes[0] == "vn" {
		ramAddr = " (exec_mode == 1'b1 && vn_state == FETCH) ? _pc : " + ramAddr
	}

	if arch.OnlyOne(op.Op_get_name(), []string{"r2mri", "r2m", "m2r", "m2rri"}) {
		result += "\tassign ram_addr = " + ramAddr + ";\n"
	}

	return result
}

func (op R2m) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	ramDepth := int(arch.L)

	reg_num := 1 << arch.R

	if len(words) != 2 {
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

	if partial, err := Process_number(words[1]); err == nil {
		result += zeros_prefix(ramDepth, partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opbits + int(arch.R) + ramDepth; i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op R2m) Disassembler(arch *Arch, instr string) (string, error) {
	ramDepth := int(arch.L)
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	value := get_id(instr[arch.R : int(arch.R)+ramDepth])
	result += strconv.Itoa(value)
	return result, nil
}

// The simulation does nothing
func (op R2m) Simulate(vm *VM, instr string) error {
	// TODO
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op R2m) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op R2m) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op R2m) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op R2m) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op R2m) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	result := ""
	//result += "\t\t\t\twr_int_ram <= #1 1'b0;\n"
	return result
}

func (Op R2m) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op R2m) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op R2m) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	// TODO Partial
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "r2m", I_NIL}
	result[0] = newnot
	return result, nil
}

func (Op R2m) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op R2m) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "r2m::*--type=reg::*--type=ram--ramaddressing=immediate")
	result = append(result, "r2m::*--type=reg::*--type=ram--ramaddressing=symbol")
	result = append(result, "mov::*--type=ram--ramaddressing=immediate::*--type=reg")
	result = append(result, "mov::*--type=ram--ramaddressing=symbol::*--type=reg")
	return result
}
func (Op R2m) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "r2m":
		regVal := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		return line, nil
	case "mov":
		regVal := line.Elements[1].GetValue()
		addressing := line.Elements[0].GetMeta("ramaddressing")
		switch addressing {
		case "immediate":
			location := line.Elements[0].GetMeta("location")
			rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
			if regVal != "" && location != "" {
				newLine := new(bmline.BasmLine)
				newOp := new(bmline.BasmElement)
				newOp.SetValue("r2m")
				newLine.Operation = newOp
				newArgs := make([]*bmline.BasmElement, 2)
				newArg0 := new(bmline.BasmElement)
				newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
				newArg0.SetValue(regVal)
				newArgs[0] = newArg0
				newArg1 := new(bmline.BasmElement)
				newArg1.SetValue(location)
				newArg1.BasmMeta = newArg1.SetMeta("type", "number")
				newArgs[1] = newArg1
				newLine.Elements = newArgs
				return newLine, nil
			}
		case "symbol":
			// The mov is a r2m, the symbol is kept unchanged because it will be resolved by the symbol resolver
			symbol := line.Elements[0].GetMeta("symbol")
			rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
			if regVal != "" && symbol != "" {
				newLine := new(bmline.BasmLine)
				newOp := new(bmline.BasmElement)
				newOp.SetValue("r2m")
				newLine.Operation = newOp
				newArgs := make([]*bmline.BasmElement, 2)
				newArg0 := new(bmline.BasmElement)
				newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
				newArg0.SetValue(regVal)
				newArgs[0] = newArg0
				newArg1 := new(bmline.BasmElement)
				newArg1.SetValue(symbol)
				newArg1.BasmMeta = newArg1.SetMeta("type", "ram")
				newArg1.BasmMeta = newArg1.SetMeta("ramaddressing", "symbol")
				newArg1.BasmMeta = newArg1.SetMeta("symbol", symbol)
				newArgs[1] = newArg1
				newLine.Elements = newArgs
				return newLine, nil
			}
		default:
			return nil, errors.New("unknown addressing mode")
		}
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op R2m) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op R2m) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
