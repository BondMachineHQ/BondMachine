package bmnumbers

import (
	"errors"
	"regexp"
	"strconv"
)

type DynFloPoCoFixedPoint struct{}

func (d DynFloPoCoFixedPoint) GetName() string {
	return "dyn_flopocofp"
}

func (d DynFloPoCoFixedPoint) MatchName(name string) bool {
	re := regexp.MustCompile("flpfps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	return re.MatchString(name)
}

func (d DynFloPoCoFixedPoint) CreateType(name string, param interface{}) (BMNumberType, error) {

	re := regexp.MustCompile("flpfps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ := strconv.Atoi(ss)
		f, _ := strconv.Atoi(fs)
		i := make(map[string]string)
		i["multop"] = "multflpfps" + ss + "f" + fs
		i["addop"] = "addflpfps" + ss + "f" + fs
		i["divop"] = "divflpfps" + ss + "f" + fs
		return FloPoCoFixedPoint{FixedPointName: name, s: s, f: f, instructions: i}, nil
	}

	return nil, errors.New("creation failed")

}
