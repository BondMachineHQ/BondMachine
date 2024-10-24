package bmmatrix

import (
	"testing"

	"github.com/mmirko/mel/pkg/mel"
)

func TestBasmExporter(t *testing.T) {

	a := new(M3BasmMatrix)
	var ep *mel.EvolutionParameters
	c := new(mel.MelConfig)
	c.Debug = false
	a.MelInit(c, ep)

	tests := []string{"m(T.json)", "m(T.json)"}

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
