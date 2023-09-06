package bmnumbers

import (
	"errors"
	"regexp"
	"strconv"
)

type DynFixedPoint struct {
}

func (d DynFixedPoint) GetName() string {
	return "dyn_fixed_point"
}

func (d DynFixedPoint) MatchName(name string) bool {
	re := regexp.MustCompile("fps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	return re.MatchString(name)
}

func (d DynFixedPoint) CreateType(name string, param interface{}) (BMNumberType, error) {

	re := regexp.MustCompile("fps(?P<s>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		fs := re.ReplaceAllString(name, "${f}")
		s, _ := strconv.Atoi(ss)
		f, _ := strconv.Atoi(fs)
		i := make(map[string]string)
		i["multop"] = "multfps" + ss + "t" + fs
		i["addop"] = "addfps" + ss + "t" + fs
		i["divop"] = "divfps" + ss + "t" + fs
		return FixedPoint{FixedPointName: name, s: s, f: f, instructions: i}, nil
	}

	return nil, errors.New("creation failed")

}
