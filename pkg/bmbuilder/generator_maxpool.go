package bmbuilder

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

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

	basmCode := `%section maxpool .romtext iomode:async
	entry _start    ; Entry point
_start:
	clr	r1
`
	for i := 0; i < mpNum; i++ {
		basmCode += `	mov	r0, i` + strconv.Itoa(i) + `
	add	r1, r0` + "\n"
	}
	basmCode += "\t" + `mov	o0, r1
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

	basmCode += `%meta bmdef	global  registersize:8` + "\n"

	if b.debug {
		fmt.Println(green("\t\t\tMaxPool Generator - Start"))
		defer fmt.Println(green("\t\t\tMaxPool Generator - End"))
	}

	l.BasmMeta = l.BasmMeta.AddMeta("basmcode", basmCode)

	return BasmGenerator(b, e, l)

}
