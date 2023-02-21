package bondmachine

import (
	//"fmt"

	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmstack"
)

// The placeholder struct

type Kbd struct{}

func (op Kbd) Shr_get_name() string {
	return "kbd"
}

func (op Kbd) Shr_get_desc() string {
	return "Kbd"
}

func (op Kbd) Shortname() string {
	return "k"
}

func (op Kbd) GV_config(element uint8) string {
	result := ""
	switch element {
	case GVNODEINPROC:
		result += "style=filled fillcolor=pink3 color=black"
	case GVNODE:
		result += "style=filled fillcolor=pink3 color=black"
	case GVEDGE:
		result += "arrowhead=none"
	case GVCLUS:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey70"
	case GVCLUSINPROC:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey70"
	}
	return result
}

func (op Kbd) Instantiate(s string) (Shared_instance, bool) {
	if strings.HasPrefix(s, "kbd:") {
		result := new(Kbd_instance)
		result.Shared_element = op
		components := strings.Split(s, ":")
		componentsN := len(components)
		if componentsN == 2 {
			if depth, ok := strconv.Atoi(components[1]); ok == nil {
				result.Depth = depth
			}
			return result, true
		}
	}
	return nil, false
}

// The instance struct

type Kbd_instance struct {
	Shared_element
	Depth int // The depth of the fifo holding the data from the keyboard
}

func (sm Kbd_instance) String() string {
	return "kbd:" + strconv.Itoa(sm.Depth)
}

func (sm Kbd_instance) Write_verilog(bmach *Bondmachine, soIndex int, kbdName string, flavor string) string {

	result := ""

	// Compute the receivers and senders of the kbd, senders will be the writers of the fifo, receivers will be the readers of the fifo
	receivers := make([]string, 0)

	for numProcessor, soList := range bmach.Shared_links {
		for _, soId := range soList {
			if soId == soIndex {
				for _, op := range bmach.Domains[bmach.Processors[numProcessor]].Op {
					switch op.Op_get_name() {
					case "k2r":
						receivers = append(receivers, "p"+strconv.Itoa(numProcessor)+"kbd_recv")
						continue
					}
				}
			}
		}
	}

	// Create the fifo where the CPs will read from
	rFifo := bmstack.CreateBasicStack()
	rFifo.ModuleName = kbdName + "rfifo"
	rFifo.DataSize = 8
	rFifo.Depth = sm.Depth
	rFifo.MemType = "FIFO"
	rFifo.Senders = []string{"kbdwriter"}
	rFifo.Receivers = receivers

	r, _ := rFifo.WriteHDL()

	result += r

	// TODO : add the kbd module

	return result

}

func (sm Kbd_instance) GetPerProcPortsWires(bmach *Bondmachine, procId int, soId int, flavor string) string {
	result := ""
	if soName, ok := bmach.Get_so_name(soId); ok {
		for _, op := range bmach.Domains[bmach.Processors[procId]].Op {
			if op.Op_get_name() == "k2r" {
				result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(procId) + soName + "receiverData;\n"
				result += "\twire p" + strconv.Itoa(procId) + soName + "receiverRead;\n"
				result += "\twire p" + strconv.Itoa(procId) + soName + "receiverAck;\n"
				result += "\n"
				break
			}
		}
	}
	return result
}

func (sm Kbd_instance) GetPerProcPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soName, ok := bmach.Get_so_name(so_id); ok {
		for _, op := range bmach.Domains[bmach.Processors[proc_id]].Op {
			if op.Op_get_name() == "k2r" {
				result += ", p" + strconv.Itoa(proc_id) + soName + "receiverData"
				result += ", p" + strconv.Itoa(proc_id) + soName + "receiverRead"
				result += ", p" + strconv.Itoa(proc_id) + soName + "receiverAck"
				break
			}
		}
	}
	return result
}

func (sm Kbd_instance) GetExternalPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Kbd_instance) GetExternalPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Kbd_instance) GetCPSharedPortsHeader(bmach *Bondmachine, soId int, flavor string) string {
	result := ""
	if soName, ok := bmach.Get_so_name(soId); ok {
		result += ", " + soName + "empty"
		result += ", " + soName + "full"
	}
	return result
}

func (sm Kbd_instance) GetCPSharedPortsWires(bmach *Bondmachine, soId int, flavor string) string {
	result := ""
	if soName, ok := bmach.Get_so_name(soId); ok {
		result += "\n"
		result += "	wire " + soName + "empty;\n"
		result += "	wire " + soName + "full\n;"
		result += "\n"
	}
	return result
}
