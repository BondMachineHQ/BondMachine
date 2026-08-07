package bmbuilder

import (
	"reflect"
	"testing"
)

func TestTokenizeLine(t *testing.T) {
	tokens, count := tokenizeLine("\tadd\t   r1, r2 ; comment")
	want := []string{"add", "r1,", "r2"}

	if count != len(want) {
		t.Fatalf("token count = %d, want %d", count, len(want))
	}

	if !reflect.DeepEqual(tokens, want) {
		t.Fatalf("tokens = %v, want %v", tokens, want)
	}
}

func TestTokenizeLineEmpty(t *testing.T) {
	tokens, count := tokenizeLine("   ; comment only")
	if count != 0 {
		t.Fatalf("token count = %d, want 0", count)
	}
	if len(tokens) != 0 {
		t.Fatalf("tokens = %v, want empty slice", tokens)
	}
}
