package bmmatrix

import (
	//"math/rand"
	//"fmt"
	"github.com/mmirko/mel/pkg/mel"
	mel3program "github.com/mmirko/mel/pkg/mel3program"
)

// Program IDs
const (
	MATRIXCONST = uint16(0) + iota
	MATRIXMULT
)

// Program types
const (
	MATRIX = uint16(0) + iota
)

const (
	MYLIBID = mel3program.LIB_M3BMMATRIX
)

// The Mel3 implementation
var Implementation = mel3program.Mel3Implementation{
	ProgramNames: map[uint16]string{
		MATRIXCONST: "m",
		MATRIXMULT:  "mult",
	},
	TypeNames: map[uint16]string{
		MATRIX: "matrix",
	},
	ProgramTypes: map[uint16]mel3program.ArgumentsTypes{
		MATRIXCONST: mel3program.ArgumentsTypes{mel3program.ArgType{MYLIBID, MATRIX, []uint64{}}},
		MATRIXMULT:  mel3program.ArgumentsTypes{mel3program.ArgType{MYLIBID, MATRIX, []uint64{}}},
	},
	NonVariadicArgs: map[uint16]mel3program.ArgumentsTypes{
		MATRIXCONST: mel3program.ArgumentsTypes{},
		MATRIXMULT:  mel3program.ArgumentsTypes{mel3program.ArgType{MYLIBID, MATRIX, []uint64{}}, mel3program.ArgType{MYLIBID, MATRIX, []uint64{}}},
	},
	IsVariadic: map[uint16]bool{
		MATRIXCONST: false,
		MATRIXMULT:  false,
	},
	VariadicType: map[uint16]mel3program.ArgType{
		MATRIXCONST: mel3program.ArgType{},
		MATRIXMULT:  mel3program.ArgType{},
	},
	ImplName: "m3bmmatrix",
}

// The effective Me3li
type M3BasmMatrix struct {
	mel3program.Mel3Object
}

// ********* Mel interface

// The Mel entry point for M3uintMe3li
func (prog *M3BasmMatrix) MelInit(c *mel.MelConfig, ep *mel.EvolutionParameters) {
	implementations := make(map[uint16]*mel3program.Mel3Implementation)
	implementations[MYLIBID] = &Implementation

	if prog.Mel3Object.DefaultCreator == nil {
		creators := make(map[uint16]mel3program.Mel3VisitorCreator)
		creators[MYLIBID] = BasmExporterCreator
		prog.Mel3Init(c, implementations, creators, ep)
	} else {
		creators := mel3program.CreateGenericCreators(&prog.Mel3Object, ep, implementations)
		prog.Mel3Init(c, implementations, creators, ep)
	}

}

func (prog *M3BasmMatrix) MelCopy() mel.Me3li {
	var result mel.Me3li
	return result
}

// The effective Me3li
type M3MatrixInfo struct {
	mel3program.Mel3Object
}

// ********* Mel interface

// The Mel entry point for M3uintMe3li
func (prog *M3MatrixInfo) MelInit(c *mel.MelConfig, ep *mel.EvolutionParameters) {
	implementations := make(map[uint16]*mel3program.Mel3Implementation)
	implementations[MYLIBID] = &Implementation

	if prog.Mel3Object.DefaultCreator == nil {
		creators := make(map[uint16]mel3program.Mel3VisitorCreator)
		creators[MYLIBID] = MatrixInfoCreator
		prog.Mel3Init(c, implementations, creators, ep)
	} else {
		creators := mel3program.CreateGenericCreators(&prog.Mel3Object, ep, implementations)
		prog.Mel3Init(c, implementations, creators, ep)
	}

}

func (prog *M3MatrixInfo) MelCopy() mel.Me3li {
	var result mel.Me3li
	return result
}
