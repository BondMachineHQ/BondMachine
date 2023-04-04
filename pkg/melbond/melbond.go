package melbond

import (
	"context"

	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/mmirko/mel/pkg/m3melbond"
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
	CodeChan  chan string
	CancelCtx context.CancelFunc
}

type MelBondProgram struct {
	*MelBondConfig
	*m3melbond.M3MelBondMe3li
	Source   string
	BasmCode string
}

func (p *MelBondProgram) WriteBasm() (string, error) {

	ctx, cancel := context.WithCancel(context.Background())
	p.CodeChan = make(chan string)
	p.CancelCtx = cancel

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case code := <-p.CodeChan:
				p.BasmCode += code
			}
		}
	}()

	defer cancel()

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
