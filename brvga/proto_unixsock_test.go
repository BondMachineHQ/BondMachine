package brvga

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

func TestUnixSocket(t *testing.T) {

	go func() {
		time.Sleep(1 * time.Second)

		c, err := net.Dial("unix", "/tmp/brvga.sock")
		if err != nil {
			panic(err)
		}
		defer c.Close()

		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		recv := &Textmemupdate{}

		proto.Unmarshal(buf, recv)
		fmt.Println(n, recv)
	}()

	// Create a new protobuf message
	msg := &Textmemupdate{
		Cpid: 1,
		Seq:  []*Textmemupdate_Byteseq{&Textmemupdate_Byteseq{Pos: 1, Payload: []byte("Hello")}},
	}

	l, err := net.Listen("unix", "/tmp/brvga.sock")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	// Marshal the message into a buffer
	out, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	fd, err := l.Accept()
	if err != nil {
		log.Fatal("accept error:", err)
	}

	_, err = fd.Write(out)
	if err != nil {
		log.Fatal("Write: ", err)
	}

	time.Sleep(2 * time.Second)
}
