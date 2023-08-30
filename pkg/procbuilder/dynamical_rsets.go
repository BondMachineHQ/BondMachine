package procbuilder

import (
	"regexp"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmconfig"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

type DynRsets struct {
}

func (d DynRsets) GetName() string {
	return "dyn_rset"
}

func (d DynRsets) MatchName(name string) bool {
	re := regexp.MustCompile("rsets(?P<s>[0-9]+)")
	return re.MatchString(name)
}

func (d DynRsets) CreateInstruction(name string) (Opcode, error) {
	var s int
	re := regexp.MustCompile("rsets(?P<s>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		s, _ = strconv.Atoi(ss)
	}

	return Rsets{rsetsName: name, s: s}, nil

}
func (d DynRsets) HLAssemblerGeneratorMatch(bmc *bmconfig.BmConfig) []string {
	result := make([]string, 0)
	if !bmc.IsActive(bmconfig.DisableDynamicalMatching) {
		result = append(result, "mov::*--type=reg::*--type=number")
	}
	return result
}

func (d DynRsets) HLAssemblerGeneratorList(bmc *bmconfig.BmConfig, line *bmline.BasmLine) []string {
	result := make([]string, 0)
	if !bmc.IsActive(bmconfig.DisableDynamicalMatching) {
		result = append(result, "rsets5")
		result = append(result, "rsets6")
		result = append(result, "rsets7")
	}
	return result
}
