package bondmachine

import (
	//"fmt"

	"strconv"
	"strings"
)

// The placeholder struct

type Queue struct{}

func (op Queue) Shr_get_name() string {
	return "queue"
}

func (op Queue) Shr_get_desc() string {
	return "Queue"
}

func (op Queue) Shortname() string {
	return "q"
}

func (op Queue) GV_config(element uint8) string {
	result := ""
	switch element {
	case GVNODEINPROC:
		result += "style=filled fillcolor=steelblue1 color=black"
	case GVNODE:
		result += "style=filled fillcolor=steelblue1 color=black"
	case GVEDGE:
		result += "arrowhead=none"
	case GVCLUS:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey70"
	case GVCLUSINPROC:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey70"
	}
	return result
}

func (op Queue) Instantiate(s string) (Shared_instance, bool) {
	if strings.HasPrefix(s, "queue:") {
		if len(s) > 6 {
			if depth, ok := strconv.Atoi(s[6:]); ok == nil {
				result := new(Queue_instance)
				result.Shared_element = op
				result.Depth = depth
				return *result, true
			}
		}
	}
	return nil, false
}

// The instance struct

type Queue_instance struct {
	Shared_element
	Depth int
}

func (sm Queue_instance) String() string {
	return "queue:" + strconv.Itoa(sm.Depth)
}

func (sm Queue_instance) Write_verilog(bmach *Bondmachine, so_index int, queue_name string, flavor string) string {

	result := ""

	return result

}

func (sm Queue_instance) GetPerProcPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
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

func (sm Queue_instance) GetPerProcPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
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

func (sm Queue_instance) GetExternalPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Queue_instance) GetExternalPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Queue_instance) GetCPSharedPortsHeader(bmach *Bondmachine, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Queue_instance) GetCPSharedPortsWires(bmach *Bondmachine, so_id int, flavor string) string {
	result := ""
	return result
}
