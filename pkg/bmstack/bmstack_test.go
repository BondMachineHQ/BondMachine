package bmstack

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	s := CreateBasicStack()
	fmt.Println(s.WriteHDL())
}
