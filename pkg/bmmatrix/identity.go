package bmmatrix

func IdentityComplex(n int) *BmMatrixSquareComplex {
	m := NewBmMatrixSquareComplex(n)
	for i := 0; i < n; i++ {
		m.Data[i][i] = Complex32{1.0, 0.0}
	}
	return m
}
