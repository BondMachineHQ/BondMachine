package bmnumbers

import (
	"testing"
)

func TestNumberToUint64(t *testing.T) {

	dataSetUint := testDataSet("uintData")

	// Loop through the inputs and expected outputs
	for input, expected := range dataSetUint.(map[string]uint64) {
		// Get the output
		output, err := ImportString(input)
		if err != nil {
			t.Errorf("Error: %s", err)
		} else {
			// Check if the output is the expected output
			if value, err := output.ExportUint64(); err != nil {
				t.Errorf("Error: %s", err)
			} else {
				if value != expected {
					t.Errorf("did not get expected output: %d != %d", value, expected)
				} else {
					t.Logf("got expected output: %d == %d", value, expected)
				}
			}
		}
	}
}

func TestNumberToBinary(t *testing.T) {

	dataSetBin := testDataSet("binData")

	// Loop through the inputs and expected outputs
	for input, expected := range dataSetBin.(map[string]string) {
		// Get the output
		output, err := ImportString(input)
		if err != nil {
			t.Errorf("Error: %s", err)
		} else {
			// Check if the output is the expected output
			if value, err := output.ExportBinary(true); err != nil {
				t.Errorf("Error: %s", err)
			} else {
				if value != expected {
					t.Errorf("did not get expected output: %s != %s", value, expected)
				} else {
					t.Logf("got expected output: %s == %s", value, expected)
				}
			}
		}
	}

	dataSetBin = testDataSet("binDataNoSize")

	// Loop through the inputs and expected outputs
	for input, expected := range dataSetBin.(map[string]string) {
		// Get the output
		output, err := ImportString(input)
		if err != nil {
			t.Errorf("Error: %s", err)
		} else {
			// Check if the output is the expected output
			if value, err := output.ExportBinary(false); err != nil {
				t.Errorf("Error: %s", err)
			} else {
				if value != expected {
					t.Errorf("did not get expected output: %s != %s", value, expected)
				} else {
					t.Logf("got expected output: %s == %s", value, expected)
				}
			}
		}
	}
}
