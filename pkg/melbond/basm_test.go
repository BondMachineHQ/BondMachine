package melbond

import (
	"testing"

	"github.com/mmirko/mel/pkg/m3number"
	"github.com/mmirko/mel/pkg/mel"
)

func TestBasmEvaluator(t *testing.T) {

	bc := new(MelBondConfig)
	a := new(m3number.M3numberMe3li)
	var ep *mel.EvolutionParameters
	c := new(mel.MelConfig)
	a.Mel3Object.DefaultCreator = bc.BasmCreator
	c.Debug = false
	a.MelInit(c, ep)

	tests := []string{`
	add(
		add(
			m3numberconst(1),
			m3numberconst(66)
		),
		add(
			m3numberconst(1),
			mult(
				m3numberconst(2),
				m3numberconst(4)
			)
		)
	)`,
	}

	for _, iString := range tests {

		if err := a.MelStringImport(iString); err != nil {
			t.Errorf("Error importing: %s", err)
		} else {
			a.Compute()
		}
	}
}
