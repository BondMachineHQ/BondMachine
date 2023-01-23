package bondmachine

import (
	"strconv"
	"strings"
)

// The placeholder struct

type Barrier struct{}

func (op Barrier) Shr_get_name() string {
	return "barrier"
}

func (op Barrier) Shr_get_desc() string {
	return "Barrier"
}

func (op Barrier) Shortname() string {
	return "br"
}

func (op Barrier) GV_config(element uint8) string {
	result := ""
	switch element {
	case GVNODEINPROC:
		result += "style=filled fillcolor=orange color=black"
	case GVNODE:
		result += "style=filled fillcolor=orange color=black"
	case GVEDGE:
		result += "arrowhead=none"
	case GVCLUS:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey65"
	case GVCLUSINPROC:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey65"
	}
	return result
}

func (op Barrier) Instantiate(s string) (Shared_instance, bool) {
	if strings.HasPrefix(s, "barrier:") {
		if len(s) > 8 {
			if timeout, ok := strconv.Atoi(s[8:]); ok == nil {
				result := new(Barrier_instance)
				result.Shared_element = op
				result.Timeout = timeout
				return *result, true
			}
		}
	}
	return nil, false
}

// The instance struct

type Barrier_instance struct {
	Shared_element
	Timeout int
}

func (sm Barrier_instance) String() string {
	return "barrier:" + strconv.Itoa(sm.Timeout)
}

func (sm Barrier_instance) Write_verilog(bmach *Bondmachine, so_index int, barrier_name string, flavor string) string {

	result := ""

	subresult := ""

	num_processors := 0

	has_tout := true

	if sm.Timeout == 0 {
		has_tout = false
	}

	binarytout := strconv.FormatInt(int64(sm.Timeout), 2)
	toutlen := len(binarytout)
	toutzero := strconv.Itoa(toutlen) + "'b" + strings.Replace(binarytout, "1", "0", -1)

	orlist := ""
	andlist := ""

	for _, solist := range bmach.Shared_links {
		for _, so_id := range solist {
			if so_id == so_index {
				subresult += ", p" + strconv.Itoa(num_processors) + "hit"
				subresult += ", p" + strconv.Itoa(num_processors) + "ishitted"
				subresult += ", p" + strconv.Itoa(num_processors) + "tout"
				orlist += "p" + strconv.Itoa(num_processors) + "hit | "
				andlist += "p" + strconv.Itoa(num_processors) + "hit & "
				num_processors++
			}
		}
	}

	orlist = orlist[0 : len(orlist)-3]
	andlist = andlist[0 : len(andlist)-3]

	result += "`timescale 1ns/1ps\n"
	result += "module " + barrier_name + "(clk, reset" + subresult + ");\n"
	result += "\n"
	result += "	//--------------Input Ports-----------------------\n"
	result += "	input clk;\n"
	result += "	input reset;\n"

	subresult_in := ""
	subresult_out := ""
	num_processors = 0

	for _, solist := range bmach.Shared_links {
		for _, so_id := range solist {
			if so_id == so_index {
				subresult_in += "	input p" + strconv.Itoa(num_processors) + "hit;\n"
				subresult_out += "	output p" + strconv.Itoa(num_processors) + "ishitted;\n"
				subresult_out += "	output p" + strconv.Itoa(num_processors) + "tout;\n"
				num_processors++
			}
		}
	}

	result += subresult_in
	result += "\n"
	result += "	//--------------Output Ports-----------------------\n"
	result += subresult_out
	result += "\n"

	result += "	reg done;\n"
	if has_tout {
		result += "	reg timeout;\n"
		result += "	reg [" + strconv.Itoa(toutlen-1) + ":0] counter;\n"
	}
	result += "\n"

	result += "	initial begin\n"
	result += "		done = 1'b0;\n"
	if has_tout {
		result += "		timeout = 1'b0;\n"
		result += "		counter = " + toutzero + ";\n"
	}
	result += "	end\n"
	result += "\n"

	result += "	always @ (posedge clock) begin\n"
	if has_tout {
		result += "		if (done || timeout) begin\n"
	} else {
		result += "		if (done) begin\n"
	}
	result += "			done <= 1'b0;\n"
	if has_tout {
		result += "			timeout <= 1'b0;\n"
		result += "			counter <= " + toutzero + ";\n"
	}
	result += "		end\n"
	result += "	end\n"
	result += "\n"

	if has_tout {
		result += "	always @ (posedge clock) begin\n"
		result += "		if (!done & !timeout & (" + orlist + "))\n"
		result += "			counter <= counter + 1'b1;\n"
		result += "	end\n"
		result += "\n"
	}

	result += "	always @(posedge clock) begin\n"
	result += "		if (" + andlist + ")\n"
	result += "			done <= 1;\n"
	if has_tout {
		result += "		else begin\n"
		result += "			if (counter == " + strconv.Itoa(toutlen) + "'b" + binarytout + ")\n"
		result += "				timeout <= 1;\n"
		result += "		end\n"
	}
	result += "	end\n"
	result += "\n"

	num_processors = 0
	for _, solist := range bmach.Shared_links {
		for _, so_id := range solist {
			if so_id == so_index {
				result += "	assign p" + strconv.Itoa(num_processors) + "ishitted = done;\n"
				if has_tout {
					result += "	assign p" + strconv.Itoa(num_processors) + "tout = timeout;\n"
				}
				num_processors++
			}
		}
	}
	result += "\n"

	result += "endmodule\n"
	result += "\n"

	return result
}

func (sm Barrier_instance) GetPerProcPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "hit;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "ishitted;\n"
		result += "\twire p" + strconv.Itoa(proc_id) + soname + "tout;\n"
		result += "\n"
	}
	return result
}

func (sm Barrier_instance) GetPerProcPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	if soname, ok := bmach.Get_so_name(so_id); ok {
		result += ", p" + strconv.Itoa(proc_id) + soname + "hit"
		result += ", p" + strconv.Itoa(proc_id) + soname + "ishitted"
		result += ", p" + strconv.Itoa(proc_id) + soname + "tout"
	}
	return result
}

func (sm Barrier_instance) GetExternalPortsHeader(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Barrier_instance) GetExternalPortsWires(bmach *Bondmachine, proc_id int, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Barrier_instance) GetCPSharedPortsHeader(bmach *Bondmachine, so_id int, flavor string) string {
	result := ""
	return result
}

func (sm Barrier_instance) GetCPSharedPortsWires(bmach *Bondmachine, so_id int, flavor string) string {
	result := ""
	return result
}
