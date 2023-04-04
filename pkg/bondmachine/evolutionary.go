package bondmachine

import (
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/simbox"
	"github.com/mmirko/mel/pkg/mel"
)

func (mach *Bondmachine) MelInit(c *mel.MelConfig, ep *mel.EvolutionParameters) {
}

func (mach *Bondmachine) MelCopy() mel.Me3li {
	newmach := new(Bondmachine)
	return newmach
}

func (bmach *Bondmachine) Fitness_default(in *simbox.Simbox, exp *simbox.Simbox, sim_interactions uint64) (float32, error) {

	vm := new(VM)
	vm.Bmach = bmach
	if err := vm.Init(); err != nil {
		return 0, err
	}

	// Build the simulation configuration
	sconfig := new(Sim_config)
	if err := sconfig.Init(in, vm, nil); err != nil {
		return 0, err
	}

	// Build the simulation driver
	sdrive := new(Sim_drive)
	if err := sdrive.Init(nil, in, vm); err != nil {
		return 0, err
	}

	out := new(simbox.Simbox)
	out.Rules = make([]simbox.Rule, 0)

	for _, rule := range exp.Rules {
		// Intercept the set rules
		if rule.Action == simbox.ACTION_SET {
			out.Rules = append(out.Rules, simbox.Rule{rule.Timec, rule.Tick, simbox.ACTION_GET, rule.Object, ""})
		}
	}

	// Build the simulation report
	srep := new(Sim_report)
	if err := srep.Init(out, vm); err != nil {
		return 0, err
	}

	if err := vm.Launch_processors(in); err != nil {
		return 0, err
	}

	for i := uint64(0); i < sim_interactions; i++ {

		// This will get actions eventually to do on this tick
		if act, exist_actions := sdrive.AbsSet[i]; exist_actions {
			for k, val := range act {
				*sdrive.Injectables[k] = val
			}
		}

		// TODO Periodic set

		if _, err := vm.Step(sconfig); err != nil {
			return 0, err
		}

		// This will get value to report on this tick
		if rep, existReports := srep.AbsGet[i]; existReports {
			for k, _ := range rep {
				rep[k] = *srep.Reportables[k]
			}
			for k, v := range rep {
				fmt.Print(k, ",", v, ",")
			}
			// fmt.Println("TEMPORARY:", rep)
		}

		// TODO Periodic get

	}

	return 0, nil
}
