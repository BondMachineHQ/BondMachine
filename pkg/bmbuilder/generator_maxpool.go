package bmbuilder

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

// TODO Needed: max,

func MaxPoolGenerator(b *BMBuilder, e *bmline.BasmElement, l *bmline.BasmLine) (*bondmachine.Bondmachine, error) {
	// Pay attention to the following: This code works only if the Generator function is called on a sequential block
	// TODO Include here a check and error if the Generator function is called outside a sequential block

	var mpNum int
	if len(l.Elements) != 1 {
		return nil, errors.New("MaxPool Generator: Wrong number of elements, it should be 1")
	} else {
		mpNumS := l.Elements[0].GetValue()
		if val, err := strconv.Atoi(mpNumS); err != nil {
			return nil, errors.New("MaxPool Generator: Error converting MaxPool number to int")
		} else {
			mpNum = val
		}
	}

	cb := b.currentBlock

	var prevIN int

	if len(b.blocks[cb].blockBMs) == 0 {
		return nil, errors.New("MaxPool Generator: No previous BMs found, it cannot be the first block")
	} else {
		prevBM := b.blocks[cb].blockBMs[len(b.blocks[cb].blockBMs)-1]
		prevIN = prevBM.Outputs
	}

	var cpNum int

	if prevIN%mpNum != 0 {
		return nil, errors.New("MaxPool Generator: The number of inputs is not divisible by the MaxPool number")
	} else {
		cpNum = prevIN / mpNum
	}

	// Get the register size, starting from the maxpool metadata and falling back to the global metadata
	// var regSize int
	var regSizeS string
	if regSizeS = l.BasmMeta.GetMeta("registersize"); regSizeS == "" {
		if regSizeS = e.BasmMeta.GetMeta("registersize"); regSizeS == "" {
			if regSizeS = b.global.BasmMeta.GetMeta("registersize"); regSizeS == "" {
				return nil, errors.New("MaxPool Generator: No register size found")
			}
		}
	}
	if _, err := strconv.Atoi(regSizeS); err != nil {
		return nil, errors.New("MaxPool Generator: Error converting register size to int")
	} else {
		// regSize = val
	}

	// Get the data type, starting from the maxpool metadata and falling back to the global metadata
	var dataType string
	if dataType = l.BasmMeta.GetMeta("datatype"); dataType == "" {
		if dataType = e.BasmMeta.GetMeta("datatype"); dataType == "" {
			if dataType = b.global.BasmMeta.GetMeta("datatype"); dataType == "" {
				return nil, errors.New("MaxPool Generator: No data type found")
			}
		}
	}

	// get the instructions from the data type using bmnumbers
	found := false
	ops := make(map[string]string)
	// typePrefix := ""
	for _, tpy := range bmnumbers.AllTypes {
		if tpy.GetName() == dataType {
			for opType, opName := range tpy.ShowInstructions() {
				ops[opType] = opName
			}
			// typePrefix = tpy.ShowPrefix()
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("MaxPool Generator: Data type not found")
	}

	// Get the maxpool type
	var mpType string
	if mpType = l.BasmMeta.GetMeta("maxpooltype"); mpType == "" {
		fmt.Println("MaxPool Generator: No maxpool type found, using default: add")
		mpType = "add"
	}

	// Type check
	switch mpType {
	case "add", "avg":
	default:
		return nil, errors.New("MaxPool Generator: MaxPool type not recognized")
	}

	// The maxpool function header
	basmCode := `%section maxpool .romtext iomode:async
	entry _start    ; Entry point
_start:
`

	// The maxpool function body
	switch mpType {
	case "add":
		var addOp string
		if op, ok := ops["addop"]; ok {
			addOp = op
		} else {
			return nil, errors.New("MaxPool Generator: Add operation not found")
		}
		basmCode += "\t" + `mov	r1, 0` + "\n"
		for i := 0; i < mpNum; i++ {
			basmCode += `	mov	r0, i` + strconv.Itoa(i) + `
	` + addOp + `	r1, r0` + "\n"
		}
		basmCode += "\t" + `mov	o0, r1
`
	case "avg":
		var addOp string
		var divOp string
		if op, ok := ops["addop"]; ok {
			addOp = op
		} else {
			return nil, errors.New("MaxPool Generator: Add operation not found")
		}
		if op, ok := ops["divop"]; ok {
			divOp = op
		} else {
			return nil, errors.New("MaxPool Generator: Div operation not found")
		}
		basmCode += "\t" + `mov	r1, 0` + "\n"
		for i := 0; i < mpNum; i++ {
			basmCode += `	mov	r0, i` + strconv.Itoa(i) + `
	` + addOp + `	r1, r0` + "\n"
		}
		basmCode += "\t" + `mov	r0, ` + strconv.Itoa(mpNum) + "\n"
		basmCode += "\t" + divOp + `	r1, r0
	mov	o0, r1
`
	}

	// The maxpool function footer
	basmCode += `
	j 	_start

%endsection
`
	for i := 0; i < cpNum; i++ {
		basmCode += `%meta cpdef	cpu` + strconv.Itoa(i) + `	romcode: maxpool, ramsize:8` + "\n"
	}

	for i := 0; i < cpNum; i++ {
		for j := 0; j < mpNum; j++ {
			basmCode += `%meta ioatt	in` + strconv.Itoa(i) + `p` + strconv.Itoa(j) + `	cp: cpu` + strconv.Itoa(i) + `, index:` + strconv.Itoa(j) + `, type:input` + "\n"
			basmCode += `%meta ioatt	in` + strconv.Itoa(i) + `p` + strconv.Itoa(j) + `	cp: bm, index:` + strconv.Itoa(i*mpNum+j) + `, type:input` + "\n"
		}
		basmCode += `%meta ioatt	out` + strconv.Itoa(i) + `	cp: cpu` + strconv.Itoa(i) + `, index:0, type:output` + "\n"
		basmCode += `%meta ioatt	out` + strconv.Itoa(i) + `	cp: bm, index:` + strconv.Itoa(i) + `, type:output` + "\n"
	}

	basmCode += `%meta bmdef	global  registersize:` + regSizeS + "\n"

	if b.debug {
		fmt.Println(green("\t\t\tMaxPool Generator - Start"))
		defer fmt.Println(green("\t\t\tMaxPool Generator - End"))
	}

	l.BasmMeta = l.BasmMeta.AddMeta("basmcode", basmCode)

	return BasmGenerator(b, e, l)

}
