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

	f, err := os.Create("serialize.v")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	if r, err := s.WriteHDL(); err != nil {
		t.Error(err)
	} else {
		f.WriteString(r)
	}

	// tb, err := os.Create("serialize_tb.v")
	// if err != nil {
	// 	t.Error(err)
	// }
	// defer tb.Close()
	// if r, err := s.WriteTestBench(); err != nil {
	// 	t.Error(err)
	// } else {
	// 	tb.WriteString(r)
	// }
}
