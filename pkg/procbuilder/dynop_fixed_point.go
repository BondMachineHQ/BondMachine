package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

const (
	FPPUT = uint8(0) + iota
	FPGET
)

// The FixedPoint opcode is both a basic instruction and a template for other instructions.
type FixedPoint struct {
	fpName   string
	s        int
	f        int
	opType   uint8
	pipeline *uint8
}

func (op FixedPoint) Op_get_name() string {
	return op.fpName
}

func (op FixedPoint) Op_get_desc() string {
	return "FixedPoint dynamical instruction " + op.fpName
}

func (op FixedPoint) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := op.fpName + " [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to the sum of its value with another register [" + strconv.Itoa(opBits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op FixedPoint) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	return opBits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op FixedPoint) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pName string) string {
	result := ""
	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + op.fpName + "_" + arch.Tag + "_input_a;\n"
	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + op.fpName + "_" + arch.Tag + "_input_b;\n"
	result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + op.fpName + "_" + arch.Tag + "_output_z;\n"

	result += "\treg	[1:0] " + op.fpName + "_" + arch.Tag + "_state;\n"
	result += "parameter " + op.fpName + "_" + arch.Tag + "_put         = 2'd0,\n"
	result += "          " + op.fpName + "_" + arch.Tag + "_get         = 2'd1;\n"

	result += "\t" + op.fpName + "_" + arch.Tag + " " + op.fpName + "_" + arch.Tag + "_inst (" + op.fpName + "_" + arch.Tag + "_input_a, " + op.fpName + "_" + arch.Tag + "_input_b,  " + op.fpName + "_" + arch.Tag + "_output_z);\n\n"

	return result
}

func (op FixedPoint) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opBits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					" + strings.ToUpper(op.fpName) + ": begin\n"
	if arch.R == 1 {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opBits-1) + "])\n"
	} else {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opBits-1) + ":" + strconv.Itoa(rom_word-opBits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {

		if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlyDestRegs)) {
			cp := arch.Tag
			req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:" + op.fpName, T: bmreqs.ObjectSet, Name: "destregs", Value: Get_register_name(i), Op: bmreqs.OpCheck})
			if req.Value == "false" {
				continue
			}
		}

		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "							case (current_instruction[" + strconv.Itoa(rom_word-opBits-int(arch.R)-1) + "])\n"
		} else {
			result += "							case (current_instruction[" + strconv.Itoa(rom_word-opBits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opBits-int(arch.R)-int(arch.R)) + "])\n"
		}

		for j := 0; j < reg_num; j++ {

			if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlySrcRegs)) {
				cp := arch.Tag
				req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:" + op.fpName, T: bmreqs.ObjectSet, Name: "sourceregs", Value: Get_register_name(j), Op: bmreqs.OpCheck})
				if req.Value == "false" {
					continue
				}
			}

			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "							case (" + op.fpName + "_" + arch.Tag + "_state)\n"
			result += "							" + op.fpName + "_" + arch.Tag + "_put : begin\n"
			result += "								" + op.fpName + "_" + arch.Tag + "_input_a <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
			result += "								" + op.fpName + "_" + arch.Tag + "_input_b <= #1 _" + strings.ToLower(Get_register_name(j)) + ";\n"
			result += "								" + op.fpName + "_" + arch.Tag + "_state <= #1 " + op.fpName + "_" + arch.Tag + "_get;\n"
			result += "							end\n"
			result += "							" + op.fpName + "_" + arch.Tag + "_get : begin\n"
			result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 " + op.fpName + "_" + arch.Tag + "_output_z;\n"
			result += "								" + op.fpName + "_" + arch.Tag + "_state <= #1 " + op.fpName + "_" + arch.Tag + "_put;\n"
			result += NextInstruction(conf, arch, 8, "_pc + 1'b1")
			result += "							end\n"
			result += "							endcase\n"
			result += "								$display(\"" + strings.ToUpper(op.fpName) + " " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "							end\n"
		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "					end\n"
	return result
}

func (op FixedPoint) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op FixedPoint) Assembler(arch *Arch, words []string) (string, error) {
	opBits := arch.Opcodes_bits()
	rom_word := arch.Max_word()

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

	partial := ""
	for i := 0; i < reg_num; i++ {
		if words[1] == strings.ToLower(Get_register_name(i)) {
			partial += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if partial == "" {
		return "", Prerror{"Unknown register name " + words[1]}
	}

	result += partial

	for i := opBits + 2*int(arch.R); i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op FixedPoint) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

func FP_mult(a, b, regsize, fsize int64) int64 {
	shifted_a := a >> regsize
	shifted_b := b >> regsize

	part1 := (shifted_a | shifted_b) * ((a * b) >> fsize)
	part2 := (1 - (shifted_a | shifted_b)) * ((a * b) >> fsize)
	multResult := part1 + part2

	return multResult
}

func FP_div(num1, num2, regsize, fsize int64) int64 {
	bit1 := (num1 >> (regsize - 1)) & 1
	bit2 := (num2 >> (regsize - 1)) & 1

	if bit1^bit2 == 1 {
		numerator := num1
		denominator := num2

		if bit1 == 1 {
			numerator = -int64(num1)
		}

		if bit2 == 1 {
			denominator = -int64(num2)
		}

		numerator <<= fsize
		result := -int64(numerator / denominator)
		return result
	} else {
		result := int64(num1 << fsize / num2)
		return result
	}
}

func (op FixedPoint) Simulate(vm *VM, instr string) error {
	regBits := vm.Mach.R
	regDest := get_id(instr[:regBits])
	regSrc := get_id(instr[regBits : regBits*2])

	s := int64(uint64(1) << (op.f))

	switch *op.pipeline {
	case FPPUT:
		*op.pipeline = LQGET
	case LQGET:
		switch op.opType {
		case LQADD:
			vm.Registers[regDest] = Int64bits(Int64FromBits(vm.Registers[regDest].(uint64)) + Int64FromBits(vm.Registers[regSrc].(uint64)))
		case LQMULT:
			// TODO Check if this is correct
			// if vm.Mach.Rsize <= 8 {
			// 	vm.Registers[regDest] = Int8bits(Int8FromBits(vm.Registers[regDest].(uint8)) * Int8FromBits(vm.Registers[regSrc].(uint8)) / int8(s))
			// } else if vm.Mach.Rsize <= 16 {
			// 	vm.Registers[regDest] = Int16bits(Int16FromBits(vm.Registers[regDest].(uint16)) * Int16FromBits(vm.Registers[regSrc].(uint16)) / int16(s))
			// } else if vm.Mach.Rsize <= 32 {
			// 	vm.Registers[regDest] = Int32bits(Int32FromBits(vm.Registers[regDest].(uint32)) * Int32FromBits(vm.Registers[regSrc].(uint32)) / int32(s))
			// } else if vm.Mach.Rsize <= 64 {
			// 	vm.Registers[regDest] = Int64bits(Int64FromBits(vm.Registers[regDest].(uint64)) * Int64FromBits(vm.Registers[regSrc].(uint64)) / int64(s))
			// } else {
			// 	return errors.New("invalid register size, must be <= 64")
			// }

			a := Int64FromBits(vm.Registers[regDest].(uint64))
			b := Int64FromBits(vm.Registers[regSrc].(uint64))
			
			// WIP: change 16 and 6
			vm.Registers[regDest] = FP_mult(a, b, op.s, op.f)


		case LQDIV:
			// TODO Check if this is correct
			a := Int64FromBits(vm.Registers[regDest].(uint64))
			b := Int64FromBits(vm.Registers[regSrc].(uint64))
			
			// WIP: change 16 and 6
			vm.Registers[regDest] = FP_div(a, b, op.s, op.f)
		}
		vm.Pc = vm.Pc + 1
		*op.pipeline = LQPUT
	}
	return nil
}

// The random genaration does nothing
func (op FixedPoint) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op FixedPoint) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op FixedPoint) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op FixedPoint) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op FixedPoint) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op FixedPoint) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op FixedPoint) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op FixedPoint) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op FixedPoint) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	result := "\n\n"
	var moduleName string
	s := strconv.Itoa(op.s - 1)
	f := strconv.Itoa(op.f)

	result += "module " + op.fpName + "_" + arch.Tag + "(\n"
	result += "input_a,\n"
	result += "input_b,\n"
	result += "output_z);\n"
	result += "	input signed    [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] input_a;\n"
	result += "	input signed    [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] input_b;\n"
	result += "	output signed   [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] output_z;\n"
	result += "\n"

	switch op.opType {
	case LQADD:
		result += "	assign output_z = input_a + input_b;\n"
		moduleName = "adder"
	case LQMULT:
		result += "	wire signed [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] mult_result;\n"
		result += "	assign mult_result = input_a * input_b;\n"
		result += "	assign output_z = (mult_result[" + s + "]) ? (mult_result >>> " + f + ") : (mult_result >> " + f + ");\n"
		moduleName = "multiplier"
	case LQDIV:
		result += "	assign output_z = (input_a[" + s + "] ^ input_b[" + s + "]) ? -((((input_a[" + s + "] ? -input_a : input_a) <<< " + f + ") / (input_b[" + s + "] ? -input_b : input_b))) : ((input_a <<< " + f + ") / input_b);\n"
		moduleName = "divider"
	}
	result += "endmodule\n"

	moduleNames := []string{moduleName}
	moduleCodes := []string{result}

	return moduleNames, moduleCodes
}

func (Op FixedPoint) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op FixedPoint) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op FixedPoint) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, Op.fpName+"::*--type=reg::*--type=reg")
	return result
}
func (Op FixedPoint) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case Op.fpName:
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regSrc, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "opcodes", Value: Op.fpName, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:" + Op.fpName, T: bmreqs.ObjectSet, Name: "destregs", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:" + Op.fpName, T: bmreqs.ObjectSet, Name: "sourceregs", Value: regSrc, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}

func (Op FixedPoint) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op FixedPoint) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
