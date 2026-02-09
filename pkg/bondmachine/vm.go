package bondmachine

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
	"github.com/BondMachineHQ/BondMachine/pkg/simbox"
)

const (
	EVENTONVALID = uint8(0) + iota
	EVENTONRECV
	EVENTONCHANGE
	EVENTONEXIT
)

type VM struct {
	Bmach                 *Bondmachine
	Processors            []*procbuilder.VM
	Inputs_regs           []interface{}
	Outputs_regs          []interface{}
	Internal_inputs_regs  []interface{}
	Internal_outputs_regs []interface{}

	InputsValid          []bool
	OutputsValid         []bool
	InternalInputsValid  []bool
	InternalOutputsValid []bool

	InputsRecv          []bool
	OutputsRecv         []bool
	InternalInputsRecv  []bool
	InternalOutputsRecv []bool

	DeferredInstructions map[string]DeferredInstruction

	SimDelayMap *simbox.SimDelays

	EmuDrivers []EmuDriver
	cmdChan    chan []byte

	Emulating bool

	send_chans   []chan int
	result_chans []chan string
	recv_chan    chan int

	wait_proc int

	abs_tick uint64
}

func (vm *VM) CopyState(vmSource *VM) error {
	// Validate inputs
	if vm == nil || vmSource == nil {
		return errors.New("cannot copy state from or to a nil VM")
	}

	// Copy processor states
	for i, pState := range vmSource.Processors {
		if err := vm.Processors[i].CopyState(pState); err != nil {
			return fmt.Errorf("failed to copy processor %d state: %w", i, err)
		}
	}

	// Copy input/output registers
	copy(vm.Inputs_regs, vmSource.Inputs_regs)
	copy(vm.Outputs_regs, vmSource.Outputs_regs)
	copy(vm.Internal_inputs_regs, vmSource.Internal_inputs_regs)
	copy(vm.Internal_outputs_regs, vmSource.Internal_outputs_regs)

	// Copy valid flags
	copy(vm.InputsValid, vmSource.InputsValid)
	copy(vm.OutputsValid, vmSource.OutputsValid)
	copy(vm.InternalInputsValid, vmSource.InternalInputsValid)
	copy(vm.InternalOutputsValid, vmSource.InternalOutputsValid)

	// Copy recv flags
	copy(vm.InputsRecv, vmSource.InputsRecv)
	copy(vm.OutputsRecv, vmSource.OutputsRecv)
	copy(vm.InternalInputsRecv, vmSource.InternalInputsRecv)
	copy(vm.InternalOutputsRecv, vmSource.InternalOutputsRecv)

	// Copy deferred instructions map
	vm.DeferredInstructions = make(map[string]DeferredInstruction)
	maps.Copy(vm.DeferredInstructions, vmSource.DeferredInstructions)

	// Copy absolute tick counter
	vm.abs_tick = vmSource.abs_tick

	return nil
}

type SimConfig struct {
	ShowTicks      bool
	ShowIoPre      bool
	ShowIoPost     bool
	GetTicks       bool
	GetAll         bool
	GetAllInternal bool
}

// Simbox rules are converted in a sim drive when the simulation starts and applied during the simulation
type SimTickSet map[int]interface{}
type SimDrive struct {
	Injectables []*interface{}
	NeedValid   map[int]int
	AbsSet      map[uint64]SimTickSet
	PerSet      map[uint64]SimTickSet
}

type simEvent struct {
	event  uint8
	object string
}

// This is initialized when the simulation starts and filled on the way
type SimTickGet map[int]interface{}
type SimTickShow map[int]struct{}
type EventPointers [2]int // [0] -> index in Reportables/Showables, [1] -> index in the event data slice
type SimReport struct {
	Reportables      []*interface{}             // Direct pointers to the elements that is possible to report
	Showables        []*interface{}             // Direct pointers to the elements that is possible to show
	EventData        []*interface{}             // Data associated to events, for example valid, recv signals
	ReportablesTypes []string                   // Types of the reportables elements
	ShowablesTypes   []string                   // Types of the showables elements
	ReportablesNames []string                   // Names of the reportables elements
	ShowablesNames   []string                   // Names of the showables elements
	AbsGet           map[uint64]SimTickGet      // Absolute tick -> map[index in Reportables]value
	PerGet           map[uint64]SimTickGet      // Periodic tick -> map[index in Reportables]value
	EventGet         map[simEvent]EventPointers // Events that trigger a get -> map to pointers in Reportables and EventData
	AbsShow          map[uint64]SimTickShow     // Absolute tick -> map[index in Showables]struct{}
	PerShow          map[uint64]SimTickShow     // Periodic tick -> map[index in Showables]struct{}
	EventShow        map[simEvent]EventPointers // Events that trigger a show -> map to pointers in Showables and EventData
}

func (vm *VM) Processor_execute(psc *procbuilder.SimConfig, instruct <-chan int, resp chan<- int, resultChan chan<- string, procId int) {
	for {
		switch <-instruct {
		case 0:
			resp <- procId
		case 1:
			result, err := vm.Processors[procId].Step(psc)
			resp <- procId
			if err == nil {
				resultChan <- result
			} else {
				resultChan <- ""
			}
		}
	}
}

func (vm *VM) Init() error {
	vm.Processors = make([]*procbuilder.VM, len(vm.Bmach.Processors))
	vm.Inputs_regs = make([]interface{}, vm.Bmach.Inputs)
	vm.Outputs_regs = make([]interface{}, vm.Bmach.Outputs)
	vm.Internal_inputs_regs = make([]interface{}, len(vm.Bmach.Internal_inputs))
	vm.Internal_outputs_regs = make([]interface{}, len(vm.Bmach.Internal_outputs))

	vm.InputsValid = make([]bool, vm.Bmach.Inputs)
	vm.OutputsValid = make([]bool, vm.Bmach.Outputs)
	vm.InputsRecv = make([]bool, vm.Bmach.Inputs)
	vm.OutputsRecv = make([]bool, vm.Bmach.Outputs)

	vm.InternalInputsValid = make([]bool, len(vm.Bmach.Internal_inputs))
	vm.InternalOutputsValid = make([]bool, len(vm.Bmach.Internal_outputs))
	vm.InternalInputsRecv = make([]bool, len(vm.Bmach.Internal_inputs))
	vm.InternalOutputsRecv = make([]bool, len(vm.Bmach.Internal_outputs))

	vm.DeferredInstructions = make(map[string]DeferredInstruction)

	vm.abs_tick = uint64(0)

	if vm.EmuDrivers == nil {
		vm.EmuDrivers = make([]EmuDriver, 0)
	}

	cmdChan := make(chan []byte)
	vm.cmdChan = cmdChan

	for _, ed := range vm.EmuDrivers {
		ed.Init()
	}

	for i, proc_dom_id := range vm.Bmach.Processors {
		pvm := new(procbuilder.VM)
		pvm.Mach = vm.Bmach.Domains[proc_dom_id]
		pvm.Emulating = vm.Emulating
		pvm.CpID = uint32(i)
		pvm.CmdChan = cmdChan

		if vm.SimDelayMap != nil {
			pvm.SimDelayArray = make([]*simbox.DelayDistribution, len(pvm.Mach.Op))
			for j, opcode := range pvm.Mach.Op {
				opName := opcode.Op_get_name()
				if delayDistr, ok := vm.SimDelayMap.OpcodeDelays[opName]; ok {
					// delayDistr.Normalize() // This will break concurrency. Must be done before starting the simulation.
					pvm.SimDelayArray[j] = &delayDistr
				} else {
					pvm.SimDelayArray[j] = nil
				}
			}
		} else {
			pvm.SimDelayArray = nil
		}

		pvm.Init()
		vm.Processors[i] = pvm
	}

	vm.send_chans = make([]chan int, len(vm.Bmach.Processors))
	vm.result_chans = make([]chan string, len(vm.Bmach.Processors))
	vm.recv_chan = make(chan int)

	vm.wait_proc = 0

	for i := 0; i < len(vm.Processors); i++ {
		vm.wait_proc = vm.wait_proc + 1
		vm.send_chans[i] = make(chan int)
		vm.result_chans[i] = make(chan string)
	}

	if vm.Bmach.Rsize <= 8 {
		for i := 0; i < vm.Bmach.Inputs; i++ {
			vm.Inputs_regs[i] = uint8(0)
		}
		for i := 0; i < vm.Bmach.Outputs; i++ {
			vm.Outputs_regs[i] = uint8(0)
		}
		for i := 0; i < len(vm.Bmach.Internal_inputs); i++ {
			vm.Internal_inputs_regs[i] = uint8(0)
		}
		for i := 0; i < len(vm.Bmach.Internal_outputs); i++ {
			vm.Internal_outputs_regs[i] = uint8(0)
		}
	} else if vm.Bmach.Rsize <= 16 {
		for i := 0; i < vm.Bmach.Inputs; i++ {
			vm.Inputs_regs[i] = uint16(0)
		}
		for i := 0; i < vm.Bmach.Outputs; i++ {
			vm.Outputs_regs[i] = uint16(0)
		}
		for i := 0; i < len(vm.Bmach.Internal_inputs); i++ {
			vm.Internal_inputs_regs[i] = uint16(0)
		}
		for i := 0; i < len(vm.Bmach.Internal_outputs); i++ {
			vm.Internal_outputs_regs[i] = uint16(0)
		}
	} else if vm.Bmach.Rsize <= 32 {
		for i := 0; i < vm.Bmach.Inputs; i++ {
			vm.Inputs_regs[i] = uint32(0)
		}
		for i := 0; i < vm.Bmach.Outputs; i++ {
			vm.Outputs_regs[i] = uint32(0)
		}
		for i := 0; i < len(vm.Bmach.Internal_inputs); i++ {
			vm.Internal_inputs_regs[i] = uint32(0)
		}
		for i := 0; i < len(vm.Bmach.Internal_outputs); i++ {
			vm.Internal_outputs_regs[i] = uint32(0)
		}
	} else if vm.Bmach.Rsize <= 64 {
		for i := 0; i < vm.Bmach.Inputs; i++ {
			vm.Inputs_regs[i] = uint64(0)
		}
		for i := 0; i < vm.Bmach.Outputs; i++ {
			vm.Outputs_regs[i] = uint64(0)
		}
		for i := 0; i < len(vm.Bmach.Internal_inputs); i++ {
			vm.Internal_inputs_regs[i] = uint64(0)
		}
		for i := 0; i < len(vm.Bmach.Internal_outputs); i++ {
			vm.Internal_outputs_regs[i] = uint64(0)
		}
	} else {
		return errors.New("invalid register size, must smaller or equal to 64 bits")
	}
	//	// Set the initial state of the internal outputs registers
	//	for i, bond := range vm.Bmach.Internal_outputs {
	//		switch bond.Map_to {
	//		case 0:
	//			vm.Internal_outputs_regs[i] = vm.Inputs_regs[bond.Res_id]
	//		case 3:
	//			vm.Internal_outputs_regs[i] = vm.Processors[bond.Res_id].Outputs[bond.Ext_id]
	//		}
	//	}

	return nil
}

func (vm *VM) EmuDriverDispatcher() {
	// TODO Complete
	// fmt.Println("EmuDriverDispatcher", vm.EmuDrivers)
	for {
		select {
		case cmd := <-vm.cmdChan:
			for _, ed := range vm.EmuDrivers {
				ed.PushCommand(cmd)
			}
		}
	}
}

func (vm *VM) Launch_processors(s *simbox.Simbox) error {
	go vm.EmuDriverDispatcher()
	for i := 0; i < len(vm.Processors); i++ {

		psc := new(procbuilder.SimConfig)
		pscerr := psc.Init(s, vm.Processors[i])
		check(pscerr)

		for _, ed := range vm.EmuDrivers {
			go ed.Run()
		}
		go vm.Processor_execute(psc, vm.send_chans[i], vm.recv_chan, vm.result_chans[i], i)
	}
	return nil
}

func (vm *VM) Step(sc *SimConfig) (string, error) {

	result := ""
	debug := false

	if sc != nil {
		if sc.ShowTicks {
			result += "Absolute tick:" + strconv.Itoa(int(vm.abs_tick)) + "\n"
		}
	}

	if sc != nil {
		if sc.ShowIoPre {
			result += "\tPre-compute IO: " + vm.DumpIO() + "\n"
		}
	}

	if debug {
		result += "\tPre-compute data movement:\n"
	}
	// Set the internal outputs registers and the relative data valid signal, for the BM inputs
	for i, bond := range vm.Bmach.Internal_outputs {
		switch bond.Map_to {
		case BMINPUT:
			if debug {
				iin, _ := vm.Bmach.GetInternalOutputName(i)
				result += "\t\t(" + vm.dumpRegister(vm.Inputs_regs[bond.Res_id]) + "," + fmt.Sprintf("%t", vm.InputsValid[bond.Res_id]) + ") BM Input " + strconv.Itoa(bond.Res_id) + "(data,valid) -> internal output: " + iin + "(data,valid)\n"
			}
			vm.Internal_outputs_regs[i] = vm.Inputs_regs[bond.Res_id]
			vm.InternalOutputsValid[i] = vm.InputsValid[bond.Res_id]
		}
	}

	// Transfer to the internal inputs registers and the relative data valids the previous outputs according the links
	for i, j := range vm.Bmach.Links {
		if j != -1 {
			if debug {
				iin, _ := vm.Bmach.GetInternalInputName(i)
				iout, _ := vm.Bmach.GetInternalOutputName(j)
				result += "\t\t(" + vm.dumpRegister(vm.Internal_outputs_regs[j]) + "," + fmt.Sprintf("%t", vm.InternalOutputsValid[j]) + ") internal output: " + iout + "(data,valid) -> internal input: " + iin + "(data,valid)\n"
			}
			vm.Internal_inputs_regs[i] = vm.Internal_outputs_regs[j]
			vm.InternalInputsValid[i] = vm.InternalOutputsValid[j]
		}
	}

	// Transfer internal inputs registers and the relative data valids to their destination in the processors
	for i, bond := range vm.Bmach.Internal_inputs {
		switch bond.Map_to {
		case CPINPUT:
			if debug {
				iin, _ := vm.Bmach.GetInternalInputName(i)
				result += "\t\t(" + vm.dumpRegister(vm.Internal_inputs_regs[i]) + "," + fmt.Sprintf("%t", vm.InternalInputsValid[i]) + ") internal input: " + iin + "(data,valid) -> CP " + strconv.Itoa(bond.Res_id) + " Input " + strconv.Itoa(bond.Ext_id) + "(data,valid)\n"
			}
			vm.Processors[bond.Res_id].Inputs[bond.Ext_id] = vm.Internal_inputs_regs[i]
			vm.Processors[bond.Res_id].InputsValid[bond.Ext_id] = vm.InternalInputsValid[i]
		}
	}

	// Set the internal input data received signals
	for i, bond := range vm.Bmach.Internal_inputs {
		switch bond.Map_to {
		case BMOUTPUT:
			if debug {
				iin, _ := vm.Bmach.GetInternalInputName(i)
				result += "\t\t(" + fmt.Sprintf("%t", vm.OutputsRecv[bond.Res_id]) + ") BM Output " + strconv.Itoa(bond.Res_id) + " (recv) -> internal input: " + iin + " (recv)\n"
			}
			vm.InternalInputsRecv[i] = vm.OutputsRecv[bond.Res_id]
		}
	}

	// Set the internal output data received signals
	dataRecv := make(map[int]bool)
	andString := make(map[int]string)
	for i, j := range vm.Bmach.Links {
		if j != -1 {
			if val, ok := dataRecv[j]; !ok {
				dataRecv[j] = vm.InternalInputsRecv[i]
				if debug {
					iin, _ := vm.Bmach.GetInternalInputName(i)
					andString[j] = fmt.Sprintf("(%t) %s", vm.InternalInputsRecv[i], iin)
				}
			} else {
				dataRecv[j] = val && vm.InternalInputsRecv[i]
				if debug {
					iin, _ := vm.Bmach.GetInternalInputName(i)
					andString[j] = andString[j] + fmt.Sprintf(" && (%t) %s", vm.InternalInputsRecv[i], iin)
				}
			}
		}
	}
	for i, _ := range vm.Bmach.Internal_outputs {
		if val, ok := dataRecv[i]; !ok {
			vm.InternalOutputsRecv[i] = false
			if debug {
				iout, _ := vm.Bmach.GetInternalOutputName(i)
				result += "\t\t(false) -> internal output: " + iout + "(recv)\n"
			}
		} else {
			vm.InternalOutputsRecv[i] = val
			if debug {
				iout, _ := vm.Bmach.GetInternalOutputName(i)
				result += "\t\t(" + fmt.Sprintf("%t", val) + ") internal inputs: " + andString[i] + "(recv) -> internal output: " + iout + "(recv)\n"
			}
		}
	}

	// Transfer internal outputs data received to their destination in the processors
	for i, bond := range vm.Bmach.Internal_outputs {
		switch bond.Map_to {
		case CPOUTPUT:
			if debug {
				iout, _ := vm.Bmach.GetInternalOutputName(i)
				result += "\t\t(" + fmt.Sprintf("%t", vm.InternalOutputsRecv[i]) + ") internal output: " + iout + "(recv) -> CP " + strconv.Itoa(bond.Res_id) + " Output " + strconv.Itoa(bond.Ext_id) + "(recv)\n"
			}
			vm.Processors[bond.Res_id].OutputsRecv[bond.Ext_id] = vm.InternalOutputsRecv[i]
		}
	}

	if debug {
		result += "\tCompute step:\n"
	}

	// Order the step to processors
	for i := 0; i < len(vm.Processors); i++ {
		vm.send_chans[i] <- 1
		vm.wait_proc = vm.wait_proc - 1
	}

	for {
		i := <-vm.recv_chan
		proc_result := <-vm.result_chans[i]
		if proc_result != "" {
			result += "\tProc: " + strconv.Itoa(i) + "\n"
			result += proc_result
		}
		vm.wait_proc = vm.wait_proc + 1
		if vm.wait_proc == len(vm.Processors) {
			break
		}
	}

	if debug {
		result += "\tPost-compute data movement:\n"
	}

	// Set the internal outputs registers
	for i, bond := range vm.Bmach.Internal_outputs {
		switch bond.Map_to {
		case CPOUTPUT:
			if debug {
				iout, _ := vm.Bmach.GetInternalOutputName(i)
				result += "\t\t(" + vm.dumpRegister(vm.Processors[bond.Res_id].Outputs[bond.Ext_id]) + "," + fmt.Sprintf("%t", vm.Processors[bond.Res_id].OutputsValid[bond.Ext_id]) + ") CP " + strconv.Itoa(bond.Res_id) + " Output " + strconv.Itoa(bond.Ext_id) + "(data,valid) -> internal output: " + iout + "(data,valid)\n"
			}
			vm.Internal_outputs_regs[i] = vm.Processors[bond.Res_id].Outputs[bond.Ext_id]
			vm.InternalOutputsValid[i] = vm.Processors[bond.Res_id].OutputsValid[bond.Ext_id]
		}
	}

	// Transfer to the internal inputs registers the previous outputs according the links
	for i, j := range vm.Bmach.Links {
		if j != -1 {
			if debug {
				iin, _ := vm.Bmach.GetInternalInputName(i)
				iout, _ := vm.Bmach.GetInternalOutputName(j)
				result += "\t\t(" + vm.dumpRegister(vm.Internal_outputs_regs[j]) + "," + fmt.Sprintf("%t", vm.InternalOutputsValid[j]) + ") internal output: " + iout + "(data,valid) -> internal input: " + iin + "(data,valid)\n"
			}
			vm.Internal_inputs_regs[i] = vm.Internal_outputs_regs[j]
			vm.InternalInputsValid[i] = vm.InternalOutputsValid[j]
		}
	}

	// Transfer internal inputs registers to their destination
	for i, bond := range vm.Bmach.Internal_inputs {
		switch bond.Map_to {
		case BMOUTPUT:
			if debug {
				iin, _ := vm.Bmach.GetInternalInputName(i)
				result += "\t\t(" + vm.dumpRegister(vm.Internal_inputs_regs[i]) + "," + fmt.Sprintf("%t", vm.InternalInputsValid[i]) + ") internal input: " + iin + "(data,valid) -> BM Output " + strconv.Itoa(bond.Res_id) + "(data,valid)\n"
			}
			vm.Outputs_regs[bond.Res_id] = vm.Internal_inputs_regs[i]
			vm.OutputsValid[bond.Res_id] = vm.InternalInputsValid[i]
		}
	}

	// Set the internal inputs registers data received signals
	for i, bond := range vm.Bmach.Internal_inputs {
		switch bond.Map_to {
		case CPINPUT:
			if debug {
				iin, _ := vm.Bmach.GetInternalInputName(i)
				result += "\t\t(" + fmt.Sprintf("%t", vm.Processors[bond.Res_id].InputsRecv[bond.Ext_id]) + ") CP " + strconv.Itoa(bond.Res_id) + " Input " + strconv.Itoa(bond.Ext_id) + "(recv) -> internal input: " + iin + " (recv)\n"
			}
			vm.InternalInputsRecv[i] = vm.Processors[bond.Res_id].InputsRecv[bond.Ext_id]
		}
	}

	// Set the internal output data received signals
	dataRecv = make(map[int]bool)
	andString = make(map[int]string)
	for i, j := range vm.Bmach.Links {
		if j != -1 {
			if val, ok := dataRecv[j]; !ok {
				dataRecv[j] = vm.InternalInputsRecv[i]
				if debug {
					iin, _ := vm.Bmach.GetInternalInputName(i)
					andString[j] = fmt.Sprintf("(%t) %s", vm.InternalInputsRecv[i], iin)
				}
			} else {
				dataRecv[j] = val && vm.InternalInputsRecv[i]
				if debug {
					iin, _ := vm.Bmach.GetInternalInputName(i)
					andString[j] = andString[j] + fmt.Sprintf(" && (%t) %s", vm.InternalInputsRecv[i], iin)
				}
			}
		}
	}
	for i, _ := range vm.Bmach.Internal_outputs {
		if val, ok := dataRecv[i]; !ok {
			vm.InternalOutputsRecv[i] = false
			if debug {
				iout, _ := vm.Bmach.GetInternalOutputName(i)
				result += "\t\t(false) -> internal output: " + iout + "(recv)\n"
			}
		} else {
			vm.InternalOutputsRecv[i] = val
			if debug {
				iout, _ := vm.Bmach.GetInternalOutputName(i)
				result += "\t\t(" + fmt.Sprintf("%t", val) + ") internal inputs: " + andString[i] + "(recv) -> internal output: " + iout + "(recv)\n"
			}
		}
	}

	// Transfer internal outputs data received to their destination in the BM inputs
	for i, bond := range vm.Bmach.Internal_outputs {
		switch bond.Map_to {
		case BMINPUT:
			if debug {
				iout, _ := vm.Bmach.GetInternalOutputName(i)
				result += "\t\t(" + fmt.Sprintf("%t", vm.InternalOutputsRecv[i]) + ") internal output: " + iout + "(recv) -> BM Input " + strconv.Itoa(bond.Res_id) + "(recv)\n"
			}
			vm.InputsRecv[bond.Res_id] = vm.InternalOutputsRecv[i]
		}
	}

	if sc != nil {
		if sc.ShowIoPost {
			result += "\tPost-compute IO: " + vm.DumpIO() + "\n"
		}
	}

	vm.abs_tick++

	return result, nil
}

func (vm *VM) dumpRegister(reg any) string {
	if vm.Bmach.Rsize <= 8 {
		return zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint8))))
	} else if vm.Bmach.Rsize <= 16 {
		return zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint16))))
	} else if vm.Bmach.Rsize <= 32 {
		return zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint32))))
	} else if vm.Bmach.Rsize <= 64 {
		return zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint64))))
	}
	return ""
}

func (vm *VM) DumpIO() string {
	result := ""
	for i, reg := range vm.Inputs_regs {
		if vm.Bmach.Rsize <= 8 {
			result = result + Get_input_name(i) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint8)))) + " "
		} else if vm.Bmach.Rsize <= 16 {
			result = result + Get_input_name(i) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint16)))) + " "
		} else if vm.Bmach.Rsize <= 32 {
			result = result + Get_input_name(i) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint32)))) + " "
		} else if vm.Bmach.Rsize <= 64 {
			result = result + Get_input_name(i) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint64)))) + " "
		} else {
			result = result + "ERROR, Rsize not supported, only <= 64 bits"
		}
		result += "(v:" + strconv.FormatBool(vm.InputsValid[i]) + " r:" + strconv.FormatBool(vm.InputsRecv[i]) + ") "
	}
	for i, reg := range vm.Outputs_regs {
		if vm.Bmach.Rsize <= 8 {
			result = result + Get_output_name(i) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint8)))) + " "
		} else if vm.Bmach.Rsize <= 16 {
			result = result + Get_output_name(i) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint16)))) + " "
		} else if vm.Bmach.Rsize <= 32 {
			result = result + Get_output_name(i) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint32)))) + " "
		} else if vm.Bmach.Rsize <= 64 {
			result = result + Get_output_name(i) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint64)))) + " "
		} else {
			result = result + "ERROR, Rsize not supported, only <= 64 bits"
		}
		result += "(v:" + strconv.FormatBool(vm.OutputsValid[i]) + " r:" + strconv.FormatBool(vm.OutputsRecv[i]) + ") "
	}
	return result
}

func (vm *VM) GetElementLocation(mnemonic string) (*interface{}, error) {
	// TODO include others

	// This function returns a pointer to the element identified by the mnemonic. The elements within
	// The VM are of the type interface{} to allow for different register sizes. So this function returns
	// A pointer to an interface{} that can be casted to the right type by the caller.
	// However, for boolean elements (like valid and recv signals) the function returns a pointer to a any
	// that can be casted to a *bool by the caller. (It is a sort of embedding a *bool into an interface{})
	// This is done to avoid returning a pointer to a bool directly, which would not match the *interface{}
	// The caller must be aware of this behavior and handle the casting accordingly.
	// The blocks where this is true are marked with a comment.

	// Input registers
	re := regexp.MustCompile("^i(?P<input>[0-9]+)$")
	if re.MatchString(mnemonic) {
		inputNum := re.ReplaceAllString(mnemonic, "${input}")
		if i, err := strconv.Atoi(inputNum); err == nil {
			if i < len(vm.Inputs_regs) {
				return &vm.Inputs_regs[i], nil
			}
		}
	}
	re = regexp.MustCompile("^i(?P<input>[0-9]+)v$")
	if re.MatchString(mnemonic) {
		inputNum := re.ReplaceAllString(mnemonic, "${input}")
		if i, err := strconv.Atoi(inputNum); err == nil {
			if i < len(vm.InputsValid) {
				var result any = &vm.InputsValid[i] // Pointer to bool embedded in interface{}
				return &result, nil
			}
		}
	}
	re = regexp.MustCompile("^i(?P<input>[0-9]+)r$")
	if re.MatchString(mnemonic) {
		inputNum := re.ReplaceAllString(mnemonic, "${input}")
		if i, err := strconv.Atoi(inputNum); err == nil {
			if i < len(vm.InputsRecv) {
				var result any = &vm.InputsRecv[i] // Pointer to bool embedded in interface{}
				return &result, nil
			}
		}
	}
	re = regexp.MustCompile("^o(?P<output>[0-9]+)$")
	if re.MatchString(mnemonic) {
		outputNum := re.ReplaceAllString(mnemonic, "${output}")
		if i, err := strconv.Atoi(outputNum); err == nil {
			if i < len(vm.Outputs_regs) {
				return &vm.Outputs_regs[i], nil
			}
		}
	}
	re = regexp.MustCompile("^o(?P<output>[0-9]+)v$")
	if re.MatchString(mnemonic) {
		outputNum := re.ReplaceAllString(mnemonic, "${output}")
		if i, err := strconv.Atoi(outputNum); err == nil {
			if i < len(vm.OutputsValid) {
				var result any = &vm.OutputsValid[i] // Pointer to bool embedded in interface{}
				return &result, nil
			}
		}
	}
	re = regexp.MustCompile("^o(?P<output>[0-9]+)r$")
	if re.MatchString(mnemonic) {
		outputNum := re.ReplaceAllString(mnemonic, "${output}")
		if i, err := strconv.Atoi(outputNum); err == nil {
			if i < len(vm.OutputsRecv) {
				var result any = &vm.OutputsRecv[i] // Pointer to bool embedded in interface{}
				return &result, nil
			}
		}
	}
	re = regexp.MustCompile("^p(?P<proc>[0-9]+)i(?P<input>[0-9]+)$")
	if re.MatchString(mnemonic) {
		procNum := re.ReplaceAllString(mnemonic, "${proc}")
		inputNum := re.ReplaceAllString(mnemonic, "${input}")
		if i, err := strconv.Atoi(procNum); err == nil {
			if i < len(vm.Processors) {
				if j, err := strconv.Atoi(inputNum); err == nil {
					if j < len(vm.Processors[i].Inputs) {
						return &vm.Processors[i].Inputs[j], nil
					}
				}
			}
		}
	}
	re = regexp.MustCompile("^p(?P<proc>[0-9]+)o(?P<output>[0-9]+)$")
	if re.MatchString(mnemonic) {
		procNum := re.ReplaceAllString(mnemonic, "${proc}")
		outputNum := re.ReplaceAllString(mnemonic, "${output}")
		if i, err := strconv.Atoi(procNum); err == nil {
			if i < len(vm.Processors) {
				if j, err := strconv.Atoi(outputNum); err == nil {
					if j < len(vm.Processors[i].Outputs) {
						return &vm.Processors[i].Outputs[j], nil
					}
				}
			}
		}
	}

	re = regexp.MustCompile("^p(?P<proc>[0-9]+)r(?P<reg>[0-9]+)$")
	if re.MatchString(mnemonic) {
		procNum := re.ReplaceAllString(mnemonic, "${proc}")
		regNum := re.ReplaceAllString(mnemonic, "${reg}")
		if i, err := strconv.Atoi(procNum); err == nil {
			if i < len(vm.Processors) {
				if j, err := strconv.Atoi(regNum); err == nil {
					if j < len(vm.Processors[i].Registers) {
						return &vm.Processors[i].Registers[j], nil
					}
				}
			}
		}
	}

	return nil, errors.New("unknown mnemonic " + mnemonic)
}

func (sc *SimConfig) Init(s *simbox.Simbox, vm *VM, conf *Config) error {

	if s != nil {

		for _, rule := range s.Rules {
			// Skip suspended rules
			if rule.Suspended {
				continue
			}
			if conf.Debug {
				fmt.Println("Loading simbox rule:", rule)
			}
			// Intercept the set rules
			if rule.Timec == simbox.TIMEC_NONE && rule.Action == simbox.ACTION_CONFIG {
				switch rule.Object {
				case "show_ticks":
					sc.ShowTicks = true
				case "show_io_pre":
					sc.ShowIoPre = true
				case "show_io_post":
					sc.ShowIoPost = true
				case "get_ticks":
					sc.GetTicks = true
				case "get_all":
					sc.GetAll = true
				case "get_all_internal":
					sc.GetAllInternal = true
				}
			}
		}
	}
	return nil
}

func (sd *SimDrive) Init(c *Config, s *simbox.Simbox, vm *VM) error {

	inj := make([]*interface{}, 0)
	needValid := make(map[int]int)
	absset := make(map[uint64]SimTickSet)
	perset := make(map[uint64]SimTickSet)

	for _, rule := range s.Rules {
		// Skip suspended rules
		if rule.Suspended {
			continue
		}
		// Intercept the set rules
		if rule.Timec == simbox.TIMEC_ABS && rule.Action == simbox.ACTION_SET {
			if loc, err := vm.GetElementLocation(rule.Object); err == nil {
				if val, err := ImportNumber(c, rule.Extra); err == nil {
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

						re := regexp.MustCompile("^i(?P<input>[0-9]+)$")
						if re.MatchString(rule.Object) {
							inIdxS := re.ReplaceAllString(rule.Object, "${input}")
							inIdx, err := strconv.Atoi(inIdxS)
							if err != nil {
								return err
							}
							needValid[ipos] = inIdx
						}
					}

					if actOnTick, ok := absset[rule.Tick]; ok {
						if vm.Bmach.Rsize <= 8 {
							actOnTick[ipos] = uint8(val)
						} else if vm.Bmach.Rsize <= 16 {
							actOnTick[ipos] = uint16(val)
						} else if vm.Bmach.Rsize <= 32 {
							actOnTick[ipos] = uint32(val)
						} else if vm.Bmach.Rsize <= 64 {
							actOnTick[ipos] = uint64(val)
						} else {
							return errors.New("unsupported register size, <= 64 are supported")
						}
					} else {
						actOnTick := make(map[int]interface{})
						if vm.Bmach.Rsize <= 8 {
							actOnTick[ipos] = uint8(val)
						} else if vm.Bmach.Rsize <= 16 {
							actOnTick[ipos] = uint16(val)
						} else if vm.Bmach.Rsize <= 32 {
							actOnTick[ipos] = uint32(val)
						} else if vm.Bmach.Rsize <= 64 {
							actOnTick[ipos] = uint64(val)
						} else {
							return errors.New("unsupported register size, <= 64 are supported")
						}
						absset[rule.Tick] = actOnTick
					}
				} else {
					return err
				}
			} else {
				return err
			}
		}
		// Intercept the periodic set rules
		if rule.Timec == simbox.TIMEC_REL && rule.Action == simbox.ACTION_SET {
			if loc, err := vm.GetElementLocation(rule.Object); err == nil {
				if val, err := ImportNumber(c, rule.Extra); err == nil {
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

					if actOnTick, ok := perset[rule.Tick]; ok {
						if vm.Bmach.Rsize <= 8 {
							actOnTick[ipos] = uint8(val)
						} else if vm.Bmach.Rsize <= 16 {
							actOnTick[ipos] = uint16(val)
						} else if vm.Bmach.Rsize <= 32 {
							actOnTick[ipos] = uint32(val)
						} else if vm.Bmach.Rsize <= 64 {
							actOnTick[ipos] = uint64(val)
						} else {
							return errors.New("unsupported register size, <= 64 are supported")
						}
					} else {
						actOnTick := make(map[int]interface{})
						if vm.Bmach.Rsize <= 8 {
							actOnTick[ipos] = uint8(val)
						} else if vm.Bmach.Rsize <= 16 {
							actOnTick[ipos] = uint16(val)
						} else if vm.Bmach.Rsize <= 32 {
							actOnTick[ipos] = uint32(val)
						} else if vm.Bmach.Rsize <= 64 {
							actOnTick[ipos] = uint64(val)
						} else {
							return errors.New("unsupported register size, <= 64 are supported")
						}
						perset[rule.Tick] = actOnTick
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
	sd.NeedValid = needValid
	sd.AbsSet = absset
	sd.PerSet = perset
	return nil
}

func (sd *SimReport) Init(s *simbox.Simbox, vm *VM) error {

	rep := make([]*interface{}, 0)
	sho := make([]*interface{}, 0)
	ev := make([]*interface{}, 0)
	repTypes := make([]string, 0)
	shoTypes := make([]string, 0)
	repNames := make([]string, 0)
	shoNames := make([]string, 0)
	absget := make(map[uint64]SimTickGet)
	perget := make(map[uint64]SimTickGet)
	absshow := make(map[uint64]SimTickShow)
	pershow := make(map[uint64]SimTickShow)
	eventget := make(map[simEvent]EventPointers)
	eventshow := make(map[simEvent]EventPointers)

	for _, rule := range s.Rules {
		// Skip suspended rules
		if rule.Suspended {
			continue
		}
		// Intercept the get rules from config action
		if rule.Timec == simbox.TIMEC_NONE && rule.Action == simbox.ACTION_CONFIG {
			objects := make([]string, 0)
			switch rule.Object {
			case "get_all":
				for i, _ := range vm.Bmach.Internal_inputs {
					objects = append(objects, vm.Bmach.Internal_inputs[i].String())
				}
				for i, _ := range vm.Bmach.Internal_outputs {
					objects = append(objects, vm.Bmach.Internal_outputs[i].String())
				}
			case "get_all_internal":
				for i, _ := range vm.Bmach.Internal_inputs {
					objects = append(objects, vm.Bmach.Internal_inputs[i].String())
				}
				for i, _ := range vm.Bmach.Internal_outputs {
					objects = append(objects, vm.Bmach.Internal_outputs[i].String())
				}
				for i, procVm := range vm.Processors {
					for j, _ := range procVm.Registers {
						objects = append(objects, fmt.Sprintf("p%dr%d", i, j))
					}
				}
			}
			for _, obj := range objects {
				if loc, err := vm.GetElementLocation(obj); err == nil {
					ipos := -1
					for i, iloc := range rep {
						if iloc == loc {
							ipos = i
							break
						}
					}
					if ipos == -1 {
						rep = append(rep, loc)
						repNames = append(repNames, obj)
						if rule.Extra == "" {
							repTypes = append(repTypes, "unsigned")
						} else {
							repTypes = append(repTypes, rule.Extra)
						}
					}
				}
			}
		}

		// Intercept the show rules from config action
		if rule.Timec == simbox.TIMEC_NONE && rule.Action == simbox.ACTION_CONFIG {
			objects := make([]string, 0)
			switch rule.Object {
			case "show_all":
				for i, _ := range vm.Bmach.Internal_inputs {
					objects = append(objects, vm.Bmach.Internal_inputs[i].String())
				}
				for i, _ := range vm.Bmach.Internal_outputs {
					objects = append(objects, vm.Bmach.Internal_outputs[i].String())
				}
			case "show_all_internal":
				for i, _ := range vm.Bmach.Internal_inputs {
					objects = append(objects, vm.Bmach.Internal_inputs[i].String())
				}
				for i, _ := range vm.Bmach.Internal_outputs {
					objects = append(objects, vm.Bmach.Internal_outputs[i].String())
				}
				for i, procVm := range vm.Processors {
					for j, _ := range procVm.Registers {
						objects = append(objects, fmt.Sprintf("p%dr%d", i, j))
					}
				}
			}
			for _, obj := range objects {
				if loc, err := vm.GetElementLocation(obj); err == nil {
					ipos := -1
					for i, iloc := range sho {
						if iloc == loc {
							ipos = i
							break
						}
					}
					if ipos == -1 {
						sho = append(sho, loc)
						shoNames = append(shoNames, obj)
						if rule.Extra == "" {
							shoTypes = append(shoTypes, "unsigned")
						} else {
							shoTypes = append(shoTypes, rule.Extra)
						}
					}
				}
			}
		}

		// Intercept the get rules on valid time
		if rule.Timec == simbox.TIMEC_ON_VALID && rule.Action == simbox.ACTION_GET {
			// Getting the location of the object to report
			ipos := -1
			if loc, err := vm.GetElementLocation(rule.Object); err == nil {
				for i, iloc := range rep {
					if iloc == loc {
						ipos = i
						break
					}
				}
				if ipos == -1 {
					ipos = len(rep)
					rep = append(rep, loc)
					repNames = append(repNames, rule.Object)

					if rule.Extra == "" {
						repTypes = append(repTypes, "unsigned")
					} else {
						repTypes = append(repTypes, rule.Extra)
					}
				}
			}
			// Getting the location of the valid signal of the object to report
			iposv := -1
			if loc, err := vm.GetElementLocation(rule.Object + "v"); err == nil {
				for i, iloc := range ev {
					if iloc == loc {
						iposv = i
						break
					}
				}
				if iposv == -1 {
					iposv = len(ev)
					ev = append(ev, loc)
				}
			}
			if (ipos != -1) && (iposv != -1) {
				eventget[simEvent{event: EVENTONVALID, object: rule.Object}] = [2]int{ipos, iposv}
			}
		}

		// Intercept the show rules on valid time
		if rule.Timec == simbox.TIMEC_ON_VALID && rule.Action == simbox.ACTION_SHOW {
			ipos := -1
			if loc, err := vm.GetElementLocation(rule.Object); err == nil {
				for i, iloc := range sho {
					if iloc == loc {
						ipos = i
						break
					}
				}
				if ipos == -1 {
					ipos = len(sho)
					sho = append(sho, loc)
					shoNames = append(shoNames, rule.Object)
					if rule.Extra == "" {
						shoTypes = append(shoTypes, "unsigned")
					} else {
						shoTypes = append(shoTypes, rule.Extra)
					}
				}
			}
			iposv := -1
			if locv, err := vm.GetElementLocation(rule.Object + "v"); err == nil {
				for i, iloc := range ev {
					if iloc == locv {
						iposv = i
						break
					}
				}
				if iposv == -1 {
					iposv = len(ev)
					ev = append(ev, locv)
				}
			}
			if (ipos != -1) && (iposv != -1) {
				eventshow[simEvent{event: EVENTONVALID, object: rule.Object}] = [2]int{ipos, iposv}
			}
		}
		// Intercept the get rules on exit time
		if rule.Timec == simbox.TIMEC_ON_EXIT && rule.Action == simbox.ACTION_GET {
			ipos := -1
			if loc, err := vm.GetElementLocation(rule.Object); err == nil {
				for i, iloc := range rep {
					if iloc == loc {
						ipos = i
						break
					}
				}
				if ipos == -1 {
					ipos = len(rep)
					rep = append(rep, loc)
					repNames = append(repNames, rule.Object)
					if rule.Extra == "" {
						repTypes = append(repTypes, "unsigned")
					} else {
						repTypes = append(repTypes, rule.Extra)
					}
				}
			}

			// The exit event does not need other signals
			if ipos != -1 {
				eventget[simEvent{event: EVENTONEXIT, object: rule.Object}] = [2]int{ipos, -1}
			}
		}

		// Intercept the show rules on exit time
		if rule.Timec == simbox.TIMEC_ON_EXIT && rule.Action == simbox.ACTION_SHOW {
			ipos := -1
			if loc, err := vm.GetElementLocation(rule.Object); err == nil {
				for i, iloc := range sho {
					if iloc == loc {
						ipos = i
						break
					}
				}
				if ipos == -1 {
					ipos = len(sho)
					sho = append(sho, loc)
					shoNames = append(shoNames, rule.Object)
					if rule.Extra == "" {
						shoTypes = append(shoTypes, "unsigned")
					} else {
						shoTypes = append(shoTypes, rule.Extra)
					}
				}
			}

			// The exit event does not need other signals
			if ipos != -1 {
				eventshow[simEvent{event: EVENTONEXIT, object: rule.Object}] = [2]int{ipos, -1}
			}
		}

		// Intercept the get rules in absolute time
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
					repNames = append(repNames, rule.Object)
					if rule.Extra == "" {
						repTypes = append(repTypes, "unsigned")
					} else {
						repTypes = append(repTypes, rule.Extra)
					}
				}

				if strOnTick, ok := absget[rule.Tick]; ok {
					if vm.Bmach.Rsize <= 8 {
						strOnTick[ipos] = uint8(0)
					} else if vm.Bmach.Rsize <= 16 {
						strOnTick[ipos] = uint16(0)
					} else if vm.Bmach.Rsize <= 32 {
						strOnTick[ipos] = uint32(0)
					} else if vm.Bmach.Rsize <= 64 {
						strOnTick[ipos] = uint64(0)
					} else {
						return errors.New("unsupported register size, <= 64 are supported")
					}
				} else {
					strOnTick := make(map[int]interface{})
					if vm.Bmach.Rsize <= 8 {
						strOnTick[ipos] = uint8(0)
					} else if vm.Bmach.Rsize <= 16 {
						strOnTick[ipos] = uint16(0)
					} else if vm.Bmach.Rsize <= 32 {
						strOnTick[ipos] = uint32(0)
					} else if vm.Bmach.Rsize <= 64 {
						strOnTick[ipos] = uint64(0)
					} else {
						return errors.New("unsupported register size, <= 64 are supported")
					}
					absget[rule.Tick] = strOnTick
				}
			} else {
				return err
			}
		}
		// Intercept the get rules in relative time
		if rule.Timec == simbox.TIMEC_REL && rule.Action == simbox.ACTION_GET {
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
					repNames = append(repNames, rule.Object)
					if rule.Extra == "" {
						repTypes = append(repTypes, "unsigned")
					} else {
						repTypes = append(repTypes, rule.Extra)
					}
				}

				if strOnTick, ok := perget[rule.Tick]; ok {
					if vm.Bmach.Rsize <= 8 {
						strOnTick[ipos] = uint8(0)
					} else if vm.Bmach.Rsize <= 16 {
						strOnTick[ipos] = uint16(0)
					} else if vm.Bmach.Rsize <= 32 {
						strOnTick[ipos] = uint32(0)
					} else if vm.Bmach.Rsize <= 64 {
						strOnTick[ipos] = uint64(0)
					} else {
						return errors.New("unsupported register size, <= 64 are supported")
					}
				} else {
					strOnTick := make(map[int]interface{})
					if vm.Bmach.Rsize <= 8 {
						strOnTick[ipos] = uint8(0)
					} else if vm.Bmach.Rsize <= 16 {
						strOnTick[ipos] = uint16(0)
					} else if vm.Bmach.Rsize <= 32 {
						strOnTick[ipos] = uint32(0)
					} else if vm.Bmach.Rsize <= 64 {
						strOnTick[ipos] = uint64(0)
					} else {
						return errors.New("unsupported register size, <= 64 are supported")
					}
					perget[rule.Tick] = strOnTick
				}
			} else {
				return err
			}
		}
		// Intercept the show rules in absolute time
		if rule.Timec == simbox.TIMEC_ABS && rule.Action == simbox.ACTION_SHOW {
			if loc, err := vm.GetElementLocation(rule.Object); err == nil {
				ipos := -1
				for i, iloc := range sho {
					if iloc == loc {
						ipos = i
						break
					}
				}
				if ipos == -1 {
					ipos = len(sho)
					sho = append(sho, loc)
					shoNames = append(shoNames, rule.Object)
					if rule.Extra == "" {
						shoTypes = append(shoTypes, "unsigned")
					} else {
						shoTypes = append(shoTypes, rule.Extra)
					}
				}

				if strOnTick, ok := absshow[rule.Tick]; ok {
					strOnTick[ipos] = struct{}{}
				} else {
					str_on_tick := make(map[int]struct{})
					str_on_tick[ipos] = struct{}{}
					absshow[rule.Tick] = str_on_tick
				}
			} else {
				return err
			}
		}
		// Intercept the show rules in relative time
		if rule.Timec == simbox.TIMEC_REL && rule.Action == simbox.ACTION_SHOW {
			if loc, err := vm.GetElementLocation(rule.Object); err == nil {
				ipos := -1
				for i, iloc := range sho {
					if iloc == loc {
						ipos = i
						break
					}
				}
				if ipos == -1 {
					ipos = len(sho)
					sho = append(sho, loc)
					shoNames = append(shoNames, rule.Object)
					if rule.Extra == "" {
						shoTypes = append(shoTypes, "unsigned")
					} else {
						shoTypes = append(shoTypes, rule.Extra)
					}
				}

				if strOnTick, ok := pershow[rule.Tick]; ok {
					strOnTick[ipos] = struct{}{}
				} else {
					strOnTick := make(map[int]struct{})
					strOnTick[ipos] = struct{}{}
					pershow[rule.Tick] = strOnTick
				}
			} else {
				return err
			}
		}
	}

	sd.Reportables = rep
	sd.Showables = sho
	sd.EventData = ev
	sd.ReportablesTypes = repTypes
	sd.ShowablesTypes = shoTypes
	sd.ReportablesNames = repNames
	sd.ShowablesNames = shoNames
	sd.AbsGet = absget
	sd.PerGet = perget
	sd.AbsShow = absshow
	sd.PerShow = pershow
	sd.EventGet = eventget
	sd.EventShow = eventshow

	return nil
}
