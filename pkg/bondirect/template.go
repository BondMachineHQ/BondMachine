package bondirect

import (
	"text/template"
)

type TData struct {
	// Define the fields for Tdata
	Rsize int
}

func (be *BondirectElement) InitTData() {
	be.TData = &TData{
		Rsize: int(be.Config.Rsize),
	}
}

var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
	"dec": func(i int) int {
		return i - 1
	},
	"next": func(i int, max int) int {
		if i < max-1 {
			return i + 1
		} else {
			return 0
		}
	},
	"bits": func(i int) int {
		return NeededBits(i)
	},
}

func NeededBits(num int) int {
	if num > 0 {
		for bits := 1; true; bits++ {
			if 1<<uint8(bits) >= num {
				return bits
			}
		}
	}
	return 0
}
