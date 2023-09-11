package procbuilder

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bcof"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

const (
	letters = "0123456789abcdef"
)

// The CPU
type Conproc struct {
	CpID  uint32
	Rsize uint8
	R     uint8 // Number of n-bit registers
	N     uint8 // Number of n-bit inputs
	M     uint8 // Number of n-bit outputs
	Op    []Opcode
}

type Config struct {
	*bmreqs.ReqRoot
	*bcof.BCOFEntry
	HwOptimizations
	Debug             bool
	Commented_verilog bool
	Runinfo           *RuntimeInfo
}

type RuntimeInfo struct {
	HeaderFlags map[string]bool
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (ri *RuntimeInfo) Init() {
	ri.HeaderFlags = make(map[string]bool)
}

func (ri *RuntimeInfo) Check(flag string) bool {
	if _, present := ri.HeaderFlags[flag]; present {
		return false
	} else {
		ri.HeaderFlags[flag] = true
		return true
	}
}

func (proc *Conproc) String() string {
	reg_num := 1 << proc.R

	result := "Opcodes: "
	for _, op := range proc.Op {
		result += op.Op_get_name() + " "
	}

	result += "\nRegisters: "
	for i := 0; i < reg_num; i++ {
		result += strings.ToUpper(Get_register_name(i)) + " "
	}

	return result
}

func (proc *Conproc) Decode_opcode(intr string) (int, error) {
	opbits := proc.Opcodes_bits()
	result := get_id(intr[0:opbits])
	return result, nil
}

func (proc *Conproc) Opcodes_bits() int {
	served := 1
	for bits := 1; bits < 16; bits++ {
		if served<<uint8(bits) >= len(proc.Op) {
			return bits
		}
	}
	return 1
}

func (proc *Conproc) Inputs_bits() int {
	served := 1
	for bits := 1; bits < 16; bits++ {
		if served<<uint8(bits) >= int(proc.N) {
			return bits
		}
	}
	return 1
}

func (proc *Conproc) Outputs_bits() int {
	served := 1
	for bits := 1; bits < 16; bits++ {
		if served<<uint8(bits) >= int(proc.M) {
			return bits
		}
	}
	return 1
}

func (proc *Conproc) Write_opcodes_verilog() string {
	opbits := proc.Opcodes_bits()

	result := ""
	result += "			// Opcodes in the istructions, lenght accourding the number of the selected.\n"

	for i, op := range proc.Op {
		if i == 0 {
			if len(proc.Op) == 1 {
				result += "	localparam	" + strings.ToUpper(op.Op_get_name()) + "=" + strconv.Itoa(opbits) + "'b" + zeros_prefix(opbits, get_binary(i)) + ";          // " + op.Op_get_desc() + "\n"
			} else {
				result += "	localparam	" + strings.ToUpper(op.Op_get_name()) + "=" + strconv.Itoa(opbits) + "'b" + zeros_prefix(opbits, get_binary(i)) + ",          // " + op.Op_get_desc() + "\n"
			}
		} else if i == len(proc.Op)-1 {
			result += "			" + strings.ToUpper(op.Op_get_name()) + "=" + strconv.Itoa(opbits) + "'b" + zeros_prefix(opbits, get_binary(i)) + ";          // " + op.Op_get_desc() + "\n"
		} else {
			result += "			" + strings.ToUpper(op.Op_get_name()) + "=" + strconv.Itoa(opbits) + "'b" + zeros_prefix(opbits, get_binary(i)) + ",          // " + op.Op_get_desc() + "\n"
		}
	}

	result += "\n"
	return result
}

func NextInstruction(conf *Config, arch *Arch, tabs int, jumpTo string) string {
	result := ""
	tabS := ""
	for i := 0; i < tabs; i++ {
		tabS += "\t"
	}
	switch arch.Modes[0] {
	case "ha":
		result += tabS + "_pc <= #1 " + jumpTo + ";\n"
	case "hy":
		result += tabS + "if (exec_mode == 1'b1) begin\n"
		result += tabS + "\tvn_state <= FETCH;\n"
		result += tabS + "end\n"
		result += tabS + "_pc <= #1 " + jumpTo + ";\n"
	case "vn":
		result += tabS + "vn_state <= FETCH;\n"
		result += tabS + "_pc <= #1 " + jumpTo + ";\n"
	}
	return result
}

func ExecutionCase(conf *Config, arch *Arch, tabs int, open bool) string {
	result := ""
	tabS := ""
	for i := 0; i < tabs; i++ {
		tabS += "\t"
	}

	switch arch.Modes[0] {
	case "ha":
	case "hy":
		if open {
			result += tabS + "if (exec_mode == 1'b0 || vn_state == EXECUTE) begin\n"
		} else {
			result += tabS + "end\n"
		}
	case "vn":
		if open {
			result += tabS + "if (vn_state == EXECUTE) begin\n"
		} else {
			result += tabS + "end\n"
		}
	}
	return result
}

func (proc *Conproc) Write_verilog(conf *Config, arch *Arch, processor_module_name string, flavor string) string {
	regsize := int(proc.Rsize)
	rom_word := arch.Max_word()
	opbits := proc.Opcodes_bits()
	inbits := proc.Inputs_bits()
	outbits := proc.Outputs_bits()

	reg_num := 1 << proc.R

	result := ""

	arch.Tag = fmt.Sprint(proc.CpID)

	// Module header
	result += "`timescale 1ns/1ps\n"
	result += "module " + processor_module_name + "(clock_signal, reset_signal"

	ramh := ""

	if int(arch.L) != 0 {
		ramh += ", ram_din, ram_dout, ram_addr, ram_wren, ram_en"
	}
	result += ", rom_bus, rom_value" + ramh

	for i := 0; i < int(arch.N); i++ {
		result += ", " + Get_input_name(i) + ", " + Get_input_name(i) + "_valid, " + Get_input_name(i) + "_received"
	}

	for i := 0; i < int(arch.M); i++ {
		result += ", " + Get_output_name(i) + ", " + Get_output_name(i) + "_valid, " + Get_output_name(i) + "_received"
	}

	if arch.Shared_constraints != "" {
		header := ""
		seq := make(map[string]int)
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if _, ok := seq[soname]; ok {
				seq[soname]++
			} else {
				seq[soname] = 0
			}

			for _, so := range Allshared {
				if so.Shr_get_name() == soname {
					header += so.GetArchHeader(arch, constraint, seq[soname])
				}
			}
		}
		result += header
	}

	result += ");\n\n"

	// Header variables declarations
	result += "\tinput clock_signal;\n"
	result += "\tinput reset_signal;\n"

	result += "\toutput  [" + strconv.Itoa(int(arch.O)-1) + ":0] rom_bus;\n"
	result += "\tinput  [" + strconv.Itoa(int(rom_word)-1) + ":0] rom_value;\n"

	if int(arch.L) != 0 {
		result += "\tinput  [" + strconv.Itoa(regsize-1) + ":0] ram_dout;\n"
		result += "\toutput [" + strconv.Itoa(regsize-1) + ":0] ram_din;\n"
		result += "\toutput  [" + strconv.Itoa(int(arch.L)-1) + ":0] ram_addr;\n"
		result += "\toutput ram_wren, ram_en;\n"
	}

	result += "\n"

	for i := 0; i < int(arch.N); i++ {
		result += "\tinput [" + strconv.Itoa(regsize-1) + ":0] " + Get_input_name(i) + ";\n"
		result += "\tinput " + Get_input_name(i) + "_valid;\n"
		result += "\toutput " + Get_input_name(i) + "_received;\n"
	}

	for i := 0; i < int(arch.M); i++ {
		result += "\toutput [" + strconv.Itoa(regsize-1) + ":0] " + Get_output_name(i) + ";\n"
		result += "\toutput " + Get_output_name(i) + "_valid;\n"
		result += "\tinput " + Get_output_name(i) + "_received;\n"
	}

	result += "\n"

	if arch.Shared_constraints != "" {
		params := ""
		seq := make(map[string]int)
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if _, ok := seq[soname]; ok {
				seq[soname]++
			} else {
				seq[soname] = 0
			}

			for _, so := range Allshared {
				if so.Shr_get_name() == soname {
					params += so.GetCPParams(arch, constraint, seq[soname])
					params += "\n"
				}
			}
		}
		result += params
	}

	//	// Module header
	//	result += "module " + processor_module_name + "(input clock_signal, input reset_signal, output [" + strconv.Itoa(int(arch.O)-1) + ":0] rom_bus, input [" + strconv.Itoa(rom_word-1) + ":0] rom_value"
	//
	//	for i := 0; i < int(proc.N); i++ {
	//		result += ", input [" + strconv.Itoa(regsize-1) + ":0] " + Get_input_name(i)
	//
	//	}
	//
	//	for i := 0; i < int(proc.M); i++ {
	//		result += ", output [" + strconv.Itoa(regsize-1) + ":0] " + Get_output_name(i)
	//
	//	}
	//
	//	result += ");\n"

	// Opcodes generation

	result += proc.Write_opcodes_verilog()

	// Registers generation
	result += "	localparam	" + strings.ToUpper(Get_register_name(0)) + "=" + strconv.Itoa(int(proc.R)) + "'b" + zeros_prefix(int(proc.R), get_binary(0)) + ",		// Registers in the intructions\n"

	for i := 1; i < reg_num; i++ {
		if i == reg_num-1 {
			result += "			" + strings.ToUpper(Get_register_name(i)) + "=" + strconv.Itoa(int(proc.R)) + "'b" + zeros_prefix(int(proc.R), get_binary(i)) + ";\n"
		} else {
			result += "			" + strings.ToUpper(Get_register_name(i)) + "=" + strconv.Itoa(int(proc.R)) + "'b" + zeros_prefix(int(proc.R), get_binary(i)) + ",\n"
		}
	}

	for i := 0; i < int(proc.N); i++ {
		if i == 0 {
			result += "	localparam"
		}
		if i == int(proc.N)-1 {
			result += "			" + strings.ToUpper(Get_input_name(i)) + "=" + strconv.Itoa(int(inbits)) + "'b" + zeros_prefix(int(inbits), get_binary(i)) + ";\n"
		} else {
			result += "			" + strings.ToUpper(Get_input_name(i)) + "=" + strconv.Itoa(int(inbits)) + "'b" + zeros_prefix(int(inbits), get_binary(i)) + ",\n"
		}
	}

	for i := 0; i < int(proc.M); i++ {
		if i == 0 {
			result += "	localparam"
		}
		if i == int(proc.M)-1 {
			result += "			" + strings.ToUpper(Get_output_name(i)) + "=" + strconv.Itoa(int(outbits)) + "'b" + zeros_prefix(int(outbits), get_binary(i)) + ";\n"
		} else {
			result += "			" + strings.ToUpper(Get_output_name(i)) + "=" + strconv.Itoa(int(outbits)) + "'b" + zeros_prefix(int(outbits), get_binary(i)) + ",\n"
		}
	}

	for i := 0; i < int(proc.M); i++ {
		result += "	reg [" + strconv.Itoa(regsize-1) + ":0] _aux" + Get_output_name(i) + ";\n"
	}

	result += "\n"
	result += "	reg [" + strconv.Itoa(regsize-1) + ":0] _ram [0:" + strconv.Itoa((1<<arch.L)-1) + "];		// Internal processor RAM\n"
	result += "\n"
	switch arch.Modes[0] {
	case "ha":
		result += "	(* KEEP = \"TRUE\" *) reg [" + strconv.Itoa(int(arch.O)-1) + ":0] _pc;		// Program counter\n"
	case "hy":
		if arch.L > arch.O {
			result += "	(* KEEP = \"TRUE\" *) reg [" + strconv.Itoa(int(arch.L)-1) + ":0] _pc;		// Program counter\n"
		} else {
			result += "	(* KEEP = \"TRUE\" *) reg [" + strconv.Itoa(int(arch.O)-1) + ":0] _pc;		// Program counter\n"
		}
	case "vn":
		result += "	(* KEEP = \"TRUE\" *) reg [" + strconv.Itoa(int(arch.L)-1) + ":0] _pc;		// Program counter\n"
	}
	result += "\n"
	result += "	// The number of registers are 2^R, two letters and an unserscore as identifier , maximum R=8 and 265 rigisters\n"

	for i := 0; i < reg_num; i++ {
		result += "	(* KEEP = \"TRUE\" *) reg [" + strconv.Itoa(regsize-1) + ":0] _" + strings.ToLower(Get_register_name(i)) + ";\n"
	}

	// modes handling
	switch arch.Modes[0] {
	case "ha":
		result += "\n"
		result += "	wire [" + strconv.Itoa(int(rom_word)-1) + ":0] current_instruction;\n"
		result += "	assign current_instruction=rom_value;\n"
		result += "\n"
	case "hy":
		result += "\n"
		result += "	wire [" + strconv.Itoa(int(rom_word)-1) + ":0] current_instruction;\n"
		result += "	reg [" + strconv.Itoa(int(rom_word)-1) + ":0] ram_instruction;\n"
		result += "	reg exec_mode; // 0 = harvard , 1=VN\n"
		result += "	reg [1:0] vn_state;\n"
		result += "	localparam FETCH=2'b00, WAIT=2'b10, EXECUTE=2'b01;\n"
		result += "	assign current_instruction= (exec_mode==1'b0) ? rom_value : ram_instruction;\n"
		result += "\n"
		if !arch.HasOp("r2m") && !arch.HasOp("m2r") && !arch.HasOp("r2mri") && !arch.HasOp("m2rri") {
			result += "	assign ram_addr = _pc;\n"
			result += "	assign ram_en = 1'b1;\n"
		}
	case "vn":
		result += "\n"
		result += "	wire [" + strconv.Itoa(int(rom_word)-1) + ":0] current_instruction;\n"
		result += "	reg [" + strconv.Itoa(int(rom_word)-1) + ":0] ram_instruction;\n"
		result += "	reg [1:0] vn_state;\n"
		result += "	localparam FETCH=2'b00, WAIT=2'b10, EXECUTE=2'b01;\n"
		result += "	assign current_instruction=ram_instruction;\n"
		result += "\n"
		if !arch.HasOp("r2m") && !arch.HasOp("m2r") && !arch.HasOp("r2mri") && !arch.HasOp("m2rri") {
			result += "	assign ram_addr = _pc;\n"
			result += "	assign ram_en = 1'b1;\n"
		}
	}

	for _, op := range proc.Op {
		if conf.Commented_verilog {
			result += "\n// Start of the component \"header\" for the opcode " + op.Op_get_name() + "\n\n"
		}
		result += op.OpInstructionVerilogHeader(conf, arch, flavor, processor_module_name)
	}

	result += "\n"
	result += "	always @(posedge clock_signal, posedge reset_signal)\n"
	result += "	begin\n"
	result += "		if(reset_signal)\n"
	result += "		begin\n"
	switch arch.Modes[0] {
	case "ha":
		result += "			_pc <= #1 " + strconv.Itoa(int(arch.O)) + "'h0;\n"
	case "hy":
		if arch.L > arch.O {
			result += "			_pc <= #1 " + strconv.Itoa(int(arch.L)) + "'h0;\n"
		} else {
			result += "			_pc <= #1 " + strconv.Itoa(int(arch.O)) + "'h0;\n"
		}
	case "vn":
		result += "			_pc <= #1 " + strconv.Itoa(int(arch.L)) + "'h0;\n"
	}

	for i := 0; i < reg_num; i++ {
		result += "			_" + strings.ToLower(Get_register_name(i)) + " <= #1 " + strconv.Itoa(int(arch.Rsize)) + "'h0;\n"
	}

	for _, op := range proc.Op {
		if conf.Commented_verilog {
			result += "\n// Start of the component \"reset\" for the opcode " + op.Op_get_name() + "\n\n"
		}
		result += op.Op_instruction_verilog_reset(arch, flavor)
	}

	result += "		end\n"
	result += "		else begin\n"

	switch arch.Modes[0] {
	case "ha":
		result += "			// ha placeholder\n"
	case "hy":
		result += "			if (exec_mode == 1'b1 && vn_state == FETCH) begin\n"
		result += "				vn_state <= WAIT;\n"
		result += "			end\n"
		result += "			else if (exec_mode == 1'b1 && vn_state == WAIT) begin\n"
		result += "				vn_state <= EXECUTE;\n"
		result += "				ram_instruction <= ram_dout;\n"
		result += "			end\n"
		result += "			else begin\n"
	case "vn":
		result += "			case (vn_state)\n"
		result += "			FETCH: begin\n"
		result += "				vn_state <= WAIT;\n"
		result += "			end\n"
		result += "			WAIT: begin\n"
		result += "				vn_state <= EXECUTE;\n"
		result += "				ram_instruction <= ram_dout;\n"
		result += "			end\n"
		result += "			EXECUTE: begin\n"
	}

	result += "			$display(\"Program Counter:%d\", _pc);\n"
	result += "			$display(\"Instruction:%b\", rom_value);\n"

	format := ""
	list := ""
	for i := 0; i < reg_num; i++ {
		format = format + strings.ToLower(Get_register_name(i)) + ":%b "
		list = list + ", _" + strings.ToLower(Get_register_name(i))
	}

	result += "			$display(\"Registers " + format + "\"" + list + ");\n"

	for _, op := range proc.Op {
		if conf.Commented_verilog {
			result += "\n// Start of the component \"internal state\" for the opcode " + op.Op_get_name() + "\n\n"
		}
		result += op.Op_instruction_verilog_internal_state(arch, flavor)
	}

	// TODO What are they, maybe needed by some opcode ?
	//	result += "			else\n"
	//	result += "				begin\n"

	for _, op := range proc.Op {
		if conf.Commented_verilog {
			result += "\n// Start of the component \"default state\" for the opcode " + op.Op_get_name() + "\n\n"
		}
		result += op.Op_instruction_verilog_default_state(arch, flavor)
	}

	if opbits == 1 {
		result += "				case(current_instruction[" + strconv.Itoa(rom_word-1) + "])\n"
	} else {
		result += "				case(current_instruction[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "])\n"
	}

	for _, op := range proc.Op {
		if conf.Commented_verilog {
			result += "\n// Start of the component of the \"state machine\" for the opcode " + op.Op_get_name() + "\n\n"
		}
		result += op.Op_instruction_verilog_state_machine(conf, arch, conf.ReqRoot, flavor)
	}

	result += "					default : begin\n"
	result += "						$display(\"Unknown Opcode\");\n"
	result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
	result += "					end\n"

	result += "				endcase\n"

	// TODO What are they, maybe needed by some opcode ?
	//	result += "			end\n"

	switch arch.Modes[0] {
	case "ha":
		result += "			// ha placeholder\n"
	case "hy":
		result += "			end\n"
	case "vn":
		result += "				end\n"
		result += "			endcase\n"
	}

	result += "		end\n"
	result += "	end\n"
	result += "	assign rom_bus = _pc;\n"

	for _, op := range proc.Op {
		result += op.Op_instruction_verilog_footer(arch, flavor)
	}

	for i := 0; i < int(proc.N); i++ {
		result += "	assign " + Get_input_name(i) + "_received = " + Get_input_name(i) + "_recv;\n"
	}

	for i := 0; i < int(proc.M); i++ {
		result += "	assign " + Get_output_name(i) + " = _aux" + Get_output_name(i) + ";\n"
		result += "	assign " + Get_output_name(i) + "_valid = " + Get_output_name(i) + "_val;\n"
	}

	result += "endmodule\n"

	doneextramod := make(map[string]bool)

	for _, op := range proc.Op {
		modlist, modcode := op.Op_instruction_verilog_extra_modules(arch, flavor)
		for i, module := range modlist {
			if _, ok := doneextramod[module]; !ok {
				result += modcode[i]
				doneextramod[module] = true
			}
		}

		files, filesCode := op.ExtraFiles(arch)
		for i, file := range files {
			f, _ := os.Create(file)
			f.WriteString(filesCode[i])
			f.Close()
		}
	}

	return result
}
