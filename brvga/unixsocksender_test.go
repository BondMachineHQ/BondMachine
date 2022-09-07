package brvga

import (
	"testing"
	"time"

	"context"
)

func TestUnixSocketSender(t *testing.T) {
	vga, _ := NewBrvgaTextMemory("1:0:0:80:25")

	ctx, _ := context.WithCancel(context.Background())

	// Create a new protobuf message
	msg := &Textmemupdate{
		Cpid: 1,
		Seq:  []*Textmemupdate_Byteseq{&Textmemupdate_Byteseq{Pos: 1, Payload: []byte("Hello")}},
	}

	vga.UNIXSockSender(ctx, "/tmp/brvga.sock", msg)

	time.Sleep(10 * time.Second)

}
