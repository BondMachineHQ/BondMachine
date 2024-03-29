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
	group  string
	groups *Groups
}

func (ev *BasmEvaluator) GetName() string {
	return "dump"
}

func (ev *BasmEvaluator) GetMel3Object() *mel3program.Mel3Object {
	return ev.Mel3Object
}

func (ev *BasmEvaluator) SetMel3Object(mel3o *mel3program.Mel3Object) {
	ev.Mel3Object = mel3o
	groupsPI := mel3o.Environment
	groupsI := *groupsPI
	groups := groupsI.(Groups)
	ev.groups = &groups
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

	if iProg.LibraryID == mel3program.BUILTINS {
		switch iProg.ProgramID {
		case mel3program.B_IN_INPUT:
			result := new(mel3program.Mel3Program)
			result.LibraryID = mel3program.BUILTINS
			result.ProgramID = mel3program.B_IN_INPUT
			result.ProgramValue = iProg.ProgramValue
			ev.Result = result
			return nil
		case mel3program.B_IN_OUTPUT:
			arg_num := len(iProg.NextPrograms) // It will be 1
			evaluators := make([]mel3program.Mel3Visitor, arg_num)
			names := make([]string, arg_num)
			for i, prog := range iProg.NextPrograms {
				evaluators[i] = mel3program.ProgMux(ev, prog)
				names[i] = nodeName + string(byte(97+i))
				evaluators[i].(*BasmEvaluator).index = ev.index + string(byte(97+i))
				evaluators[i].(*BasmEvaluator).groups = ev.groups
				evaluators[i].(*BasmEvaluator).group = ev.group
				evaluators[i].Visit(prog)
				if evaluators[i].GetError() != nil {
					ev.error = evaluators[i].GetError()
					return nil
				}
				result := evaluators[i].GetResult()
				ev.MelBondConfig.CodeChan <- fmt.Sprintln("%meta filinkdef " + nodeName + "_" + names[i] + " type:fl")
				if result != nil {
					if !(result.LibraryID == mel3program.BUILTINS && result.ProgramID == mel3program.B_IN_INPUT) {
						ev.MelBondConfig.CodeChan <- fmt.Sprintln("%meta filinkatt " + nodeName + "_" + names[i] + " fi:ext, type:output, index:" + iProg.ProgramValue)
						ev.MelBondConfig.CodeChan <- fmt.Sprintln("%meta filinkatt " + nodeName + "_" + names[i] + " fi:" + names[i] + ", type:output, index:0")
					} else {
						ev.error = errors.New("an input is directly connected to an output")
						return nil
					}
				} else {
					ev.error = errors.New("no result")
					return nil
				}
				ev.Result = result
				return nil
			}

		case mel3program.B_IN_GROUP:
			if ev.group != "" {
				ev.error = errors.New("group already defined")
				return nil
			}

			ev.group = iProg.ProgramValue
			arg_num := len(iProg.NextPrograms) // It will be 1
			evaluators := make([]mel3program.Mel3Visitor, arg_num)
			for i, prog := range iProg.NextPrograms {
				evaluators[i] = mel3program.ProgMux(ev, prog)
				evaluators[i].(*BasmEvaluator).index = ev.index
				evaluators[i].(*BasmEvaluator).groups = ev.groups
				evaluators[i].(*BasmEvaluator).group = ev.group
				evaluators[i].Visit(prog)
				if evaluators[i].GetError() != nil {
					ev.error = evaluators[i].GetError()
					return nil
				}
				result := evaluators[i].GetResult()
				if result == nil {
					ev.error = errors.New("no result")
					return nil
				}
				ev.Result = result
				return nil
			}

		case mel3program.B_IN_UNGROUP:
			if ev.group == "" {
				ev.error = errors.New("no group defined")
				return nil
			}
			ev.group = ""
			arg_num := len(iProg.NextPrograms) // It will be 1
			evaluators := make([]mel3program.Mel3Visitor, arg_num)
			for i, prog := range iProg.NextPrograms {
				evaluators[i] = mel3program.ProgMux(ev, prog)
				evaluators[i].(*BasmEvaluator).index = ev.index
				evaluators[i].(*BasmEvaluator).groups = ev.groups
				evaluators[i].(*BasmEvaluator).group = ev.group
				evaluators[i].Visit(prog)
				if evaluators[i].GetError() != nil {
					ev.error = evaluators[i].GetError()
					return nil
				}
				result := evaluators[i].GetResult()
				if result == nil {
					ev.error = errors.New("no result")
					return nil
				}
				ev.Result = result
				return nil
			}
		default:
			ev.error = errors.New("unknown builtin")
			return nil
		}
	} else {
		implementation := ev.Implementation[iProg.LibraryID]

		nodeCodeName := implementation.ImplName + "_" + implementation.ProgramNames[iProg.ProgramID]

		if neuron, ok := ev.MelBondConfig.Neurons[nodeCodeName]; ok {

			ev.MelBondConfig.CodeChan <- fmt.Sprint("%meta fidef " + nodeName + " fragment:" + nodeCodeName)

			// groups management
			myGroup := nodeName
			if ev.group != "" {
				myGroup = ev.group
			}

			if ev.groups == nil {
				newGroups := make(Groups)
				ev.groups = &newGroups
			}
			groups := *ev.groups
			if _, ok := groups[myGroup]; !ok {
				groups[myGroup] = make([]string, 1)
				groups[myGroup][0] = nodeName
			} else {
				groups[myGroup] = append(groups[myGroup], nodeName)
			}

			isFunctional := true

			if len(implementation.NonVariadicArgs[iProg.ProgramID]) == 0 && !implementation.IsVariadic[iProg.ProgramID] {
				isFunctional = false
			}

			for _, param := range neuron.Params {
				if !isFunctional && param == implementation.ProgramNames[iProg.ProgramID] {
					ev.MelBondConfig.CodeChan <- fmt.Sprintf(", %s:%s", implementation.ProgramNames[iProg.ProgramID], iProg.ProgramValue)
					continue
				}
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

			arg_num := len(iProg.NextPrograms)
			evaluators := make([]mel3program.Mel3Visitor, arg_num)
			names := make([]string, arg_num)
			for i, prog := range iProg.NextPrograms {
				evaluators[i] = mel3program.ProgMux(ev, prog)
				names[i] = nodeName + string(byte(97+i))
				evaluators[i].(*BasmEvaluator).index = ev.index + string(byte(97+i))
				evaluators[i].(*BasmEvaluator).groups = ev.groups
				evaluators[i].(*BasmEvaluator).group = ev.group
				evaluators[i].Visit(prog)
				if evaluators[i].GetError() != nil {
					ev.error = evaluators[i].GetError()
					return nil
				}
				result := evaluators[i].GetResult()
				ev.MelBondConfig.CodeChan <- fmt.Sprintln("%meta filinkdef " + nodeName + "_" + names[i] + " type:fl")
				if result != nil {
					if !(result.LibraryID == mel3program.BUILTINS && result.ProgramID == mel3program.B_IN_INPUT) {
						ev.MelBondConfig.CodeChan <- fmt.Sprintln("%meta filinkatt " + nodeName + "_" + names[i] + " fi:" + nodeName + ", type:input, index:" + fmt.Sprint(i))
						ev.MelBondConfig.CodeChan <- fmt.Sprintln("%meta filinkatt " + nodeName + "_" + names[i] + " fi:" + names[i] + ", type:output, index:0")
					} else {
						ev.MelBondConfig.CodeChan <- fmt.Sprintln("%meta filinkatt " + nodeName + "_" + names[i] + " fi:" + nodeName + ", type:input, index:" + fmt.Sprint(i))
						ev.MelBondConfig.CodeChan <- fmt.Sprintln("%meta filinkatt " + nodeName + "_" + names[i] + " fi:ext, type:input, index:" + result.ProgramValue)
					}
				} else {
					ev.error = errors.New("no result")
					return nil
				}
			}
		} else {
			ev.error = errors.New("Unknown neuron " + nodeCodeName)
			return nil
		}
	}
	ev.Result = iProg
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
