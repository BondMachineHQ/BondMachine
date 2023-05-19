package procbuilder

import (
	"strconv"
	"strings"
)

type Uart struct{}

func (op Uart) Shr_get_name() string {
	return "uart"
}

func (op Uart) Shortname() string {
	return "u"
}

func (op Uart) GetArchHeader(arch *Arch, shared_constraint string, seq int) string {
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
	result += ", " + uartName + "rempty, " + uartName + "rfull, " + uartName + "wempty, " + uartName + "wfull"
	return result
}

func (op Uart) GetArchParams(arch *Arch, shared_constraint string, seq int) string {
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

	result += "	input " + uartName + "rempty;\n"
	result += "	input " + uartName + "rfull;\n"
	result += "	input " + uartName + "wempty;\n"
	result += "	input " + uartName + "wfull;\n"

	return result
}

func (op Uart) GetCPParams(arch *Arch, shared_constraint string, seq int) string {

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

	result += "	input " + uartName + "rempty;\n"
	result += "	input " + uartName + "rfull;\n"
	result += "	input " + uartName + "wempty;\n"
	result += "	input " + uartName + "wfull;\n"

	return result
}
