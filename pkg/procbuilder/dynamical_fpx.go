package procbuilder

import (
	"regexp"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmconfig"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

const (
	FPXADD = uint8(0) + iota
	FPXMULT
	FPXDIV
)

type DynFPX struct {
}

func (d DynFPX) GetName() string {
	return "dyn_fpx"
}

func (d DynFPX) MatchName(name string) bool {
	re := regexp.MustCompile("multfpxs(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("addfpxs(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("divfpxs(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	return re.MatchString(name)
}

func (d DynFPX) CreateInstruction(name string) (Opcode, error) {
	var f, s int
	var opType uint8
	re := regexp.MustCompile("multfpxs(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ = strconv.Atoi(ss)
		f, _ = strconv.Atoi(fs)
		opType = FPXMULT
	}
	re = regexp.MustCompile("addfpxs(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ = strconv.Atoi(ss)
		f, _ = strconv.Atoi(fs)
		opType = FPXADD
	}
	re = regexp.MustCompile("divfpxs(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ = strconv.Atoi(ss)
		f, _ = strconv.Atoi(fs)
		opType = FPXDIV
	}

	return FPX{fpName: name, s: s, f: f, opType: opType, pipeline: new(uint8)}, nil

}

func (d DynFPX) HLAssemblerGeneratorMatch(bmc *bmconfig.BmConfig) []string {
	result := make([]string, 0)
	return result
}

func (d DynFPX) HLAssemblerGeneratorList(bmc *bmconfig.BmConfig, bl *bmline.BasmLine) []string {
	result := make([]string, 0)
	return result
}
