package brvga

import (
	"testing"
	"time"

	"context"
)

func TestUnixSocketSender(t *testing.T) {
	vga, _ := NewBrvgaTextMemory("vtextmem:0:0:0:16:16:1:25:25:16:16")

	ctx, _ := context.WithCancel(context.Background())

	// Create a new protobuf message
	msg := &Textmemupdate{
		Cpid: 1,
		Seq:  []*Textmemupdate_Byteseq{&Textmemupdate_Byteseq{Pos: 0, Payload: []byte{0, 0, 0, 0, 0}}},
	}

	vga.UNIXSockSender(ctx, "/tmp/brvga.sock", msg)

	time.Sleep(10 * time.Second)

}
