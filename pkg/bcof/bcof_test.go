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

	bcofData := &BCOFData{
		Id:        1,
		Rsize:     8,
		Signature: "data",
		Payload:   []byte("hello world"),
	}

	list := make([]*BCOFEntrySubentry, 2)
	list[0] = new(BCOFEntrySubentry)
	bin1 := new(BCOFEntrySubentry_Binary)
	bin1.Binary = bcofData
	list[0].Pl = bin1

	list[1] = new(BCOFEntrySubentry)
	bin2 := new(BCOFEntrySubentry_Binary)
	bin2.Binary = bcofData
	list[1].Pl = bin2

	bcof := &BCOFEntry{
		Id:        1,
		Rsize:     8,
		Signature: "sub",
		Data:      list,
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
