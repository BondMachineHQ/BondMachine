package bcof

// A test to see if we can use the generated code to read and write a BCOF file.

import (
	"fmt"
	"os"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestBCOF(t *testing.T) {
	// Create a BCOF file.
	bcof := &BCOFEntry{
		Id:        1,
		Signature: "sub",
		Payload: &BCOFEntry_Data{
			Data: []byte("hello world"),
		},
	}

	fmt.Println(bcof)

	bcofBytes, err := proto.Marshal(bcof)
	if err != nil {
		t.Fatalf("failed to marshal BCOF: %v", err)
	}
	if err := os.WriteFile("test.bcof", bcofBytes, 0644); err != nil {
		t.Fatalf("failed to write BCOF file: %v", err)
	}
	defer os.Remove("test.bcof")

	// Read the BCOF file.
	bcofBytes, err = os.ReadFile("test.bcof")
	if err != nil {
		t.Fatalf("failed to read BCOF file: %v", err)
	}
	if err := proto.Unmarshal(bcofBytes, bcof); err != nil {
		t.Fatalf("failed to unmarshal BCOF: %v", err)
	}

	fmt.Println(bcof)
}
