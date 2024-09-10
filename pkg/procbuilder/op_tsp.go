package procbuilder

// TODO This is the ROM, change it to halndle also the RAM case

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Tsp struct{}

func (op Tsp) Op_get_name() string {
	return "tsp"
}

func (op Tsp) Op_get_desc() string {
	return "Start a thread at a specific location"
}

func (op Tsp) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := ""
	switch arch.Modes[0] {
	case "ha":
		result = "j [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.O)) + "(Location)] [ 8 (nice)] 	// Start a thread at a specific location [" + strconv.Itoa(opBits+int(arch.O)+8) + "]\n"
	case "vn":
		result = "j [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.L)) + "(Location)] [ 8 (nice)] // Start a thread at a specific location [" + strconv.Itoa(opBits+int(arch.L)+8) + "]\n"
	case "hy":
		if arch.O > arch.L {
			result = "j [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.O)) + "(Location)] [ 8 (nice)]	// Start a thread at a specific location [" + strconv.Itoa(opBits+int(arch.O)+8) + "]\n"
		} else {
			result = "j [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.L)) + "(Location)] [ 8 (nice)]	// Start a thread at a specific location [" + strconv.Itoa(opBits+int(arch.L)+8) + "]\n"
		}
	}
	return result
}

func (op Tsp) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	switch arch.Modes[0] {
	case "ha":
		return opbits + int(arch.R) + int(arch.O) + 8 // The bits for the opcode + bits for the register + bits for a location + 8 bits for the nice
	case "vn":
		return opbits + int(arch.R) + int(arch.L) + 8 // The bits for the opcode + bits for the register + bits for a location + 8 bits for the nice
	case "hy":
		if arch.O > arch.L {
			return opbits + int(arch.R) + int(arch.O) + 8 // The bii for the opcode + bits for the register + bits for a location + 8 bits for the nice
		} else {
			return opbits + int(arch.R) + int(arch.L) + 8 // The bits for the opcode + bits for the register + bits for a location + 8 bits for the nice
		}
	}
	return 0
}

func (op Tsp) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	return ""
}

func (op Tsp) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (op Tsp) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op Tsp) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Tsp) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {

	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()
	tabsNum := 5
	regNum := 1 << arch.R
	locationBits := arch.O

	threadStack := "threadStack" + strconv.Itoa(arch.Threaded)
	sm := threadStack + "SM"

	switch arch.Modes[0] {
	case "ha":
		locationBits = arch.O
	case "vn":
		locationBits = arch.L
	case "hy":
		if arch.O > arch.L {
			locationBits = arch.O
		} else {
			locationBits = arch.L
		}
	}

	result := ""
	result += tabs(tabsNum) + "TSP: begin\n"
	if arch.R == 1 {
		result += tabs(tabsNum+1) + "case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + "])\n"
	} else {
		result += tabs(tabsNum+1) + "case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)) + "])\n"
	}
	for i := 0; i < regNum; i++ {

		if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlyDestRegs)) {
			cp := arch.Tag
			req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:tsp", T: bmreqs.ObjectSet, Name: "destregs", Value: Get_register_name(i), Op: bmreqs.OpCheck})
			if req.Value == "false" {
				continue
			}
		}

		result += tabs(tabsNum+2) + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		result += tabs(tabsNum+3) + "case (" + sm + ")\n"
		result += tabs(tabsNum+3) + "CTXEXE: begin\n"
		result += tabs(tabsNum+4) + sm + " <= CTXSEND;\n"
		result += tabs(tabsNum+4) + "provpc <= current_instruction[" + strconv.Itoa(romWord-opBits-1-int(arch.R)) + ":" + strconv.Itoa(romWord-opBits-int(locationBits)-int(arch.R)) + "];\n"

		result += tabs(tabsNum+3) + "end\n"
		result += tabs(tabsNum+3) + "CTXSEND: begin\n"
		result += tabs(tabsNum+5) + "if (!" + threadStack + "senderAck) begin\n"
		regList := ""
		for j := 0; j < regNum; j++ {
			regList += "_" + strings.ToLower(Get_register_name(j)) + ", "
		}
		result += tabs(tabsNum+6) + threadStack + "senderData <= {ThreadID, 1'b0, provpc, " + regList + "current_instruction[" + strconv.Itoa(romWord-opBits-1-int(arch.R)-int(locationBits)) + ":" + strconv.Itoa(romWord-opBits-int(locationBits)-int(arch.R)-8) + "]};\n"
		result += tabs(tabsNum+6) + threadStack + "senderWrite <= 1'b1;\n"
		result += tabs(tabsNum+6) + sm + " <= CTXWSEND;\n"
		result += tabs(tabsNum+5) + "end\n"
		result += tabs(tabsNum+4) + "end\n"
		result += tabs(tabsNum+3) + "CTXWSEND: begin\n"
		result += tabs(tabsNum+4) + "if (" + threadStack + "senderAck) begin\n"
		result += tabs(tabsNum+5) + threadStack + "senderWrite <= 1'b0;\n"
		result += tabs(tabsNum+5) + sm + " <= CTXEXE;\n"
		result += NextInstruction(conf, arch, tabsNum+5, "_pc + 1'b1")

		result += tabs(tabsNum+2) + "\t\tend\n"
		result += tabs(tabsNum+2) + "\tend\n"
		result += tabs(tabsNum+2) + "\tendcase\n"

		result += tabs(tabsNum+2) + "end\n"
	}
	result += tabs(tabsNum+1) + "endcase\n"
	result += tabs(tabsNum) + "end\n"
	return result
}

func (op Tsp) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Tsp) Assembler(arch *Arch, words []string) (string, error) {
	reg_num := 1 << arch.R

	locationBits := arch.O

	switch arch.Modes[0] {
	case "ha":
		locationBits = arch.O
	case "vn":
		locationBits = arch.L
	case "hy":
		if arch.O > arch.L {
			locationBits = arch.O
		} else {
			locationBits = arch.L
		}
	}

	if len(words) != 3 {
		return "", Prerror{"Wrong arguments number"}
	}

	result := ""
	for i := 0; i < reg_num; i++ {
		if words[0] == strings.ToLower(Get_register_name(i)) {
			result += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if partial, err := Process_number(words[1]); err == nil {
		result += zeros_prefix(int(locationBits), partial)
	} else {
		return "", Prerror{err.Error()}
	}

	if partial, err := Process_number(words[2]); err == nil {
		result += zeros_prefix(8, partial)
	} else {
		return "", Prerror{err.Error()}
	}

	return result, nil
}

func (op Tsp) Disassembler(arch *Arch, instr string) (string, error) {
	regId := get_id(instr[:arch.R])

	locationBits := arch.O

	switch arch.Modes[0] {
	case "ha":
		locationBits = arch.O
	case "vn":
		locationBits = arch.L
	case "hy":
		if arch.O > arch.L {
			locationBits = arch.O
		} else {
			locationBits = arch.L
		}
	}
	result := strings.ToLower(Get_register_name(regId)) + " "
	value := get_id(instr[arch.R : int(arch.R)+int(locationBits)])
	result += strconv.Itoa(value) + " "
	value = get_id(instr[int(arch.R)+int(locationBits) : int(arch.R)+int(locationBits)+8])
	result += strconv.Itoa(value)
	return result, nil
}

func (op Tsp) Simulate(vm *VM, instr string) error {
	value := get_id(instr[:vm.Mach.O])
	if value < len(vm.Mach.Slocs) {
		vm.Pc = uint64(value)
	} else {
		vm.Pc = vm.Pc + 1
	}
	return nil
}

func (op Tsp) Generate(arch *Arch) string {
	max_value := 1 << arch.O
	value := rand.Intn(max_value)
	return zeros_prefix(int(arch.O), get_binary(value))
}

func (op Tsp) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Tsp) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Tsp) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Tsp) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (op Tsp) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "tsp", I_NIL}
	result[0] = newnot
	return result, nil
}

func (op Tsp) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (op Tsp) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 4)
	result[0] = "tsp::*--type=reg::*--type=number::*--type=number"
	result[1] = "tsp::*--type=reg::*--type=symbol::*--type=number"
	return result
}
func (op Tsp) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "tsp":
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (op Tsp) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (op Tsp) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
