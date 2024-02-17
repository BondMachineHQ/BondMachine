package bmqsim

import (
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmatrix"
)

// QasmToBmMatrices converts a QASM file to a list of BmMatrixSquareComplex, the input is a BasmBody with all the metadata and the list of quantum instructions
func QasmToBmMatrices(qasm *bmline.BasmBody) ([]*bmmatrix.BmMatrixSquareComplex, error) {
	return nil, nil
}
