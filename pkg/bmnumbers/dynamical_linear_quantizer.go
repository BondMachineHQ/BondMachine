package bmnumbers

import (
	"errors"
	"regexp"
)

type DynLinearQuantizer struct{}

func (d DynLinearQuantizer) MatchName(name string) bool {
	re := regexp.MustCompile("lq(?P<s>[0-9]+)")
	if re.MatchString(name) {
		return true
	}

	return false
}

func (d DynLinearQuantizer) CreateType(name string, param interface{}) (BMNumberType, error) {

	re := regexp.MustCompile("lq(?P<s>[0-9]+)")
	if re.MatchString(name) {
		// es := re.ReplaceAllString(name, "${e}")
		return LinearQuantizer{linearQuantizerName: name}, nil
	}

	return nil, errors.New("creation failed")

}
