package bondgo

import (
	"fmt"
	"testing"
)

func TestActions(t *testing.T) {
	if test1, err := Type_from_string("uint8"); err == nil {
		fmt.Println(test1)
		if test2, err := Type_from_string(" *   bool"); err == nil {
			fmt.Println(test2)
			if test3, err := Type_from_string(" uint8"); err == nil {
				fmt.Println(test3)
				fmt.Println("test1 test1", Same_Type(test1, test1))
				fmt.Println("test1 test2", Same_Type(test1, test2))
				fmt.Println("test1 test3", Same_Type(test1, test3))
				fmt.Println("test2 test1", Same_Type(test2, test1))
				fmt.Println("test2 test2", Same_Type(test2, test2))
				fmt.Println("test2 test3", Same_Type(test2, test3))
				fmt.Println("test3 test1", Same_Type(test3, test1))
				fmt.Println("test3 test2", Same_Type(test3, test2))
				fmt.Println("test3 test3", Same_Type(test3, test3))
			}
		}
	}
}

func TestMemused(t *testing.T) {
	ttype, _ := Type_from_string("uint8")
	a := VarCell{ttype, REGISTER, 0, 0, 0, 0, 0, 0}
	b := VarCell{ttype, REGISTER, 0, 0, 0, 0, 0, 0}
	fmt.Println(&a == &b)
}
