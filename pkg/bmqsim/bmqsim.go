package bmqsim

import (
	"fmt"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmatrix"
)

type BmQSimulator struct {
	verbose  bool
	debug    bool
	qbits    []string
	qbitsNum map[string]int
}

// BmQSimulatorInit initializes the BmQSimulator
func (sim *BmQSimulator) BmQSimulatorInit() {
	sim.verbose = false
	sim.debug = false
	sim.qbits = make([]string, 0)
	sim.qbitsNum = make(map[string]int)
}

func (sim *BmQSimulator) Dump() string {
	return fmt.Sprintf("BmQSimulator: verbose=%t, debug=%t, qbits=%v, qbitsNum=%v", sim.verbose, sim.debug, sim.qbits, sim.qbitsNum)
}

func (sim *BmQSimulator) SetVerbose() {
	sim.verbose = true
}

func (sim *BmQSimulator) SetDebug() {
	sim.debug = true
}

// QasmToBmMatrices converts a QASM file to a list of BmMatrixSquareComplex, the input is a BasmBody with all the metadata and the list of quantum instructions
func (sim *BmQSimulator) QasmToBmMatrices(qasm *bmline.BasmBody) ([]*bmmatrix.BmMatrixSquareComplex, error) {

	// Get the qbits and their names
	qbits := qasm.GetMeta("qbits")

	if qbits == "" {
		return nil, fmt.Errorf("no qbits defined")
	}

	for i, qbit := range strings.Split(qbits, ":") {
		if _, ok := sim.qbitsNum[qbit]; ok {
			return nil, fmt.Errorf("qbit %s already defined", qbit)
		} else {
			sim.qbits = append(sim.qbits, qbit)
			sim.qbitsNum[qbit] = i
		}
	}

	currOp := make([]*bmline.BasmLine, 0)
	currQbits := make(map[int]struct{})

	for i, line := range qasm.Lines {
		op := line.Operation.GetValue()

		if sim.debug {
			fmt.Printf("Processing line %d: %s\n", i, line.String())
		}

		// Check if the operation is ready to form a matrix
		nextOp := false

		// Include the qbits in the operation to the currQbits map
		for _, arg := range line.Elements {
			argName := arg.GetValue()
			// Check if the argument is a qbit, otherwise ignore it
			if qbitN, ok := sim.qbitsNum[argName]; ok {
				if _, ok := currQbits[qbitN]; ok {
					nextOp = true
					break
				} else {
					currQbits[sim.qbitsNum[argName]] = struct{}{}
				}
			}

		}

		if op == "nextop" {
			nextOp = true
		}

		singleLast := false

		if i == len(qasm.Lines)-1 {
			// If the last line is not already a nextOp lets put it in current operation, otherwise we will set singleLast to true
			// and process it alone later on
			if !nextOp {
				currOp = append(currOp, line)
			} else {
				singleLast = true
			}
			nextOp = true
		}

		// If the operation is ready to form a matrix, create the matrix
		if nextOp {
			// Create the matrix
			if len(currOp) > 0 {
				// Create the matrix
				_, err := sim.BmMatrixFromOperation(currOp)
				if err != nil {
					return nil, fmt.Errorf("error creating matrix from operation: %v", err)
				}
			}

			// Reset the operation and the qbits
			currOp = make([]*bmline.BasmLine, 0)
			currQbits = make(map[int]struct{})
		}

		// If the last line is not already been added to the last matrix, lets put it alone in a new matrix
		// Otherwise (not the last or already done) we will put it in the current operations list and set the involved qbits into the currQbits map
		if singleLast {
			currOp = append(currOp, line)
			// Create the matrix
			_, err := sim.BmMatrixFromOperation(currOp)
			if err != nil {
				return nil, fmt.Errorf("error creating matrix from operation: %v", err)
			}
		} else {
			currOp = append(currOp, line)

			// Include the qbits in the operation to the currQbits map
			for _, arg := range line.Elements {
				argName := arg.GetValue()
				// Check if the argument is a qbit, otherwise ignore it
				if qbitN, ok := sim.qbitsNum[argName]; ok {
					currQbits[qbitN] = struct{}{}
				}

			}

		}
	}
	return nil, nil
}

func (sim *BmQSimulator) BmMatrixFromOperation(op []*bmline.BasmLine) (*bmmatrix.BmMatrixSquareComplex, error) {
	// Let prepare the matrix sequence that will be tensor-producted to form the final matrix
	// mSeq := make([]bmmatrix.BmMatrixSquareComplex, 0)

	// loop over the qbits each qbit will have only one (or zero) operation within the op list
	// Every line where the qbits in not in sequence will be reordered swapping the qbits
	// and swapped back after the matrix is created
	// for _, qbit := range sim.qbits {
	// }
	// TODO: continue here
	return nil, nil
}
