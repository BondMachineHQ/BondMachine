package bmqsim

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmatrix"
)

type StateArray struct {
	Vector []bmmatrix.Complex32
}

type IOmap struct {
	Assoc map[string]string
}
type BmQSimulator struct {
	verbose  bool
	debug    bool
	qbits    []string
	qbitsNum map[string]int
	Mtx      []*bmmatrix.BmMatrixSquareComplex
	Inputs   []StateArray
	Outputs  []StateArray
}

// BmQSimulatorInit initializes the BmQSimulator
func (sim *BmQSimulator) BmQSimulatorInit() {
	sim.verbose = false
	sim.debug = false
	sim.qbits = make([]string, 0)
	sim.qbitsNum = make(map[string]int)
}

func (sim *BmQSimulator) QbitsNum() int {
	return len(sim.qbits)
}

func (sim *BmQSimulator) StateSize() int {
	return int(math.Pow(2, float64(len(sim.qbits))))
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

	result := make([]*bmmatrix.BmMatrixSquareComplex, 0)

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

	curOp := make([]*bmline.BasmLine, 0)
	curQBits := make(map[int]struct{})

	for i := 0; i < len(qasm.Lines); {
		line := qasm.Lines[i]
		op := line.Operation.GetValue()

		// Check if the operation is ready to form a matrix
		nextOp := false

		// Include the qbits in the operation to the currQbits map
		for _, arg := range line.Elements {
			argName := arg.GetValue()
			// Check if the argument is a qbit, otherwise ignore it
			if qbitN, ok := sim.qbitsNum[argName]; ok {
				if _, ok := curQBits[qbitN]; ok {
					nextOp = true
					break
				} else {
					curQBits[sim.qbitsNum[argName]] = struct{}{}
				}
			}

		}

		if op == "nextop" {
			nextOp = true
		}

		if i == len(qasm.Lines)-1 && !nextOp {
			if sim.debug {
				fmt.Println("\tProcessing line", i, ":", line.String())
			}
			curOp = append(curOp, line)
			i++
			nextOp = true
		}

		// If the operation is ready to form a matrix, create the matrix
		if nextOp {

			if sim.debug {
				fmt.Println(red("\tNew operation") + " (ready to form a matrix)")
			}

			// Create the matrix
			if len(curOp) > 0 {
				// Create the matrix
				if m, err := sim.BmMatrixFromOperation(curOp); err != nil {
					return nil, fmt.Errorf("error creating matrix from operation: %v", err)
				} else {
					if m != nil {
						result = append(result, m)
					}
				}
			}

			// Reset the operation and the qbits
			curOp = make([]*bmline.BasmLine, 0)
			curQBits = make(map[int]struct{})
		} else {
			if sim.debug {
				fmt.Println("\tProcessing line", i, ":", line.String())
			}
			curOp = append(curOp, line)
			i++
		}
	}
	return result, nil
}

type swap struct {
	s1 int
	s2 int
}

func (sim *BmQSimulator) BmMatrixFromOperation(op []*bmline.BasmLine) (*bmmatrix.BmMatrixSquareComplex, error) {
	// Let prepare the matrix that will be tensor-producted to form the final matrix
	var result *bmmatrix.BmMatrixSquareComplex

	swaps := make([]swap, 0)
	// loop over the qbits each qbit will have only one (or zero) operation within the op list
	// Every line where the qbits in not in sequence will be reordered swapping the qbits
	// and swapped back after the matrix is created

	localQBits := make([]string, len(sim.qbits))
	copy(localQBits, sim.qbits)

	if sim.debug {
		ops := ""
		for _, l := range op {
			ops += l.String() + ", "
		}
		fmt.Println("\t\tBuilding matrix from operation" + ops)
	}

	for q := 0; q < len(localQBits); q++ {

		if sim.debug {
			fmt.Println("\t\tAnalyzing qbit number:", q, "name:", localQBits[q])
			fmt.Println("\t\t\tCurrent QBits order:", localQBits)
			fmt.Println("\t\t\tCurrent QBits swaps:", swaps)
		}

		qbit := localQBits[q]
		// Find the operation for the qbit
		found := false
		fundLine := -1
		for i, line := range op {
			for _, arg := range line.Elements {
				argName := arg.GetValue()
				if argName == qbit {
					found = true
					fundLine = i
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			if sim.debug {
				fmt.Println("\t\t\tNo operation for qbit", q, "adding identity matrix")
			}
			// No operation for the qbit, lets add an identity matrix
			// Create the identity matrix
			ident := bmmatrix.IdentityComplex(2)
			if result == nil {
				result = ident
			} else {
				result = bmmatrix.TensorProductComplex(result, ident)
			}
			if sim.debug {
				fmt.Println("\t\t\tPartial Matrix:")
				fmt.Println(result)
			}
		} else {
			if sim.debug {
				fmt.Print("\t\t\tOperation for qbit ", localQBits[q], " found:"+op[fundLine].String())
			}
			argNumQBits := 0
			for _, arg := range op[fundLine].Elements {
				if _, ok := sim.qbitsNum[arg.GetValue()]; ok {
					argNumQBits++
				}
			}
			if sim.debug {
				fmt.Println(" the operation has", argNumQBits, "qbits")
			}

			if argNumQBits == 1 {
				// Single qbit operation
				// Create the matrix
				if matrix, err := sim.MatrixFromOp(op[fundLine]); err != nil {
					return nil, fmt.Errorf("error creating matrix from operation: %v", err)
				} else {
					if result == nil {
						result = matrix
					} else {
						result = bmmatrix.TensorProductComplex(result, matrix)
					}
					if sim.debug {
						fmt.Println("\t\t\tPartial Matrix:")
						fmt.Println(result)
					}
				}
			} else {
				// Multi qbit operation
				// The order of the qbits in the operation is important, we need to reorder the qbits in the operation
				// if they are not in sequence by swapping the qbits and then swapping them back after the matrix is created

				localOrder := make([]int, len(op[fundLine].Elements))
				for i, arg := range op[fundLine].Elements {
					argName := arg.GetValue()
					if _, ok := sim.qbitsNum[argName]; ok {
						localOrder[i] = sim.qbitsNum[argName]
					} else {
						// Leaving out the arguments that are not qbits
						localOrder[i] = -1
					}
				}

				for i, lq := range localOrder {
					if lq != q && lq != -1 {
						// Swap the qbits
						localQBits[q], localQBits[lq] = localQBits[lq], localQBits[q]
						// Add the swap to the list
						swaps = append(swaps, swap{q, lq})
						if sim.debug {
							fmt.Printf("Swapping qbits %d and %d\n", q, lq)
						}
						// Swap the localOrder if needed
						for j, lq2 := range localOrder {
							if lq2 == q {
								localOrder[j] = lq
							} else if lq2 == lq {
								localOrder[j] = q
							}
						}
						if sim.debug {
							fmt.Println("newLocalOrder:", localOrder)
						}

					}
					if i != len(localOrder)-1 {
						q++
					}
				}

				// Create the matrix
				if matrix, err := sim.MatrixFromOp(op[fundLine]); err != nil {
					return nil, fmt.Errorf("error creating matrix from operation: %v", err)
				} else {
					if result == nil {
						result = matrix
					} else {
						result = bmmatrix.TensorProductComplex(result, matrix)
					}
					if sim.debug {
						fmt.Println("\t\t\tPartial Matrix:")
						fmt.Println(result)
					}
				}
			}
		}
	}

	if len(swaps) == 0 {
		if sim.debug {
			fmt.Println("\t\tNo swaps needed")
		}
	} else {
		if sim.debug {
			fmt.Println("\t\tPerforming the swaps back")
		}

		for i := range len(swaps) {
			s := swaps[len(swaps)-1-i]
			baseSwaps := swaps2baseSwaps(s, len(sim.qbits))
			for _, bs := range baseSwaps {
				result = bmmatrix.SwapRowsColsComplex(result, bs.s1, bs.s2)
			}
			if sim.debug {
				fmt.Println("Swapping qbits", s.s1, "and", s.s2)
				fmt.Println("\t\t\tPartial Matrix:")
				fmt.Println(result)
			}
		}
	}
	return result, nil
}

func (sim *BmQSimulator) SwapQbits(m *bmmatrix.BmMatrixSquareComplex, s1, s2 int) *bmmatrix.BmMatrixSquareComplex {
	s := swap{s1, s2}
	baseSwaps := swaps2baseSwaps(s, len(sim.qbits))
	// fmt.Println("baseSwaps:", baseSwaps)
	for _, bs := range baseSwaps {
		m = bmmatrix.SwapRowsColsComplex(m, bs.s1, bs.s2)
	}
	return m
}

func swaps2baseSwaps(s swap, n int) []swap {
	baseNum := uint64(1 << n)
	// fmt.Println("baseNum:", baseNum)

	iDone := make(map[uint64]struct{})
	result := make([]swap, 0)

	s1 := s.s1
	s2 := s.s2

	max := uint64(1 << (n - 1))

	mask1 := uint64(max >> s1)
	mask2 := uint64(max >> s2)

	// fmt.Println("mask1:", int2bin(int(mask1), n))
	// fmt.Println("mask2:", int2bin(int(mask2), n))

	for i := uint64(0); i < baseNum; i++ {
		// fmt.Println(i, int2bin(int(i), n), i&mask1>>(n-1-s1), i&mask2>>(n-1-s2))
		if _, ok := iDone[i]; !ok {
			if i&mask1>>(n-1-s1) != i&mask2>>(n-1-s2) {
				num1 := i
				num2 := i ^ mask1 ^ mask2
				iDone[num1] = struct{}{}
				iDone[num2] = struct{}{}
				result = append(result, swap{int(num1), int(num2)})
			}
		}
	}

	return result
}

// func int2bin(i, n int) string {
// 	result := ""
// 	for a := n - 1; a >= 0; a-- {
// 		bit := i & (1 << uint(a))
// 		result += fmt.Sprintf("%d", bit>>uint(a))
// 	}
// 	return result
// }

func (sim *BmQSimulator) MatrixFromOp(line *bmline.BasmLine) (*bmmatrix.BmMatrixSquareComplex, error) {
	op := line.Operation.GetValue()
	op = strings.ToLower(op)
	switch op {
	case "h", "hadamard":
		return bmmatrix.Hadamard(), nil
	case "x", "paulix":
		return bmmatrix.PauliX(), nil
	case "y", "pauliy":
		return bmmatrix.PauliY(), nil
	case "z", "pauliz":
		return bmmatrix.PauliZ(), nil
	case "cx", "cnot", "xor":
		return bmmatrix.CNot(), nil
	case "s", "p":
		return bmmatrix.S(), nil
	case "v", "sx":
		return bmmatrix.Sx(), nil
	case "t":
		return bmmatrix.T(), nil
	case "xnor":
		return bmmatrix.XNor(), nil
	case "cz", "cphase", "csign", "cpf":
		return bmmatrix.Cphase(), nil
	case "dcnot":
		return bmmatrix.Dcnot(), nil
	case "swap":
		return bmmatrix.Swap(), nil
	case "iswap":
		return bmmatrix.Iswap(), nil
	case "phase", "ph":
		return sim.Phase(line)
	case "r":
		return sim.P(line)
	case "zero":
	// Ignore the zero operation
	default:
		return nil, fmt.Errorf("unknown operation %s", op)
	}
	return nil, nil
}

func (sim *BmQSimulator) EmitBMAPIMaps(hwflavor string) (string, error) {
	if sim != nil && sim.qbits != nil {
		var ioNum int
		switch hwflavor {
		case "seq_hardcoded_real":
			ioNum = int(math.Pow(float64(2), float64(len(sim.qbits))))
		case "seq_hardcoded_complex", "seq_hardcoded_addtree_complex":
			ioNum = int(math.Pow(float64(2), float64(len(sim.qbits)))) * 2
		}
		newMap := new(IOmap)
		newMap.Assoc = make(map[string]string)
		for i := 0; i < ioNum; i++ {
			newMap.Assoc[fmt.Sprintf("i%d", i)] = fmt.Sprintf("%d", i)
			newMap.Assoc[fmt.Sprintf("o%d", i)] = fmt.Sprintf("%d", i)
		}
		if jData, err := json.Marshal(newMap); err != nil {
			return "", fmt.Errorf("error marshalling the map: %v", err)
		} else {
			return string(jData), nil
		}
	}
	return "", fmt.Errorf("no qbits defined")
}

func (sim *BmQSimulator) RunSoftwareSimulation() error {
	// TODO Include the checks
	sim.Outputs = make([]StateArray, len(sim.Inputs))
	for i, input := range sim.Inputs {
		if len(input.Vector) != sim.StateSize() {
			return fmt.Errorf("input %d has %d vectors, expected %d", i, len(input.Vector), sim.StateSize())
		}

		curState := make([]bmmatrix.Complex32, sim.StateSize())
		copy(curState, input.Vector)

		for j, mtx := range sim.Mtx {
			if mtx == nil {
				return fmt.Errorf("matrix %d is nil", j)
			}
			if len(mtx.Data) != sim.StateSize() {
				return fmt.Errorf("matrix %d has %d vectors, expected %d", j, len(mtx.Data), sim.StateSize())
			}

			if newState, err := bmmatrix.MatrixVectorProductComplex(mtx, curState); err != nil {
				return fmt.Errorf("error running the software simulation: %v", err)
			} else {
				curState = newState
			}
		}

		sim.Outputs[i] = StateArray{Vector: curState}

	}
	return nil
}
