package bmmatrix

import "math"

func IdentityComplex(n int) *BmMatrixSquareComplex {
	m := NewBmMatrixSquareComplex(n)
	for i := 0; i < n; i++ {
		m.Data[i][i] = Complex32{1.0, 0.0}
	}
	return m
}

func GlobalPhase(n int, phase float32) *BmMatrixSquareComplex {

	// The real part is the cosine of the phase
	realPart := float32(math.Cos(float64(phase)))
	imagPart := float32(math.Sin(float64(phase)))
	m := NewBmMatrixSquareComplex(n)
	for i := 0; i < n; i++ {
		m.Data[i][i] = Complex32{realPart, imagPart}
	}
	return m
}

func PhaseShift(phase float32) *BmMatrixSquareComplex {
	m := NewBmMatrixSquareComplex(2)
	m.Data[0][0] = Complex32{1.0, 0.0}
	m.Data[1][1] = Complex32{float32(math.Cos(float64(phase))), float32(math.Sin(float64(phase)))}
	return m
}

func T() *BmMatrixSquareComplex {
	return PhaseShift(math.Pi / 4)
}
