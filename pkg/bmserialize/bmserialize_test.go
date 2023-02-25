package bmserialize

import (
	"os"
	"testing"
)

func TestStack(t *testing.T) {
	s := CreateBasicSerializer()
	s.TerminalDataSize = 32
	s.Depth = 8
	s.MemType = "FIFO"
	s.Senders = []string{"sender1", "sender2"}
	s.Receivers = []string{"receiver1"}

	s.Pushes = []Push{
		Push{"sender1", 200, "32'd1"},
	}

	s.Pops = []Pop{
		Pop{"receiver1", 60},
	}

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

	tb, err := os.Create("stack_tb.v")
	if err != nil {
		t.Error(err)
	}
	defer tb.Close()
	if r, err := s.WriteTestBench(); err != nil {
		t.Error(err)
	} else {
		tb.WriteString(r)
	}
}
