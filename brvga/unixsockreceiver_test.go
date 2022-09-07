package brvga

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"context"
)

func TestUnixSocketReceiver(t *testing.T) {
	vga, err := NewBrvgaTextMemory("textvga:1:1:1:16:16")

	if err != nil {
		t.Fatal(err)
	}

	ctx, _ := context.WithCancel(context.Background())

	go vga.UNIXSockReceiver(ctx, "/tmp/brvga.sock")

	for {
		time.Sleep(1 * time.Second)
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		fmt.Println(vga.Dump())
	}

}
