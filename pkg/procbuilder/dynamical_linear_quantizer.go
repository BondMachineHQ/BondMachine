package procbuilder

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
)

const (
	LQADD = uint8(0) + iota
	LQMULT
	LQDIV
)

type DynLinearQuantizer struct {
	Ranges *map[int]bmnumbers.LinearDataRange
}

func (d DynLinearQuantizer) GetName() string {
	return "dyn_linear_quantizer"
}

func (d DynLinearQuantizer) MatchName(name string) bool {
	re := regexp.MustCompile("multlqs(?P<s>[0-9]+)t(?P<t>[0-9]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("addlqs(?P<s>[0-9]+)t(?P<t>[0-9]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("divlqs(?P<s>[0-9]+)t(?P<t>[0-9]+)")
	return re.MatchString(name)
}

func (d DynLinearQuantizer) CreateInstruction(name string) (Opcode, error) {
	if d.Ranges == nil {
		return nil, errors.New("Ranges not initialized")
	}
	var t, s int
	var opType uint8
	var max float64
	re := regexp.MustCompile("multlqs(?P<s>[0-9]+)t(?P<t>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		ts := re.ReplaceAllString(name, "${t}")
		s, _ = strconv.Atoi(ss)
		t, _ = strconv.Atoi(ts)
		opType = LQMULT
		if val, ok := (*d.Ranges)[t]; ok {
			max = val.Max
		} else {
			return nil, errors.New("Invalid range for index " + strconv.Itoa(t))
		}
	}
	re = regexp.MustCompile("addlqs(?P<s>[0-9]+)t(?P<t>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		ts := re.ReplaceAllString(name, "${t}")
		s, _ = strconv.Atoi(ss)
		t, _ = strconv.Atoi(ts)
		opType = LQADD
		if val, ok := (*d.Ranges)[t]; ok {
			max = val.Max
		} else {
			return nil, errors.New("Invalid range for index " + strconv.Itoa(t))
		}
	}
	re = regexp.MustCompile("divlqs(?P<s>[0-9]+)t(?P<t>[0-9]+)")
	if re.MatchString(name) {
		ss := re.ReplaceAllString(name, "${s}")
		ts := re.ReplaceAllString(name, "${t}")
		s, _ = strconv.Atoi(ss)
		t, _ = strconv.Atoi(ts)
		opType = LQDIV
		if val, ok := (*d.Ranges)[t]; ok {
			max = val.Max
		} else {
			return nil, errors.New("Invalid range for index " + strconv.Itoa(t))
		}
	}

	return LinearQuantizer{lqName: name, s: s, t: t, opType: opType, max: max, pipeline: new(uint8)}, nil

}

func (Op DynLinearQuantizer) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
