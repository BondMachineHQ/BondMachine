package bmline

// Body lines and elements
import (
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
)

type BasmBody struct {
	*bmmeta.BasmMeta
	Lines []*BasmLine
}

type BasmLine struct {
	*bmmeta.BasmMeta
	Operation *BasmElement
	Elements  []*BasmElement
}

type BasmElement struct {
	*bmmeta.BasmMeta
	string
}

type ByName []*BasmElement

func (s ByName) Len() int {
	return len(s)
}
func (s ByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByName) Less(i, j int) bool {
	return s[i].string < s[j].string
}

func (be *BasmElement) SetValue(val string) {
	if be != nil {
		be.string = val
	}
}

func (be *BasmElement) GetValue() string {
	if be != nil {
		return be.string
	}
	return ""
}
