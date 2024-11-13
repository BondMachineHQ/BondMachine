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
	MINPUT = uint8(0) + iota // Input matrix
	MSTD                     // Standard matrix with values
	MREF                     // Reference to a matrix, there is not values but it depends other matrix
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
	basmCode *string
}

func newExporterEnv() exporterEnv {
	result := new(exporterEnv)
	l := new(lists)
	l.ls = make(map[uint64]listInfo)
	result.lists = l
	result.listId = 0
	result.basmCode = new(string)
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
		fmt.Printf("Mel3Object at %p\n", obj)
		fmt.Printf("Environment at %p\n", &env)
	}

	listId := env.listId + 1
	env.listId = listId
	implementations := obj.Implementation
	var envI interface{} = env
	obj.Environment = envI
	ev.SetMel3Object(obj)

	defer func() {
		if debug {
			fmt.Printf("Environment at %p\n", &env)
			fmt.Printf("Mel3Object at %p\n", obj)
			fmt.Println("Tree so far:", env.ls)
		}

		var envI interface{} = env
		obj.Environment = envI
		ev.SetMel3Object(obj)
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

		if arg_num == 2 {
			res0 := evaluators[0].GetResult()
			res1 := evaluators[1].GetResult()
			var value0 string
			var lInfo0 listInfo
			if res0 != nil && res0.LibraryID == libraryId && res0.ProgramID == MATRIXCONST {
				value0 = res0.ProgramValue
				if _, lI, err := getMatrixData(value0, 0); err == nil {
					lInfo0 = lI
				} else {
					ev.error = err
					return nil
				}
			} else {
				ev.error = errors.New("wrong argument type")
				return nil
			}

			var value1 string
			var lInfo1 listInfo
			if res1 != nil && res1.LibraryID == libraryId && res1.ProgramID == MATRIXCONST {
				value1 = res1.ProgramValue
				if _, lI, err := getMatrixData(value1, 0); err == nil {
					lInfo1 = lI
				} else {
					ev.error = err
					return nil
				}
			} else {
				ev.error = errors.New("wrong argument type")
				return nil
			}

			opResult := ""

			col0 := strconv.Itoa(lInfo0.cols)
			row0 := strconv.Itoa(lInfo0.rows)
			col1 := strconv.Itoa(lInfo1.cols)
			row1 := strconv.Itoa(lInfo1.rows)

			if col0 == row1 {
				rowS, _ := strconv.Atoi(row0)
				colS, _ := strconv.Atoi(col1)
				lInfo := listInfo{listId: listId, listName: "", mType: MREF, rows: rowS, cols: colS, values: nil}
				env.lists.ls[listId] = lInfo
				opResult = fmt.Sprintf("ref:%s:%s", row0, col1)
			} else {
				ev.error = errors.New("wrong argument, matrix dimensions do not match")
				return nil
			}

			*env.basmCode += fmt.Sprintf("; entering MATRIXMULT with %s and %s\n", value0, value1)

			templ := ev.createBasicTemplateData2M()

			templ.Mtx1 = make([][]string, lInfo0.rows)
			switch lInfo0.mType {
			case MSTD:
				for i := 0; i < lInfo0.rows; i++ {
					templ.Mtx1[i] = make([]string, lInfo0.cols)
					for j := 0; j < lInfo0.cols; j++ {
						templ.Mtx1[i][j] = fmt.Sprintf("%f", lInfo0.values[i][j])
					}
				}
			case MINPUT, MREF:
				label := strconv.Itoa(int(lInfo0.listId))
				if lInfo0.mType == MINPUT {
					label = "in_" + lInfo0.listName + "_" + label
				} else {
					label = "ref_" + label
				}
				for i := 0; i < lInfo0.rows; i++ {
					templ.Mtx1[i] = make([]string, lInfo0.cols)
					for j := 0; j < lInfo0.cols; j++ {
						templ.Mtx1[i][j] = fmt.Sprintf("%s_el_%d_%d", label, i, j)
					}
				}
			}

			templ.Mtx2 = make([][]string, lInfo1.rows)
			switch lInfo1.mType {
			case MSTD:
				for i := 0; i < lInfo1.rows; i++ {
					templ.Mtx2[i] = make([]string, lInfo1.cols)
					for j := 0; j < lInfo1.cols; j++ {
						templ.Mtx2[i][j] = fmt.Sprintf("%f", lInfo1.values[i][j])
					}
				}
			case MINPUT, MREF:
				label := strconv.Itoa(int(lInfo1.listId))
				if lInfo1.mType == MINPUT {
					label = "in_" + lInfo1.listName + "_" + label
				} else {
					label = "ref_" + label
				}
				for i := 0; i < lInfo1.rows; i++ {
					templ.Mtx2[i] = make([]string, lInfo1.cols)
					for j := 0; j < lInfo1.cols; j++ {
						templ.Mtx2[i][j] = fmt.Sprintf("%s_el_%d_%d", label, i, j)
					}
				}
			}

			switch in_prog.LibraryID {
			case MYLIBID:
				switch in_prog.ProgramID {
				case MATRIXMULT:
					if debug {
						fmt.Println("Processing MATRIXMULT")
					}

					if code, err := ev.ApplyTemplate2M(templ, "mult", templateMult); err == nil {
						*env.basmCode += code
					} else {
						ev.error = err
						return nil
					}

					result := new(mel3program.Mel3Program)
					result.LibraryID = libraryId
					result.ProgramID = MATRIXCONST
					result.ProgramValue = opResult
					result.NextPrograms = nil
					ev.Result = result
					return nil
				default:
					ev.error = errors.New("unknown LibraryID")
					return nil
				}
			}
		}

		ev.error = errors.New("wrong number of arguments")
		return nil

	} else {

		switch in_prog.LibraryID {
		case MYLIBID:
			switch in_prog.ProgramID {
			case MATRIXCONST:
				if debug {
					fmt.Println("Processing MATRIXCONST")
				}
				if m, lInfo, err := getMatrixData(in_prog.ProgramValue, listId); err != nil {
					ev.error = err
					return nil
				} else {
					*env.basmCode += fmt.Sprintf("; entering MATRIXCONST with %s\n", m)
					env.lists.ls[listId] = lInfo
					result := new(mel3program.Mel3Program)
					result.LibraryID = libraryId
					result.ProgramID = MATRIXCONST
					result.ProgramValue = m
					result.NextPrograms = nil
					ev.Result = result
					return nil
				}
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

func getMatrixData(programValue string, listId uint64) (string, listInfo, error) {
	m := programValue
	ref := regexp.MustCompile(`^ref:([0-9]+):([0-9]+)$`)
	// Match an alredy ref matrix
	if ref.MatchString(m) {
		rowsS := ref.FindStringSubmatch(m)[1]
		colsS := ref.FindStringSubmatch(m)[2]
		rows, _ := strconv.Atoi(rowsS)
		cols, _ := strconv.Atoi(colsS)
		lInfo := listInfo{listId: listId, listName: "", mType: MREF, rows: rows, cols: cols, values: nil}
		return m, lInfo, nil
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
		// Replace in with ref
		newM := "ref:" + in.FindStringSubmatch(m)[2] + ":" + in.FindStringSubmatch(m)[3]
		return newM, lInfo, nil
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
				var lInfo listInfo
				major := len(matrix)
				minor := len(matrix[0])
				for i := 1; i < major; i++ {
					if len(matrix[i]) != minor {
						return "", listInfo{}, errors.New("matrix rows have different length")
					}
				}
				if rowM {
					lInfo = listInfo{listId: listId, listName: "", mType: MSTD, rows: major, cols: minor, values: matrix}
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
					lInfo = listInfo{listId: listId, listName: "", mType: MSTD, rows: minor, cols: major, values: invMatrix}
				}
				var newM string
				if rowM {
					newM = fmt.Sprintf("ref:%d:%d", major, minor)
				} else {
					newM = fmt.Sprintf("ref:%d:%d", minor, major)
				}
				return newM, lInfo, nil
			}
		}
	}
	return "", listInfo{}, errors.New("wrong argument")
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
