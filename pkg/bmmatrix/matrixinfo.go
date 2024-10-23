package bmmatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"

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
					res0 := evaluators[0].GetResult()
					res1 := evaluators[1].GetResult()
					value0 := ""
					if res0 != nil && res0.LibraryID == libraryId && res0.ProgramID == MATRIXCONST {
						value0 = res0.ProgramValue
					} else {
						ev.error = errors.New("wrong argument type")
						return nil
					}

					value1 := ""
					if res1 != nil && res1.LibraryID == libraryId && res1.ProgramID == MATRIXCONST {
						value1 = res1.ProgramValue
					} else {
						ev.error = errors.New("wrong argument type")
						return nil
					}

					opResult := ""

					values0 := regexp.MustCompile(`^ref:([0-9]+):([0-9]+)$`)
					values1 := regexp.MustCompile(`^ref:([0-9]+):([0-9]+)$`)
					if values0.MatchString(value0) && values1.MatchString(value1) {
						row0 := values0.FindStringSubmatch(value0)[1]
						col0 := values0.FindStringSubmatch(value0)[2]
						row1 := values1.FindStringSubmatch(value1)[1]
						col1 := values1.FindStringSubmatch(value1)[2]
						if col0 == row1 {
							opResult = fmt.Sprintf("ref:%s:%s", row0, col1)
						} else {
							ev.error = errors.New("wrong argument, matrix dimensions do not match")
							return nil
						}
					} else {
						ev.error = errors.New("wrong argument, matrix dimensions do not match")
						return nil
					}

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
				m := in_prog.ProgramValue
				ref := regexp.MustCompile(`^ref:([0-9]+):([0-9]+)$`)
				// Match an alredy ref matrix
				if ref.MatchString(m) {
					result := new(mel3program.Mel3Program)
					result.LibraryID = libraryId
					result.ProgramID = MATRIXCONST
					result.ProgramValue = m
					result.NextPrograms = nil
					ev.Result = result
					return nil
				}

				in := regexp.MustCompile(`^in:([a-zA-Z0-9]+):([0-9]+):([0-9]+)$`)
				// Match an inpput matrix
				if in.MatchString(m) {
					result := new(mel3program.Mel3Program)
					result.LibraryID = libraryId
					result.ProgramID = MATRIXCONST
					// Replace in with ref
					result.ProgramValue = "ref:" + in.FindStringSubmatch(m)[2] + ":" + in.FindStringSubmatch(m)[3]
					result.NextPrograms = nil
					ev.Result = result
					return nil
				}

				// Match a matrix json file (a file m exists in the filesystem)
				fileName := m
				rowM := true
				rowMajor := regexp.MustCompile(`^rowmajor:(.+)$`)
				if rowMajor.MatchString(m) {
					fileName = rowMajor.FindStringSubmatch(m)[1]
				}
				colMajor := regexp.MustCompile(`^colmajor:(.+)$`)
				if colMajor.MatchString(m) {
					fileName = colMajor.FindStringSubmatch(m)[1]
					rowM = false
				}

				if _, err := os.Stat(fileName); err == nil {
					// Load the file
					file, err := os.ReadFile(fileName)
					if err == nil {
						// Unmarshal the file
						var matrix [][]float64
						err = json.Unmarshal(file, &matrix)
						if err == nil {
							major := len(matrix)
							minor := len(matrix[0])
							for i := 1; i < major; i++ {
								if len(matrix[i]) != minor {
									ev.error = errors.New("matrix rows have different length")
									return nil
								}
							}
							result := new(mel3program.Mel3Program)
							result.LibraryID = libraryId
							result.ProgramID = MATRIXCONST
							if rowM {
								result.ProgramValue = fmt.Sprintf("ref:%d:%d", major, minor)
							} else {
								result.ProgramValue = fmt.Sprintf("ref:%d:%d", minor, major)
							}
							result.NextPrograms = nil
							ev.Result = result
							return nil
						}
					}
				}

				ev.error = errors.New("wrong argument")
				return nil
			default:
				ev.error = errors.New("unknown ProgramID")
				return nil
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
