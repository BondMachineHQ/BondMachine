package bondmachine

import (
	//"fmt"

	"fmt"
	"strconv"
	"strings"
)

// The placeholder struct

type Stack struct{}

func (op Stack) Shr_get_name() string {
	return "stack"
}

func (op Stack) Shr_get_desc() string {
	return "Stack"
}

func (op Stack) Shortname() string {
	return "st"
}

func (op Stack) GV_config(element uint8) string {
	result := ""
	switch element {
	case GVNODEINPROC:
		result += "style=filled fillcolor=pink1 color=black"
	case GVNODE:
		result += "style=filled fillcolor=pink1 color=black"
	case GVEDGE:
		result += "arrowhead=none"
	case GVCLUS:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey70"
	case GVCLUSINPROC:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey70"
	}
	return result
}

func (op Stack) Instantiate(s string) (Shared_instance, bool) {
	if strings.HasPrefix(s, "stack:") {
		if len(s) > 6 {
			if depth, ok := strconv.Atoi(s[6:]); ok == nil {
				result := new(Stack_instance)
				result.Shared_element = op
				result.Depth = depth
				return *result, true
			}
		}
	}
	return nil, false
}

// The instance struct

type Stack_instance struct {
	Shared_element
	Depth int
}

func (sm Stack_instance) String() string {
	return "stack:" + strconv.Itoa(sm.Depth)
}

func (sm Stack_instance) Write_verilog(bmach *Bondmachine, soIndex int, stackName string, flavor string) string {

	result := ""

	receivers := make([]string, 0)
	senders := make([]string, 0)

	for numProcessor, soList := range bmach.Shared_links {
		for _, soId := range soList {
			if soId == soIndex {
				for _, op := range bmach.Domains[bmach.Processors[numProcessor]].Op {
					switch op.Op_get_name() {
					case "t2r":
						receivers = append(receivers, "p"+strconv.Itoa(numProcessor)+"stack_recv")
						continue
					case "r2t":
						senders = append(senders, "p"+strconv.Itoa(numProcessor)+"stack_send")
						continue
					}
				}
			}
		}
	}

	fmt.Println("Stack", stackName, "receivers", receivers, "senders", senders)

	return result

}

func (sm Stack_instance) GetPerProcPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(proc_id) + soname + "din;\n"
		result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(proc_id) + soname + "dout;\n"
		result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] p" + strconv.Itoa(proc_id) + soname + "addr;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "wren;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "en;\n"
		result += "\n"
	}
	return result
}

func (sm Stack_instance) GetPerProcPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += ", p" + strconv.Itoa(proc_id) + soname + "din"
		result += ", p" + strconv.Itoa(proc_id) + soname + "dout"
		result += ", p" + strconv.Itoa(proc_id) + soname + "addr"
		result += ", p" + strconv.Itoa(proc_id) + soname + "wren"
		result += ", p" + strconv.Itoa(proc_id) + soname + "en"
	}
	return result
}

func (sm Stack_instance) GetExternalPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Stack_instance) GetExternalPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}
