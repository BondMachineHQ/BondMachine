package melbond

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/mmirko/mel/pkg/m3melbond"
)

const (
	ASYNC = uint8(0) + iota
	SYNC
)

type Group []string

type Groups map[string]Group

type Neuron struct {
	Params []string
}

type MelBondConfig struct {
	*bminfo.BMinfo
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
	CodeChan      chan string
	CancelCtx     context.CancelFunc
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

	//Create the groups map and copy it to the environment

	var groups interface{}
	groups = make(Groups)
	p.Mel3Object.Environment = &groups

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

		regSize := p.RegisterSize
		p.CodeChan <- fmt.Sprintf("%%meta bmdef     global registersize:%d\n", regSize)
		switch p.IOMode {
		case ASYNC:
			p.CodeChan <- fmt.Sprintf("%%meta bmdef     global iomode:async\n")
		case SYNC:
			p.CodeChan <- fmt.Sprintf("%%meta bmdef     global iomode:sync\n")
		}

		if err := p.Compute(); err != nil {
			return "", err
		} else {
			for gName, group := range groups.(Groups) {
				nodeList := ""
				for _, nName := range group {
					nodeList += ":" + nName
				}
				p.CodeChan <- fmt.Sprintln("%meta cpdef  " + gName + " fragcollapse" + nodeList)
			}
			time.Sleep(1 * time.Second)
			return p.BasmCode, nil
		}
	}
}
