package bmmatrix

import (
	"testing"

	"github.com/mmirko/mel/pkg/mel"
)

func TestM3numberEvaluator(t *testing.T) {

	a := new(M3numberMe3li)
	var ep *mel.EvolutionParameters
	c := new(mel.MelConfig)
	c.Debug = false
	a.MelInit(c, ep)

	tests := []string{"m3numberconst(45.2)", "m3numberconst(45.2)"}
	tests = append(tests, "add(m3numberconst(1E+1),m3numberconst(2))", "m3numberconst(1.2E+01)")
	tests = append(tests, "mult(m3numberconst(3.2),m3numberconst(5))", "m3numberconst(1.6E+01)")

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
