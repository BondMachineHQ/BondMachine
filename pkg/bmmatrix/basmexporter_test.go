package bmmatrix

import (
	"testing"

	"github.com/mmirko/mel/pkg/mel"
)

func TestBasmExporter(t *testing.T) {

	a := new(M3BasmMatrix)
	var env interface{}
	env = newExporterEnv()
	a.Mel3Object.Environment = &env

	var ep *mel.EvolutionParameters
	c := new(mel.MelConfig)
	c.Debug = true
	a.MelInit(c, ep)

	tests := []string{"m(T.json)", "m(ref:3:2)"}
	tests = append(tests, "mult(m(T.json),m(T.json))", "m(ref:3:2)")

	for i, iString := range tests {

		if i%2 == 1 {
			continue
		}

		a.MelStringImport(iString)
		a.Compute()
		if a.Inspect() != tests[i+1] {
			// a.MelDump(nil)
			t.Errorf("Expected %s, got %s", tests[i+1], a.Inspect())
		}

	}
}
