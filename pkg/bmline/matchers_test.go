package bmline

import (
	"testing"
)

func TestMatchers(t *testing.T) {
	// TODO Fare il vero
	tests := []string{"prova::*--type=reg",
		"prova2::*--type=reg",
		"prova3--exec=async::*--type=reg"}

	for _, test := range tests {
		if bline, err := Text2BasmLine(test); err != nil {
			t.Fail()
		} else {
			if dcheck, err := BasmLine2Text(bline); err != nil {
				t.Fail()
			} else {
				if dcheck != test {
					t.Fail()
				}
			}
		}
	}
}
