package procbuilder

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// The architecture
type Arch struct {
	Modes []string
	Conproc
	Rom
	Ram
	Shared_constraints string
	Tag                string
}

func (arch *Arch) Shared_num(soname string) int {
	so := strings.Split(arch.Shared_constraints, ",")
	counted := 0
	for _, sos := range so {
		splitted := strings.Split(sos, ":")
		if len(splitted) == 2 {
			if splitted[0] == soname {
				counted++
			}
		}
	}
	return counted
}

func (arch *Arch) Shared_bits(soname string) int {
	counted := arch.Shared_num(soname)
	if counted == 0 {
		return 0
	}
	served := 1
	for bits := 1; bits < 256; bits++ {
		if served<<uint8(bits) >= counted {
			return bits
		}
	}
	return 0
}

func (arch *Arch) Shared_depth(soname string, so_id int) int {
	// TODO
	served := 1
	for bits := 1; bits < 16; bits++ {
		if served<<uint8(bits) >= int(arch.N) {
			return bits
		}
	}
	return 1
}

func (arch *Arch) String() string {
	result := ""
	result += arch.Conproc.String()
	result += arch.Rom.String()
	return result
}

func (arch *Arch) Max_word() int {
	now := 1
	for _, op := range arch.Op {
		neww := op.Op_get_instruction_len(arch)
		if neww > now {
			now = neww
		}
	}

	//	served := 1
	//	for bits := 1; bits < 16; bits++ {
	//		if served<<uint8(bits) >= now {
	//			return int(served << uint8(bits))
	//		}
	//	}

	return now
}

func (arch *Arch) Write_verilog(arch_module_name string, modules_names map[string]string, flavor string) string {
	regsize := int(arch.Rsize)
	rom_word := arch.Max_word()
	//inbits := arch.Inputs_bits()
	//outbits := arch.Outputs_bits()

	result := ""

	// Module header
	result += "`timescale 1ns/1ps\n"
	result += "module " + arch_module_name + "(clock_signal, reset_signal"

	ioh := ""
	for i := 0; i < int(arch.N); i++ {
		ioh += ", " + Get_input_name(i) + ", " + Get_input_name(i) + "_valid , " + Get_input_name(i) + "_received"
	}

	for i := 0; i < int(arch.M); i++ {
		ioh += ", " + Get_output_name(i) + ", " + Get_output_name(i) + "_valid, " + Get_output_name(i) + "_received"
	}

	result += ioh

	header := ""
	if arch.Shared_constraints != "" {
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
					header += so.Get_header(arch, constraint, seq[soname])
				}
			}
		}
		result += header
	}

	result += ");\n\n"

	// Header variables declarations
	result += "\tinput clock_signal;\n"
	result += "\tinput reset_signal;\n"

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
					params += so.Get_params(arch, constraint, seq[soname])
					params += "\n"
				}
			}
		}
		result += params
	}

	// Local parameters
	result += "\twire [" + strconv.Itoa(int(arch.O)-1) + ":0] rom_bus;\n"
	result += "\twire [" + strconv.Itoa(int(rom_word)-1) + ":0] rom_value;\n"

	result += "\n"

	ramh := ""
	if int(arch.L) != 0 {
		for _, ramcomp := range []string{"din", "dout", "addr", "wren", "en"} {
			ramh += ", " + arch_module_name + ramcomp
		}
		result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + arch_module_name + "din;\n"
		result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + arch_module_name + "dout;\n"
		result += "\twire [" + strconv.Itoa(int(arch.L)-1) + ":0] " + arch_module_name + "addr;\n"
		result += "\twire " + arch_module_name + "wren;\n"
		result += "\twire " + arch_module_name + "en;\n"
	}

	// Architecture components
	procname := modules_names["processor"]
	romname := modules_names["rom"]
	ramname := modules_names["ram"]

	result += "\n"

	// Processor
	result += "\t" + procname + " " + procname + "_instance(clock_signal, reset_signal, rom_bus, rom_value" + ramh + ioh + header + ");\n"

	// Rom
	result += "\t" + romname + " " + romname + "_instance(rom_bus, rom_value);\n"

	// Ram
	if int(arch.L) != 0 {
		result += "\t" + ramname + " " + ramname + "_instance(clock_signal, reset_signal" + ramh + ");\n"
	}

	result += "\n"

	result += "endmodule\n"
	return result
}

func (arch *Arch) Show_assembler() string {
	result := ""

	for _, op := range arch.Op {
		result += op.Op_show_assembler(arch)
	}

	return result
}

func (arch *Arch) Assembler_process_line(line []byte) (string, error) {
	currline := strings.ToLower(string(line))
	words := strings.Fields(currline)
	opbits := arch.Opcodes_bits()

	if len(words) != 0 {
		if words[0][0] != '#' {
			for i, op := range arch.Op {
				if op.Op_get_name() == words[0] {
					if result, err := op.Assembler(arch, words[1:]); err == nil {
						return zeros_prefix(opbits, get_binary(i)) + result, nil
					} else {
						return "", Prerror{err.Error() + ", error processing " + op.Op_get_name()}
					}
				}
			}
		} else {
			return "", nil
		}
	} else {
		return "", nil
	}

	return "", Prerror{"Unknown Opcode"}
}

func (arch *Arch) Assembler(inp []byte) (Program, error) {
	// TODO keep in mind this
	currline := make([]byte, 256)

	maxlines := make([]string, int(1<<arch.O))
	iline := 0
	j := 0
	impline := 1
	for _, ch := range inp {
		if ch == 10 {
			//fmt.Println(currline[0:iline])
			//fmt.Println(string(currline[0:iline]))
			if result, err := arch.Assembler_process_line(currline[0:iline]); err == nil {
				if result != "" {
					maxlines[j] = result
					j = j + 1
					impline = impline + 1
				}
				iline = 0
			} else {
				return Program{}, Prerror{err.Error() + " on line " + strconv.Itoa(impline)}
			}
		} else {
			currline[iline] = ch
			iline = iline + 1
		}
	}

	lines := make([]string, impline-1)
	for i, _ := range lines {
		lines[i] = maxlines[i]
	}

	//	for i := j; i < int(1<<arch.O); i++ {
	//		result, _ := arch.Assembler_process_line([]byte{'n', 'o', 'p'})
	//		maxlines[i] = result
	//	}

	return Program{lines}, nil
}

func (mach *Machine) Disassembler() (string, error) {
	result := ""
	opbits := mach.Opcodes_bits()
	for _, instr := range mach.Program.Slocs {

		curline := ""

		if opcode_id, err := mach.Conproc.Decode_opcode(instr); err == nil {
			op := mach.Arch.Conproc.Op[opcode_id]
			curline = curline + op.Op_get_name() + " "
			if part, err := op.Disassembler(&mach.Arch, instr[opbits:]); err == nil {
				curline = curline + part
				result += curline + "\n"
			} else {
				return "", Prerror{"Error dissasembling"}
			}
		} else {
			return "", Prerror{"Unknown opecode"}
		}
	}
	return result, nil
}

func (mach *Machine) Program_alias() (string, error) {
	aliases := make(map[string]string)

	result := ""
	opbits := mach.Opcodes_bits()
	for _, instr := range mach.Program.Slocs {

		curline := ""

		if opcode_id, err := mach.Conproc.Decode_opcode(instr); err == nil {
			op := mach.Arch.Conproc.Op[opcode_id]
			curline = curline + op.Op_get_name() + " "
			if part, err := op.Disassembler(&mach.Arch, instr[opbits:]); err == nil {
				curline = curline + part
				if _, ok := aliases[instr]; !ok {
					result += instr + " " + curline + "\n"
					aliases[instr] = curline
				}
			} else {
				return "", Prerror{"Error dissasembling"}
			}
		} else {
			return "", Prerror{"Unknown opcode"}
		}
	}
	return result, nil
}

func (mach *Machine) Instructions_alias() (string, error) {
	opbits := mach.Opcodes_bits()

	result := ""
	for i, op := range mach.Op {
		result += zeros_prefix(opbits, get_binary(i)) + " " + op.Op_get_name() + "\n"
	}

	return result, nil
}

func (arch *Arch) Program_generate() Program {
	mem_size := 1 << arch.O
	progr_lenght := rand.Intn(mem_size-1) + 1
	lines := make([]string, progr_lenght)
	opcodes := len(arch.Op)
	opbits := arch.Opcodes_bits()
	word_size := arch.Max_word()
	fmt.Println(word_size)
	for i := 0; i < progr_lenght; i++ {
		opcode := rand.Intn(opcodes)
		prefix := zeros_prefix(opbits, get_binary(opcode))
		opcode_code := arch.Op[opcode].Generate(arch)
		line := zeros_suffix(word_size, prefix+opcode_code)
		lines[i] = line
	}

	return Program{lines}
}
