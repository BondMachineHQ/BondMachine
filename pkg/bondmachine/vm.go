package bondmachine

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
	"github.com/BondMachineHQ/BondMachine/pkg/simbox"
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

	EmuDrivers []EmuDriver
	cmdChan    chan []byte

	send_chans   []chan int
	result_chans []chan string
	recv_chan    chan int

	wait_proc int

	abs_tick uint64
}

func (vm *VM) CopyState(vmsource *VM) {
	for i, pstate := range vmsource.Processors {
		vm.Processors[i].CopyState(pstate)
	}
	// TODO Finish
}

type Sim_config struct {
	Show_ticks     bool
	Show_io_pre    bool
	Show_io_post   bool
	GetTicks       bool
	GetAll         bool
	GetAllInternal bool
}

// Simbox rules are converted in a sim drive when the simulation starts and applied during the simulation
type Sim_tick_set map[int]interface{}
type Sim_drive struct {
	Injectables []*interface{}
	NeedValid   map[int]int
	AbsSet      map[uint64]Sim_tick_set
	PerSet      map[uint64]Sim_tick_set
}

// This is initializated when the simulation starts and filled on the way
type Sim_tick_get map[int]interface{}
type Sim_tick_show map[int]bool
type Sim_report struct {
	Reportables      []*interface{}
	Showables        []*interface{}
	ReportablesTypes []string
	ShowablesTypes   []string
	ReportablesNames []string
	AbsGet           map[uint64]Sim_tick_get
	PerGet           map[uint64]Sim_tick_get
	AbsShow          map[uint64]Sim_tick_show
	PerShow          map[uint64]Sim_tick_show
}

func (vm *VM) Processor_execute(psc *procbuilder.Sim_config, instruct <-chan int, resp chan<- int, result_chan chan<- string, proc_id int) {
	for {
		switch <-instruct {
		case 0:
			resp <- proc_id
			break
		case 1:
			result, err := vm.Processors[proc_id].Step(psc)
			resp <- proc_id
			if err == nil {
				result_chan <- result
			} else {
				result_chan <- ""
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
		pvm.CpID = uint32(i)
		pvm.CmdChan = cmdChan
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

		psc := new(procbuilder.Sim_config)
		pscerr := psc.Init(s, vm.Processors[i])
		check(pscerr)

		for _, ed := range vm.EmuDrivers {
			go ed.Run()
		}
		go vm.Processor_execute(psc, vm.send_chans[i], vm.recv_chan, vm.result_chans[i], i)
	}
	return nil
}

func (vm *VM) Step(sc *Sim_config) (string, error) {

	result := ""

	if sc != nil {
		if sc.Show_ticks {
			result += "Absolute tick:" + strconv.Itoa(int(vm.abs_tick)) + "\n"
		}
	}

	// Set the internal outputs registers and the relative data valid signal, for the BM inputs
	for i, bond := range vm.Bmach.Internal_outputs {
		switch bond.Map_to {
		case BMINPUT:
			vm.Internal_outputs_regs[i] = vm.Inputs_regs[bond.Res_id]
			vm.InternalOutputsValid[i] = vm.InputsValid[bond.Res_id]
		}
	}

	// Transfer to the internal inputs registers and the relative data valids the previous outputs according the links
	for i, j := range vm.Bmach.Links {
		if j != -1 {
			vm.Internal_inputs_regs[i] = vm.Internal_outputs_regs[j]
			vm.InternalInputsValid[i] = vm.InternalOutputsValid[j]
		}
	}

	// Transfer internal inputs registers and the relative data valids to their destination in the processors
	for i, bond := range vm.Bmach.Internal_inputs {
		switch bond.Map_to {
		case CPINPUT:
			vm.Processors[bond.Res_id].Inputs[bond.Ext_id] = vm.Internal_inputs_regs[i]
			vm.Processors[bond.Res_id].InputsValid[bond.Ext_id] = vm.InternalInputsValid[i]
		}
	}

	// Set the internal input data received signals
	for i, bond := range vm.Bmach.Internal_inputs {
		switch bond.Map_to {
		case BMOUTPUT:
			vm.InternalInputsRecv[i] = vm.OutputsRecv[bond.Res_id]
		}
	}

	// Set the internal output data received signals
	dataRecv := make(map[int]bool)
	for i, j := range vm.Bmach.Links {
		if j != -1 {
			if val, ok := dataRecv[j]; !ok {
				dataRecv[j] = vm.InternalInputsRecv[i]
			} else {
				dataRecv[j] = val && vm.InternalInputsRecv[i]
			}
		}
	}
	for i, _ := range vm.Bmach.Internal_outputs {
		if val, ok := dataRecv[i]; !ok {
			vm.InternalOutputsRecv[i] = false
		} else {
			vm.InternalOutputsRecv[i] = val
		}
	}

	// Transfer internal outputd data received to their destination in the processors
	for i, bond := range vm.Bmach.Internal_outputs {
		switch bond.Map_to {
		case CPOUTPUT:
			vm.Processors[bond.Res_id].OutputsRecv[bond.Ext_id] = vm.InternalOutputsRecv[i]
		}
	}

	if sc != nil {
		if sc.Show_io_pre {
			result += "\tPre-compute IO: " + vm.DumpIO() + "\n"
		}
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

	if sc != nil {
		if sc.Show_io_post {
			result += "\tPost-compute IO: " + vm.DumpIO() + "\n"
		}
	}

	// Set the internal outputs registers
	for i, bond := range vm.Bmach.Internal_outputs {
		switch bond.Map_to {
		case CPOUTPUT:
			vm.Internal_outputs_regs[i] = vm.Processors[bond.Res_id].Outputs[bond.Ext_id]
			vm.InternalOutputsValid[i] = vm.Processors[bond.Res_id].OutputsValid[bond.Ext_id]
		}
	}

	// Set the internal inputs registers data received signals
	for i, bond := range vm.Bmach.Internal_inputs {
		switch bond.Map_to {
		case CPINPUT:
			vm.InternalInputsRecv[i] = vm.Processors[bond.Res_id].InputsRecv[bond.Ext_id]
		}
	}

	// Transfer to the internal inputs registers the previous outputs according the links
	for i, j := range vm.Bmach.Links {
		if j != -1 {
			vm.Internal_inputs_regs[i] = vm.Internal_outputs_regs[j]
			vm.InternalInputsValid[i] = vm.InternalOutputsValid[j]
		}
	}

	// Transfer internal inputs registers to their destination
	for i, bond := range vm.Bmach.Internal_inputs {
		switch bond.Map_to {
		case BMOUTPUT:
			vm.Outputs_regs[bond.Res_id] = vm.Internal_inputs_regs[i]
			vm.OutputsValid[bond.Res_id] = vm.InternalInputsValid[i]
		}
	}

	vm.abs_tick++

	return result, nil
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

	re := regexp.MustCompile("^i(?P<input>[0-9]+)$")
	if re.MatchString(mnemonic) {
		inputNum := re.ReplaceAllString(mnemonic, "${input}")
		if i, err := strconv.Atoi(inputNum); err == nil {
			if i < len(vm.Inputs_regs) {
				return &vm.Inputs_regs[i], nil
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

func (sc *Sim_config) Init(s *simbox.Simbox, vm *VM, conf *Config) error {

	if s != nil {

		for _, rule := range s.Rules {
			if conf.Debug {
				fmt.Println("Loading simbox rule:", rule)
			}
			// Intercept the set rules
			if rule.Timec == simbox.TIMEC_NONE && rule.Action == simbox.ACTION_CONFIG {
				switch rule.Object {
				case "show_ticks":
					sc.Show_ticks = true
				case "show_io_pre":
					sc.Show_io_pre = true
				case "show_io_post":
					sc.Show_io_post = true
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

func (sd *Sim_drive) Init(c *Config, s *simbox.Simbox, vm *VM) error {

	inj := make([]*interface{}, 0)
	needValid := make(map[int]int)
	absset := make(map[uint64]Sim_tick_set)
	perset := make(map[uint64]Sim_tick_set)

	for _, rule := range s.Rules {
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

func (sd *Sim_report) Init(s *simbox.Simbox, vm *VM) error {

	rep := make([]*interface{}, 0)
	sho := make([]*interface{}, 0)
	repTypes := make([]string, 0)
	shoTypes := make([]string, 0)
	repNames := make([]string, 0)
	absget := make(map[uint64]Sim_tick_get)
	perget := make(map[uint64]Sim_tick_get)
	absshow := make(map[uint64]Sim_tick_show)
	pershow := make(map[uint64]Sim_tick_show)

	for _, rule := range s.Rules {
		// Intercept the get_all rules
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
					if rule.Extra == "" {
						shoTypes = append(shoTypes, "unsigned")
					} else {
						shoTypes = append(shoTypes, rule.Extra)
					}
				}

				if strOnTick, ok := absshow[rule.Tick]; ok {
					strOnTick[ipos] = true
				} else {
					str_on_tick := make(map[int]bool)
					str_on_tick[ipos] = true
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
					if rule.Extra == "" {
						shoTypes = append(shoTypes, "unsigned")
					} else {
						shoTypes = append(shoTypes, rule.Extra)
					}
				}

				if strOnTick, ok := pershow[rule.Tick]; ok {
					strOnTick[ipos] = true
				} else {
					strOnTick := make(map[int]bool)
					strOnTick[ipos] = true
					pershow[rule.Tick] = strOnTick
				}
			} else {
				return err
			}
		}
	}

	sd.Reportables = rep
	sd.Showables = sho
	sd.ReportablesTypes = repTypes
	sd.ShowablesTypes = shoTypes
	sd.ReportablesNames = repNames
	sd.AbsGet = absget
	sd.PerGet = perget
	sd.AbsShow = absshow
	sd.PerShow = pershow

	return nil
}
