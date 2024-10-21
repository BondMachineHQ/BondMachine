package bmmatrix

import (
	"errors"
	"fmt"

	mel3program "github.com/mmirko/mel/pkg/mel3program"
)

type MatrixInfo struct {
	*mel3program.Mel3Object
	error
	Result *mel3program.Mel3Program
}

func MatrixInfoCreator() mel3program.Mel3Visitor {
	return new(MatrixInfo)
}

func (ev *MatrixInfo) GetName() string {
	return "m3bmmatrix"
}

func (ev *MatrixInfo) GetMel3Object() *mel3program.Mel3Object {
	return ev.Mel3Object
}

func (ev *MatrixInfo) SetMel3Object(mel3o *mel3program.Mel3Object) {
	ev.Mel3Object = mel3o
}

func (ev *MatrixInfo) GetError() error {
	return ev.error
}

func (ev *MatrixInfo) GetResult() *mel3program.Mel3Program {
	return ev.Result
}

func (ev *MatrixInfo) Visit(in_prog *mel3program.Mel3Program) mel3program.Mel3Visitor {

	debug := ev.Config.Debug

	if debug {
		fmt.Println("m3bmmatrix enter: ", in_prog)
		defer fmt.Println("m3bmmatrix exit")
	}

	checkEv := mel3program.ProgMux(ev, in_prog)

	if ev.GetName() != checkEv.GetName() {
		return checkEv.Visit(in_prog)
	}

	obj := ev.GetMel3Object()
	implementations := obj.Implementation

	programId := in_prog.ProgramID
	libraryId := in_prog.LibraryID

	implementation := implementations[libraryId]

	isFunctional := true

	if len(implementation.NonVariadicArgs[programId]) == 0 && !implementation.IsVariadic[programId] {
		isFunctional = false
	}

	if isFunctional {
		arg_num := len(in_prog.NextPrograms)
		evaluators := make([]mel3program.Mel3Visitor, arg_num)
		for i, prog := range in_prog.NextPrograms {
			evaluators[i] = mel3program.ProgMux(ev, prog)
			evaluators[i].Visit(prog)
		}

		switch in_prog.LibraryID {
		case MYLIBID:
			switch in_prog.ProgramID {
			case MATRIXMULT:
				if arg_num == 2 {
					// res0 := evaluators[0].GetResult()
					// res1 := evaluators[1].GetResult()
					// value0 := ""
					// if res0 != nil && res0.LibraryID == libraryId && res0.ProgramID == M3NUMBERCONST {
					// 	value0 = res0.ProgramValue
					// } else {
					// 	ev.error = errors.New("wrong argument type")
					// 	return nil
					// }

					// value1 := ""
					// if res1 != nil && res1.LibraryID == libraryId && res1.ProgramID == M3NUMBERCONST {
					// 	value1 = res1.ProgramValue
					// } else {
					// 	ev.error = errors.New("wrong argument type")
					// 	return nil
					// }

					opResult := ""

					// if value0n64, err := strconv.ParseFloat(value0, 32); err == nil {
					// 	if value1n64, err := strconv.ParseFloat(value1, 32); err == nil {
					// 		value0n := float32(value0n64)
					// 		value1n := float32(value1n64)

					// 		var opResultN float32

					// 		switch in_prog.ProgramID {
					// 		case ADD:
					// 			opResultN = value0n + value1n
					// 		case SUB:
					// 			opResultN = value0n - value1n
					// 		case MULT:
					// 			opResultN = value0n * value1n
					// 		case DIV:
					// 			opResultN = value0n / value1n
					// 		}

					// 		opResult = strconv.FormatFloat(float64(opResultN), 'E', -1, 32)

					// 	} else {
					// 		ev.error = errors.New("convert to number failed")
					// 		return nil
					// 	}
					// } else {
					// 	ev.error = errors.New("convert to number failed")
					// 	return nil
					// }

					result := new(mel3program.Mel3Program)
					result.LibraryID = libraryId
					result.ProgramID = MATRIXCONST
					result.ProgramValue = opResult
					result.NextPrograms = nil
					ev.Result = result
					return nil
				} else {
					ev.error = errors.New("wrong argument number")
					return nil
				}
			}
		default:
			ev.error = errors.New("unknown LibraryID")
			return nil
		}
	} else {

		switch in_prog.LibraryID {
		case MYLIBID:
			switch in_prog.ProgramID {
			case MATRIXCONST:
				switch in_prog.ProgramValue {
				default:
					result := new(mel3program.Mel3Program)
					result.LibraryID = libraryId
					result.ProgramID = programId
					result.ProgramValue = in_prog.ProgramValue
					result.NextPrograms = nil
					ev.Result = result
					return nil
				}
			}
		default:
			ev.error = errors.New("unknown LibraryID")
			return nil
		}
	}

	return ev
}

func (ev *MatrixInfo) Inspect() string {
	obj := ev.GetMel3Object()
	implementations := obj.Implementation
	if ev.error == nil {
		if dump, err := mel3program.ProgDump(implementations, ev.Result); err == nil {
			return "Evaluation ok: " + dump
		} else {
			return "Result export failed:" + fmt.Sprint(err)
		}
	} else {
		return fmt.Sprint(ev.error)
	}
}
