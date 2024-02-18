package bmmatrix

func Cnot() *BmMatrixSquareComplex {
	m := NewBmMatrixSquareComplex(4)
	m.Data[0][0] = Complex32{1.0, 0.0}
	m.Data[0][1] = Complex32{0.0, 0.0}
	m.Data[0][2] = Complex32{0.0, 0.0}
	m.Data[0][3] = Complex32{0.0, 0.0}
	m.Data[1][0] = Complex32{0.0, 0.0}
	m.Data[1][1] = Complex32{1.0, 0.0}
	m.Data[1][2] = Complex32{0.0, 0.0}
	m.Data[1][3] = Complex32{0.0, 0.0}
	m.Data[2][0] = Complex32{0.0, 0.0}
	m.Data[2][1] = Complex32{0.0, 0.0}
	m.Data[2][2] = Complex32{0.0, 0.0}
	m.Data[2][3] = Complex32{1.0, 0.0}
	m.Data[3][0] = Complex32{0.0, 0.0}
	m.Data[3][1] = Complex32{0.0, 0.0}
	m.Data[3][2] = Complex32{1.0, 0.0}
	m.Data[3][3] = Complex32{0.0, 0.0}
	return m
}
