package bmmatrix

import (
	"fmt"
	"testing"
)

// Print test

func TestPrint(t *testing.T) {
	m1 := NewBmMatrixSquareReal(3)
	m1.Data[0][0] = 1.0
	m1.Data[0][1] = 0.0
	m1.Data[0][2] = 0.0
	m1.Data[1][0] = 0.0
	m1.Data[1][1] = 1.0
	m1.Data[1][2] = 0.0
	m1.Data[2][0] = 0.0
	m1.Data[2][1] = 0.0
	m1.Data[2][2] = 1.0
	fmt.Println(m1)

	m2 := NewBmMatrixSquareReal(2)
	m2.Data[0][0] = 1.0
	m2.Data[0][1] = 1.0
	m2.Data[1][0] = 1.0
	m2.Data[1][1] = 1.0
	fmt.Println(m2)

	tp := TensorProductReal(m1, m2)
	fmt.Println(tp)

	c1 := NewBmMatrixSquareComplex(3)
	c1.Data[0][0] = Complex32{1.0, 0.0}
	c1.Data[0][1] = Complex32{0.0, 0.0}
	c1.Data[0][2] = Complex32{0.0, 0.0}
	c1.Data[1][0] = Complex32{0.0, 0.0}
	c1.Data[1][1] = Complex32{1.0, 0.0}
	c1.Data[1][2] = Complex32{0.0, 0.0}
	c1.Data[2][0] = Complex32{0.0, 0.0}
	c1.Data[2][1] = Complex32{0.0, 0.0}
	c1.Data[2][2] = Complex32{1.0, 0.0}
	fmt.Println(c1)

	c2 := NewBmMatrixSquareComplex(2)
	c2.Data[0][0] = Complex32{1.0, 0.0}
	c2.Data[0][1] = Complex32{1.0, 0.0}
	c2.Data[1][0] = Complex32{1.0, 0.0}
	c2.Data[1][1] = Complex32{1.0, 0.0}
	fmt.Println(c2)

	tpc := TensorProductComplex(c1, c2)
	fmt.Println(tpc)

}
