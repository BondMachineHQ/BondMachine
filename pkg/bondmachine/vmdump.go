package bondmachine

import "fmt"

func (vm *VM) String() string {
	if vm == nil {
		return "Empty VM"
	}

	result := "BondMachine VM State\n"
	result += "====================\n\n"

	// Absolute tick
	result += fmt.Sprintf("Absolute Tick: %d\n\n", vm.abs_tick)

	// BondMachine info
	if vm.Bmach != nil {
		result += fmt.Sprintln("BondMachine:")
		result += fmt.Sprintf("  Register Size: %d bits\n", vm.Bmach.Rsize)
		result += fmt.Sprintf("  Processors: %d\n", len(vm.Bmach.Processors))
		result += fmt.Sprintf("  Inputs: %d\n", vm.Bmach.Inputs)
		result += fmt.Sprintf("  Outputs: %d\n", vm.Bmach.Outputs)
		result += fmt.Sprintf("  Internal Inputs: %d\n", len(vm.Bmach.Internal_inputs))
		result += fmt.Sprintf("  Internal Outputs: %d\n\n", len(vm.Bmach.Internal_outputs))
	}

	// BM Inputs
	if len(vm.Inputs_regs) > 0 {
		result += "BM Inputs:\n"
		for i, reg := range vm.Inputs_regs {
			result += fmt.Sprintf("  i%d: %s (v:%t r:%t)\n",
				i,
				vm.dumpRegister(reg),
				vm.InputsValid[i],
				vm.InputsRecv[i])
		}
		result += "\n"
	}

	// BM Outputs
	if len(vm.Outputs_regs) > 0 {
		result += "BM Outputs:\n"
		for i, reg := range vm.Outputs_regs {
			result += fmt.Sprintf("  o%d: %s (v:%t r:%t)\n",
				i,
				vm.dumpRegister(reg),
				vm.OutputsValid[i],
				vm.OutputsRecv[i])
		}
		result += "\n"
	}

	// Internal Inputs
	if len(vm.Internal_inputs_regs) > 0 {
		result += "Internal Inputs:\n"
		for i, reg := range vm.Internal_inputs_regs {
			name := ""
			if vm.Bmach != nil && i < len(vm.Bmach.Internal_inputs) {
				name = vm.Bmach.Internal_inputs[i].String()
			}
			result += fmt.Sprintf("  [%d] %s: %s (v:%t r:%t)\n",
				i,
				name,
				vm.dumpRegister(reg),
				vm.InternalInputsValid[i],
				vm.InternalInputsRecv[i])
		}
		result += "\n"
	}

	// Internal Outputs
	if len(vm.Internal_outputs_regs) > 0 {
		result += "Internal Outputs:\n"
		for i, reg := range vm.Internal_outputs_regs {
			name := ""
			if vm.Bmach != nil && i < len(vm.Bmach.Internal_outputs) {
				name = vm.Bmach.Internal_outputs[i].String()
			}
			result += fmt.Sprintf("  [%d] %s: %s (v:%t r:%t)\n",
				i,
				name,
				vm.dumpRegister(reg),
				vm.InternalOutputsValid[i],
				vm.InternalOutputsRecv[i])
		}
		result += "\n"
	}

	// Processors
	if len(vm.Processors) > 0 {
		result += "Processors:\n"
		for i, proc := range vm.Processors {
			if proc != nil && proc.Mach != nil {
				result += fmt.Sprintf("  Processor %d:\n", i)
				result += fmt.Sprintf("    PC: %d\n", proc.Pc)

				// Processor Registers
				if len(proc.Registers) > 0 {
					result += "    Registers:\n"
					for j, reg := range proc.Registers {
						result += fmt.Sprintf("      r%d: %s\n", j, vm.dumpRegister(reg))
					}
				}

				// Processor Inputs
				if len(proc.Inputs) > 0 {
					result += "    Inputs:\n"
					for j, input := range proc.Inputs {
						result += fmt.Sprintf("      i%d: %s (v:%t r:%t)\n",
							j,
							vm.dumpRegister(input),
							proc.InputsValid[j],
							proc.InputsRecv[j])
					}
				}

				// Processor Outputs
				if len(proc.Outputs) > 0 {
					result += "    Outputs:\n"
					for j, output := range proc.Outputs {
						result += fmt.Sprintf("      o%d: %s (v:%t r:%t)\n",
							j,
							vm.dumpRegister(output),
							proc.OutputsValid[j],
							proc.OutputsRecv[j])
					}
				}

				// Processor Memory (if any)
				if len(proc.Memory) > 0 {
					result += fmt.Sprintf("    Memory: %d locations\n", len(proc.Memory))
				}

				result += "\n"
			}
		}
	}

	// Deferred Instructions
	if len(vm.DeferredInstructions) > 0 {
		result += fmt.Sprintf("Deferred Instructions: %d\n", len(vm.DeferredInstructions))
		for name := range vm.DeferredInstructions {
			result += fmt.Sprintf("  - %s\n", name)
		}
		result += "\n"
	}

	// Emulation Drivers
	if len(vm.EmuDrivers) > 0 {
		result += fmt.Sprintf("Emulation Drivers: %d\n", len(vm.EmuDrivers))
		for i, ed := range vm.EmuDrivers {
			result += fmt.Sprintf("  [%d] %T\n", i, ed)
		}
		result += "\n"
	}

	// Links (if BondMachine exists)
	if vm.Bmach != nil && len(vm.Bmach.Links) > 0 {
		result += "Internal Links:\n"
		for i, link := range vm.Bmach.Links {
			if link != -1 {
				inName := ""
				outName := ""
				if i < len(vm.Bmach.Internal_inputs) {
					inName = vm.Bmach.Internal_inputs[i].String()
				}
				if link < len(vm.Bmach.Internal_outputs) {
					outName = vm.Bmach.Internal_outputs[link].String()
				}
				result += fmt.Sprintf("  %s -> %s\n", outName, inName)
			}
		}
		result += "\n"
	}

	return result
}
