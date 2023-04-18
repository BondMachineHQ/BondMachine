package bmnumbers

import (
	"errors"
	"regexp"
	"strconv"
)

type LinearDataRange struct {
	Max float64
}

type DynLinearQuantizer struct {
	Ranges *map[int]LinearDataRange
}

func (d DynLinearQuantizer) GetName() string {
	return "dyn_linear_quantizer"
}

func (d DynLinearQuantizer) MatchName(name string) bool {
	re := regexp.MustCompile("lqs(?P<s>[0-9]+)t(?P<t>[0-9]+)")
	if re.MatchString(name) {
		return true
	}

	return false
}

func (d DynLinearQuantizer) CreateType(name string, param interface{}) (BMNumberType, error) {

	re := regexp.MustCompile("lqs(?P<s>[0-9]+)t(?P<t>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		ts := re.ReplaceAllString(name, "${t}")
		s, _ := strconv.Atoi(ss)
		t, _ := strconv.Atoi(ts)
		i := make(map[string]string)
		i["multop"] = "multlqs" + ss + "t" + ts
		i["addop"] = "addlqs" + ss + "t" + ts
		i["divop"] = "divlqs" + ss + "t" + ts
		return LinearQuantizer{linearQuantizerName: name, s: s, t: t, instructions: i}, nil
	}

	return nil, errors.New("creation failed")

}
