package main

import (
	"flag"

	"context"

	"github.com/BondMachineHQ/BondMachine/brvga"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var constraint = flag.String("constraints", "", "Use a BM constraint file to set the VGA up")
var sockFile = flag.String("sock", "/tmp/brvga.sock", "Socket file to use")
var cpId = flag.Int("cpid", 0, "CPU ID")
var pos = flag.Int("pos", 0, "Position in the CPU")
var payload = flag.String("payload", "", "Payload to send")

func init() {
	flag.Parse()
	if *constraint == "" {
		panic("Must specify a constraint string")
	}
}

func main() {
	vga, _ := brvga.NewBrvgaTextMemory(*constraint)

	ctx, _ := context.WithCancel(context.Background())

	// Create a new protobuf message
	msg := &brvga.Textmemupdate{
		Cpid: uint32(*cpId),
		Seq:  []*brvga.Textmemupdate_Byteseq{&brvga.Textmemupdate_Byteseq{Pos: uint32(*pos), Payload: []byte(*payload)}},
	}

	vga.UNIXSockSender(ctx, *sockFile, msg)
}
