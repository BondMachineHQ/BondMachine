package procbuilder

import (
	"strconv"
	"strings"
)

type Kbd struct{}

func (op Kbd) Shr_get_name() string {
	return "kbd"
}

func (op Kbd) Shortname() string {
	return "k"
}

func (op Kbd) GetArchHeader(arch *Arch, shared_constraint string, seq int) string {
	result := ""
	kbdName := "k" + strconv.Itoa(seq)
	for _, op := range arch.Op {
		if op.Op_get_name() == "k2r" {
			result += ", " + kbdName + "receiverData, " + kbdName + "receiverRead, " + kbdName + "receiverAck"
			break
		}
	}
	result += ", " + kbdName + "empty, " + kbdName + "full"
	return result
}

func (op Kbd) GetArchParams(arch *Arch, shared_constraint string, seq int) string {
	kbdName := "k" + strconv.Itoa(seq)
	result := ""

	for _, op := range arch.Op {
		if op.Op_get_name() == "k2r" {
			result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + kbdName + "receiverData;\n"
			result += "	output " + kbdName + "receiverRead;\n"
			result += "	input " + kbdName + "receiverAck;\n"
			break
		}
	}

	result += "	input " + kbdName + "empty;\n"
	result += "	input " + kbdName + "full;\n"

	return result
}

func (op Kbd) GetCPParams(arch *Arch, shared_constraint string, seq int) string {

	result := ""

	kbdNum := 0
	if arch.Shared_constraints != "" {
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if soname == "kbd" {
				kbdNum++
			}
		}
	}

	kbdName := "u" + strconv.Itoa(seq)

	for _, op := range arch.Op {
		if op.Op_get_name() == "k2r" {
			result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + kbdName + "receiverData;\n"
			result += "	output reg " + kbdName + "receiverRead;\n"
			result += "	input " + kbdName + "receiverAck;\n"
			break
		}
	}

	result += "	input " + kbdName + "empty;\n"
	result += "	input " + kbdName + "full;\n"

	return result
}
