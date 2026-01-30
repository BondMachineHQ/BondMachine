package bmnumbers

import (
	"errors"
	"regexp"
	"strconv"
)

type DynFPX struct {
}

func (d DynFPX) GetName() string {
	return "dyn_fpx"
}

func (d DynFPX) MatchName(name string) bool {
	re := regexp.MustCompile("fpxs(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	return re.MatchString(name)
}

func (d DynFPX) CreateType(name string, param interface{}) (BMNumberType, error) {

	re := regexp.MustCompile("fpxs(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ := strconv.Atoi(ss)
		f, _ := strconv.Atoi(fs)
		i := make(map[string]string)
		i["multop"] = "multfpxs" + ss + "f" + fs
		i["addop"] = "addfpxs" + ss + "f" + fs
		i["divop"] = "divfpxs" + ss + "f" + fs
		return FPX{FixedPointName: name, s: s, f: f, instructions: i}, nil
	}

	return nil, errors.New("creation failed")

}
