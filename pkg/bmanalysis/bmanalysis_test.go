package bmanalysis

import (
	"os"
	"testing"
)

func TestStack(t *testing.T) {
	s := CreateBasicStack()
	s.DataSize = 32
	s.Depth = 8
	s.MemType = "FIFO"
	s.Senders = []string{"sender1", "sender2"}
	s.Receivers = []string{"receiver1"}

	// s.SaveJSON("stack.json")
	// s.LoadJSON("stack.json")

	f, err := os.Create("stack.v")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	if r, err := s.WriteHDL(); err != nil {
		t.Error(err)
	} else {
		f.WriteString(r)
	}
}
