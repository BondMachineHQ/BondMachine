package procbuilder

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

const (
	FXPPUT = uint8(0) + iota
	FXPGET
)

// The FXP opcode is both a basic instruction and a template for other instructions.
type FXP struct {
	fpName   string
	s        int
	f        int
	opType   uint8
	pipeline *uint8
}

func (op FXP) Op_get_name() string {
	return op.fpName
}

func (op FXP) Op_get_desc() string {
	return "FXP dynamical instruction " + op.fpName
}

func (op FXP) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := op.fpName + " [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to the sum of its value with another register [" + strconv.Itoa(opBits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op FXP) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	return opBits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op FXP) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pName string) string {
	result := ""

	dri := op.fpName + "_" + arch.Tag

	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + dri + "_input_a;\n"
	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + dri + "_input_b;\n"
	result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + dri + "_output_z;\n"

	result += "\treg	[1:0] " + dri + "_state;\n"
	result += "parameter " + dri + "_put         = 2'd0,\n"
	result += "          " + dri + "_get         = 2'd1;\n"

	var fxpModule string
	switch op.opType {
	case FPXADD:
		fxpModule = "fxp_add"
	case FPXMULT:
		fxpModule = "fxp_mul"
	case FPXDIV:
		fxpModule = "fxp_div"
	default:
		return result
	}

	s := strconv.Itoa(op.s - op.f)
	f := strconv.Itoa(op.f)

	result += "\t" + fxpModule + " #(\n"
	result += "\t\t.WIIA(" + s + "),\n"
	result += "\t\t.WIFA(" + f + "),\n"
	result += "\t\t.WIIB(" + s + "),\n"
	result += "\t\t.WIFB(" + f + "),\n"
	result += "\t\t.WOI(" + s + "),\n"
	result += "\t\t.WOF(" + f + ")\n"
	result += "\t) " + dri + "_inst (\n"
	switch op.opType {
	case FXPDIV:
		result += "\t\t.dividend(" + dri + "_input_a),\n"
		result += "\t\t.divisor(" + dri + "_input_b),\n"
	default:
		result += "\t\t.ina(" + dri + "_input_a),\n"
		result += "\t\t.inb(" + dri + "_input_b),\n"
	}
	result += "\t\t.out(" + dri + "_output_z)\n"
	result += "\t);\n\n"

	return result
}

func (op FXP) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
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

func (op FXP) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op FXP) Assembler(arch *Arch, words []string) (string, error) {
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

func (op FXP) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

func fxpMult(a, b, regsize, fsize int64) int64 {
	shifted_a := a >> regsize
	shifted_b := b >> regsize

	part1 := (shifted_a | shifted_b) * ((a * b) >> fsize)
	part2 := (1 - (shifted_a | shifted_b)) * ((a * b) >> fsize)
	multResult := part1 + part2

	return multResult
}

func fxpDiv(num1, num2, regsize, fsize int64) int64 {
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

func (op FXP) Simulate(vm *VM, instr string) error {
	regBits := vm.Mach.R
	regDest := get_id(instr[:regBits])
	regSrc := get_id(instr[regBits : regBits*2])

	switch *op.pipeline {
	case FXPPUT:
		*op.pipeline = FXPGET
	case FXPGET:
		var dest int64
		var src int64
		var res int64
		if vm.Mach.Rsize <= 8 {
			dest = int64(Int8FromBits(vm.Registers[regDest].(uint8)))
			src = int64(Int8FromBits(vm.Registers[regSrc].(uint8)))
		} else if vm.Mach.Rsize <= 16 {
			dest = int64(Int16FromBits(vm.Registers[regDest].(uint16)))
			src = int64(Int16FromBits(vm.Registers[regSrc].(uint16)))
		} else if vm.Mach.Rsize <= 32 {
			dest = int64(Int32FromBits(vm.Registers[regDest].(uint32)))
			src = int64(Int32FromBits(vm.Registers[regSrc].(uint32)))
		} else if vm.Mach.Rsize <= 64 {
			dest = int64(Int64FromBits(vm.Registers[regDest].(uint64)))
			src = int64(Int64FromBits(vm.Registers[regSrc].(uint64)))
		} else {
			return errors.New("invalid register size, must be <= 64")
		}
		switch op.opType {
		case LQADD:
			res = dest + src
		case LQMULT:
			res = fxpMult(dest, src, int64(op.s), int64(op.f))
		case LQDIV:
			res = fxpDiv(dest, src, int64(op.s), int64(op.f))
		}
		if vm.Mach.Rsize <= 8 {
			vm.Registers[regDest] = uint8(Int8bits(int8(res)))
		} else if vm.Mach.Rsize <= 16 {
			vm.Registers[regDest] = uint16(Int16bits(int16(res)))
		} else if vm.Mach.Rsize <= 32 {
			vm.Registers[regDest] = uint32(Int32bits(int32(res)))
		} else if vm.Mach.Rsize <= 64 {
			vm.Registers[regDest] = uint64(Int64bits(int64(res)))
		} else {
			return errors.New("invalid register size, must be <= 64")
		}
		vm.Pc = vm.Pc + 1
		*op.pipeline = LQPUT
	}
	return nil
}

// The random genaration does nothing
func (op FXP) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op FXP) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op FXP) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op FXP) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op FXP) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op FXP) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op FXP) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op FXP) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op FXP) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op FXP) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op FXP) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op FXP) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, Op.fpName+"::*--type=reg::*--type=reg")
	return result
}
func (Op FXP) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
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

func (Op FXP) ExtraFiles(arch *Arch) ([]string, []string) {
	fileList := make([]string, 0)
	switch Op.opType {
	case FXPADD:
		fileList = append(fileList, "fxp_zoom.v")
		fileList = append(fileList, "fxp_add.v")
	case FXPMULT:
		fileList = append(fileList, "fxp_zoom.v")
		fileList = append(fileList, "fxp_mul.v")
	case FXPDIV:
		fileList = append(fileList, "fxp_zoom.v")
		fileList = append(fileList, "fxp_div.v")
	}

	names := make([]string, 0)
	codes := make([]string, 0)

	for _, fileName := range fileList {
		if !stringInSlice(fileName, strings.Split(arch.Conproc.SharedHDLOps, ",")) {

			if data, err := os.ReadFile("/tmp/fxpcode/" + fileName); err != nil {
				log.Fatal("unable to load file")
			} else {

				if len(arch.Conproc.SharedHDLOps) > 0 {
					arch.Conproc.SharedHDLOps += ","
				}
				arch.Conproc.SharedHDLOps += fileName
				names = append(names, fileName)
				codes = append(codes, string(data))
			}
		}
	}
	return names, codes
}

func (Op FXP) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
