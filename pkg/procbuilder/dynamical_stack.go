package procbuilder

import (
	"regexp"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmconfig"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

const (
	OP_PUSH = uint8(0) + iota
	OP_PULL
)

type DynStack struct {
}

func (d DynStack) GetName() string {
	return "dyn_stack"
}

func (d DynStack) MatchName(name string) bool {
	re := regexp.MustCompile("push(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("pull(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)")
	return re.MatchString(name)
}

func (d DynStack) CreateInstruction(name string) (Opcode, error) {
	var s int
	var sn string
	var opType uint8
	re := regexp.MustCompile("push(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${stacksize}")
		sn = re.ReplaceAllString(name, "${stackname}")
		s, _ = strconv.Atoi(ss)
		opType = OP_PUSH
	}
	re = regexp.MustCompile("pull(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${stacksize}")
		sn = re.ReplaceAllString(name, "${stackname}")
		s, _ = strconv.Atoi(ss)
		opType = OP_PULL
	}

	return DynOpStack{callName: name, s: s, sn: sn, opType: opType}, nil

}

func (d DynStack) HLAssemblerGeneratorMatch(bmc *bmconfig.BmConfig) []string {
	result := make([]string, 0)
	return result
}

func (d DynStack) HLAssemblerGeneratorList(bmc *bmconfig.BmConfig, bl *bmline.BasmLine) []string {
	result := make([]string, 0)
	return result
}
