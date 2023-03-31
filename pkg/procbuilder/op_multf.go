package procbuilder

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Multf struct{}

func (op Multf) Op_get_name() string {
	return "multf"
}

func (op Multf) Op_get_desc() string {
	return "Register float32 multiplcation"
}

func (op Multf) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "multf [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to the product of its value with another register [" + strconv.Itoa(opbits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op Multf) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op Multf) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	result += "\treg [31:0] multiplier_" + arch.Tag + "_input_a;\n"
	result += "\treg [31:0] multiplier_" + arch.Tag + "_input_b;\n"
	result += "\treg multiplier_" + arch.Tag + "_input_a_stb;\n"
	result += "\treg multiplier_" + arch.Tag + "_input_b_stb;\n"
	result += "\treg multiplier_" + arch.Tag + "_output_z_ack;\n\n"

	result += "\twire [31:0] multiplier_" + arch.Tag + "_output_z;\n"
	result += "\twire multiplier_" + arch.Tag + "_output_z_stb;\n"
	result += "\twire multiplier_" + arch.Tag + "_input_a_ack;\n"
	result += "\twire multiplier_" + arch.Tag + "_input_b_ack;\n\n"

	result += "\treg	[1:0] multiplier_" + arch.Tag + "_state;\n"
	result += "parameter multiplier_" + arch.Tag + "_put_a         = 2'd0,\n"
	result += "          multiplier_" + arch.Tag + "_put_b         = 2'd1,\n"
	result += "          multiplier_" + arch.Tag + "_get_z         = 2'd2;\n"

	result += "\tmultiplier_" + arch.Tag + " multiplier_" + arch.Tag + "_inst (multiplier_" + arch.Tag + "_input_a, multiplier_" + arch.Tag + "_input_b, multiplier_" + arch.Tag + "_input_a_stb, multiplier_" + arch.Tag + "_input_b_stb, multiplier_" + arch.Tag + "_output_z_ack, clock_signal, reset_signal, multiplier_" + arch.Tag + "_output_z, multiplier_" + arch.Tag + "_output_z_stb, multiplier_" + arch.Tag + "_input_a_ack, multiplier_" + arch.Tag + "_input_b_ack);\n\n"

	return result
}

func (op Multf) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					MULTF: begin\n"
	if arch.R == 1 {
		result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {

		if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlyDestRegs)) {
			cp := arch.Tag
			req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:multf", T: bmreqs.ObjectSet, Name: "destregs", Value: Get_register_name(i), Op: bmreqs.OpCheck})
			if req.Value == "false" {
				continue
			}
		}

		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "							case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
		} else {
			result += "							case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.R)) + "])\n"
		}

		for j := 0; j < reg_num; j++ {

			if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlySrcRegs)) {
				cp := arch.Tag
				req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:multf", T: bmreqs.ObjectSet, Name: "sourceregs", Value: Get_register_name(j), Op: bmreqs.OpCheck})
				if req.Value == "false" {
					continue
				}
			}

			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "							case (multiplier_" + arch.Tag + "_state)\n"
			result += "							multiplier_" + arch.Tag + "_put_a : begin\n"
			result += "								if (multiplier_" + arch.Tag + "_input_a_ack) begin\n"
			result += "									multiplier_" + arch.Tag + "_input_a <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
			result += "									multiplier_" + arch.Tag + "_input_a_stb <= #1 1;\n"
			result += "									multiplier_" + arch.Tag + "_output_z_ack <= #1 0;\n"
			result += "									multiplier_" + arch.Tag + "_state <= #1 multiplier_" + arch.Tag + "_put_b;\n"
			result += "								end\n"
			result += "							end\n"
			result += "							multiplier_" + arch.Tag + "_put_b : begin\n"
			result += "								if (multiplier_" + arch.Tag + "_input_b_ack) begin\n"
			result += "									multiplier_" + arch.Tag + "_input_b <= #1 _" + strings.ToLower(Get_register_name(j)) + ";\n"
			result += "									multiplier_" + arch.Tag + "_input_b_stb <= #1 1;\n"
			result += "									multiplier_" + arch.Tag + "_output_z_ack <= #1 0;\n"
			result += "									multiplier_" + arch.Tag + "_state <= #1 multiplier_" + arch.Tag + "_get_z;\n"
			result += "									multiplier_" + arch.Tag + "_input_a_stb <= #1 0;\n"
			result += "								end\n"
			result += "							end\n"
			result += "							multiplier_" + arch.Tag + "_get_z : begin\n"
			result += "								if (multiplier_" + arch.Tag + "_output_z_stb) begin\n"
			result += "									_" + strings.ToLower(Get_register_name(i)) + " <= #1 multiplier_" + arch.Tag + "_output_z;\n"
			result += "									multiplier_" + arch.Tag + "_output_z_ack <= #1 1;\n"
			result += "									multiplier_" + arch.Tag + "_state <= #1 multiplier_" + arch.Tag + "_put_a;\n"
			result += "									multiplier_" + arch.Tag + "_input_b_stb <= #1 0;\n"
			result += "									_pc <= #1 _pc + 1'b1 ;\n"
			result += "								end\n"
			result += "							end\n"
			result += "							endcase\n"
			result += "								$display(\"ADDF " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "							end\n"
		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "					end\n"
	return result
}

func (op Multf) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Multf) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
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

	for i := opbits + 2*int(arch.R); i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op Multf) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

func (op Multf) Simulate(vm *VM, instr string) error {
	reg_bits := vm.Mach.R
	regDest := get_id(instr[:reg_bits])
	regSrc := get_id(instr[reg_bits : reg_bits*2])
	switch vm.Mach.Rsize {
	case 32:
		var floatDest float32
		var floatSrc float32
		if v, ok := vm.Registers[regDest].(uint32); ok {
			floatDest = math.Float32frombits(v)
		} else {
			floatDest = float32(0.0)
		}
		if v, ok := vm.Registers[regSrc].(uint32); ok {
			floatSrc = math.Float32frombits(v)
		} else {
			floatSrc = float32(0.0)
		}
		vm.Registers[regDest] = math.Float32bits(floatDest * floatSrc)
	default:
		return errors.New("invalid register size, for float registers has to be 32 bits")
	}
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Multf) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Multf) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Multf) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Multf) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Multf) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Multf) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Multf) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Multf) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Multf) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	result := "\n\n"

	result += "//IEEE Floating Point Multiplier (Single Precision)\n"
	result += "//Copyright (C) Jonathan P Dawson 2013\n"
	result += "//2013-12-12\n"
	result += "module multiplier_" + arch.Tag + "(\n"
	result += "        input_a,\n"
	result += "        input_b,\n"
	result += "        input_a_stb,\n"
	result += "        input_b_stb,\n"
	result += "        output_z_ack,\n"
	result += "        clk,\n"
	result += "        rst,\n"
	result += "        output_z,\n"
	result += "        output_z_stb,\n"
	result += "        input_a_ack,\n"
	result += "        input_b_ack);\n"
	result += "\n"
	result += "  input     clk;\n"
	result += "  input     rst;\n"
	result += "\n"
	result += "  input     [31:0] input_a;\n"
	result += "  input     input_a_stb;\n"
	result += "  output    input_a_ack;\n"
	result += "\n"
	result += "  input     [31:0] input_b;\n"
	result += "  input     input_b_stb;\n"
	result += "  output    input_b_ack;\n"
	result += "\n"
	result += "  output    [31:0] output_z;\n"
	result += "  output    output_z_stb;\n"
	result += "  input     output_z_ack;\n"
	result += "\n"
	result += "  reg       s_output_z_stb;\n"
	result += "  reg       [31:0] s_output_z;\n"
	result += "  reg       s_input_a_ack;\n"
	result += "  reg       s_input_b_ack;\n"
	result += "\n"
	result += "  reg       [3:0] state;\n"
	result += "  parameter get_a         = 4'd0,\n"
	result += "            get_b         = 4'd1,\n"
	result += "            unpack        = 4'd2,\n"
	result += "            special_cases = 4'd3,\n"
	result += "            normalise_a   = 4'd4,\n"
	result += "            normalise_b   = 4'd5,\n"
	result += "            multiply_0    = 4'd6,\n"
	result += "            multiply_1    = 4'd7,\n"
	result += "            normalise_1   = 4'd8,\n"
	result += "            normalise_2   = 4'd9,\n"
	result += "            round         = 4'd10,\n"
	result += "            pack          = 4'd11,\n"
	result += "            put_z         = 4'd12;\n"
	result += "\n"
	result += "  reg       [31:0] a, b, z;\n"
	result += "  reg       [23:0] a_m, b_m, z_m;\n"
	result += "  reg       [9:0] a_e, b_e, z_e;\n"
	result += "  reg       a_s, b_s, z_s;\n"
	result += "  reg       guard, round_bit, sticky;\n"
	result += "  reg       [49:0] product;\n"
	result += "\n"
	result += "  always @(posedge clk)\n"
	result += "  begin\n"
	result += "\n"
	result += "    case(state)\n"
	result += "\n"
	result += "      get_a:\n"
	result += "      begin\n"
	result += "        s_input_a_ack <= 1;\n"
	result += "        if (s_input_a_ack && input_a_stb) begin\n"
	result += "          a <= input_a;\n"
	result += "          s_input_a_ack <= 0;\n"
	result += "          state <= get_b;\n"
	result += "        end\n"
	result += "      end\n"
	result += "\n"
	result += "      get_b:\n"
	result += "      begin\n"
	result += "        s_input_b_ack <= 1;\n"
	result += "        if (s_input_b_ack && input_b_stb) begin\n"
	result += "          b <= input_b;\n"
	result += "          s_input_b_ack <= 0;\n"
	result += "          state <= unpack;\n"
	result += "        end\n"
	result += "      end\n"
	result += "\n"
	result += "      unpack:\n"
	result += "      begin\n"
	result += "        a_m <= a[22 : 0];\n"
	result += "        b_m <= b[22 : 0];\n"
	result += "        a_e <= a[30 : 23] - 127;\n"
	result += "        b_e <= b[30 : 23] - 127;\n"
	result += "        a_s <= a[31];\n"
	result += "        b_s <= b[31];\n"
	result += "        state <= special_cases;\n"
	result += "      end\n"
	result += "\n"
	result += "      special_cases:\n"
	result += "      begin\n"
	result += "        //if a is NaN or b is NaN return NaN \n"
	result += "        if ((a_e == 128 && a_m != 0) || (b_e == 128 && b_m != 0)) begin\n"
	result += "          z[31] <= 1;\n"
	result += "          z[30:23] <= 255;\n"
	result += "          z[22] <= 1;\n"
	result += "          z[21:0] <= 0;\n"
	result += "          state <= put_z;\n"
	result += "        //if a is inf return inf\n"
	result += "        end else if (a_e == 128) begin\n"
	result += "          z[31] <= a_s ^ b_s;\n"
	result += "          z[30:23] <= 255;\n"
	result += "          z[22:0] <= 0;\n"
	result += "          state <= put_z;\n"
	result += "           //if b is zero return NaN\n"
	result += "          if ($signed(b_e == -127) && (b_m == 0)) begin\n"
	result += "            z[31] <= 1;\n"
	result += "            z[30:23] <= 255;\n"
	result += "            z[22] <= 1;\n"
	result += "            z[21:0] <= 0;\n"
	result += "            state <= put_z;\n"
	result += "          end\n"
	result += "        //if b is inf return inf\n"
	result += "        end else if (b_e == 128) begin\n"
	result += "          z[31] <= a_s ^ b_s;\n"
	result += "          z[30:23] <= 255;\n"
	result += "          z[22:0] <= 0;\n"
	result += "          state <= put_z;\n"
	result += "        //if a is zero return zero\n"
	result += "        end else if (($signed(a_e) == -127) && (a_m == 0)) begin\n"
	result += "          z[31] <= a_s ^ b_s;\n"
	result += "          z[30:23] <= 0;\n"
	result += "          z[22:0] <= 0;\n"
	result += "          state <= put_z;\n"
	result += "        //if b is zero return zero\n"
	result += "        end else if (($signed(b_e) == -127) && (b_m == 0)) begin\n"
	result += "          z[31] <= a_s ^ b_s;\n"
	result += "          z[30:23] <= 0;\n"
	result += "          z[22:0] <= 0;\n"
	result += "          state <= put_z;\n"
	result += "        end else begin\n"
	result += "          //Denormalised Number\n"
	result += "          if ($signed(a_e) == -127) begin\n"
	result += "            a_e <= -126;\n"
	result += "          end else begin\n"
	result += "            a_m[23] <= 1;\n"
	result += "          end\n"
	result += "          //Denormalised Number\n"
	result += "          if ($signed(b_e) == -127) begin\n"
	result += "            b_e <= -126;\n"
	result += "          end else begin\n"
	result += "            b_m[23] <= 1;\n"
	result += "          end\n"
	result += "          state <= normalise_a;\n"
	result += "        end\n"
	result += "      end\n"
	result += "\n"
	result += "      normalise_a:\n"
	result += "      begin\n"
	result += "        if (a_m[23]) begin\n"
	result += "          state <= normalise_b;\n"
	result += "        end else begin\n"
	result += "          a_m <= a_m << 1;\n"
	result += "          a_e <= a_e - 1;\n"
	result += "        end\n"
	result += "      end\n"
	result += "\n"
	result += "      normalise_b:\n"
	result += "      begin\n"
	result += "        if (b_m[23]) begin\n"
	result += "          state <= multiply_0;\n"
	result += "        end else begin\n"
	result += "          b_m <= b_m << 1;\n"
	result += "          b_e <= b_e - 1;\n"
	result += "        end\n"
	result += "      end\n"
	result += "\n"
	result += "      multiply_0:\n"
	result += "      begin\n"
	result += "        z_s <= a_s ^ b_s;\n"
	result += "        z_e <= a_e + b_e + 1;\n"
	result += "        product <= a_m * b_m * 4;\n"
	result += "        state <= multiply_1;\n"
	result += "      end\n"
	result += "\n"
	result += "      multiply_1:\n"
	result += "      begin\n"
	result += "        z_m <= product[49:26];\n"
	result += "        guard <= product[25];\n"
	result += "        round_bit <= product[24];\n"
	result += "        sticky <= (product[23:0] != 0);\n"
	result += "        state <= normalise_1;\n"
	result += "      end\n"
	result += "\n"
	result += "      normalise_1:\n"
	result += "      begin\n"
	result += "        if (z_m[23] == 0) begin\n"
	result += "          z_e <= z_e - 1;\n"
	result += "          z_m <= z_m << 1;\n"
	result += "          z_m[0] <= guard;\n"
	result += "          guard <= round_bit;\n"
	result += "          round_bit <= 0;\n"
	result += "        end else begin\n"
	result += "          state <= normalise_2;\n"
	result += "        end\n"
	result += "      end\n"
	result += "\n"
	result += "      normalise_2:\n"
	result += "      begin\n"
	result += "        if ($signed(z_e) < -126) begin\n"
	result += "          z_e <= z_e + 1;\n"
	result += "          z_m <= z_m >> 1;\n"
	result += "          guard <= z_m[0];\n"
	result += "          round_bit <= guard;\n"
	result += "          sticky <= sticky | round_bit;\n"
	result += "        end else begin\n"
	result += "          state <= round;\n"
	result += "        end\n"
	result += "      end\n"
	result += "\n"
	result += "      round:\n"
	result += "      begin\n"
	result += "        if (guard && (round_bit | sticky | z_m[0])) begin\n"
	result += "          z_m <= z_m + 1;\n"
	result += "          if (z_m == 24'hffffff) begin\n"
	result += "            z_e <=z_e + 1;\n"
	result += "          end\n"
	result += "        end\n"
	result += "        state <= pack;\n"
	result += "      end\n"
	result += "\n"
	result += "      pack:\n"
	result += "      begin\n"
	result += "        z[22 : 0] <= z_m[22:0];\n"
	result += "        z[30 : 23] <= z_e[7:0] + 127;\n"
	result += "        z[31] <= z_s;\n"
	result += "        if ($signed(z_e) == -126 && z_m[23] == 0) begin\n"
	result += "          z[30 : 23] <= 0;\n"
	result += "        end\n"
	result += "        //if overflow occurs, return inf\n"
	result += "        if ($signed(z_e) > 127) begin\n"
	result += "          z[22 : 0] <= 0;\n"
	result += "          z[30 : 23] <= 255;\n"
	result += "          z[31] <= z_s;\n"
	result += "        end\n"
	result += "        state <= put_z;\n"
	result += "      end\n"
	result += "\n"
	result += "      put_z:\n"
	result += "      begin\n"
	result += "        s_output_z_stb <= 1;\n"
	result += "        s_output_z <= z;\n"
	result += "        if (s_output_z_stb && output_z_ack) begin\n"
	result += "          s_output_z_stb <= 0;\n"
	result += "          state <= get_a;\n"
	result += "        end\n"
	result += "      end\n"
	result += "\n"
	result += "    endcase\n"
	result += "\n"
	result += "    if (rst == 1) begin\n"
	result += "      state <= get_a;\n"
	result += "      s_input_a_ack <= 0;\n"
	result += "      s_input_b_ack <= 0;\n"
	result += "      s_output_z_stb <= 0;\n"
	result += "    end\n"
	result += "\n"
	result += "  end\n"
	result += "  assign input_a_ack = s_input_a_ack;\n"
	result += "  assign input_b_ack = s_input_b_ack;\n"
	result += "  assign output_z_stb = s_output_z_stb;\n"
	result += "  assign output_z = s_output_z;\n"
	result += "\n"
	result += "endmodule\n"

	return []string{"multiplier"}, []string{result}
}

func (Op Multf) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Multf) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Multf) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "multf::*--type=reg::*--type=reg")
	return result
}
func (Op Multf) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "multf":
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regSrc, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "opcodes", Value: "multf", Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:multf", T: bmreqs.ObjectSet, Name: "destregs", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:multf", T: bmreqs.ObjectSet, Name: "sourceregs", Value: regSrc, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Multf) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
