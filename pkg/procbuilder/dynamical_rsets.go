package procbuilder

import (
	"regexp"
	"strconv"
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
