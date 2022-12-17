package bmnumbers

import (
	"errors"
	"regexp"
	"strconv"
)

type DynFloPoCo struct{}

func (d DynFloPoCo) MatchName(name string) bool {
	re := regexp.MustCompile("flpe(?P<e>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		return true
	}

	return false
}

func (d DynFloPoCo) CreateType(name string, param interface{}) (BMNumberType, error) {

	re := regexp.MustCompile("flpe(?P<e>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		es := re.ReplaceAllString(name, "${e}")
		fs := re.ReplaceAllString(name, "${f}")
		e, _ := strconv.Atoi(es)
		f, _ := strconv.Atoi(fs)
		return FloPoCo{floPoCoName: name, e: e, f: f}, nil
	}

	return nil, errors.New("creation failed")

}
