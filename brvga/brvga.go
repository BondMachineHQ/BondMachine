package brvga

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
)

type cpTextMemory struct {
	CpId    int
	Width   int
	Height  int
	TopPos  int
	LeftPos int
	mem     []byte
}

type BrvgaTextMemory struct {
	constraintString string
	Cps              []cpTextMemory
}

func NewBrvgaTextMemory(constraint string) (*BrvgaTextMemory, error) {
	components := strings.Split(constraint, ":")
	componentsN := len(components)
	if (componentsN-1)%5 == 0 {
		boxes := make([]cpTextMemory, 0)

		for i := 1; i < componentsN; i = i + 5 {
			newBox := cpTextMemory{}
			if newCP, err := strconv.Atoi(components[i]); err == nil {
				newBox.CpId = newCP
			} else {
				return nil, err
			}
			if newLeft, err := strconv.Atoi(components[i+1]); err == nil {
				newBox.LeftPos = newLeft
			} else {
				return nil, err
			}
			if newTop, err := strconv.Atoi(components[i+2]); err == nil {
				newBox.TopPos = newTop
			} else {
				return nil, err
			}
			if newWidth, err := strconv.Atoi(components[i+3]); err == nil {
				newBox.Width = newWidth
			} else {
				return nil, err
			}
			if newHeight, err := strconv.Atoi(components[i+4]); err == nil {
				newBox.Height = newHeight
			} else {
				return nil, err
			}

			memL := newBox.Width * newBox.Height
			newBox.mem = make([]byte, memL)
			for j := 0; j < memL; j++ {
				newBox.mem[j] = 0x00
			}

			boxes = append(boxes, newBox)
		}

		result := new(BrvgaTextMemory)
		result.constraintString = constraint
		result.Cps = boxes
		return result, nil
	}

	return nil, errors.New("invalid constraint string")
}

func (b *BrvgaTextMemory) Dump() string {
	result := ""
	result += fmt.Sprintf("Constraint %s\n\n", b.constraintString)
	for _, cp := range b.Cps {
		result += fmt.Sprintf("cp %d: %d x %d at %d, %d\n", cp.CpId, cp.Width, cp.Height, cp.LeftPos, cp.TopPos)
		for i := 0; i < cp.Height*cp.Width; i++ {
			if i%cp.Width == 0 {
				result += "\n"
			}
			result += fmt.Sprintf("%02x ", cp.mem[i])
		}
		result += "\n\n"
	}
	return result
}

func (b *BrvgaTextMemory) UNIXSockReceiver(ctx context.Context, path string) {

	// Create the socket
	sock, err := net.Listen("unix", path)
	if err != nil {
		panic(err)
	}

	// Wait for a connection
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := sock.Accept()
			if err != nil {
				panic(err)
			}
			go b.handleConnection(ctx, conn)
		}
	}
}

func (b *BrvgaTextMemory) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			b.handleMessage(ctx, buf[:n])
		}
	}
}

func (b *BrvgaTextMemory) handleMessage(ctx context.Context, buf []byte) {
	recv := &Textmemupdate{}

	proto.Unmarshal(buf, recv)
	fmt.Println(recv)
}

func (b *BrvgaTextMemory) UNIXSockSender(ctx context.Context, path string, buf *Textmemupdate) error {

	// Marshal the message into a buffer
	out, err := proto.Marshal(buf)
	if err != nil {
		return err
	}

	// Connect to the socket
	conn, err := net.Dial("unix", path)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send the message
	_, err = conn.Write(out)
	return err
}
