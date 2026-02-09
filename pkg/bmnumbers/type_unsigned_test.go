package bmnumbers

import (
	"testing"
)

func TestUnsignedWithSize(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		wantValue   uint64
		wantBits    int
	}{
		{
			name:        "8-bit unsigned max value",
			input:       "0u<8>255",
			expectError: false,
			wantValue:   255,
			wantBits:    8,
		},
		{
			name:        "8-bit unsigned min value",
			input:       "0u<8>0",
			expectError: false,
			wantValue:   0,
			wantBits:    8,
		},
		{
			name:        "8-bit unsigned overflow",
			input:       "0u<8>256",
			expectError: true,
		},
		{
			name:        "16-bit unsigned",
			input:       "0u<16>1000",
			expectError: false,
			wantValue:   1000,
			wantBits:    16,
		},
		{
			name:        "16-bit unsigned max value",
			input:       "0u<16>65535",
			expectError: false,
			wantValue:   65535,
			wantBits:    16,
		},
		{
			name:        "16-bit unsigned overflow",
			input:       "0u<16>65536",
			expectError: true,
		},
		{
			name:        "32-bit unsigned",
			input:       "0u<32>123456",
			expectError: false,
			wantValue:   123456,
			wantBits:    32,
		},
		{
			name:        "32-bit unsigned max value",
			input:       "0u<32>4294967295",
			expectError: false,
			wantValue:   4294967295,
			wantBits:    32,
		},
		{
			name:        "64-bit unsigned",
			input:       "0u<64>9876543210",
			expectError: false,
			wantValue:   9876543210,
			wantBits:    64,
		},
		{
			name:        "decimal notation 8-bit",
			input:       "0d<8>200",
			expectError: false,
			wantValue:   200,
			wantBits:    8,
		},
		{
			name:        "decimal notation 16-bit",
			input:       "0d<16>5000",
			expectError: false,
			wantValue:   5000,
			wantBits:    16,
		},
		{
			name:        "4-bit unsigned",
			input:       "0u<4>15",
			expectError: false,
			wantValue:   15,
			wantBits:    4,
		},
		{
			name:        "4-bit unsigned overflow",
			input:       "0u<4>16",
			expectError: true,
		},
		{
			name:        "1-bit unsigned (0)",
			input:       "0u<1>0",
			expectError: false,
			wantValue:   0,
			wantBits:    1,
		},
		{
			name:        "1-bit unsigned (1)",
			input:       "0u<1>1",
			expectError: false,
			wantValue:   1,
			wantBits:    1,
		},
		{
			name:        "1-bit unsigned overflow",
			input:       "0u<1>2",
			expectError: true,
		},
		{
			name:        "invalid size (0)",
			input:       "0u<0>10",
			expectError: true,
		},
		{
			name:        "invalid size (>64)",
			input:       "0u<65>10",
			expectError: true,
		},
		{
			name:        "12-bit unsigned",
			input:       "0u<12>2048",
			expectError: false,
			wantValue:   2048,
			wantBits:    12,
		},
		{
			name:        "12-bit unsigned max",
			input:       "0u<12>4095",
			expectError: false,
			wantValue:   4095,
			wantBits:    12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := ImportString(tt.input)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none for input %s", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}

			// Check the value
			value, err := output.ExportUint64()
			if err != nil {
				t.Errorf("error exporting value: %s", err)
				return
			}

			if value != tt.wantValue {
				t.Errorf("got value %d, want %d", value, tt.wantValue)
			}

			// Check the bit size
			if output.bits != tt.wantBits {
				t.Errorf("got bits %d, want %d", output.bits, tt.wantBits)
			}

			// Verify type
			if output.nType.GetName() != "unsigned" {
				t.Errorf("got type %s, want unsigned", output.nType.GetName())
			}
		})
	}
}

func TestUnsignedWithSizeExport(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantString   string
		wantBinary   string
	}{
		{
			name:         "8-bit value export",
			input:        "0u<8>42",
			wantString:   "42",
			wantBinary:   "0b<8>101010",
		},
		{
			name:         "16-bit value export",
			input:        "0u<16>1000",
			wantString:   "1000",
			wantBinary:   "0b<16>1111101000",
		},
		{
			name:         "4-bit value export",
			input:        "0u<4>7",
			wantString:   "7",
			wantBinary:   "0b<4>111",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := ImportString(tt.input)
			if err != nil {
				t.Fatalf("error importing: %s", err)
			}

			// Check string export
			strValue, err := output.ExportString(nil)
			if err != nil {
				t.Errorf("error exporting string: %s", err)
			} else if strValue != tt.wantString {
				t.Errorf("got string %s, want %s", strValue, tt.wantString)
			}

			// Check binary export
			binValue, err := output.ExportBinary(true)
			if err != nil {
				t.Errorf("error exporting binary: %s", err)
			} else if binValue != tt.wantBinary {
				t.Errorf("got binary %s, want %s", binValue, tt.wantBinary)
			}
		})
	}
}

// BenchmarkUnsignedImportWithSize benchmarks the unsigned import with size function
func BenchmarkUnsignedImportWithSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ImportString("0u<16>1234")
	}
}

// BenchmarkUnsignedImportNoSize benchmarks the unsigned import without size function
func BenchmarkUnsignedImportNoSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ImportString("1234")
	}
}
