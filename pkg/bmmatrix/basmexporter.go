package bmmatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	mel3program "github.com/mmirko/mel/pkg/mel3program"
)

const (
	MINPUT = uint8(0) + iota
	MSTD
	MREF
)

// Every list (in the lisp sense) has a listInfo struct and a listId (a unique identifier)
type listInfo struct {
	listId   uint64
	listName string
	mType    uint8
	rows     int
	cols     int
	values   [][]float64
}

type lists struct {
	ls map[uint64]listInfo
}

type BasmExporter struct {
	*mel3program.Mel3Object
	error
	Result *mel3program.Mel3Program
}

type exporterEnv struct {
	*lists
	listId   uint64
	basmCode string
}

func newExporterEnv() exporterEnv {
	result := new(exporterEnv)
	l := new(lists)
	l.ls = make(map[uint64]listInfo)
	result.lists = l
	result.listId = 0
	return *result
}

func BasmExporterCreator() mel3program.Mel3Visitor {
	return new(BasmExporter)
}

func (ev *BasmExporter) GetName() string {
	return "m3bmmatrix"
}

func (ev *BasmExporter) GetMel3Object() *mel3program.Mel3Object {
	return ev.Mel3Object
}

func (ev *BasmExporter) SetMel3Object(mel3o *mel3program.Mel3Object) {
	ev.Mel3Object = mel3o
}

func (ev *BasmExporter) GetError() error {
	return ev.error
}

func (ev *BasmExporter) GetResult() *mel3program.Mel3Program {
	return ev.Result
}

func (ev *BasmExporter) Visit(in_prog *mel3program.Mel3Program) mel3program.Mel3Visitor {

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
	env := (obj.Environment).(exporterEnv)
	if debug {
		fmt.Printf("Get Mel3Object at %p\n", obj)
		fmt.Printf("env at %p\n", &env)
	}

	listId := env.listId + 1
	env.listId = listId
	implementations := obj.Implementation
	var envI interface{} = env
	obj.Environment = envI
	ev.SetMel3Object(obj)

	defer func() {
		if debug {
			fmt.Printf("Put env at %p\n", &env)
			fmt.Printf("Put Mel3Object at %p\n", obj)
		}

		var envI interface{} = env
		obj.Environment = envI
		ev.SetMel3Object(obj)
		fmt.Println(env.ls)
	}()

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
				if debug {
					fmt.Println("MATRIXMULT")
				}
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

					// templ := ev.createBasicTemplateData2M()

					opResult := ""

					values0 := regexp.MustCompile(`^ref:([0-9]+):([0-9]+)$`)
					values1 := regexp.MustCompile(`^ref:([0-9]+):([0-9]+)$`)
					if values0.MatchString(value0) && values1.MatchString(value1) {
						row0 := values0.FindStringSubmatch(value0)[1]
						col0 := values0.FindStringSubmatch(value0)[2]
						row1 := values1.FindStringSubmatch(value1)[1]
						col1 := values1.FindStringSubmatch(value1)[2]
						if col0 == row1 {
							rows, _ := strconv.Atoi(row0)
							cols, _ := strconv.Atoi(col1)
							lInfo := listInfo{listId: listId, listName: "", mType: MREF, rows: rows, cols: cols, values: nil}
							env.lists.ls[listId] = lInfo
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
				if debug {
					fmt.Println("MATRIXCONST")
				}
				m := in_prog.ProgramValue
				ref := regexp.MustCompile(`^ref:([0-9]+):([0-9]+)$`)
				// Match an alredy ref matrix
				if ref.MatchString(m) {
					rowsS := ref.FindStringSubmatch(m)[1]
					colsS := ref.FindStringSubmatch(m)[2]
					rows, _ := strconv.Atoi(rowsS)
					cols, _ := strconv.Atoi(colsS)
					lInfo := listInfo{listId: listId, listName: "", mType: MREF, rows: rows, cols: cols, values: nil}
					env.lists.ls[listId] = lInfo
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
					rowsS := in.FindStringSubmatch(m)[2]
					colsS := in.FindStringSubmatch(m)[3]
					name := in.FindStringSubmatch(m)[1]
					rows, _ := strconv.Atoi(rowsS)
					cols, _ := strconv.Atoi(colsS)
					lInfo := listInfo{listId: listId, listName: name, mType: MINPUT, rows: rows, cols: cols, values: nil}
					env.lists.ls[listId] = lInfo
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
							if rowM {
								lInfo := listInfo{listId: listId, listName: "", mType: MSTD, rows: major, cols: minor, values: matrix}
								env.lists.ls[listId] = lInfo
							} else {
								invMatrix := make([][]float64, minor)
								for i := 0; i < minor; i++ {
									invMatrix[i] = make([]float64, major)
								}
								for i := 0; i < major; i++ {
									for j := 0; j < minor; j++ {
										invMatrix[j][i] = matrix[i][j]
									}
								}
								lInfo := listInfo{listId: listId, listName: "", mType: MSTD, rows: minor, cols: major, values: invMatrix}
								env.lists.ls[listId] = lInfo
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

func (ev *BasmExporter) Inspect() string {
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
