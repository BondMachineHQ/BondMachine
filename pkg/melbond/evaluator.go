package melbond

import (
	"errors"
	"fmt"

	"github.com/mmirko/mel/pkg/mel3program"
)

type BasmEvaluator struct {
	*mel3program.Mel3Object
	MelBondConfig *MelBondConfig
	error
	Result *mel3program.Mel3Program
	index  string
}

func (ev *BasmEvaluator) GetName() string {
	return "dump"
}

func (ev *BasmEvaluator) GetMel3Object() *mel3program.Mel3Object {
	return ev.Mel3Object
}

func (ev *BasmEvaluator) SetMel3Object(mel3o *mel3program.Mel3Object) {
	ev.Mel3Object = mel3o
}

func (ev *BasmEvaluator) GetError() error {
	return ev.error
}

func (ev *BasmEvaluator) GetResult() *mel3program.Mel3Program {
	return ev.Result
}

func (ev *BasmEvaluator) Visit(iProg *mel3program.Mel3Program) mel3program.Mel3Visitor {

	debug := ev.Config.Debug

	if debug {
		fmt.Println("basm: Visit: ", iProg)
	}

	nodeName := "node" + ev.index
	implementation := ev.Implementation[iProg.LibraryID]

	nodeCodeName := implementation.ImplName + "_" + implementation.ProgramNames[iProg.ProgramID]

	if neuron, ok := ev.MelBondConfig.Neurons[nodeCodeName]; ok {

		ev.MelBondConfig.CodeChan <- fmt.Sprint("%meta fidef " + nodeName + " fragment:" + nodeCodeName)
		for _, param := range neuron.Params {
			switch param {
			default:
				if value, ok := ev.MelBondConfig.Params[param]; ok {
					ev.MelBondConfig.CodeChan <- fmt.Sprintf(", %s:%s", param, value)
				} else {
					ev.error = errors.New("Unknown parameter " + param)
					return nil
				}
			}
		}
		ev.MelBondConfig.CodeChan <- "\n"

		isFunctional := true

		if len(implementation.NonVariadicArgs[iProg.ProgramID]) == 0 && !implementation.IsVariadic[iProg.ProgramID] {
			isFunctional = false
		}

		if !isFunctional {
			ev.MelBondConfig.CodeChan <- fmt.Sprintln(", " + implementation.ProgramNames[iProg.ProgramID] + ":" + iProg.ProgramValue)
		} else {
			ev.MelBondConfig.CodeChan <- fmt.Sprintln()
		}

		arg_num := len(iProg.NextPrograms)
		evaluators := make([]mel3program.Mel3Visitor, arg_num)
		names := make([]string, arg_num)
		for i, prog := range iProg.NextPrograms {
			evaluators[i] = mel3program.ProgMux(ev, prog)
			names[i] = nodeName + string(byte(97+i))
			evaluators[i].(*BasmEvaluator).index = ev.index + string(byte(97+i))
			ev.MelBondConfig.CodeChan <- fmt.Sprintln("%meta filinkatt " + nodeName + "_" + names[i] + " fi:" + nodeName + ", type:input, index:" + fmt.Sprint(i))
			ev.MelBondConfig.CodeChan <- fmt.Sprintln("%meta filinkatt " + nodeName + "_" + names[i] + " fi:" + names[i] + ", type:output, index:0")
			evaluators[i].Visit(prog)
		}
	} else {
		ev.error = errors.New("Unknown neuron " + nodeCodeName)
	}

	return nil
}

func (ev *BasmEvaluator) Inspect() string {
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

func (c *MelBondConfig) BasmCreator() mel3program.Mel3Visitor {
	result := new(BasmEvaluator)
	result.MelBondConfig = c
	return result
}
