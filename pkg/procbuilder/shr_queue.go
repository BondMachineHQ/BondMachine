package procbuilder

import (
	"strconv"
	"strings"
)

type Queue struct{}

func (op Queue) Shr_get_name() string {
	return "queue"
}

func (op Queue) Shortname() string {
	return "q"
}

func (op Queue) GetArchHeader(arch *Arch, shared_constraint string, seq int) string {
	result := ""
	queueName := "q" + strconv.Itoa(seq)
	for _, op := range arch.Op {
		if op.Op_get_name() == "r2q" {
			result += ", " + queueName + "senderData, " + queueName + "senderWrite, " + queueName + "senderAck"
			break
		}
	}
	for _, op := range arch.Op {
		if op.Op_get_name() == "q2r" {
			result += ", " + queueName + "receiverData, " + queueName + "receiverRead, " + queueName + "receiverAck"
			break
		}
	}
	result += ", " + queueName + "empty, " + queueName + "full"
	return result
}

func (op Queue) GetArchParams(arch *Arch, shared_constraint string, seq int) string {
	queueName := "q" + strconv.Itoa(seq)
	result := ""

	for _, op := range arch.Op {
		if op.Op_get_name() == "r2q" {
			result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + queueName + "senderData;\n"
			result += "	output " + queueName + "senderWrite;\n"
			result += "	input " + queueName + "senderAck;\n"
			break
		}
	}
	for _, op := range arch.Op {
		if op.Op_get_name() == "q2r" {
			result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + queueName + "receiverData;\n"
			result += "	output " + queueName + "receiverRead;\n"
			result += "	input " + queueName + "receiverAck;\n"
			break
		}
	}

	result += "	input " + queueName + "empty;\n"
	result += "	input " + queueName + "full;\n"

	return result
}

func (op Queue) GetCPParams(arch *Arch, shared_constraint string, seq int) string {

	result := ""

	queueNum := 0
	if arch.Shared_constraints != "" {
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if soname == "queue" {
				queueNum++
			}
		}
	}

	queueName := "q" + strconv.Itoa(seq)

	for _, op := range arch.Op {
		if op.Op_get_name() == "r2q" {
			result += "	output reg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + queueName + "senderData;\n"
			result += "	output reg " + queueName + "senderWrite;\n"
			result += "	input " + queueName + "senderAck;\n"
			break
		}
	}
	for _, op := range arch.Op {
		if op.Op_get_name() == "q2r" {
			result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + queueName + "receiverData;\n"
			result += "	output reg " + queueName + "receiverRead;\n"
			result += "	input " + queueName + "receiverAck;\n"
			break
		}
	}

	result += "	input " + queueName + "empty;\n"
	result += "	input " + queueName + "full;\n"

	return result
}
