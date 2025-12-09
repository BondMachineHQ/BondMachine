package simbox

import (
	"strings"
	"testing"
)

func TestRuleSuspension(t *testing.T) {
	// Create a simbox with some rules
	sb := &Simbox{
		Rules: []Rule{
			{Timec: TIMEC_ABS, Tick: 10, Action: ACTION_SET, Object: "i0", Extra: "5", Suspended: false},
			{Timec: TIMEC_REL, Tick: 5, Action: ACTION_GET, Object: "o0", Extra: "unsigned", Suspended: false},
			{Timec: TIMEC_NONE, Action: ACTION_CONFIG, Object: "show_ticks", Extra: "", Suspended: false},
		},
	}

	// Test initial state - all rules should be active
	for i, rule := range sb.Rules {
		if rule.Suspended {
			t.Errorf("Rule %d should not be suspended initially", i)
		}
	}

	// Test suspending a rule
	err := sb.Suspend(1)
	if err != nil {
		t.Errorf("Failed to suspend rule: %v", err)
	}
	if !sb.Rules[1].Suspended {
		t.Error("Rule 1 should be suspended after Suspend() call")
	}

	// Test that other rules are not affected
	if sb.Rules[0].Suspended || sb.Rules[2].Suspended {
		t.Error("Other rules should remain active")
	}

	// Test reactivating a rule
	err = sb.Reactivate(1)
	if err != nil {
		t.Errorf("Failed to reactivate rule: %v", err)
	}
	if sb.Rules[1].Suspended {
		t.Error("Rule 1 should not be suspended after Reactivate() call")
	}

	// Test suspending with invalid index
	err = sb.Suspend(99)
	if err == nil {
		t.Error("Expected error when suspending rule with invalid index")
	}

	// Test reactivating with invalid index
	err = sb.Reactivate(99)
	if err == nil {
		t.Error("Expected error when reactivating rule with invalid index")
	}
}

func TestPrintWithSuspendedRules(t *testing.T) {
	sb := &Simbox{
		Rules: []Rule{
			{Timec: TIMEC_ABS, Tick: 10, Action: ACTION_SET, Object: "i0", Extra: "5", Suspended: false},
			{Timec: TIMEC_REL, Tick: 5, Action: ACTION_GET, Object: "o0", Extra: "unsigned", Suspended: true},
		},
	}

	output := sb.Print()

	// Check that suspended marker appears for suspended rule
	if !strings.Contains(output, "[SUSPENDED]") {
		t.Error("Print() should show [SUSPENDED] marker for suspended rules")
	}

	// Check that the output contains rule information
	if !strings.Contains(output, "000 -") || !strings.Contains(output, "001 -") {
		t.Error("Print() should show rule indices")
	}
}

func TestAddRuleDefaultsToNotSuspended(t *testing.T) {
	sb := &Simbox{
		Rules: []Rule{},
	}

	// Add a rule
	err := sb.Add("absolute:10:set:i0:5")
	if err != nil {
		t.Errorf("Failed to add rule: %v", err)
	}

	// Check that the new rule is not suspended by default
	if len(sb.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(sb.Rules))
	}
	if sb.Rules[0].Suspended {
		t.Error("Newly added rule should not be suspended by default")
	}
}
