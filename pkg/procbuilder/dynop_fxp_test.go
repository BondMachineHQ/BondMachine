package procbuilder

import (
	"testing"
)

func TestFxpMult(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		regsize int64
		fsize   int64
		want    int64
	}{
		{
			name:    "simple multiplication 2*3 with 8-bit regsize and 4-bit frac",
			a:       32, // 2.0 in Q4.4 format (2 << 4)
			b:       48, // 3.0 in Q4.4 format (3 << 4)
			regsize: 8,
			fsize:   4,
			want:    96, // 6.0 in Q4.4 format (6 << 4)
		},
		{
			name:    "multiplication with fractional parts 1.5*2.5",
			a:       24, // 1.5 in Q4.4 format (1.5 * 16 = 24)
			b:       40, // 2.5 in Q4.4 format (2.5 * 16 = 40)
			regsize: 8,
			fsize:   4,
			want:    60, // 3.75 in Q4.4 format (3.75 * 16 = 60)
		},
		{
			name:    "multiplication with zero",
			a:       0,
			b:       48,
			regsize: 8,
			fsize:   4,
			want:    0,
		},
		{
			name:    "multiplication both zero",
			a:       0,
			b:       0,
			regsize: 8,
			fsize:   4,
			want:    0,
		},
		{
			name:    "16-bit multiplication 4*5 with 8-bit frac",
			a:       1024, // 4.0 in Q8.8 format (4 << 8)
			b:       1280, // 5.0 in Q8.8 format (5 << 8)
			regsize: 16,
			fsize:   8,
			want:    5120, // 20.0 in Q8.8 format (20 << 8)
		},
		{
			name:    "small fractional multiplication 0.5*0.5",
			a:       8, // 0.5 in Q4.4 format (0.5 * 16 = 8)
			b:       8, // 0.5 in Q4.4 format
			regsize: 8,
			fsize:   4,
			want:    4, // 0.25 in Q4.4 format (0.25 * 16 = 4)
		},
		{
			name:    "negative number multiplication -2*3",
			a:       -32, // -2.0 in Q4.4 format
			b:       48,  // 3.0 in Q4.4 format
			regsize: 8,
			fsize:   4,
			want:    -96, // -6.0 in Q4.4 format
		},
		{
			name:    "both negative -2*-3",
			a:       -32, // -2.0 in Q4.4 format
			b:       -48, // -3.0 in Q4.4 format
			regsize: 8,
			fsize:   4,
			want:    96, // 6.0 in Q4.4 format
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fxpMult(tt.a, tt.b, tt.regsize, tt.fsize)
			if got != tt.want {
				t.Errorf("fxpMult(%d, %d, %d, %d) = %d, want %d",
					tt.a, tt.b, tt.regsize, tt.fsize, got, tt.want)
			}
		})
	}
}

func TestFxpDiv(t *testing.T) {
	tests := []struct {
		name    string
		num1    int64
		num2    int64
		regsize int64
		fsize   int64
		want    int64
	}{
		{
			name:    "simple division 6/2 with 8-bit regsize and 4-bit frac",
			num1:    96, // 6.0 in Q4.4 format (6 << 4)
			num2:    32, // 2.0 in Q4.4 format (2 << 4)
			regsize: 8,
			fsize:   4,
			want:    48, // 3.0 in Q4.4 format (3 << 4)
		},
		{
			name:    "division with fractional result 5/2",
			num1:    80, // 5.0 in Q4.4 format (5 << 4)
			num2:    32, // 2.0 in Q4.4 format (2 << 4)
			regsize: 8,
			fsize:   4,
			want:    40, // 2.5 in Q4.4 format (2.5 * 16 = 40)
		},
		{
			name:    "division resulting in less than 1, (1/2)",
			num1:    16, // 1.0 in Q4.4 format (1 << 4)
			num2:    32, // 2.0 in Q4.4 format (2 << 4)
			regsize: 8,
			fsize:   4,
			want:    8, // 0.5 in Q4.4 format (0.5 * 16 = 8)
		},
		{
			name:    "division by 1",
			num1:    64, // 4.0 in Q4.4 format
			num2:    16, // 1.0 in Q4.4 format
			regsize: 8,
			fsize:   4,
			want:    64, // 4.0 in Q4.4 format
		},
		{
			name:    "16-bit division 20/4 with 8-bit frac",
			num1:    5120, // 20.0 in Q8.8 format (20 << 8)
			num2:    1024, // 4.0 in Q8.8 format (4 << 8)
			regsize: 16,
			fsize:   8,
			want:    1280, // 5.0 in Q8.8 format (5 << 8)
		},
		{
			name:    "fractional division 1.5/0.5",
			num1:    24, // 1.5 in Q4.4 format (1.5 * 16 = 24)
			num2:    8,  // 0.5 in Q4.4 format (0.5 * 16 = 8)
			regsize: 8,
			fsize:   4,
			want:    48, // 3.0 in Q4.4 format (3 << 4)
		},
		{
			name:    "negative dividend -6/2",
			num1:    -96, // -6.0 in Q4.4 format
			num2:    32,  // 2.0 in Q4.4 format
			regsize: 8,
			fsize:   4,
			want:    -48, // -3.0 in Q4.4 format
		},
		{
			name:    "negative divisor 6/-2",
			num1:    96,  // 6.0 in Q4.4 format
			num2:    -32, // -2.0 in Q4.4 format
			regsize: 8,
			fsize:   4,
			want:    -48, // -3.0 in Q4.4 format
		},
		{
			name:    "both negative -6/-2",
			num1:    -96, // -6.0 in Q4.4 format
			num2:    -32, // -2.0 in Q4.4 format
			regsize: 8,
			fsize:   4,
			want:    48, // 3.0 in Q4.4 format
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fxpDiv(tt.num1, tt.num2, tt.regsize, tt.fsize)
			if got != tt.want {
				t.Errorf("fxpDiv(%d, %d, %d, %d) = %d, want %d",
					tt.num1, tt.num2, tt.regsize, tt.fsize, got, tt.want)
			}
		})
	}
}

// TestFxpMultDiv tests multiplication followed by division to verify round-trip accuracy
func TestFxpMultDiv(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		regsize int64
		fsize   int64
	}{
		{
			name:    "round-trip 3*4/4 should equal 3",
			a:       48, // 3.0 in Q4.4 format
			b:       64, // 4.0 in Q4.4 format
			regsize: 8,
			fsize:   4,
		},
		{
			name:    "round-trip with fractional 2.5*2.0/2.0 should equal 2.5",
			a:       40, // 2.5 in Q4.4 format
			b:       32, // 2.0 in Q4.4 format
			regsize: 8,
			fsize:   4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Multiply then divide
			multResult := fxpMult(tt.a, tt.b, tt.regsize, tt.fsize)
			divResult := fxpDiv(multResult, tt.b, tt.regsize, tt.fsize)

			// Should get back the original value (within rounding tolerance)
			if divResult != tt.a {
				t.Logf("Round-trip test: (%d * %d) / %d = %d, expected %d",
					tt.a, tt.b, tt.b, divResult, tt.a)
				// Allow small rounding differences
				diff := divResult - tt.a
				if diff < -1 || diff > 1 {
					t.Errorf("Round-trip error too large: got %d, want %d (diff: %d)",
						divResult, tt.a, diff)
				}
			}
		})
	}
}

// BenchmarkFxpMult benchmarks the fixed-point multiplication function
func BenchmarkFxpMult(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fxpMult(48, 32, 8, 4)
	}
}

// BenchmarkFxpDiv benchmarks the fixed-point division function
func BenchmarkFxpDiv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fxpDiv(96, 32, 8, 4)
	}
}
