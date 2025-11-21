package bondmachine

import (
	"testing"

	"github.com/BondMachineHQ/BondMachine/pkg/simbox"
)

func TestSuspendedRulesAreSkipped(t *testing.T) {
	// Create a simbox with both active and suspended rules
	sbox := &simbox.Simbox{
		Rules: []simbox.Rule{
			{Timec: simbox.TIMEC_NONE, Action: simbox.ACTION_CONFIG, Object: "show_ticks", Suspended: false},
			{Timec: simbox.TIMEC_NONE, Action: simbox.ACTION_CONFIG, Object: "show_io_pre", Suspended: true},
			{Timec: simbox.TIMEC_NONE, Action: simbox.ACTION_CONFIG, Object: "show_io_post", Suspended: false},
		},
	}

	// Create a minimal VM structure for testing
	vm := &VM{
		Bmach: &Bondmachine{
			Rsize: 8,
		},
	}

	// Initialize simulation config
	sc := &Sim_config{}
	conf := &Config{Debug: false}
	err := sc.Init(sbox, vm, conf)
	if err != nil {
		t.Fatalf("Failed to initialize Sim_config: %v", err)
	}

	// Verify that active rules are applied
	if !sc.Show_ticks {
		t.Error("show_ticks should be enabled (rule was not suspended)")
	}
	if !sc.Show_io_post {
		t.Error("show_io_post should be enabled (rule was not suspended)")
	}

	// Verify that suspended rules are NOT applied
	if sc.Show_io_pre {
		t.Error("show_io_pre should NOT be enabled (rule was suspended)")
	}
}

func TestSuspendedSetRulesAreSkipped(t *testing.T) {
	// Create a simbox with suspended SET rules
	sbox := &simbox.Simbox{
		Rules: []simbox.Rule{
			{Timec: simbox.TIMEC_ABS, Tick: 10, Action: simbox.ACTION_SET, Object: "i0", Extra: "5", Suspended: false},
			{Timec: simbox.TIMEC_ABS, Tick: 20, Action: simbox.ACTION_SET, Object: "i1", Extra: "10", Suspended: true},
		},
	}

	// Create a minimal VM structure for testing
	vm := &VM{
		Bmach: &Bondmachine{
			Rsize:  8,
			Inputs: 2,
		},
		Inputs_regs: []interface{}{uint8(0), uint8(0)},
	}

	// Initialize simulation drive
	sd := &Sim_drive{}
	conf := &Config{}
	err := sd.Init(conf, sbox, vm)
	if err != nil {
		t.Fatalf("Failed to initialize Sim_drive: %v", err)
	}

	// Verify that the active rule's tick is present
	if _, ok := sd.AbsSet[10]; !ok {
		t.Error("Tick 10 should be present (rule was not suspended)")
	}

	// Verify that the suspended rule's tick is NOT present
	if _, ok := sd.AbsSet[20]; ok {
		t.Error("Tick 20 should NOT be present (rule was suspended)")
	}

	// Verify that we have only 1 injectable (from the non-suspended rule)
	if len(sd.Injectables) != 1 {
		t.Errorf("Expected 1 injectable, got %d", len(sd.Injectables))
	}
}
