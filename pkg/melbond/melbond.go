package melbond

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/mmirko/mel/pkg/m3melbond"
)

type Group []string

type Neuron struct {
	Params []string
}

type MelBondConfig struct {
	RegisterSize  uint8
	DataType      string
	TypePrefix    string
	Params        map[string]string
	Neurons       map[string]*Neuron
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

func (n *MelBondConfig) Init() error {
	n.Neurons = make(map[string]*Neuron)

	// List nb files in the neuron library path and load them
	neuronFiles, err := ioutil.ReadDir(n.NeuronLibPath)
	if err != nil {
		return err
	}
	for _, f := range neuronFiles {
		if len(f.Name()) > 3 && f.Name()[len(f.Name())-3:] == ".nb" {
			neuronFile, err := ioutil.ReadFile(n.NeuronLibPath + "/" + f.Name())
			if err != nil {
				return err
			}
			neuron := new(Neuron)
			if err := json.Unmarshal(neuronFile, neuron); err != nil {
				return err
			}
			n.Neurons[f.Name()[0:len(f.Name())-3]] = neuron
		}
	}

	return nil
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
