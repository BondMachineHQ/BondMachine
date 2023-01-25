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

func (op Stack) GetArchHeader(arch *Arch, shared_constraint string, seq int) string {
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
	result += ", " + stackName + "empty, " + stackName + "full"
	return result
}

func (op Stack) GetArchParams(arch *Arch, shared_constraint string, seq int) string {
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
			result += "	output " + stackName + "receiverRead;\n"
			result += "	input " + stackName + "receiverAck;\n"
			break
		}
	}

	result += "	input " + stackName + "empty;\n"
	result += "	input " + stackName + "full;\n"

	return result
}

func (op Stack) GetCPParams(arch *Arch, shared_constraint string, seq int) string {

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
			result += "	output reg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + stackName + "senderData;\n"
			result += "	output reg " + stackName + "senderWrite;\n"
			result += "	input " + stackName + "senderAck;\n"
			break
		}
	}
	for _, op := range arch.Op {
		if op.Op_get_name() == "t2r" {
			result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + stackName + "receiverData;\n"
			result += "	output reg " + stackName + "receiverRead;\n"
			result += "	input " + stackName + "receiverAck;\n"
			break
		}
	}

	result += "	input " + stackName + "empty;\n"
	result += "	input " + stackName + "full;\n"

	return result
}
