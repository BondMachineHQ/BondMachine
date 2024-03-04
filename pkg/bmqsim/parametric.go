package bmqsim

import (
	"fmt"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmatrix"
)

func (sim *BmQSimulator) Phase(line *bmline.BasmLine) (*bmmatrix.BmMatrixSquareComplex, error) {
	// The first parameter is the qbit, we can safely ignore it since it has been already treated
	// The second parameter is the phase
	if len(line.Elements) != 2 {
		return nil, fmt.Errorf("Phase: wrong number of parameters")
	}

	phase := line.Elements[1].GetValue()

	// Convert the phase to a float32
	if phaseFloat, err := strconv.ParseFloat(phase, 32); err != nil {
		return nil, fmt.Errorf("Phase: error parsing phase %s", phase)
	} else {
		return bmmatrix.GlobalPhase(2, float32(phaseFloat)), nil
	}
}

func (sim *BmQSimulator) P(line *bmline.BasmLine) (*bmmatrix.BmMatrixSquareComplex, error) {
	// The first parameter is the qbit, we can safely ignore it since it has been already treated
	// The second parameter is the phase
	if len(line.Elements) != 2 {
		return nil, fmt.Errorf("P: wrong number of parameters")
	}

	phase := line.Elements[1].GetValue()

	// Convert the phase to a float32
	if phaseFloat, err := strconv.ParseFloat(phase, 32); err != nil {
		return nil, fmt.Errorf("P: error parsing phase %s", phase)
	} else {
		return bmmatrix.PhaseShift(float32(phaseFloat)), nil
	}
}
