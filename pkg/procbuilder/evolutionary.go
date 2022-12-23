package procbuilder

import (
	//"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/mmirko/mel"
)

func (mach *Machine) MelInit(ep *mel.EvolutionParameters) {

	myarch := new(Arch)

	var eops []string

	if value, ok := ep.GetValue("procbuilder:opcodes"); ok {
		eops = strings.Split(value, ",")
	} else {
		// If not specified populate the opcodes with only the simulable (the ones not needing shared objects)
		eops = make([]string, 0)
		for _, op := range Allopcodes {
			if unsim, _ := op.Required_shared(); !unsim {
				eops = append(eops, op.Op_get_name())
			}
		}
	}

	opcodes := make([]Opcode, len(eops))

	for i, opname := range eops {
		for _, op := range Allopcodes {
			if op.Op_get_name() == opname {
				opcodes[i] = op
			}
		}
	}

	sort.Sort(ByName(opcodes))

	myarch.Op = opcodes

	if value, ok := ep.GetValue("procbuilder:rsize"); ok {
		valuei, _ := strconv.Atoi(value)
		myarch.Rsize = uint8(valuei)
	} else {
		myarch.Rsize = uint8(8)
	}

	if value, ok := ep.GetValue("procbuilder:r"); ok {
		valuei, _ := strconv.Atoi(value)
		myarch.R = uint8(valuei)
	} else {
		myarch.R = uint8(3)
	}

	if value, ok := ep.GetValue("procbuilder:l"); ok {
		valuei, _ := strconv.Atoi(value)
		myarch.L = uint8(valuei)
	} else {
		myarch.L = uint8(8)
	}

	if value, ok := ep.GetValue("procbuilder:n"); ok {
		valuei, _ := strconv.Atoi(value)
		myarch.N = uint8(valuei)
	} else {
		myarch.N = uint8(0)
	}

	if value, ok := ep.GetValue("procbuilder:m"); ok {
		valuei, _ := strconv.Atoi(value)
		myarch.M = uint8(valuei)
	} else {
		myarch.M = uint8(0)
	}

	if value, ok := ep.GetValue("procbuilder:o"); ok {
		valuei, _ := strconv.Atoi(value)
		myarch.O = uint8(valuei)
	} else {
		myarch.O = uint8(8)
	}

	mach.Arch = *myarch
	//fmt.Println(*myarch)
}

func (mach *Machine) MelCopy() mel.Me3li {
	newmach := new(Machine)
	newmach.Arch = mach.Arch
	newmach.Arch.Op = make([]Opcode, len(mach.Arch.Op))
	for i, op := range mach.Arch.Op {
		newmach.Arch.Op[i] = op
	}
	newmach.Program.Slocs = make([]string, len(mach.Program.Slocs))
	for i, sloc := range mach.Program.Slocs {
		newmach.Program.Slocs[i] = sloc
	}
	return newmach
}

// These 3 are the main functions to handle program evaolution, in these cases the machines have the same ISA

func Machine_Program_Generate(ep *mel.EvolutionParameters) mel.Me3li {
	var result mel.Me3li
	var eobj *Machine
	eobj = new(Machine)
	eobj.MelInit(ep)
	eobj.Program = eobj.Program_generate()
	result = eobj
	return result
}

func Machine_Program_Mutate(p mel.Me3li, ep *mel.EvolutionParameters) mel.Me3li {
	var result mel.Me3li
	return result
}

func Machine_Program_Crossover(p mel.Me3li, q mel.Me3li, ep *mel.EvolutionParameters) mel.Me3li {
	var result mel.Me3li
	return result
}
