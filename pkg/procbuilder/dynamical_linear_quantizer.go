package procbuilder

import (
	"regexp"
	"strconv"
)

type DynLinearQuantizer struct{}

func (d DynLinearQuantizer) MatchName(name string) bool {
	re := regexp.MustCompile("multlqs(?P<s>[0-9]+)t(?P<t>[0-9]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("addlqs(?P<s>[0-9]+)t(?P<t>[0-9]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("divflqs(?P<s>[0-9]+)t(?P<t>[0-9]+)")
	return re.MatchString(name)
}

func (d DynLinearQuantizer) CreateInstruction(name string) (Opcode, error) {
	var t, s int
	re := regexp.MustCompile("multlqs(?P<e>[0-9]+)t(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		ts := re.ReplaceAllString(name, "${t}")
		s, _ = strconv.Atoi(ss)
		t, _ = strconv.Atoi(ts)
	}
	re = regexp.MustCompile("addlqs(?P<e>[0-9]+)t(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		ts := re.ReplaceAllString(name, "${t}")
		s, _ = strconv.Atoi(ss)
		t, _ = strconv.Atoi(ts)
	}
	re = regexp.MustCompile("divlqs(?P<e>[0-9]+)t(?P<f>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		ts := re.ReplaceAllString(name, "${t}")
		s, _ = strconv.Atoi(ss)
		t, _ = strconv.Atoi(ts)
	}

	return LinearQuantizer{lqName: name, s: s, t: t}, nil

}
