package bondmachine

import (
	"context"

	"github.com/BondMachineHQ/BondMachine/pkg/brvga"
)

type Vga800x600Emu struct {
	Socket      string
	Constraints string
	vga         *brvga.BrvgaTextMemory
	context     context.Context
}

func (ed *Vga800x600Emu) Init() error {
	// TODO Remove this
	ed.Constraints = "vtextmem:0:3:3:16:16:1:25:25:16:16"
	ed.Socket = "/tmp/brvga.sock"
	vga, _ := brvga.NewBrvgaTextMemory(ed.Constraints)
	ctx, _ := context.WithCancel(context.Background())

	ed.vga = vga
	ed.context = ctx
	return nil
}

func (ed *Vga800x600Emu) Run() error {
	return nil
}

func (ed *Vga800x600Emu) PushCommand(cmd []byte) ([]byte, error) {
	// fmt.Println("PushCommand", cmd)

	cpId := uint32(cmd[0])
	pos := uint32(cmd[1])
	payload := []byte{cmd[2]}

	// fmt.Println("cpId", cpId)
	// fmt.Println("pos", pos)
	// fmt.Println("payload", payload)

	// Create a new protobuf message
	msg := &brvga.Textmemupdate{
		Cpid: uint32(cpId),
		Seq:  []*brvga.Textmemupdate_Byteseq{&brvga.Textmemupdate_Byteseq{Pos: pos, Payload: payload}},
	}

	ed.vga.UNIXSockSender(ed.context, ed.Socket, msg)

	return nil, nil
}
