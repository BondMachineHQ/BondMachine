package brvga

import (
	"fmt"
	"log"
	"os"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestFile(t *testing.T) {

	// Create a new protobuf message
	msg := &Textmemupdate{
		Cpid: 1,
		Seq:  []*Textmemupdate_Byteseq{&Textmemupdate_Byteseq{Pos: 1, Payload: []byte("Hello")}},
	}

	// Marshal the message into a buffer
	out, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	// Create a file and write the buffer to it
	f, err := os.Create("/tmp/brvga")
	if err != nil {
		log.Fatal(err)
	}

	f.Write(out)
	f.Close()

	// Reopen the file and read the buffer from it
	f, err = os.Open("/tmp/brvga")
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new protobuf message and unmarshal the buffer into it
	recv := &Textmemupdate{}

	proto.Unmarshal(buf, recv)

	fmt.Println(n, recv)

}
