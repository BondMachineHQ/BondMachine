package procbuilder

import (
	"regexp"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmconfig"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

const (
	OP_CALLO = uint8(0) + iota
	OP_CALLA
	OP_RET
)

type DynCall struct {
}

func (d DynCall) GetName() string {
	return "dyn_call"
}

func (d DynCall) MatchName(name string) bool {
	re := regexp.MustCompile("callo(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("calla(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("ret(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)")
	return re.MatchString(name)
}

func (d DynCall) CreateInstruction(name string) (Opcode, error) {
	var s int
	var sn string
	var opType uint8
	re := regexp.MustCompile("callo(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${stacksize}")
		sn = re.ReplaceAllString(name, "${stackname}")
		s, _ = strconv.Atoi(ss)
		opType = OP_CALLO
	}
	re = regexp.MustCompile("calla(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${stacksize}")
		sn = re.ReplaceAllString(name, "${stackname}")
		s, _ = strconv.Atoi(ss)
		opType = OP_CALLA
	}
	re = regexp.MustCompile("ret(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${stacksize}")
		sn = re.ReplaceAllString(name, "${stackname}")
		s, _ = strconv.Atoi(ss)
		opType = OP_RET
	}

	return Call{callName: name, s: s, sn: sn, opType: opType}, nil

}

func (d DynCall) HLAssemblerGeneratorMatch(bmc *bmconfig.BmConfig) []string {
	result := make([]string, 0)
	return result
}

func (d DynCall) HLAssemblerGeneratorList(bmc *bmconfig.BmConfig, bl *bmline.BasmLine) []string {
	result := make([]string, 0)
	return result
}
