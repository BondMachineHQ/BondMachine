package bmmatrix

import (
	"testing"

	"github.com/mmirko/mel/pkg/mel"
)

func TestMatrixInfo(t *testing.T) {

	a := new(M3MatrixInfo)
	var ep *mel.EvolutionParameters
	c := new(mel.MelConfig)
	c.Debug = false
	a.MelInit(c, ep)

	tests := []string{"m(ref:3:4)", "m(ref:3:4)"}
	tests = append(tests, "m(in:test:4:7)", "m(ref:4:7)")
	tests = append(tests, "m(T.json)", "m(ref:3:2)")
	tests = append(tests, "m(rowmajor:T.json)", "m(ref:3:2)")
	tests = append(tests, "m(colmajor:T.json)", "m(ref:2:3)")
	tests = append(tests, "mult(m(ref:3:2),m(ref:2:1))", "m(ref:3:1)")

	for i, iString := range tests {

		if i%2 == 1 {
			continue
		}

		a.MelStringImport(iString)
		a.Compute()
		if a.Inspect() != tests[i+1] {
			t.Errorf("Expected %s, got %s", tests[i+1], a.Inspect())
		}

	}
}
