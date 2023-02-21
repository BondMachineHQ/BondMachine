package procbuilder

import (
	"strconv"
	"strings"
)

type Kbd struct{}

func (op Kbd) Shr_get_name() string {
	return "uart"
}

func (op Kbd) Shortname() string {
	return "u"
}

func (op Kbd) GetArchHeader(arch *Arch, shared_constraint string, seq int) string {
	result := ""
	uartName := "u" + strconv.Itoa(seq)
	for _, op := range arch.Op {
		if op.Op_get_name() == "r2u" {
			result += ", " + uartName + "senderData, " + uartName + "senderWrite, " + uartName + "senderAck"
			break
		}
	}
	for _, op := range arch.Op {
		if op.Op_get_name() == "u2r" {
			result += ", " + uartName + "receiverData, " + uartName + "receiverRead, " + uartName + "receiverAck"
			break
		}
	}
	result += ", " + uartName + "empty, " + uartName + "full"
	return result
}

func (op Kbd) GetArchParams(arch *Arch, shared_constraint string, seq int) string {
	uartName := "u" + strconv.Itoa(seq)
	result := ""

	for _, op := range arch.Op {
		if op.Op_get_name() == "r2u" {
			result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + uartName + "senderData;\n"
			result += "	output " + uartName + "senderWrite;\n"
			result += "	input " + uartName + "senderAck;\n"
			break
		}
	}
	for _, op := range arch.Op {
		if op.Op_get_name() == "u2r" {
			result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + uartName + "receiverData;\n"
			result += "	output " + uartName + "receiverRead;\n"
			result += "	input " + uartName + "receiverAck;\n"
			break
		}
	}

	result += "	input " + uartName + "empty;\n"
	result += "	input " + uartName + "full;\n"

	return result
}

func (op Kbd) GetCPParams(arch *Arch, shared_constraint string, seq int) string {

	result := ""

	uartNum := 0
	if arch.Shared_constraints != "" {
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if soname == "uart" {
				uartNum++
			}
		}
	}

	uartName := "u" + strconv.Itoa(seq)

	for _, op := range arch.Op {
		if op.Op_get_name() == "r2u" {
			result += "	output reg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + uartName + "senderData;\n"
			result += "	output reg " + uartName + "senderWrite;\n"
			result += "	input " + uartName + "senderAck;\n"
			break
		}
	}
	for _, op := range arch.Op {
		if op.Op_get_name() == "u2r" {
			result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + uartName + "receiverData;\n"
			result += "	output reg " + uartName + "receiverRead;\n"
			result += "	input " + uartName + "receiverAck;\n"
			break
		}
	}

	result += "	input " + uartName + "empty;\n"
	result += "	input " + uartName + "full;\n"

	return result
}
