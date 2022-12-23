package procbuilder

import (
	"fmt"
	"testing"
)

func TestProcessNumber(t *testing.T) {
	fmt.Println(Process_number("56"))
	fmt.Println(Process_number("0x901"))
	fmt.Println(Process_number("0b10101"))
	fmt.Println(Process_number("0d56"))
	fmt.Println(Process_number("0f56"))
	fmt.Println(Process_number("0finfinity"))
	fmt.Println(Process_number("0fNaN"))
	fmt.Println(Process_number("0f4e-4"))
}
