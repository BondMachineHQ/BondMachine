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
