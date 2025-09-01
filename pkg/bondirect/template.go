package bondirect

import (
	"fmt"
	"text/template"
)

const (
	ActionsNum = 3
)

type TData struct {
	// Define the fields for Tdata
	NodeName     string
	EdgeName     string
	Rsize        int // Register size
	NodeNum      int // Number of nodes
	NodeBits     int // Bits needed for node addressing
	IONum        int // Maximum number of inputs or outputs Among all nodes
	IOBits       int // Bits needed for IO addressing
	InnerMessLen int // Length of inner messages
}

func (be *BondirectElement) InitTData() {
	be.TData = &TData{
		Rsize:    int(be.Config.Rsize),
		NodeNum:  len(be.Cluster.Peers),
		NodeBits: NeededBits(len(be.Cluster.Peers)),
	}

	maxIO := 0
	for _, node := range be.Cluster.Peers {
		if len(node.Inputs) > maxIO {
			maxIO = len(node.Inputs)
		}
		if len(node.Outputs) > maxIO {
			maxIO = len(node.Outputs)
		}
	}
	be.TData.IONum = maxIO
	be.TData.IOBits = NeededBits(maxIO)
	be.TData.InnerMessLen = be.TData.NodeBits + be.TData.IOBits + NeededBits(ActionsNum) + be.TData.Rsize
}

func (be *BondirectElement) DumpTemplateData() string {
	result := ""
	result += fmt.Sprintf("Register Size: %d\n", be.TData.Rsize)
	result += fmt.Sprintf("Node Number: %d\n", be.TData.NodeNum)
	result += fmt.Sprintf("Node Bits: %d\n", be.TData.NodeBits)
	result += fmt.Sprintf("IO Number: %d\n", be.TData.IONum)
	result += fmt.Sprintf("IO Bits: %d\n", be.TData.IOBits)
	result += fmt.Sprintf("Inner Message Length: %d\n", be.TData.InnerMessLen)
	return result
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
