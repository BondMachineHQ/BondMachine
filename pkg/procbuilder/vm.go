package procbuilder

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/simbox"
)

type VM struct {
	CpID      uint32
	Mach      *Machine
	Registers []interface{}
	Memory    []interface{}
	Inputs    []interface{}
	Outputs   []interface{}

	InputsValid  []bool
	OutputsValid []bool

	InputsRecv  []bool
	OutputsRecv []bool

	Pc           uint64
	Extra_states map[string]interface{}
	CmdChan      chan []byte
}

func (vm *VM) CopyState(vmsource *VM) {
	for i, reg := range vmsource.Registers {
		vm.Registers[i] = reg
	}
	for i, inp := range vmsource.Inputs {
		vm.Inputs[i] = inp
	}
	for i, outp := range vmsource.Outputs {
		vm.Outputs[i] = outp
	}
	// TODO FINISH
}

// Simbox rules are converted in a sim drive when the simulation starts and applied during the simulation
type Sim_tick_set map[int]interface{}
type Sim_drive struct {
	Injectables []*interface{}
	AbsSet      map[uint64]Sim_tick_set
}

// This is initializated when the simulation starts and filled on the way
type Sim_tick_get map[int]interface{}
type Sim_report struct {
	Reportables []*interface{}
	AbsGet      map[uint64]Sim_tick_get
}

type Sim_config struct {
	Show_pc          bool
	Show_instruction bool
	Show_disasm      bool
	Show_regs_pre    bool
	Show_regs_post   bool
	Show_io_pre      bool
	Show_io_post     bool
}

func (vm *VM) Init() error {
	if vm.Mach == nil {
		return Prerror{"Cannot initialize the nil machine"}
	}

	if vm.Mach.R == 0 {
		return Prerror{"Cannot initialize a machine with 0 registers"}
	}

	if len(vm.Mach.Program.Slocs) == 0 {
		return Prerror{"Cannot initialize a machine with 0 lines of program"}
	}

	//TODO Other checks

	reg_num := 1 << vm.Mach.R

	mem_num := 1 << vm.Mach.L

	vm.Registers = make([]interface{}, reg_num)
	vm.Memory = make([]interface{}, mem_num)
	vm.Inputs = make([]interface{}, vm.Mach.N)
	vm.Outputs = make([]interface{}, vm.Mach.M)

	vm.InputsValid = make([]bool, vm.Mach.N)
	vm.OutputsValid = make([]bool, vm.Mach.M)
	vm.InputsRecv = make([]bool, vm.Mach.N)
	vm.OutputsRecv = make([]bool, vm.Mach.M)

	vm.Pc = 0

	if vm.Mach.Rsize <= 8 {
		for i := 0; i < reg_num; i++ {
			vm.Registers[i] = uint8(0)
		}
		for i := 0; i < mem_num; i++ {
			vm.Memory[i] = uint8(0)
		}
		for i := 0; i < int(vm.Mach.N); i++ {
			vm.Inputs[i] = uint8(0)
		}
		for i := 0; i < int(vm.Mach.M); i++ {
			vm.Outputs[i] = uint8(0)
		}
	} else if vm.Mach.Rsize <= 16 {
		for i := 0; i < reg_num; i++ {
			vm.Registers[i] = uint16(0)
		}
		for i := 0; i < mem_num; i++ {
			vm.Memory[i] = uint16(0)
		}
		for i := 0; i < int(vm.Mach.N); i++ {
			vm.Inputs[i] = uint16(0)
		}
		for i := 0; i < int(vm.Mach.M); i++ {
			vm.Outputs[i] = uint16(0)
		}
	} else if vm.Mach.Rsize <= 32 {
		for i := 0; i < reg_num; i++ {
			vm.Registers[i] = uint32(0)
		}
		for i := 0; i < mem_num; i++ {
			vm.Memory[i] = uint32(0)
		}
		for i := 0; i < int(vm.Mach.N); i++ {
			vm.Inputs[i] = uint32(0)
		}
		for i := 0; i < int(vm.Mach.M); i++ {
			vm.Outputs[i] = uint32(0)
		}
	} else if vm.Mach.Rsize <= 64 {
		for i := 0; i < reg_num; i++ {
			vm.Registers[i] = uint64(0)
		}
		for i := 0; i < mem_num; i++ {
			vm.Memory[i] = uint64(0)
		}
		for i := 0; i < int(vm.Mach.N); i++ {
			vm.Inputs[i] = uint64(0)
		}
		for i := 0; i < int(vm.Mach.M); i++ {
			vm.Outputs[i] = uint64(0)
		}
	} else {
		return Prerror{"Cannot initialize a machine with a register size greater than 64"}
	}

	vm.Extra_states = make(map[string]interface{})

	return nil
}

func (vm *VM) Step(psc *Sim_config) (string, error) {
	result := ""

	if psc != nil {
		if psc.Show_pc {
			result += "\t\tPC: " + strconv.Itoa(int(vm.Pc)) + "\n"
		}
	}

	//	reg_num := 1 << vm.Mach.R
	num_instr := len(vm.Mach.Program.Slocs)
	opBits := vm.Mach.Opcodes_bits()

	if int(vm.Pc) > num_instr {
		return "", Prerror{"Program counter outside limits"}
	}

	if int(vm.Pc) == num_instr {
		// Halt computation
		// vm.Pc = 0
	} else {
		instr := vm.Mach.Program.Slocs[vm.Pc]

		if psc != nil {
			if psc.Show_instruction {
				result += "\t\tInstr: " + instr + "\n"
			}
		}

		if opcode_id, err := vm.Mach.Conproc.Decode_opcode(instr); err == nil {
			op := vm.Mach.Arch.Conproc.Op[opcode_id]

			if psc != nil {
				if psc.Show_disasm {
					curline := "\t\tDisasm: " + op.Op_get_name() + " "
					if disas, err := op.Disassembler(&vm.Mach.Arch, instr[opBits:]); err != nil {
						return "", Prerror{"Disassembling falied"}
					} else {
						result += curline + disas + "\n"
					}
				}
			}
			if psc != nil {
				if psc.Show_io_pre {
					result += "\t\tPre-compute IO: " + vm.DumpIO() + "\n"
				}
				if psc.Show_regs_pre {
					result += "\t\tPre-compute Regs: " + vm.DumpRegisters() + "\n"
				}
			}

			if err := op.Simulate(vm, instr[opBits:]); err != nil {
				return "", Prerror{"Simulation failed"}
			}

			if psc != nil {
				if psc.Show_io_pre {
					result += "\t\tPost-compute IO: " + vm.DumpIO() + "\n"
				}
				if psc.Show_regs_pre {
					result += "\t\tPost-compute Regs: " + vm.DumpRegisters() + "\n"
				}
			}

		} else {
			return "", Prerror{"Unknown opcode"}
		}
	}

	return result, nil
}

func (vm *VM) DumpRegisters() string {
	result := ""
	for i, reg := range vm.Registers {
		if vm.Mach.Rsize <= 8 {
			result = result + Get_register_name(i) + ": " + zeros_prefix(int(vm.Mach.Rsize), get_binary(int(reg.(uint8)))) + " "
		} else if vm.Mach.Rsize <= 16 {
			result = result + Get_register_name(i) + ": " + zeros_prefix(int(vm.Mach.Rsize), get_binary(int(reg.(uint16)))) + " "
		} else if vm.Mach.Rsize <= 32 {
			result = result + Get_register_name(i) + ": " + zeros_prefix(int(vm.Mach.Rsize), get_binary(int(reg.(uint32)))) + " "
		} else if vm.Mach.Rsize <= 64 {
			result = result + Get_register_name(i) + ": " + zeros_prefix(int(vm.Mach.Rsize), get_binary(int(reg.(uint64)))) + " "
		} else {
			result = "Cannot dump registers for a machine with a register size greater than 64"
		}
	}
	return result
}

func (vm *VM) DumpIO() string {
	result := ""
	for i, reg := range vm.Inputs {
		switch vm.Mach.Rsize {
		case 8:
			result = result + Get_input_name(i) + ": " + zeros_prefix(int(vm.Mach.Rsize), get_binary(int(reg.(uint8)))) + " "
		case 16:
			result = result + Get_input_name(i) + ": " + zeros_prefix(int(vm.Mach.Rsize), get_binary(int(reg.(uint16)))) + " "
		case 32:
			result = result + Get_input_name(i) + ": " + zeros_prefix(int(vm.Mach.Rsize), get_binary(int(reg.(uint32)))) + " "
		case 64:
			result = result + Get_input_name(i) + ": " + zeros_prefix(int(vm.Mach.Rsize), get_binary(int(reg.(uint64)))) + " "
		default:
			result = "unsupported register size"
		}
	}
	for i, reg := range vm.Outputs {
		switch vm.Mach.Rsize {
		case 8:
			result = result + Get_output_name(i) + ": " + zeros_prefix(int(vm.Mach.Rsize), get_binary(int(reg.(uint8)))) + " "
		case 16:
			result = result + Get_output_name(i) + ": " + zeros_prefix(int(vm.Mach.Rsize), get_binary(int(reg.(uint16)))) + " "
		case 32:
			result = result + Get_output_name(i) + ": " + zeros_prefix(int(vm.Mach.Rsize), get_binary(int(reg.(uint32)))) + " "
		case 64:
			result = result + Get_output_name(i) + ": " + zeros_prefix(int(vm.Mach.Rsize), get_binary(int(reg.(uint64)))) + " "
		default:
			result = "unsupported register size"
		}
	}
	return result
}

func (vm *VM) GetElementLocation(mnemonic string) (*interface{}, error) {
	// TODO include others
	if len(mnemonic) > 1 && mnemonic[0] == 'i' {
		if i, err := strconv.Atoi(mnemonic[1:]); err == nil {
			if i < len(vm.Inputs) {
				return &vm.Inputs[i], nil
			}
		}
	}
	return nil, Prerror{mnemonic + " unknown"}
}

func (sc *Sim_config) Init(s *simbox.Simbox, vm *VM) error {

	if s != nil {

		for _, rule := range s.Rules {
			// Intercept the set rules
			if rule.Timec == simbox.TIMEC_NONE && rule.Action == simbox.ACTION_CONFIG {
				switch rule.Object {
				case "show_pc":
					sc.Show_pc = true
				case "show_instruction":
					sc.Show_instruction = true
				case "show_disasm":
					sc.Show_disasm = true
				case "show_proc_regs_pre":
					sc.Show_regs_pre = true
				case "show_proc_regs_post":
					sc.Show_regs_post = true
				case "show_proc_io_pre":
					sc.Show_io_pre = true
				case "show_proc_io_post":
					sc.Show_io_post = true
				}
			}
		}
	}
	return nil
}

func (sd *Sim_drive) Init(s *simbox.Simbox, vm *VM) error {

	if s != nil {

		inj := make([]*interface{}, 0)
		act := make(map[uint64]Sim_tick_set)

		for _, rule := range s.Rules {
			fmt.Println(rule)
			// Intercept the set rules
			if rule.Timec == simbox.TIMEC_ABS && rule.Action == simbox.ACTION_SET {
				if loc, err := vm.GetElementLocation(rule.Object); err == nil {
					if val, err := strconv.Atoi(rule.Extra); err == nil {
						ipos := -1
						for i, iloc := range inj {
							if iloc == loc {
								ipos = i
								break
							}
						}
						if ipos == -1 {
							ipos = len(inj)
							inj = append(inj, loc)
						}

						if act_on_tick, ok := act[rule.Tick]; ok {
							if vm.Mach.Rsize <= 8 {
								act_on_tick[ipos] = uint8(val)
							} else if vm.Mach.Rsize <= 16 {
								act_on_tick[ipos] = uint16(val)
							} else if vm.Mach.Rsize <= 32 {
								act_on_tick[ipos] = uint32(val)
							} else if vm.Mach.Rsize <= 64 {
								act_on_tick[ipos] = uint64(val)
							} else {
								return errors.New("unsupported register size")
							}
						} else {
							act_on_tick := make(map[int]interface{})
							if vm.Mach.Rsize <= 8 {
								act_on_tick[ipos] = uint8(val)
							} else if vm.Mach.Rsize <= 16 {
								act_on_tick[ipos] = uint16(val)
							} else if vm.Mach.Rsize <= 32 {
								act_on_tick[ipos] = uint32(val)
							} else if vm.Mach.Rsize <= 64 {
								act_on_tick[ipos] = uint64(val)
							} else {
								return errors.New("unsupported register size")

							}
							act[rule.Tick] = act_on_tick
						}
					} else {
						return err
					}
				} else {
					return err
				}
			}
		}

		sd.Injectables = inj
		sd.AbsSet = act
	}

	return nil
}

func (sd *Sim_report) Init(s *simbox.Simbox, vm *VM) error {

	if s != nil {

		rep := make([]*interface{}, 0)
		str := make(map[uint64]Sim_tick_get)

		for _, rule := range s.Rules {
			fmt.Println(rule)
			// Intercept the set rules
			if rule.Timec == simbox.TIMEC_ABS && rule.Action == simbox.ACTION_GET {
				if loc, err := vm.GetElementLocation(rule.Object); err == nil {
					ipos := -1
					for i, iloc := range rep {
						if iloc == loc {
							ipos = i
							break
						}
					}
					if ipos == -1 {
						ipos = len(rep)
						rep = append(rep, loc)
					}

					if str_on_tick, ok := str[rule.Tick]; ok {
						if vm.Mach.Rsize <= 8 {
							str_on_tick[ipos] = uint8(0)
						} else if vm.Mach.Rsize <= 16 {
							str_on_tick[ipos] = uint16(0)
						} else if vm.Mach.Rsize <= 32 {
							str_on_tick[ipos] = uint32(0)
						} else if vm.Mach.Rsize <= 64 {
							str_on_tick[ipos] = uint64(0)
						} else {
							return errors.New("unsupported register size")
						}
					} else {
						str_on_tick := make(map[int]interface{})
						if vm.Mach.Rsize <= 8 {
							str_on_tick[ipos] = uint8(0)
						} else if vm.Mach.Rsize <= 16 {
							str_on_tick[ipos] = uint16(0)
						} else if vm.Mach.Rsize <= 32 {
							str_on_tick[ipos] = uint32(0)
						} else if vm.Mach.Rsize <= 64 {
							str_on_tick[ipos] = uint64(0)
						} else {
							return errors.New("unsupported register size")
						}
						str[rule.Tick] = str_on_tick
					}
				} else {
					return err
				}
			}
		}

		sd.Reportables = rep
		sd.AbsGet = str

	}

	return nil
}
