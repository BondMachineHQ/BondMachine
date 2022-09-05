package brvga

// Test unit to read from a UNIX socket and decode the data
// into a protobuf message.

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

// TestReceiver tests the receiver
func TestReceiver(t *testing.T) {
	// Create a UNIX socket
	socket, err := net.Listen("unix", "/tmp/brvga.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	// Create a goroutine to accept the connection
	go func() {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		time.Sleep(1 * time.Second) // Wait for the client to connect

		// Create a buffer to read the data into
		buf := make([]byte, 1024)

		// Read the data
		_, err = conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new protobuf message
		msg := &Textmemupdate{}

		// Read the data into the protobuf message
		err = proto.Unmarshal(buf, msg)
		if err != nil {
			log.Fatal(err)
		}

		// Print the message
		fmt.Println(msg)
	}()

	// Create a new protobuf message
	msg := &Textmemupdate{
		Cpid: 1,
		Seq:  []*Textmemupdate_Byteseq{&Textmemupdate_Byteseq{Pos: 1, Payload: []byte("Hello")}},
	}

	// Write the data
	out, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	// Create a connection to the UNIX socket
	conn, err := net.Dial("unix", "/tmp/brvga.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Write the data
	_, err = conn.Write(out)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

}
