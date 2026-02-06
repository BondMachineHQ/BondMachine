package procbuilder

import (
	"regexp"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmconfig"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

const (
	FXPADD = uint8(0) + iota
	FXPMULT
	FXPDIV
)

type DynFXP struct {
}

func (d DynFXP) GetName() string {
	return "dyn_fxp"
}

func (d DynFXP) MatchName(name string) bool {
	re := regexp.MustCompile("multfxps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("addfxps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("divfxps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	return re.MatchString(name)
}

func (d DynFXP) CreateInstruction(name string) (Opcode, error) {
	var f, s int
	var opType uint8
	re := regexp.MustCompile("multfxps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ = strconv.Atoi(ss)
		f, _ = strconv.Atoi(fs)
		opType = FXPMULT
	}
	re = regexp.MustCompile("addfxps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ = strconv.Atoi(ss)
		f, _ = strconv.Atoi(fs)
		opType = FXPADD
	}
	re = regexp.MustCompile("divfxps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ = strconv.Atoi(ss)
		f, _ = strconv.Atoi(fs)
		opType = FXPDIV
	}

	return FXP{fpName: name, s: s, f: f, opType: opType, pipeline: new(uint8)}, nil

}

func (d DynFXP) HLAssemblerGeneratorMatch(bmc *bmconfig.BmConfig) []string {
	result := make([]string, 0)
	return result
}

func (d DynFXP) HLAssemblerGeneratorList(bmc *bmconfig.BmConfig, bl *bmline.BasmLine) []string {
	result := make([]string, 0)
	return result
}
