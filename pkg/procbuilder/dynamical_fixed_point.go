package procbuilder

import (
	"regexp"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmconfig"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

const (
	FPADD = uint8(0) + iota
	FPMULT
	FPDIV
)

type DynFixedPoint struct {
}

func (d DynFixedPoint) GetName() string {
	return "dyn_fixed_point"
}

func (d DynFixedPoint) MatchName(name string) bool {
	re := regexp.MustCompile("multfps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("addfps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("divfps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	return re.MatchString(name)
}

func (d DynFixedPoint) CreateInstruction(name string) (Opcode, error) {
	var f, s int
	var opType uint8
	re := regexp.MustCompile("multfps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ = strconv.Atoi(ss)
		f, _ = strconv.Atoi(fs)
		opType = FPMULT
	}
	re = regexp.MustCompile("addfps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ = strconv.Atoi(ss)
		f, _ = strconv.Atoi(fs)
		opType = FPADD
	}
	re = regexp.MustCompile("divfps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ = strconv.Atoi(ss)
		f, _ = strconv.Atoi(fs)
		opType = FPDIV
	}

	return FixedPoint{fpName: name, s: s, f: f, opType: opType, pipeline: new(uint8)}, nil

}

func (d DynFixedPoint) HLAssemblerGeneratorMatch(bmc *bmconfig.BmConfig) []string {
	result := make([]string, 0)
	return result
}

func (d DynFixedPoint) HLAssemblerGeneratorList(bmc *bmconfig.BmConfig, bl *bmline.BasmLine) []string {
	result := make([]string, 0)
	return result
}
