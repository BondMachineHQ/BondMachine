package bmmatrix

const (
	SQRT2 = 1.4142135623730951
)

func Hadamard() *BmMatrixSquareComplex {
	m := NewBmMatrixSquareComplex(2)
	m.Data[0][0] = Complex32{1.0 / SQRT2, 0.0}
	m.Data[0][1] = Complex32{1.0 / SQRT2, 0.0}
	m.Data[1][0] = Complex32{1.0 / SQRT2, 0.0}
	m.Data[1][1] = Complex32{-1.0 / SQRT2, 0.0}
	return m
}

func PauliX() *BmMatrixSquareComplex {
	m := NewBmMatrixSquareComplex(2)
	m.Data[0][0] = Complex32{0.0, 0.0}
	m.Data[0][1] = Complex32{1.0, 0.0}
	m.Data[1][0] = Complex32{1.0, 0.0}
	m.Data[1][1] = Complex32{0.0, 0.0}
	return m
}

func X() *BmMatrixSquareComplex {
	return PauliX()
}

func PauliY() *BmMatrixSquareComplex {
	m := NewBmMatrixSquareComplex(2)
	m.Data[0][0] = Complex32{0.0, 0.0}
	m.Data[0][1] = Complex32{0.0, -1.0}
	m.Data[1][0] = Complex32{0.0, 1.0}
	m.Data[1][1] = Complex32{0.0, 0.0}
	return m
}

func Y() *BmMatrixSquareComplex {
	return PauliY()
}

func PauliZ() *BmMatrixSquareComplex {
	m := NewBmMatrixSquareComplex(2)
	m.Data[0][0] = Complex32{1.0, 0.0}
	m.Data[0][1] = Complex32{0.0, 0.0}
	m.Data[1][0] = Complex32{0.0, 0.0}
	m.Data[1][1] = Complex32{-1.0, 0.0}
	return m
}

func Z() *BmMatrixSquareComplex {
	return PauliZ()
}
