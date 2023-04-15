package procbuilder

import (
	"math/rand"
	"time"
)

// The machine is an architecture provided with and execution code and an intial state
type Machine struct {
	Arch
	Program
}

type Machine_json struct {
	Modes              []string
	Rsize              uint8
	R                  uint8 // Number of n-bit registers
	N                  uint8 // Number of n-bit inputs
	M                  uint8 // Number of n-bit outputs
	L                  uint8 // Number of n-bit memory banks internal to the processor 2^
	O                  uint8 // Number of ROM cells 2^
	Shared_constraints string
	Op                 []string
	Slocs              []string
}

var Allopcodes []Opcode
var Allshared []Sharedel
var AllDynamicalInstructions []DynamicInstruction

func init() {
	rand.Seed(int64(time.Now().Unix()))

	// Keep in mind the lists og opcodes has to be kept ordered by name
	Allopcodes = make([]Opcode, 0)

	Allopcodes = append(Allopcodes, Adc{})
	Allopcodes = append(Allopcodes, Add{})
	Allopcodes = append(Allopcodes, Addf{})
	Allopcodes = append(Allopcodes, Addi{})
	Allopcodes = append(Allopcodes, Addp{})
	Allopcodes = append(Allopcodes, And{})
	Allopcodes = append(Allopcodes, Chc{})
	Allopcodes = append(Allopcodes, Chw{})
	Allopcodes = append(Allopcodes, Cil{})
	Allopcodes = append(Allopcodes, Cilc{})
	Allopcodes = append(Allopcodes, Cir{})
	Allopcodes = append(Allopcodes, Cirn{})
	Allopcodes = append(Allopcodes, Clc{})
	Allopcodes = append(Allopcodes, Clr{})
	Allopcodes = append(Allopcodes, Cpy{})
	Allopcodes = append(Allopcodes, Cset{})
	Allopcodes = append(Allopcodes, Dec{})
	Allopcodes = append(Allopcodes, Div{})
	Allopcodes = append(Allopcodes, Divf{})
	Allopcodes = append(Allopcodes, Divp{})
	Allopcodes = append(Allopcodes, Dpc{})
	Allopcodes = append(Allopcodes, Expf{})
	Allopcodes = append(Allopcodes, Hit{})
	Allopcodes = append(Allopcodes, Hlt{})
	Allopcodes = append(Allopcodes, I2r{})
	Allopcodes = append(Allopcodes, I2rw{})
	Allopcodes = append(Allopcodes, Incc{})
	Allopcodes = append(Allopcodes, Inc{})
	Allopcodes = append(Allopcodes, Jc{})
	Allopcodes = append(Allopcodes, Je{})
	Allopcodes = append(Allopcodes, Jgt0f{})
	// TODO: planned Allopcodes = append(Allopcodes, Jlt{})
	// TODO: planned Allopcodes = append(Allopcodes, Jlte{})
	// TODO: planned Allopcodes = append(Allopcodes, Jr{})
	Allopcodes = append(Allopcodes, Jz{})
	Allopcodes = append(Allopcodes, J{})
	Allopcodes = append(Allopcodes, Lfsr82r{})
	// TODO: planned Allopcodes = append(Allopcodes, Lfsr162r{})
	Allopcodes = append(Allopcodes, M2r{})
	Allopcodes = append(Allopcodes, Mod{})
	Allopcodes = append(Allopcodes, Mulc{})
	Allopcodes = append(Allopcodes, Mult{})
	Allopcodes = append(Allopcodes, Multf{})
	Allopcodes = append(Allopcodes, Multp{})
	Allopcodes = append(Allopcodes, Nand{})
	Allopcodes = append(Allopcodes, Nop{})
	Allopcodes = append(Allopcodes, Nor{})
	Allopcodes = append(Allopcodes, Not{})
	Allopcodes = append(Allopcodes, Or{})
	Allopcodes = append(Allopcodes, Q2r{})
	Allopcodes = append(Allopcodes, R2m{})
	Allopcodes = append(Allopcodes, R2o{})
	Allopcodes = append(Allopcodes, R2owa{})
	Allopcodes = append(Allopcodes, R2owaa{})
	Allopcodes = append(Allopcodes, R2q{})
	Allopcodes = append(Allopcodes, R2s{})
	Allopcodes = append(Allopcodes, R2v{})
	Allopcodes = append(Allopcodes, R2vri{})
	Allopcodes = append(Allopcodes, R2t{})
	Allopcodes = append(Allopcodes, R2u{})
	Allopcodes = append(Allopcodes, Ro2r{})
	Allopcodes = append(Allopcodes, Ro2rri{})
	Allopcodes = append(Allopcodes, Rsc{})
	Allopcodes = append(Allopcodes, Rset{})
	Allopcodes = append(Allopcodes, Sic{})
	Allopcodes = append(Allopcodes, S2r{})
	Allopcodes = append(Allopcodes, Saj{})
	Allopcodes = append(Allopcodes, Sbc{})
	Allopcodes = append(Allopcodes, Sub{})
	Allopcodes = append(Allopcodes, T2r{})
	Allopcodes = append(Allopcodes, U2r{})
	Allopcodes = append(Allopcodes, Wrd{})
	Allopcodes = append(Allopcodes, Wwr{})
	Allopcodes = append(Allopcodes, Xnor{})
	Allopcodes = append(Allopcodes, Xor{})

	AllDynamicalInstructions = make([]DynamicInstruction, 0)
	AllDynamicalInstructions = append(AllDynamicalInstructions, DynFloPoCo{})

	Allshared = make([]Sharedel, 0)
	Allshared = append(Allshared, Sharedmem{})
	Allshared = append(Allshared, Channel{})
	Allshared = append(Allshared, Barrier{})
	Allshared = append(Allshared, Lfsr8{})
	Allshared = append(Allshared, Queue{})
	Allshared = append(Allshared, Stack{})
	Allshared = append(Allshared, Uart{})
	Allshared = append(Allshared, Kbd{})
	Allshared = append(Allshared, Vtextmem{})
}

func (mach *Machine) String() string {
	result := ""
	result += mach.Arch.String()
	result += mach.Program.String()
	return result
}

func (mach *Machine) Descr() string {
	result := ""
	result += mach.Arch.String()
	return result
}

func (mach *Machine) Jsoner() *Machine_json {
	result := new(Machine_json)
	result.Modes = make([]string, len(mach.Modes))
	for i, val := range mach.Modes {
		result.Modes[i] = val
	}
	result.Rsize = mach.Rsize
	result.R = mach.R
	result.N = mach.N
	result.M = mach.M
	result.L = mach.L
	result.O = mach.O
	result.Shared_constraints = mach.Shared_constraints
	result.Slocs = make([]string, len(mach.Slocs))
	for i, val := range mach.Slocs {
		result.Slocs[i] = val
	}
	result.Op = make([]string, len(mach.Op))
	for i, val := range mach.Op {
		result.Op[i] = val.Op_get_name()
	}
	return result
}

func (machj *Machine_json) Dejsoner() *Machine {
	result := new(Machine)
	result.Modes = make([]string, len(machj.Modes))
	for i, val := range machj.Modes {
		result.Modes[i] = val
	}
	result.Rsize = machj.Rsize
	result.R = machj.R
	result.N = machj.N
	result.M = machj.M
	result.L = machj.L
	result.O = machj.O
	result.Shared_constraints = machj.Shared_constraints
	result.Slocs = make([]string, len(machj.Slocs))
	for i, val := range machj.Slocs {
		result.Slocs[i] = val
	}
	result.Op = make([]Opcode, len(machj.Op))
	for i, opname := range machj.Op {

		EventuallyCreateInstruction(opname)

		for _, op := range Allopcodes {
			if op.Op_get_name() == opname {
				result.Op[i] = op
			}
		}
	}
	return result
}

func (mach *Machine) Constraint_check() (string, bool) {
	result := ""
	shared := make([]string, 0)
	required := make([]string, 0)
	forbidden := make([]string, 0)

	for _, op := range mach.Op {
		if present, modes := op.Required_modes(); present {
			for _, mode := range modes {
				checkpres := false
				for _, curr := range required {
					if curr == mode {
						checkpres = true
						break
					}
				}
				if !checkpres {
					required = append(required, mode)
					result += "Added required mode: " + mode + "\n"
				}
			}
		}
		if present, modes := op.Forbidden_modes(); present {
			for _, mode := range modes {
				checkpres := false
				for _, curr := range forbidden {
					if curr == mode {
						checkpres = true
						break
					}
				}
				if !checkpres {
					forbidden = append(forbidden, mode)
					result += "Added forbidden mode: " + mode + "\n"

				}
			}

		}
		if present, sos := op.Required_shared(); present {
			for _, so := range sos {
				checkpres := false
				for _, curr := range shared {
					if curr == so {
						checkpres = true
						break
					}
				}
				if !checkpres {
					shared = append(shared, so)
					result += "Added required shared object: " + so + "\n"

				}
			}

		}
	}

	// Checking required modes constraints
	// TODO Finish
	for _, mode := range required {
		switch mode {
		case "ha":
		case "vn":
		case "hy":
		case "ramabs":
		case "romabs":
		case "ramind":
			if mach.Rsize != mach.L {
				return "Register size nd RAM depth mismatch", false
			}
		case "romind":
			if mach.Rsize != mach.O {
				return "Register size nd ROM depth mismatch", false
			}
		}
	}

	// Checking conflitting modes
	for _, mode := range forbidden {
		for _, req := range required {
			if mode == req {
				return "", false
			}
		}
	}

	// TODO Finish shared checks

	return result, true
}
