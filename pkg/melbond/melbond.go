package melbond

import (
	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/mmirko/mel/pkg/m3number"
)

type Group []string

type MelBondConfig struct {
	RegisterSize  uint8
	DataType      string
	TypePrefix    string
	Params        map[string]string
	Pruned        []string
	Collapsed     []Group
	Debug         bool
	Verbose       bool
	NeuronLibPath string
	IOMode        uint8
	*bminfo.BMinfo
}

type MelBondProgram struct {
	*MelBondConfig
	*m3number.M3numberMe3li
	Source   string
	BasmCode string
}

func (p *MelBondProgram) WriteBasm() (string, error) {
	if err := p.MelStringImport(p.Source); err != nil {
		return "", err
	} else {
		if err := p.Compute(); err != nil {
			return "", err
		} else {
			return p.BasmCode, nil
		}
	}
}
