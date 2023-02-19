package bondmachine

import (
	//"fmt"

	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmstack"
)

// The placeholder struct

type Uart struct{}

func (op Uart) Shr_get_name() string {
	return "uart"
}

func (op Uart) Shr_get_desc() string {
	return "Uart"
}

func (op Uart) Shortname() string {
	return "u"
}

func (op Uart) GV_config(element uint8) string {
	result := ""
	switch element {
	case GVNODEINPROC:
		result += "style=filled fillcolor=pink2 color=black"
	case GVNODE:
		result += "style=filled fillcolor=pink2 color=black"
	case GVEDGE:
		result += "arrowhead=none"
	case GVCLUS:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey70"
	case GVCLUSINPROC:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey70"
	}
	return result
}

func (op Uart) Instantiate(s string) (Shared_instance, bool) {
	if strings.HasPrefix(s, "uart:") {
		result := new(Uart_instance)
		result.Shared_element = op
		components := strings.Split(s, ":")
		componentsN := len(components)
		if componentsN == 3 {
			if baudRate, ok := strconv.Atoi(components[1]); ok == nil {
				result.BaudRate = baudRate
			}
			if depth, ok := strconv.Atoi(components[2]); ok == nil {
				result.Depth = depth
			}
			return result, true
		}
	}
	return nil, false
}

// The instance struct

type Uart_instance struct {
	Shared_element
	Depth    int
	BaudRate int
}

func (sm Uart_instance) String() string {
	return "uart:" + strconv.Itoa(sm.BaudRate) + ":" + strconv.Itoa(sm.Depth)
}

func (sm Uart_instance) Write_verilog(bmach *Bondmachine, soIndex int, uartName string, flavor string) string {

	result := ""

	// Compute the receivers and senders of the uart
	receivers := make([]string, 0)
	senders := make([]string, 0)

	for numProcessor, soList := range bmach.Shared_links {
		for _, soId := range soList {
			if soId == soIndex {
				for _, op := range bmach.Domains[bmach.Processors[numProcessor]].Op {
					switch op.Op_get_name() {
					case "t2r":
						receivers = append(receivers, "p"+strconv.Itoa(numProcessor)+"uart_recv")
						continue
					case "r2t":
						senders = append(senders, "p"+strconv.Itoa(numProcessor)+"uart_send")
						continue
					}
				}
			}
		}
	}

	s := bmstack.CreateBasicStack()
	s.ModuleName = uartName
	s.DataSize = int(bmach.Rsize)
	s.Depth = sm.Depth
	s.MemType = "LIFO"
	s.Senders = senders
	s.Receivers = receivers

	r, _ := s.WriteHDL()

	result += r

	return result

}

func (sm Uart_instance) GetPerProcPortsWires(bmach *Bondmachine, procId int, soId int, flavor string) string {
	result := ""
	if soName, ok := bmach.Get_so_name(soId); ok {
		for _, op := range bmach.Domains[bmach.Processors[procId]].Op {
			if op.Op_get_name() == "r2u" {
				result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(procId) + soName + "senderData;\n"
				result += "\twire p" + strconv.Itoa(procId) + soName + "senderWrite;\n"
				result += "\twire p" + strconv.Itoa(procId) + soName + "senderAck;\n"
				result += "\n"
				break
			}
		}
		for _, op := range bmach.Domains[bmach.Processors[procId]].Op {
			if op.Op_get_name() == "u2r" {
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

func (sm Uart_instance) GetPerProcPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soName, ok := bmach.Get_so_name(so_id); ok {
		for _, op := range bmach.Domains[bmach.Processors[proc_id]].Op {
			if op.Op_get_name() == "r2u" {
				result += ", p" + strconv.Itoa(proc_id) + soName + "senderData"
				result += ", p" + strconv.Itoa(proc_id) + soName + "senderWrite"
				result += ", p" + strconv.Itoa(proc_id) + soName + "senderAck"
				break
			}
		}
		for _, op := range bmach.Domains[bmach.Processors[proc_id]].Op {
			if op.Op_get_name() == "u2r" {
				result += ", p" + strconv.Itoa(proc_id) + soName + "receiverData"
				result += ", p" + strconv.Itoa(proc_id) + soName + "receiverRead"
				result += ", p" + strconv.Itoa(proc_id) + soName + "receiverAck"
				break
			}
		}
	}
	return result
}

func (sm Uart_instance) GetExternalPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Uart_instance) GetExternalPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Uart_instance) GetCPSharedPortsHeader(bmach *Bondmachine, soId int, flavor string) string {
	result := ""
	if soName, ok := bmach.Get_so_name(soId); ok {
		result += ", " + soName + "empty"
		result += ", " + soName + "full"
	}
	return result
}

func (sm Uart_instance) GetCPSharedPortsWires(bmach *Bondmachine, soId int, flavor string) string {
	result := ""
	if soName, ok := bmach.Get_so_name(soId); ok {
		result += "\n"
		result += "	wire " + soName + "empty;\n"
		result += "	wire " + soName + "full\n;"
		result += "\n"
	}
	return result
}
