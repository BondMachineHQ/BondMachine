package bondmachine

import (
	"testing"

	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

func TestCopyState(t *testing.T) {
	// Create a source VM with state
	sourceVM := &VM{
		Bmach: &Bondmachine{
			Rsize:            8,
			Inputs:           2,
			Outputs:          2,
			Internal_inputs:  []Bond{{Map_to: CPINPUT, Res_id: 0, Ext_id: 0}},
			Internal_outputs: []Bond{{Map_to: CPOUTPUT, Res_id: 0, Ext_id: 0}},
			Processors:       []int{0},
			Domains:          []*procbuilder.Machine{{}},
		},
	}

	// Initialize source VM
	err := sourceVM.Init()
	if err != nil {
		t.Fatalf("Failed to initialize source VM: %v", err)
	}

	// Set some state in source VM
	sourceVM.Inputs_regs[0] = uint8(42)
	sourceVM.Inputs_regs[1] = uint8(123)
	sourceVM.Outputs_regs[0] = uint8(99)
	sourceVM.Outputs_regs[1] = uint8(88)
	sourceVM.Internal_inputs_regs[0] = uint8(77)
	sourceVM.Internal_outputs_regs[0] = uint8(66)

	sourceVM.InputsValid[0] = true
	sourceVM.InputsValid[1] = false
	sourceVM.OutputsValid[0] = true
	sourceVM.OutputsValid[1] = false
	sourceVM.InternalInputsValid[0] = true
	sourceVM.InternalOutputsValid[0] = false

	sourceVM.InputsRecv[0] = true
	sourceVM.InputsRecv[1] = false
	sourceVM.OutputsRecv[0] = false
	sourceVM.OutputsRecv[1] = true
	sourceVM.InternalInputsRecv[0] = false
	sourceVM.InternalOutputsRecv[0] = true

	sourceVM.DeferredInstructions["test"] = func(*VM) bool { return true }
	sourceVM.abs_tick = uint64(100)

	// Create a destination VM
	destVM := &VM{
		Bmach: sourceVM.Bmach,
	}

	// Initialize destination VM
	err = destVM.Init()
	if err != nil {
		t.Fatalf("Failed to initialize dest VM: %v", err)
	}

	// Copy state
	err = destVM.CopyState(sourceVM)
	if err != nil {
		t.Fatalf("Failed to copy state: %v", err)
	}

	// Verify all register copies
	if destVM.Inputs_regs[0].(uint8) != 42 {
		t.Errorf("Inputs_regs[0]: expected 42, got %v", destVM.Inputs_regs[0])
	}
	if destVM.Inputs_regs[1].(uint8) != 123 {
		t.Errorf("Inputs_regs[1]: expected 123, got %v", destVM.Inputs_regs[1])
	}
	if destVM.Outputs_regs[0].(uint8) != 99 {
		t.Errorf("Outputs_regs[0]: expected 99, got %v", destVM.Outputs_regs[0])
	}
	if destVM.Outputs_regs[1].(uint8) != 88 {
		t.Errorf("Outputs_regs[1]: expected 88, got %v", destVM.Outputs_regs[1])
	}
	if destVM.Internal_inputs_regs[0].(uint8) != 77 {
		t.Errorf("Internal_inputs_regs[0]: expected 77, got %v", destVM.Internal_inputs_regs[0])
	}
	if destVM.Internal_outputs_regs[0].(uint8) != 66 {
		t.Errorf("Internal_outputs_regs[0]: expected 66, got %v", destVM.Internal_outputs_regs[0])
	}

	// Verify all valid flags
	if destVM.InputsValid[0] != true {
		t.Error("InputsValid[0]: expected true")
	}
	if destVM.InputsValid[1] != false {
		t.Error("InputsValid[1]: expected false")
	}
	if destVM.OutputsValid[0] != true {
		t.Error("OutputsValid[0]: expected true")
	}
	if destVM.OutputsValid[1] != false {
		t.Error("OutputsValid[1]: expected false")
	}
	if destVM.InternalInputsValid[0] != true {
		t.Error("InternalInputsValid[0]: expected true")
	}
	if destVM.InternalOutputsValid[0] != false {
		t.Error("InternalOutputsValid[0]: expected false")
	}

	// Verify all recv flags
	if destVM.InputsRecv[0] != true {
		t.Error("InputsRecv[0]: expected true")
	}
	if destVM.InputsRecv[1] != false {
		t.Error("InputsRecv[1]: expected false")
	}
	if destVM.OutputsRecv[0] != false {
		t.Error("OutputsRecv[0]: expected false")
	}
	if destVM.OutputsRecv[1] != true {
		t.Error("OutputsRecv[1]: expected true")
	}
	if destVM.InternalInputsRecv[0] != false {
		t.Error("InternalInputsRecv[0]: expected false")
	}
	if destVM.InternalOutputsRecv[0] != true {
		t.Error("InternalOutputsRecv[0]: expected true")
	}

	// Verify deferred instructions
	if _, ok := destVM.DeferredInstructions["test"]; !ok {
		t.Error("DeferredInstructions not copied correctly")
	}

	// Verify absolute tick
	if destVM.abs_tick != 100 {
		t.Errorf("abs_tick: expected 100, got %d", destVM.abs_tick)
	}
}

func TestCopyStateIndependence(t *testing.T) {
	// Create a source VM with state
	sourceVM := &VM{
		Bmach: &Bondmachine{
			Rsize:            8,
			Inputs:           1,
			Outputs:          1,
			Internal_inputs:  []Bond{},
			Internal_outputs: []Bond{},
			Processors:       []int{},
			Domains:          []*procbuilder.Machine{},
		},
	}

	// Initialize source VM
	err := sourceVM.Init()
	if err != nil {
		t.Fatalf("Failed to initialize source VM: %v", err)
	}

	// Set some state in source VM
	sourceVM.Inputs_regs[0] = uint8(50)
	sourceVM.InputsValid[0] = true
	sourceVM.abs_tick = uint64(50)

	// Create a destination VM
	destVM := &VM{
		Bmach: sourceVM.Bmach,
	}

	// Initialize destination VM
	err = destVM.Init()
	if err != nil {
		t.Fatalf("Failed to initialize dest VM: %v", err)
	}

	// Copy state
	err = destVM.CopyState(sourceVM)
	if err != nil {
		t.Fatalf("Failed to copy state: %v", err)
	}

	// Modify source VM after copy
	sourceVM.Inputs_regs[0] = uint8(100)
	sourceVM.InputsValid[0] = false
	sourceVM.abs_tick = uint64(200)

	// Verify destination VM is unchanged
	if destVM.Inputs_regs[0].(uint8) != 50 {
		t.Errorf("Inputs_regs[0]: expected 50, got %v (should be independent of source changes)", destVM.Inputs_regs[0])
	}
	if destVM.InputsValid[0] != true {
		t.Error("InputsValid[0]: expected true (should be independent of source changes)")
	}
	if destVM.abs_tick != 50 {
		t.Errorf("abs_tick: expected 50, got %d (should be independent of source changes)", destVM.abs_tick)
	}
}

func TestCopyStateDifferentRegisterSizes(t *testing.T) {
	testCases := []struct {
		name  string
		rsize uint8
		value interface{}
	}{
		{"8-bit", 8, uint8(42)},
		{"16-bit", 16, uint16(1000)},
		{"32-bit", 32, uint32(100000)},
		{"64-bit", 64, uint64(10000000000)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sourceVM := &VM{
				Bmach: &Bondmachine{
					Rsize:            tc.rsize,
					Inputs:           1,
					Outputs:          1,
					Internal_inputs:  []Bond{},
					Internal_outputs: []Bond{},
					Processors:       []int{},
					Domains:          []*procbuilder.Machine{},
				},
			}

			err := sourceVM.Init()
			if err != nil {
				t.Fatalf("Failed to initialize source VM: %v", err)
			}

			sourceVM.Inputs_regs[0] = tc.value

			destVM := &VM{
				Bmach: sourceVM.Bmach,
			}

			err = destVM.Init()
			if err != nil {
				t.Fatalf("Failed to initialize dest VM: %v", err)
			}

			err = destVM.CopyState(sourceVM)
			if err != nil {
				t.Fatalf("Failed to copy state: %v", err)
			}

			if destVM.Inputs_regs[0] != tc.value {
				t.Errorf("Register value not copied correctly: expected %v, got %v", tc.value, destVM.Inputs_regs[0])
			}
		})
	}
}
