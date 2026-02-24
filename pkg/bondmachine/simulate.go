package bondmachine

import (
	"fmt"
	"log"
	"sort"

	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/simbox"
)

func (bm *Bondmachine) SinglePipelineSimulate(dataType string, input []string, sDelay *simbox.SimDelays) ([]string, error) {

	if bm.Outputs == 0 {
		return nil, fmt.Errorf("Bondmachine must have at least one output")
	}

	result := make([]string, 0)

	// Build the simulation box
	sbox := new(simbox.Simbox)

	// Populate the simulation box with the input values
	for inId, inVal := range input {
		// fmt.Println("absolute:0:set:i" + fmt.Sprintf("%d", inId) + ":" + fmt.Sprintf("%v", inVal))
		if err := sbox.Add("absolute:0:set:i" + fmt.Sprintf("%d", inId) + ":" + fmt.Sprintf("%v", inVal)); err != nil {
			return nil, err
		}
	}

	// Populate the simulation box with the output placeholders
	for outId := 0; outId < bm.Outputs; outId++ {
		if outId == bm.Outputs-1 {
			// fmt.Println("onexit:show:o" + fmt.Sprintf("%d", outId) + ":unsigned")
			if err := sbox.Add("onexit:show:o" + fmt.Sprintf("%d", outId) + ":unsigned"); err != nil {
				return nil, err
			}
		} else {
			// fmt.Println("onexit:show:o" + fmt.Sprintf("%d", outId) + ":" + dataType)
			if err := sbox.Add("onexit:show:o" + fmt.Sprintf("%d", outId) + ":" + dataType); err != nil {
				return nil, err
			}
		}
	}

	c := new(Config)

	// Build the simulation VM
	vm := new(VM)
	vm.Bmach = bm
	vm.SimDelayMap = sDelay

	err := vm.Init()
	check(err)

	oldVm := new(VM)
	oldVm.Bmach = bm
	oldVm.SimDelayMap = vm.SimDelayMap
	if err := oldVm.Init(); err != nil {
		return nil, err
	}

	// Build the simulation configuration
	sconfig := new(SimConfig)
	if err := sconfig.Init(sbox, vm, c); err != nil {
		return nil, err
	}

	// Build the simulation driver
	sdrive := new(SimDrive)
	if err := sdrive.Init(c, sbox, vm); err != nil {
		return nil, err
	}

	// Build the simulation report
	srep := new(SimReport)
	if err := srep.Init(sbox, vm); err != nil {
		return nil, err
	}

	// Build the simulation report for the old VM
	srepOld := new(SimReport)
	if err := srepOld.Init(sbox, oldVm); err != nil {
		return nil, err
	}

	// Launch the processors
	if err := vm.Launch_processors(sbox); err != nil {
		return nil, err
	}

	// Main simulation loop, tick by tick
	for i := uint64(0); i < uint64(1000000000); i++ {

		shutDownSim := false

		if vm.OutputsValid[bm.Outputs-1] {
			shutDownSim = true
			// The simulation will stop, but we want to show the last values
			// so we execute the show/report part after this block
		}

		if !shutDownSim {
			// Manage the valid/recv states of the inputs
			for inIdx, inRecv := range vm.InputsRecv {
				if inRecv {
					vm.InputsValid[inIdx] = false
				}
			}

			// This will get actions eventually to do on this tick
			if act, exist_actions := sdrive.AbsSet[i]; exist_actions {
				for k, val := range act {
					*sdrive.Injectables[k] = val
					if inIdx, ok := sdrive.NeedValid[k]; ok {
						vm.InputsValid[inIdx] = true
					}
				}
			}

			if _, err := vm.Step(sconfig); err != nil {
				return nil, err
			}

			// Manage the valid/recv states of the outputs
			for outIdx, outValid := range vm.OutputsValid {
				if outValid {
					vm.OutputsRecv[outIdx] = true
				} else {
					vm.OutputsRecv[outIdx] = false
				}
			}
		}

		showList := make([]int, 0, len(srep.Showables))

		// This will get value to show on this tick
		if slist, exist_shows := srep.AbsShow[i]; exist_shows {
			for k, _ := range slist {
				showList = append(showList, k)
			}
		}

		// This will get value to show on periodic ticks
		for j, slist := range srep.PerShow {
			if i%j == 0 {
				for k, _ := range slist {
					alredtIn := false
					for _, v := range showList {
						if v == k {
							alredtIn = true
							break
						}
					}
					if !alredtIn {
						showList = append(showList, k)
					}
				}
			}
		}

		// This will get value to show on events
		slist, err := EventListShow(shutDownSim, srep, srepOld, vm, oldVm)
		if err != nil {
			log.Fatal(err)
		}
		for k, _ := range slist {
			alredtIn := false
			for _, v := range showList {
				if v == k {
					alredtIn = true
					break
				}
			}
			if !alredtIn {
				showList = append(showList, k)
			}
		}

		sort.Ints(showList)

		// Show the tick values
		for _, k := range showList {
			nType := srep.ShowablesTypes[k]
			if _, err := bmnumbers.EventuallyCreateType(nType, nil); err != nil {
				return nil, err
			}
			if v := bmnumbers.GetType(nType); v == nil {
				return nil, fmt.Errorf("Type " + nType + " not found")
			} else {
				bits := v.GetSize()
				if number, err := bmnumbers.ImportUint(*srep.Showables[k], bits); err != nil {
					return nil, err
				} else {
					if err := bmnumbers.CastType(number, v); err != nil {
						return nil, err
					} else {
						if numberS, err := number.ExportString(nil); err != nil {
							return nil, err
						} else {
							result = append(result, numberS)
						}
					}
				}
			}
		}

		if shutDownSim {
			break
		}

		err = oldVm.CopyState(vm)
		check(err)
	}

	return result, nil
}
