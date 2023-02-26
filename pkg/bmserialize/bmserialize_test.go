package bmserialize

import (
	"os"
	"testing"
)

func TestStack(t *testing.T) {
	s := CreateBasicSerializer()
	s.TerminalDataSize = 32
	s.SerialDataSize = 32
	s.Terminals = 2
	s.Direction = "serialize"
	s.Depth = 8

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
