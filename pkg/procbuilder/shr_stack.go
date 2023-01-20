package procbuilder

import (
	"strconv"
	"strings"
)

type Stack struct{}

func (op Stack) Shr_get_name() string {
	return "stack"
}

func (op Stack) Shortname() string {
	return "st"
}

func (op Stack) Get_header(arch *Arch, shared_constraint string, seq int) string {
	result := ""
	stackName := "st" + strconv.Itoa(seq)
	for _, op := range arch.Op {
		if op.Op_get_name() == "r2t" {
			result += ", " + stackName + "senderData, " + stackName + "senderWrite, " + stackName + "senderAck"
			break
		}
	}
	for _, op := range arch.Op {
		if op.Op_get_name() == "t2r" {
			result += ", " + stackName + "receiverData, " + stackName + "receiverRead, " + stackName + "receiverAck"
			break
		}
	}
	return result
}

func (op Stack) Get_params(arch *Arch, shared_constraint string, seq int) string {
	stackName := "st" + strconv.Itoa(seq)
	result := ""

	for _, op := range arch.Op {
		if op.Op_get_name() == "r2t" {
			result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + stackName + "senderData;\n"
			result += "	output " + stackName + "senderWrite;\n"
			result += "	input " + stackName + "senderAck;\n"
			break
		}
	}
	for _, op := range arch.Op {
		if op.Op_get_name() == "t2r" {
			result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + stackName + "receiverData;\n"
			result += "	input " + stackName + "receiverRead;\n"
			result += "	output " + stackName + "receiverAck;\n"
			break
		}
	}

	return result
}

func (op Stack) Get_internal_params(arch *Arch, shared_constraint string, seq int) string {

	result := ""

	stackNum := 0
	if arch.Shared_constraints != "" {
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if soname == "stack" {
				stackNum++
			}
		}
	}

	stackName := "st" + strconv.Itoa(seq)

	for _, op := range arch.Op {
		if op.Op_get_name() == "r2t" {
			result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + stackName + "senderData;\n"
			result += "	output " + stackName + "senderWrite;\n"
			result += "	input " + stackName + "senderAck;\n"
			break
		}
	}
	for _, op := range arch.Op {
		if op.Op_get_name() == "t2r" {
			result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + stackName + "receiverData;\n"
			result += "	input " + stackName + "receiverRead;\n"
			result += "	output " + stackName + "receiverAck;\n"
			break
		}
	}

	return result
}
