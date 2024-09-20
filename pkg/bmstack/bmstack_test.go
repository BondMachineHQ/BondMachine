package bmstack

import (
	"os"
	"testing"
)

func TestStack(t *testing.T) {
	s := CreateBasicStack()
	s.DataSize = 32
	s.Depth = 8
	s.MemType = "FIFO"
	s.Senders = []string{"sender1", "sender2", "sender3"}
	s.Receivers = []string{"receiver1", "receiver2"}

	s.Pushes = []Push{
		Push{"sender1", 100, "32'd1"},
		Push{"sender2", 150, "32'd2"},
		Push{"sender3", 200, "32'd3"},
		Push{"sender1", 250, "32'd4"},
		Push{"sender2", 300, "32'd5"},
		Push{"sender3", 350, "32'd6"},
	}

	s.Pops = []Pop{
		Pop{"receiver1", 110},
		Pop{"receiver2", 160},
		Pop{"receiver1", 210},
		Pop{"receiver2", 260},
		Pop{"receiver1", 310},
		Pop{"receiver2", 360},
	}

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
