package brvgasdl

import "github.com/BondMachineHQ/BondMachine/brvga"

type BrvgaSdl struct {
	*brvga.BrvgaTextMemory
}

func NewBrvgaSdl(constraint string) (*BrvgaSdl, error) {
	result := new(BrvgaSdl)
	textMem, err := brvga.NewBrvgaTextMemory(constraint)
	if err != nil {
		return nil, err
	}
	result.BrvgaTextMemory = textMem
	return result, nil
}
