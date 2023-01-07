package bmstack

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	s := CreateBasicStack()
	s.DataSize = 32
	s.Senders = []string{"sender1", "sender2"}
	s.Receivers = []string{"receiver1", "receiver2"}

	fmt.Println(s.WriteHDL())
}
