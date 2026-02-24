package bmnumbers

import (
	"errors"
	"regexp"
	"strconv"
)

type DynFXP struct {
}

func (d DynFXP) GetName() string {
	return "dyn_fxp"
}

func (d DynFXP) MatchName(name string) bool {
	re := regexp.MustCompile("fxps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	return re.MatchString(name)
}

func (d DynFXP) CreateType(name string, param interface{}) (BMNumberType, error) {

	re := regexp.MustCompile("fxps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ := strconv.Atoi(ss)
		f, _ := strconv.Atoi(fs)
		i := make(map[string]string)
		i["multop"] = "multfxps" + ss + "f" + fs
		i["addop"] = "addfxps" + ss + "f" + fs
		i["divop"] = "divfxps" + ss + "f" + fs
		return FXP{FixedPointName: name, s: s, f: f, instructions: i}, nil
	}

	return nil, errors.New("creation failed")

}
