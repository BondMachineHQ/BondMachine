package bondmachine

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
	"github.com/BondMachineHQ/BondMachine/pkg/simbox"
)

func nth_assoc(assoc string, seq int) string {
	re := regexp.MustCompile("\\[(?P<to>[0-9]+):(?P<from>[0-9]+)\\] +(?P<name>[a-zA-Z0-9]+)")
	if re.MatchString(assoc) {
		to_string := re.ReplaceAllString(assoc, "${to}")
		from_string := re.ReplaceAllString(assoc, "${from}")
		name_string := re.ReplaceAllString(assoc, "${name}")
		to, _ := strconv.Atoi(to_string)
		from, _ := strconv.Atoi(from_string)
		step := 1
		if from > to {
			step = -1
		}
		pos := from + seq*step

		return name_string + "[" + strconv.Itoa(pos) + "]"
	}
	return assoc
}

func (bmach *Bondmachine) Write_verilog(conf *Config, flavor string, iomaps *IOmap, extramods []ExtraModule, sbox *simbox.Simbox) error {
	if len(bmach.Domains) != 0 {
		pConf := conf.ProcbuilderConfig()

		// Check if the bmapi module is present, needed for the exclusion of the bondmachine_main module from the accelerators
		var bmapiModuleAXIStream bool
		for _, mod := range extramods {
			if mod.Get_Name() == "bmapi" {
				bmapiParams := mod.Get_Params().Params
				if bmapiFlavor, ok := bmapiParams["bmapi_flavor"]; ok {
					switch bmapiFlavor {
					case "axist":
						bmapiModuleAXIStream = true
					}
				}
			}
		}

		//Instatiation of the Processor
		for i, dom_id := range bmach.Processors {

			ri := new(procbuilder.RuntimeInfo)
			ri.Init()

			pConf.Runinfo = ri
			dom := bmach.Domains[dom_id]

			sharedlist := ""
			solist := bmach.Shared_links[i]
			for j, so_id := range solist {
				sharedlist += bmach.Shared_objects[so_id].String()
				if j != len(solist)-1 {
					sharedlist += ","
				}
			}

			dom.Arch.Shared_constraints = sharedlist

			arch_filename := "arch_" + strconv.Itoa(i)
			arch_mod_name := "a" + strconv.Itoa(i)

			arch_names := map[string]string{"processor": "p" + strconv.Itoa(i), "rom": "p" + strconv.Itoa(i) + "rom", "ram": "p" + strconv.Itoa(i) + "ram"}

			// Set the arch tag to a deterministic value
			dom.Conproc.CpID = uint32(i)

			if _, err := os.Stat(arch_filename + ".v"); os.IsNotExist(err) {
				f, err := os.Create(arch_filename + ".v")
				check(err)
				//defer f.Close()
				_, err = f.WriteString(dom.Arch.Write_verilog(arch_mod_name, arch_names, flavor))
				check(err)
				f.Close()
			}

			if _, err := os.Stat(arch_names["processor"] + ".v"); os.IsNotExist(err) {
				f, err := os.Create(arch_names["processor"] + ".v")
				check(err)
				//defer f.Close()
				_, err = f.WriteString(dom.Arch.Conproc.Write_verilog(pConf, &dom.Arch, arch_names["processor"], flavor))
				check(err)
				f.Close()
			}
			if _, err := os.Stat(arch_names["rom"] + ".v"); os.IsNotExist(err) {
				f, err := os.Create(arch_names["rom"] + ".v")
				check(err)
				//defer f.Close()
				_, err = f.WriteString(dom.Arch.Rom.Write_verilog(dom, arch_names["rom"], flavor))
				check(err)
				f.Close()
			}

			if int(dom.L) != 0 {
				if _, err := os.Stat(arch_names["ram"] + ".v"); os.IsNotExist(err) {
					f, err := os.Create(arch_names["ram"] + ".v")
					check(err)
					//defer f.Close()
					_, err = f.WriteString(dom.Arch.Ram.Write_verilog(pConf, dom, arch_names["ram"], flavor))
					check(err)
					f.Close()
				}
			}
		}

		if len(bmach.Shared_objects) > 0 {

			seq := make(map[string]int)

			for i, so := range bmach.Shared_objects {
				sname := so.Shortname()
				if _, ok := seq[sname]; !ok {
					seq[sname] = 0
				}

				if _, err := os.Stat(sname + strconv.Itoa(seq[sname]) + ".v"); os.IsNotExist(err) {
					f, err := os.Create(sname + strconv.Itoa(seq[sname]) + ".v")
					check(err)
					defer f.Close()
					_, err = f.WriteString(so.Write_verilog(bmach, i, sname+strconv.Itoa(seq[sname]), flavor))
					check(err)
				}

				seq[sname]++
			}
		}

		if _, err := os.Stat("bondmachine.v"); os.IsNotExist(err) {
			f, err := os.Create("bondmachine.v")
			check(err)
			defer f.Close()
			_, err = f.WriteString(bmach.Write_verilog_main(conf, "bondmachine", flavor))
			check(err)
		}

		for _, mod := range extramods {
			files, filescode := mod.ExtraFiles()
			for i, file := range files {
				f, err := os.Create(file)
				check(err)
				_, err = f.WriteString(filescode[i])
				check(err)
				f.Close()
			}
		}

		switch flavor {
		case "iverilog_simulation", "iverilog":
			if _, err := os.Stat("bondmachine_tb.v"); os.IsNotExist(err) {
				f, err := os.Create("bondmachine_tb.v")
				check(err)
				defer f.Close()
				_, err = f.WriteString(bmach.Write_verilog_testbench("bondmachine", flavor, iomaps, extramods, sbox))
				check(err)
			}
		case "alveou50", "basys3", "kc705", "zedboard", "ebaz4205", "zc702", "ice40lp1k", "icefun", "icebreaker", "de10nano", "max1000", "icesugarnano":

			// Create the board file only if it doesn't belong to an AXIStream accelerator
			if !bmapiModuleAXIStream {
				if _, err := os.Stat("bondmachine_main.v"); os.IsNotExist(err) {
					f, err := os.Create("bondmachine_main.v")
					check(err)
					defer f.Close()
					_, err = f.WriteString(bmach.Write_verilog_board(conf, "bondmachine", flavor, iomaps, extramods))
					check(err)
				}
			}
		default:
			return Prerror{"Verilog flavor unknown"}
		}
	} else {
		return Prerror{"No defined domains"}
	}
	return nil
}

func (bmach *Bondmachine) Write_verilog_main(conf *Config, module_name string, flavor string) string {

	result := ""
	result += "module " + module_name + "(clk, reset"

	// The External_inputs connected are defined in the module as port
	for i := 0; i < bmach.Inputs; i++ {
		result += ", i" + strconv.Itoa(i)
		result += ", i" + strconv.Itoa(i) + "_valid"
		result += ", i" + strconv.Itoa(i) + "_received"
	}

	// The External_inputs connected are defined in the module as port
	for i := 0; i < bmach.Outputs; i++ {
		result += ", o" + strconv.Itoa(i)
		result += ", o" + strconv.Itoa(i) + "_valid"
		result += ", o" + strconv.Itoa(i) + "_received"
	}

	// The External ports defined in Shared objects
	if len(bmach.Shared_objects) > 0 {
		for proc_id, solist := range bmach.Shared_links {
			for _, so_id := range solist {
				//fmt.Println(proc_id, so_id)
				result += bmach.Shared_objects[so_id].GetExternalPortsHeader(bmach, proc_id, so_id, flavor)
			}
		}
	}

	result += ");\n\n"

	if conf.CommentedVerilog {
		result += "\t// Clock and reset input ports\n"
	}
	result += "\tinput clk, reset;\n"

	// The External_inputs connected are defined as input port
	for i := 0; i < bmach.Inputs; i++ {
		result += "	input [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] i" + strconv.Itoa(i) + ";\n"
		result += "	input i" + strconv.Itoa(i) + "_valid;\n"
		result += "	output i" + strconv.Itoa(i) + "_received;\n"
	}

	result += "	//--------------Output Ports-----------------------\n"

	// The External_inputs connected are defined as input port
	for i := 0; i < bmach.Outputs; i++ {
		result += "	output [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] o" + strconv.Itoa(i) + ";\n"
		result += "	output o" + strconv.Itoa(i) + "_valid;\n"
		result += "	input o" + strconv.Itoa(i) + "_received;\n"
	}

	// The External ports defined in Shared objects
	if len(bmach.Shared_objects) > 0 {
		for proc_id, solist := range bmach.Shared_links {
			for _, so_id := range solist {
				//fmt.Println(proc_id, so_id)
				result += bmach.Shared_objects[so_id].GetExternalPortsWires(bmach, proc_id, so_id, flavor)
			}
		}
	}

	result += "\n\n"

	// The Internal_inputs connected are wire the unconnected are registers, the machine outputs are always wire handled with assign
	for i, linked := range bmach.Links {
		if linked == -1 {
			result += "	wire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] " + strings.ToLower(bmach.Internal_inputs[i].String()) + ";\n"
		}
	}

	result += "\n"

	// The Internal_outputs are always wire except the machine inputs, but here we write only the unconnected
	for ioID, bond := range bmach.Internal_outputs {
		if conf.CommentedVerilog {
			result += "\t//Analyzing Internal output " + strings.ToLower(bond.String()) + "\n"
			for iiID, linked := range bmach.Links {
				if linked == ioID {
					result += "\t//Internal output " + strings.ToLower(bond.String()) + " is connected to " + strings.ToLower(bmach.Internal_inputs[iiID].String()) + "\n"
				}
			}
		}
		if bond.Map_to == 0 {
			//TODO Check if every possibility is taken into account
			result += "	wire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] " + strings.ToLower(bond.String()) + ";\n"
			result += "	wire " + strings.ToLower(bond.String()) + "_valid;\n"
			result += "	wire " + strings.ToLower(bond.String()) + "_received;\n"
		} else {
			//TODO Check if every possibility is taken into account
			result += "	wire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] " + strings.ToLower(bond.String()) + ";\n"
			result += "	wire " + strings.ToLower(bond.String()) + "_valid;\n"
			result += "	wire " + strings.ToLower(bond.String()) + "_received;\n"
		}

		// Count the number of internal inputs connected to me
		var connected int = 0
		for _, linked := range bmach.Links {
			if linked == ioID {
				connected++
			}
		}

		// Create wires for the internal inputs connected to me
		if connected > 0 {
			for iiID, linked := range bmach.Links {
				if linked == ioID {
					//result += "	wire " + strings.ToLower(bond.String()) + "_received_from_" + strings.ToLower(bmach.Internal_inputs[iiID].String()) + ";\n"
					result += "	wire " + strings.ToLower(bmach.Internal_inputs[iiID].String()) + "_received;\n"
				}
			}
		}

	}

	result += "\n"

	// The internal wire for the connection to the SOs
	if len(bmach.Shared_objects) > 0 {
		for procId, soList := range bmach.Shared_links {
			for _, soId := range soList {
				//fmt.Println(proc_id, so_id)
				result += bmach.Shared_objects[soId].GetPerProcPortsWires(bmach, procId, soId, flavor)
			}
		}

		for siId, so := range bmach.Shared_objects {
			result += so.GetCPSharedPortsWires(bmach, siId, flavor)
		}
	}

	result += "\n"

	result += "	//Instantiation of the Processors and Shared Objects\n"

	//Instantiation of the Processor
	for i, dom_id := range bmach.Processors {
		dom := bmach.Domains[dom_id]
		arch := dom.Arch

		arch_mod_name := "a" + strconv.Itoa(i)
		arch_instance_name := arch_mod_name + "_inst"

		result += "	" + arch_mod_name + " " + arch_instance_name + "(clk, reset"

		// A Processor input is connected to a register with its name (if there is no bond) or to an internal output
		// Default: the name is the name of the internal input side if the bond
		for j := 0; j < int(arch.N); j++ {
			map_to := uint8(2)
			res_id := i
			ext_id := j

			for k, linked := range bmach.Links {
				bond := bmach.Internal_inputs[k]
				if bond.Map_to == map_to && bond.Res_id == res_id && bond.Ext_id == ext_id {
					if linked == -1 {
						result += ", " + strings.ToLower(bond.String()) + ", " + strings.ToLower(bond.String()) + "_valid, " + strings.ToLower(bond.String()) + "_received"
					} else {
						result += ", " + strings.ToLower(bmach.Internal_outputs[linked].String()) + ", " + strings.ToLower(bmach.Internal_outputs[linked].String()) + "_valid, "

						//result += strings.ToLower(bmach.Internal_outputs[linked].String()) + "_received_from_" + strings.ToLower(bond.String())
						result += strings.ToLower(bond.String()) + "_received"
					}
					break
				}
			}
		}

		// A Processor output always use its name connected whether or not
		for j := 0; j < int(arch.M); j++ {
			map_to := uint8(3)
			res_id := i
			ext_id := j

			bond := Bond{map_to, res_id, ext_id}
			result += ", " + strings.ToLower(bond.String()) + ", " + strings.ToLower(bond.String()) + "_valid, " + strings.ToLower(bond.String()) + "_received"
		}

		//memorize the name of the shared object
		// The module for the connection to the shared memory and the channel
		if len(bmach.Shared_objects) > 0 {
			soList := bmach.Shared_links[i]
			for _, soId := range soList {
				result += bmach.Shared_objects[soId].GetPerProcPortsHeader(bmach, i, soId, flavor)
				result += bmach.Shared_objects[soId].GetCPSharedPortsHeader(bmach, soId, flavor)
			}
		}
		result += ");\n"
	}

	//Instantiation of the Shared object
	if len(bmach.Shared_objects) > 0 {

		seq := make(map[string]int)

		for i, so := range bmach.Shared_objects {
			sname := so.Shortname()
			if _, ok := seq[sname]; !ok {
				seq[sname] = 0
			}

			result += "	" + sname + strconv.Itoa(seq[sname]) + " " + sname + strconv.Itoa(seq[sname]) + "_inst (clk, reset"

			for procId, solist := range bmach.Shared_links {
				for _, soId := range solist {
					if soId == i {
						result += bmach.Shared_objects[soId].GetPerProcPortsHeader(bmach, procId, soId, flavor)
						// result += bmach.Shared_objects[soId].GetExternalPortsHeader(bmach, procId, soId, flavor)
					}
				}
			}

			for procId, solist := range bmach.Shared_links {
				for _, soId := range solist {
					if soId == i {
						// result += bmach.Shared_objects[soId].GetPerProcPortsHeader(bmach, procId, soId, flavor)
						result += bmach.Shared_objects[soId].GetExternalPortsHeader(bmach, procId, soId, flavor)
					}
				}
			}
			// Include SO ports eventually shared with CORES (not core dependent)
		sharedLoop:
			for _, soList := range bmach.Shared_links {
				for _, soId := range soList {
					if soId == i {
						result += bmach.Shared_objects[soId].GetCPSharedPortsHeader(bmach, soId, flavor)
						break sharedLoop
					}
				}
			}

			result += ");\n"
			seq[sname]++
		}
	}

	result += "\n"

	for i, linked := range bmach.Links {
		if linked != -1 {
			if bmach.Internal_inputs[i].Map_to == uint8(1) {
				result += "	assign " + strings.ToLower(bmach.Internal_inputs[i].String()) + " = " + strings.ToLower(bmach.Internal_outputs[linked].String()) + ";\n"
				result += "	assign " + strings.ToLower(bmach.Internal_inputs[i].String()) + "_valid = " + strings.ToLower(bmach.Internal_outputs[linked].String()) + "_valid;\n"
				//result += "	assign " + strings.ToLower(bmach.Internal_outputs[linked].String()) + "_received = " + strings.ToLower(bmach.Internal_inputs[i].String()) + "_received;\n"
			}
		}
	}

	result += "\n"

	for ioID, bond := range bmach.Internal_outputs {

		// Count the number of internal inputs connected to me
		var connected int = 0
		for _, linked := range bmach.Links {
			if linked == ioID {
				connected++
			}
		}

		// Create assigns for the internal inputs connected to me
		if connected == 1 {
			for iiID, linked := range bmach.Links {
				if linked == ioID {
					result += "	assign " + strings.ToLower(bond.String()) + "_received = " + strings.ToLower(bmach.Internal_inputs[iiID].String()) + "_received;\n"
				}
			}
		} else if connected > 1 {
			result += "	assign " + strings.ToLower(bond.String()) + "_received = ( 1'b1 \n"
			for iiID, linked := range bmach.Links {
				if linked == ioID {
					result += "		& (" + strings.ToLower(bmach.Internal_inputs[iiID].String()) + "_received) \n"
				}
			}
			result += "		);\n"
		}
	}

	result += "\n"

	result += "endmodule\n"

	return result

}

func (bmach *Bondmachine) Write_verilog_testbench(module_name string, flavor string, iomaps *IOmap, extramods []ExtraModule, sbox *simbox.Simbox) string {

	result := ""
	result += `
module request(clk, reset, req, ack, impulse);
    input clk;
    input reset;
    output reg req;
    input ack;
    input impulse;

    reg state;

    initial begin
        state = 0;
        req = 0;
    end

    always @(posedge clk) begin
        if (reset) begin
            state <= 0;
        end else begin
            case (state)
                0: begin
                    req <= 0;
                    if (impulse) begin
                        state <= 1;
                    end
                end
                1: begin
                    req <= 1;
                    if (ack) begin
                        state <= 0;
                    end
                end
            endcase
        end
    end

endmodule

module unlock(clk, reset, valid, received);
    input clk;
    input reset;
    input valid;
    output reg received;

    always @(posedge clk) begin
        if (reset) begin
	    received <= 0;
	end else begin
	    if (valid) begin
		received <= 1;
	    end
	    else begin
		received <= 0;
	    end
	end
	end
endmodule

`

	result += "module " + module_name + "_tb;\n"
	result += "\n"
	result += "	reg clk, reset;\n"
	result += "\n"

	// The External_inputs connected are defined as input port
	for i := 0; i < bmach.Inputs; i++ {
		result += "	reg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] i" + strconv.Itoa(i) + ";\n"
		result += "	wire i" + strconv.Itoa(i) + "_valid;\n"
		result += "	wire i" + strconv.Itoa(i) + "_received;\n"
		result += "	reg i" + strconv.Itoa(i) + "_impulse;\n"
	}

	// The External_inputs connected are defined as input port
	for i := 0; i < bmach.Outputs; i++ {
		result += "	wire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] o" + strconv.Itoa(i) + ";\n"
		result += "	wire o" + strconv.Itoa(i) + "_valid;\n"
		result += "	wire o" + strconv.Itoa(i) + "_received;\n"
	}

	result += "\n"

	result += "\t" + module_name + " " + module_name + "_inst " + "(clk, reset"

	// The External_inputs connected are defined as input port
	for i := 0; i < bmach.Inputs; i++ {
		result += ", i" + strconv.Itoa(i)
		result += ", i" + strconv.Itoa(i) + "_valid"
		result += ", i" + strconv.Itoa(i) + "_received"
	}

	// The External_inputs connected are defined as input port
	for i := 0; i < bmach.Outputs; i++ {
		result += ", o" + strconv.Itoa(i)
		result += ", o" + strconv.Itoa(i) + "_valid"
		result += ", o" + strconv.Itoa(i) + "_received"
	}

	result += ");\n\n"

	for i := 0; i < bmach.Inputs; i++ {
		result += "	request i" + strconv.Itoa(i) + "_req(\n"
		result += "		.clk(clk),\n"
		result += "		.reset(reset),\n"
		result += "		.req(i" + strconv.Itoa(i) + "_valid),\n"
		result += "		.ack(i" + strconv.Itoa(i) + "_received),\n"
		result += "		.impulse(i" + strconv.Itoa(i) + "_impulse)\n"
		result += "	);\n"
	}

	for i := 0; i < bmach.Outputs; i++ {
		result += "	unlock o" + strconv.Itoa(i) + "_unlock(\n"
		result += "		.clk(clk),\n"
		result += "		.reset(reset),\n"
		result += "		.valid(o" + strconv.Itoa(i) + "_valid),\n"
		result += "		.received(o" + strconv.Itoa(i) + "_received)\n"
		result += "	);\n"
	}

	result += "\talways #1 clk = ~clk;\n"

	result += "\tinitial  begin\n"
	result += "\t\t$dumpfile (\"working_dir/bondmachine.vcd\");\n"
	result += "\t\t$dumpvars;\n"
	result += "\tend\n"

	result += "\tinitial begin\n"
	result += "\t\tclk = 1'b0;\n"
	result += "\t\treset = 1'b1;\n"

	for i := 0; i < bmach.Inputs; i++ {
		result += "\t\ti" + strconv.Itoa(i) + " = 1'b0;\n"
	}

	result += "\t\t#100;\n\n"
	result += "\t\treset = 1'b0;\n"

	for _, rule := range sbox.Rules {
		if rule.Timec == simbox.TIMEC_ABS && rule.Action == simbox.ACTION_SET {
			result += "\t\t#" + strconv.Itoa(int(rule.Tick)) + ";\n"
			result += "\t\t" + rule.Object + " = " + rule.Extra + ";\n"
			result += "\t\t" + rule.Object + "_impulse = 1'b1;\n"
			result += "\t\t#4 " + rule.Object + "_impulse = 1'b0;\n"
		}
	}

	result += "\n"
	result += "\t\t#100000;\n"
	result += "\t\t$finish;\n"
	result += "\tend\n"
	result += "endmodule\n"
	return result
}

func (bmach *Bondmachine) Write_verilog_basys3_7segment(module_name string, flavor string, iomaps *IOmap, extramods []ExtraModule) (string, error) {
	for _, mod := range extramods {
		if mod.Get_Name() == "basys3_7segment" {

			result := "\n"
			result += "module bond2seg(\n"
			result += "	input clk,\n"
			result += "	input reset,\n"
			result += "	input [15:0] value,\n"
			result += "	output [6:0] segment,\n"
			result += "	output enable_D1,\n"
			result += "	output enable_D2,\n"
			result += "	output enable_D3,\n"
			result += "	output enable_D4,\n"
			result += "	output dp\n"
			result += ");\n"
			result += "\n"
			result += "wire clk_point1hz;\n"
			result += "wire refreshClk;\n"
			result += "reg [3:0] hex;\n"
			result += "reg [3:0] reg_d0;\n"
			result += "reg [3:0] reg_d1;\n"
			result += "reg [3:0] reg_d2;\n"
			result += "reg [3:0] reg_d3;\n"
			result += "\n"
			result += "clkgen Uclkgen(\n"
			result += ".clk(clk),\n"
			result += ".refreshClk(refreshClk),\n"
			result += ".clk_point1hz(clk_point1hz)\n"
			result += ");\n"
			result += "\n"
			result += "enable_sr Uenable(\n"
			result += ".refreshClk(refreshClk),\n"
			result += ".enable_D1(enable_D1),\n"
			result += ".enable_D2(enable_D2),\n"
			result += ".enable_D3(enable_D3),\n"
			result += ".enable_D4(enable_D4)\n"
			result += ");\n"
			result += "\n"
			result += "ssd Ussd(\n"
			result += ".hex(hex),\n"
			result += ".segment(segment),\n"
			result += ".dp(dp)\n"
			result += ");\n"
			result += "\n"
			result += "always @(posedge clk) begin\n"
			result += "    reg_d0 <= value[3:0];\n"
			result += "    reg_d1 <= value[7:4];\n"
			result += "    reg_d2 <= value[11:8];\n"
			result += "    reg_d3 <= value[15:12];\n"
			result += "end\n"
			result += "\n"
			result += "always @ (*)\n"
			result += "\n"
			result += "case ({enable_D1,enable_D2,enable_D3,enable_D4})\n"
			result += "    4'b0111: hex = reg_d0;\n"
			result += "    4'b1011: hex = reg_d1;\n"
			result += "    4'b1101: hex = reg_d2;\n"
			result += "    4'b1110: hex = reg_d3;\n"
			result += "    default: hex = 0; \n"
			result += "endcase \n"
			result += "\n"
			result += "endmodule\n"
			result += "\n"
			result += "module clkgen(\n"
			result += "	input     clk, \n"
			result += "	output    refreshClk,\n"
			result += "	output    clk_point1hz\n"
			result += ");\n"
			result += "\n"
			result += "reg [26:0] 	count = 0;\n"
			result += "reg [16:0] 	refresh = 0;\n"
			result += "\n"
			result += "\n"
			result += "reg      	tmp_clk = 0;\n"
			result += "reg 		rclk = 0;\n"
			result += "\n"
			result += "\n"
			result += "assign clk_point1hz = tmp_clk;\n"
			result += "assign refreshClk = rclk;\n"
			result += "\n"
			result += "\n"
			result += "BUFG clock_buf_0(\n"
			result += "  .I(clk),\n"
			result += "  .O(clk_100mhz)\n"
			result += ");\n"
			result += "\n"
			result += "always @(posedge clk_100mhz) begin\n"
			result += "  if (count < 10000000) begin \n"
			result += "    count <= count + 1;\n"
			result += "  end\n"
			result += "  else begin\n"
			result += "    tmp_clk <= ~tmp_clk;\n"
			result += "    count <= 0;\n"
			result += "  end\n"
			result += "end\n"
			result += "\n"
			result += "always @(posedge clk_100mhz) begin\n"
			result += "	if (refresh < 100000) begin\n"
			result += "		refresh <= refresh + 1;\n"
			result += "	end else begin\n"
			result += "		refresh <= 0;\n"
			result += "		rclk <= ~rclk;\n"
			result += "	end\n"
			result += "end\n"
			result += "\n"
			result += "endmodule\n"
			result += "\n"
			result += "module enable_sr(\n"
			result += "	input         refreshClk,\n"
			result += "	output        enable_D1,\n"
			result += "	output        enable_D2,\n"
			result += "	output        enable_D3,\n"
			result += "	output        enable_D4\n"
			result += ");\n"
			result += "\n"
			result += "reg [3:0] pattern = 4'b0111; \n"
			result += "\n"
			result += "assign enable_D1 = pattern[3];\n"
			result += "assign enable_D2 = pattern[2];\n"
			result += "assign enable_D3 = pattern[1];\n"
			result += "assign enable_D4 = pattern[0];\n"
			result += "\n"
			result += "always @(posedge refreshClk) begin\n"
			result += "	pattern <= {pattern[0],pattern[3:1]};	\n"
			result += "end\n"
			result += "\n"
			result += "\n"
			result += "\n"
			result += "endmodule\n"
			result += "\n"
			result += "module ssd(\n"
			result += "	input [3:0] hex,\n"
			result += "	output reg [6:0] segment,\n"
			result += "	output dp\n"
			result += ");\n"
			result += "\n"
			result += "always @ (*)\n"
			result += "	case (hex) \n"
			result += "		0: segment = 7'b0000001;\n"
			result += "		1: segment = 7'b1001111;\n"
			result += "		2: segment = 7'b0010010;\n"
			result += "		3: segment = 7'b0000110;\n"
			result += "		4: segment = 7'b1001100;\n"
			result += "		5: segment = 7'b0100100;\n"
			result += "		6: segment = 7'b0100000;\n"
			result += "		7: segment = 7'b0001101;\n"
			result += "		8: segment = 7'b0000000;\n"
			result += "		9: segment = 7'b0000100;\n"
			result += "		10: segment = 7'b0001000;\n"
			result += "		11: segment = 7'b1100000;\n"
			result += "		12: segment = 7'b0110001;\n"
			result += "		13: segment = 7'b1000010;\n"
			result += "		14: segment = 7'b0110000;\n"
			result += "		15: segment = 7'b0111000;\n"
			result += "		default: segment = 7'b0000001;\n"
			result += "	endcase	\n"
			result += "assign dp = 4'b1111;\n"
			result += "    \n"
			result += "endmodule\n"
			result += "\n"

			return result, nil

		}
	}

	return "", errors.New("No basys3_7segment module found")
}

func (bmach *Bondmachine) WriteVerilogVgaText800x600(module_name string, flavor string, iomaps *IOmap, extramods []ExtraModule) (string, error) {
	for _, mod := range extramods {
		if mod.Get_Name() == "vga800x600" {
			result := ""
			result += `
module romfonts #(parameter ADDR_WIDTH=8, DATA_WIDTH=8, DEPTH=256, FONTSFILE="") (
    input wire [ADDR_WIDTH-1:0] addr, 
    output wire [DATA_WIDTH-1:0] data 
    );

    reg [DATA_WIDTH-1:0] fonts_array [0:DEPTH-1]; 

    initial begin
        if (FONTSFILE > 0)
        begin
            $display("Loading memory init file '" + FONTSFILE + "' into array.");
            $readmemh(FONTSFILE, fonts_array);
        end
    end

	assign data = fonts_array[addr];

endmodule

`
			return result, nil
		}
	}
	return "", errors.New("No ps2keyboard module found")
}

func (bmach *Bondmachine) WriteVerilogBMAPI(module_name string, flavor string, iomaps *IOmap, extramods []ExtraModule) (string, error) {
	for _, mod := range extramods {
		if mod.Get_Name() == "bmapi" {
			bmapiParams := mod.Get_Params().Params

			templatedata := bmach.createBasicTemplateData()
			templatedata.Inputs = make([]string, 0)
			templatedata.Outputs = make([]string, 0)
			templatedata.InputsBins = make([]string, 0)
			templatedata.OutputsBins = make([]string, 0)

			sendObj := make([]string, 0)

			sendOps := 1
			sendObj = append(sendObj, "SENDSTM_HANDSH")

			sortedKeys := make([]string, 0)
			for param, _ := range bmapiParams {
				sortedKeys = append(sortedKeys, param)
			}
			sort.Strings(sortedKeys)

			for _, param := range sortedKeys {
				if strings.HasPrefix(param, "assoc") {
					bmport := strings.Split(param, "_")[1]
					if strings.HasPrefix(bmport, "o") {
						onum, _ := strconv.Atoi(bmport[1:])
						templatedata.Outputs = append(templatedata.Outputs, "port_"+bmport)
						templatedata.OutputsBins = append(templatedata.OutputsBins, zeros_prefix(5, get_binary(onum)))
						sendObj = append(sendObj, "SENDSTM_port_"+bmport)
						sendObj = append(sendObj, "SENDSTM_port_"+bmport+"_valid")
						sendOps += 2
					} else if strings.HasPrefix(bmport, "i") {
						inum, _ := strconv.Atoi(bmport[1:])
						templatedata.Inputs = append(templatedata.Inputs, "port_"+bmport)
						templatedata.InputsBins = append(templatedata.InputsBins, zeros_prefix(5, get_binary(inum)))
						sendObj = append(sendObj, "SENDSTM_port_"+bmport+"_recv")
						sendOps += 1
					} else {
						return "", errors.New("Unknown port")
					}
				}
			}
			sendOps += 1
			sendObj = append(sendObj, "SENDSTM_KEEP")

			sendBits := Needed_bits(sendOps)
			sendBusWidth := "[" + strconv.Itoa(sendBits-1) + ":0]"
			sendBin := make([]string, 0)
			for i, _ := range sendObj {
				sendBin = append(sendBin, zeros_prefix(sendBits, get_binary(i)))
			}

			sendSMData := stateMachine{Nums: sendOps, Bits: sendBits, Names: sendObj, Binary: sendBin, Buswidth: sendBusWidth}

			templatedata.SendSM = sendSMData

			var f bytes.Buffer

			switch bmapiParams["bmapi_flavor"] {
			case "uartusb":
				t, err := template.New("bmapiuarttransceiver").Funcs(templatedata.funcmap).Parse(bmapiuarttransceiver)
				if err != nil {
					return "", err
				}

				err = t.Execute(&f, *templatedata)
				if err != nil {
					return "", err
				}
			case "aximm":
				t, err := template.New("bmapiaximmtransceiver").Funcs(templatedata.funcmap).Parse(bmapiaximmtransceiver)
				if err != nil {
					return "", err
				}

				err = t.Execute(&f, *templatedata)
				if err != nil {
					return "", err
				}
			}
			result := f.String()
			return result, nil
		}
	}
	return "", errors.New("No bmapi module found")
}

func (bmach *Bondmachine) WriteVerilogPs2Keyboard(module_name string, flavor string, iomaps *IOmap, extramods []ExtraModule) (string, error) {
	for _, mod := range extramods {
		if mod.Get_Name() == "ps2keyboard" {
			return "", nil
		}
	}
	return "", errors.New("No ps2keyboard module found")
}

func (bmach *Bondmachine) Write_verilog_etherbond(module_name string, flavor string, iomaps *IOmap, extramods []ExtraModule) (string, error) {
	for _, mod := range extramods {
		if mod.Get_Name() == "etherbond" {
			etherbond_params := mod.Get_Params().Params
			fmt.Println(etherbond_params)

			cluster_id, _ := strconv.Atoi(etherbond_params["cluster_id"])
			peer_id, _ := strconv.Atoi(etherbond_params["peer_id"])

			result := ""
			result += "module " + module_name + "_main(\n"
			result += "\n"

			clk_name := "clk"
			rst_name := "reset"

			result += "\tinput " + clk_name + ",\n"
			result += "\tinput " + rst_name + ",\n"
			result += "\toutput sck,\n"
			result += "\toutput mosi,\n"
			result += "\toutput cs_n,\n"
			result += "\tinput miso,\n"
			result += "\tinput INT_n,\n"

			inames := strings.Split(etherbond_params["inputs"], ",")
			iids := strings.Split(etherbond_params["input_ids"], ",")

			onames := strings.Split(etherbond_params["outputs"], ",")
			oids := strings.Split(etherbond_params["output_ids"], ",")
			odests := strings.Split(etherbond_params["destinations"], ",")

			for _, iname := range inames {
				if iname != "" {
					result += "\toutput [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] " + iname + ",\n"
				}
			}

			for _, oname := range onames {
				if oname != "" {
					result += "\tinput [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] " + oname + ",\n"
				}
			}

			result = result[0:len(result)-2] + "\n);\n\n"

			result += "\n"

			result += "\t//Etherbond type\n"
			result += "\tlocalparam [15:0] ETHERBOND_TYPE = 16'h8888;\n"

			result += "\n"

			result += "\t//Etherbond commands\n"
			result += "\tlocalparam [7:0] ADV_CLU = 8'h01;\n"
			result += "\tlocalparam [7:0] ADV_CH = 8'h02;\n"
			result += "\tlocalparam [7:0] ADV_IN = 8'h03;\n"
			result += "\tlocalparam [7:0] ADV_OUT = 8'h04;\n"
			result += "\tlocalparam [7:0] IO_TR = 8'h05;\n"
			result += "\tlocalparam [7:0] TAGACK = 8'hff;\n"

			result += "\n"

			result += "\t//packet definition for ethernet transmission\n"
			result += "\tlocalparam [15:0] ethertype = 16'h8888;\n"
			result += "\tlocalparam [47:0] mymac = 48'h0288" + fmt.Sprintf("%08d", peer_id) + ";\n"
			result += "\tlocalparam [31:0] mycluster_id = 32'd" + fmt.Sprintf("%08d", cluster_id) + ";\n"
			result += "\tlocalparam [31:0] mypeer_id = 32'd" + fmt.Sprintf("%08d", peer_id) + ";\n"
			result += "\treg [31:0] tag;\n"

			result += "\n"

			// 1 Advertises
			sendwires := 0

			for i, iname := range inames {
				if iname != "" {
					tmpi, _ := strconv.Atoi(iids[i])
					result += "\tlocalparam [31:0] id_resource_" + iname + " = 32'd" + fmt.Sprintf("%08d", tmpi) + ";\n"
					sendwires++
				}
			}

			result += "\n"

			for i, oname := range onames {
				if oname != "" {
					sendwires++
					tmpi, _ := strconv.Atoi(oids[i])
					result += "\tlocalparam [31:0] id_resource_" + oname + " = 32'd" + fmt.Sprintf("%08d", tmpi) + ";\n"
					ods := strings.Split(odests[i], "-")
					for j, oid := range ods {
						oid_n, _ := strconv.Atoi(oid)
						result += "\tlocalparam [31:0] " + oname + "_dest_" + strconv.Itoa(j) + "_id = 32'd" + fmt.Sprintf("%08d", oid_n) + ";\n"
						if mac, ok := etherbond_params["peer_"+oid+"_mac"]; ok {
							if mac == "auto" {
								result += "\tlocalparam [47:0] " + oname + "_dest_" + strconv.Itoa(j) + "_mac = 48'h0288" + fmt.Sprintf("%08d", oid_n) + ";\n"
							} else if mac == "adv" {
								// TODO just the case here eventually implemented logic has to be written
								result += "\treg [31:0] " + oname + "_dest_" + strconv.Itoa(j) + "_mac;\n"
							} else {
								result += "\tlocalparam [47:0] " + oname + "_dest_" + strconv.Itoa(j) + "_mac = 48'h" + mac + ";\n"
							}
						} else {
							result += "\tlocalparam [47:0] " + oname + "_dest_" + strconv.Itoa(j) + "_mac = 48'h0288" + fmt.Sprintf("%08d", oid_n) + ";\n"
						}
						sendwires++
					}
				}
			}

			result += "\n"

			for _, oname := range onames {
				if oname != "" {
					result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] " + oname + "_copy;\n"
					result += "\tassign " + oname + "_copy = " + oname + ";\n"
				}
			}

			result += "\n"

			result += "\tlocalparam LASTFRAME = " + strconv.Itoa(Needed_bits(sendwires+2)) + "'d" + strconv.Itoa(sendwires+1) + ";\n"

			framesize := 256 + 24 + int(bmach.Rsize)

			result += "\twire [" + strconv.Itoa(framesize-1) + ":0] ethtxs [0:" + strconv.Itoa(sendwires) + "];\n"

			iw := 1
			result += "\tassign ethtxs[0] = {8'hFF, 48'hFFFFFFFFFFFF, mymac, ethertype, ADV_CLU, mycluster_id, mypeer_id, " + strconv.Itoa(framesize-192) + "'b0};\n"
			for _, iname := range inames {
				if iname != "" {
					result += "\tassign ethtxs[" + strconv.Itoa(iw) + "] = {8'hFF, 48'hFFFFFFFFFFFF, mymac, ethertype, ADV_IN, mycluster_id, mypeer_id, id_resource_" + iname + ", " + strconv.Itoa(framesize-224) + "'b0};\n"
					iw++
				}
			}

			for i, oname := range onames {
				if oname != "" {
					result += "\tassign ethtxs[" + strconv.Itoa(iw) + "] = {8'hFF, 48'hFFFFFFFFFFFF, mymac, ethertype, ADV_OUT, mycluster_id, mypeer_id, id_resource_" + oname + ", " + strconv.Itoa(framesize-224) + "'b0};\n"
					iw++
					ods := strings.Split(odests[i], "-")
					for j, _ := range ods {
						dname := oname + "_dest_" + strconv.Itoa(j) + "_mac"
						result += "\tassign ethtxs[" + strconv.Itoa(iw) + "] = {8'hFF, " + dname + ", mymac, ethertype, IO_TR , tag,  mycluster_id, mypeer_id, id_resource_" + oname + ", " + oname + "_copy, 24'b0};\n"
						iw++
					}
				}
			}

			result += "\twire RESET_n = !reset;\n"

			result += "\n"

			result += "\t// Logic to start the configuration phase\n"
			result += "\t(* keep=\"true\" *) wire end_configuration;\n"
			result += "\treg start_conf;\n"
			result += "\treg [9:0] cnt4cnf;\n"
			result += "\n"
			result += "\talways@(posedge clk)\n"
			result += "\tbegin\n"
			result += "\t\tif(RESET_n==1'b0)\n"
			result += "\t\t\tcnt4cnf <= 32'b0;\n"
			result += "\t\telse\n"
			result += "\t\t\tcnt4cnf <= cnt4cnf + 1'b1;\n"
			result += "\tend\n"

			result += "\talways@(posedge clk)\n"
			result += "\tbegin\n"
			result += "\t\tif(RESET_n==1'b0 | end_configuration)\n"
			result += "\t\t\tstart_conf <= 1'b0;\n"
			result += "\t\telse if(cnt4cnf==9'h1FF)\n"
			result += "\t\t\tstart_conf <= 1'b1;\n"
			result += "\tend\n"

			result += "\n"

			result += "\t// Frame explorer\n"
			result += "\t(* keep=\"true\" *) reg [" + strconv.Itoa(Needed_bits(sendwires+1)) + ":0] frame_explorer;\n"
			result += "\t(* keep=\"true\" *) reg [28:0] frame_counter;\n"

			for i := 0; i < sendwires+1; i++ {
				result += "\treg frame_ready_" + strconv.Itoa(i) + ";\n"
			}
			result += "\treg frame_ready_rx_0;\n"

			for i := 0; i < sendwires+1; i++ {
				result += "\t(* keep=\"true\" *) wire frame_ready_reset_" + strconv.Itoa(i) + ";\n"
			}
			result += "\t(* keep=\"true\" *) wire frame_ready_reset_rx_0;\n"

			result += "\n"

			// Write enabling
			result += "\treg start_write, start_rx;\n"
			result += "\twire done_write;\n"
			result += "\treg [ETHERNET_LENGTH-1:0] EthTx;\n"
			result += "\t(* keep=\"true\" *) reg [ETHERNET_LENGTH-1:0] EthRx;\n"

			result += "\n"

			result += "\talways@(posedge clk)\n"
			result += "\tbegin\n"
			result += "\t\tif(RESET_n==1'b0) \n"
			result += "\t\t\tframe_counter <= 0;\n"
			result += "\t\telse\n"
			result += "\t\t\tframe_counter <= frame_counter + 1;\n"
			result += "\tend\n"

			result += "\n"

			subresult := "frame_ready_rx_0, "
			subresult2 := "1'b1, "
			for i := sendwires; i >= 0; i-- {
				subresult += "frame_ready_" + strconv.Itoa(i) + ", "
				subresult2 += "1'b0, "
			}
			result += "\twire [" + strconv.Itoa(sendwires+1) + ":0] frame_ready = {" + subresult[0:len(subresult)-2] + "};\n"
			result += "\twire [" + strconv.Itoa(sendwires+1) + ":0] frame_type = {" + subresult2[0:len(subresult2)-2] + "};\n"

			result += "\n"

			result += "\twire stop_frame_explorer = (frame_explorer==LASTFRAME) ? 1'b1 : 1'b0;\n"

			result += "\n"

			result += "\talways@(posedge clk)\n"
			result += "\tbegin\n"
			result += "\t\tif(RESET_n==1'b0  | end_configuration==1'b0 | stop_frame_explorer)\n"
			zerostring := ""
			for i := 0; i < Needed_bits(sendwires+2); i++ {
				zerostring += "0"
			}
			result += "\t\t\tframe_explorer <= " + strconv.Itoa(Needed_bits(sendwires+2)) + "'b" + zerostring + ";\n"
			result += "\t\telse if(start_write == 1'b0 & start_rx == 1'b0)\n"
			result += "\t\t\tframe_explorer <= frame_explorer + 1;\n"
			result += "\tend\n"

			result += "\n"

			result += "\treg [" + strconv.Itoa(sendwires+1) + ":0] frame_ready_reset;\n"
			for i := 0; i < sendwires+1; i++ {
				result += "\tassign frame_ready_reset_" + strconv.Itoa(i) + " = frame_ready_reset[" + strconv.Itoa(i) + "];\n"
			}
			result += "\tassign frame_ready_reset_rx_0 = frame_ready_reset[" + strconv.Itoa(sendwires+1) + "];\n"

			result += "\n"

			result += "\talways@(posedge clk)\n"
			result += "\tbegin\n"
			result += "\t\tif(RESET_n==1'b0)\n"
			result += "\t\t\ttag <= 32'b0;\n"
			result += "\t\telse if(done_write == 1'b1)\n"
			result += "\t\t\ttag <= tag + 1'b1;\n"
			result += "\tend\n"

			result += "\n"

			result += "\talways@(posedge clk)\n"
			result += "\tbegin\n"
			result += "\t\tif(RESET_n==1'b0) begin\n"
			result += "\t\t\tstart_write <= 1'b0;\n"
			result += "\t\t\tstart_rx <= 1'b0;\n"
			result += "\t\t\tEthTx <= 'b0;\n"
			result += "\t\tend\n"
			result += "\t\telse\n"
			result += "\t\tbegin\n"
			result += "\t\t\tframe_ready_reset <= 'b0;\n"
			result += "\t\t\tif(frame_ready[frame_explorer] & start_write==1'b0 & start_rx==1'b0 & frame_type[frame_explorer]==1'b0)\n"
			result += "\t\t\tbegin\n"
			result += "\t\t\t\tframe_ready_reset[frame_explorer] <= 1'b1;\n"
			result += "\t\t\t\tEthTx <= ethtxs[frame_explorer];\n"
			result += "\t\t\t\tstart_write <= 1'b1;\n"
			result += "\t\t\tend\n"
			result += "\t\t\telse if(done_write)\n"
			result += "\t\t\t\tstart_write <= 1'b0;\n"
			result += "\t\t\tif(frame_ready[frame_explorer] & start_write==1'b0 & start_rx==1'b0 & frame_type[frame_explorer]==1'b1)\n"
			result += "\t\t\tbegin\n"
			result += "\t\t\t\tframe_ready_reset[frame_explorer] <= 1'b1;\n"
			result += "\t\t\t\tstart_rx <= 1'b1;\n"
			result += "\t\t\tend\n"
			result += "\t\t\telse if(done_rx)\n"
			result += "\t\t\t\tstart_rx <= 1'b0;\n"
			result += "\t\tend\n"
			result += "\tend\n"

			result += "\n"

			result += "\treg old_clu_adv;\n"
			result += "\talways@(posedge clk)\n"
			result += "\tbegin\n"
			result += "\t\tif(RESET_n==1'b0) \n"
			result += "\t\t\told_clu_adv <= 0;\n"
			result += "\t\telse\n"
			result += "\t\t\told_clu_adv <= frame_counter[27];\n"
			result += "\tend\n"

			result += "\talways@(posedge clk)\n"
			result += "\tbegin\n"
			result += "\t\tif(RESET_n==1'b0 | frame_ready_reset_0) \n"
			result += "\t\t\tframe_ready_0 <= 1'b0;\n"
			result += "\t\telse if(old_clu_adv!=frame_counter[27]) \n"
			result += "\t\t\tframe_ready_0 <= 1'b1;\n"
			result += "\tend\n"

			result += "\n"

			iw = 1
			for i, iname := range inames {
				if iname != "" {
					adv_name := "old_in_adv_" + strconv.Itoa(i)
					result += "\treg " + adv_name + ";\n"
					result += "\talways@(posedge clk)\n"
					result += "\tbegin\n"
					result += "\t\tif(RESET_n==1'b0) \n"
					result += "\t\t\t" + adv_name + " <= 0;\n"
					result += "\t\telse\n"
					result += "\t\t\t" + adv_name + " <= frame_counter[27];\n"
					result += "\tend\n"

					result += "\talways@(posedge clk)\n"
					result += "\tbegin\n"
					result += "\t\tif(RESET_n==1'b0 | frame_ready_reset_" + strconv.Itoa(iw) + ") \n"
					result += "\t\t\tframe_ready_" + strconv.Itoa(iw) + " <= 1'b0;\n"
					result += "\t\telse if(" + adv_name + "!=frame_counter[27]) \n"
					result += "\t\t\tframe_ready_" + strconv.Itoa(iw) + " <= 1'b1;\n"
					result += "\tend\n"

					result += "\n"
					iw++
				}
			}

			for i, oname := range onames {
				if oname != "" {
					adv_name := "old_out_adv_" + strconv.Itoa(i)
					result += "\treg " + adv_name + ";\n"
					result += "\talways@(posedge clk)\n"
					result += "\tbegin\n"
					result += "\t\tif(RESET_n==1'b0) \n"
					result += "\t\t\t" + adv_name + " <= 'b0;\n"
					result += "\t\telse\n"
					result += "\t\t\t" + adv_name + " <= frame_counter[27];\n"
					result += "\tend\n"

					result += "\talways@(posedge clk)\n"
					result += "\tbegin\n"
					result += "\t\tif(RESET_n==1'b0 | frame_ready_reset_" + strconv.Itoa(iw) + ") \n"
					result += "\t\t\tframe_ready_" + strconv.Itoa(iw) + " <= 1'b0;\n"
					result += "\t\telse if(" + adv_name + "!=frame_counter[27]) \n"
					result += "\t\t\tframe_ready_" + strconv.Itoa(iw) + " <= 1'b1;\n"
					result += "\tend\n"

					result += "\n"
					iw++

					result += "\treg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] " + oname + "_d1;\n"
					result += "\talways@(posedge clk)\n"
					result += "\tbegin\n"
					result += "\t\tif(RESET_n==1'b0) \n"
					result += "\t\t\t" + oname + "_d1 <= 'b0;\n"
					result += "\t\telse\n"
					result += "\t\t\t" + oname + "_d1 <= " + oname + "_copy;\n"
					result += "\tend\n"

					result += "\talways@(posedge clk)\n"
					result += "\tbegin\n"
					result += "\t\tif(RESET_n==1'b0 | frame_ready_reset_" + strconv.Itoa(iw) + ") \n"
					result += "\t\t\tframe_ready_" + strconv.Itoa(iw) + " <= 1'b0;\n"
					result += "\t\telse if(" + oname + "_d1!=" + oname + "_copy) \n"
					result += "\t\tbegin\n"

					ods := strings.Split(odests[i], "-")
					for _, _ = range ods {
						result += "\t\t\tframe_ready_" + strconv.Itoa(iw) + " <= 1'b1;\n"
						iw++
					}
					result += "\t\tend\n"
					result += "\tend\n"

					result += "\n"
				}
			}

			result += "  //logic to enable the RX frame\n"

			result += "  reg enable_rx;\n"
			result += "  always@(posedge clk)\n"
			result += "  begin\n"
			result += "      if(RESET_n==1'b0)\n"
			result += "          enable_rx <= 'b0;\n"
			result += "      else\n"
			result += "          enable_rx <= frame_counter[27];\n"
			result += "  end\n"
			result += "  always@(posedge clk)\n"
			result += "  begin\n"
			result += "      if(RESET_n==1'b0 | frame_ready_reset_rx_0)\n"
			result += "          frame_ready_rx_0 <= 1'b0;\n"
			result += "      else if(enable_rx!=frame_counter[28])\n"
			result += "          frame_ready_rx_0 <= 1'b1;\n"
			result += "  end\n"
			result += "\n"
			result += "  (* keep = \"true\" *) reg [47:0] source_addess;\n"
			result += "  (* keep = \"true\" *) reg [47:0] destination_addess;\n"
			result += "  (* keep = \"true\" *) reg [15:0] ethernet_type;\n"
			result += "  (* keep = \"true\" *) reg [31:0] tagin;\n"
			result += "  (* keep = \"true\" *) reg [31:0] clusterid;\n"
			result += "  (* keep = \"true\" *) reg [31:0] nodeid;\n"

			result += "\n"
			result += "  always@(posedge clk)\n"
			result += "  begin\n"
			result += "    if(RESET_n==1'b0) begin\n"
			result += "        source_addess <= 'b0;\n"
			result += "        destination_addess <= 'b0;\n"
			result += "        ethernet_type <= 'b0;\n"
			result += "        i0_copy <= 8'b0;\n"
			result += "        tagin <= 32'b0;\n"
			result += "        nodeid <= 32'b0;\n"
			result += "        clusterid <= 32'b0;\n"
			result += "    end\n"
			result += "    else if(done_rx) begin\n"
			result += "        case (EthRx[ETHERNET_LENGTH-1-8-48*2:ETHERNET_LENGTH-8-48*2-16])\n"
			result += "            ETHERBOND_TYPE: begin\n"
			result += "                case (EthRx[ETHERNET_LENGTH-1-8-48*2-16:ETHERNET_LENGTH-8-48*2-16-8])\n"
			result += "                    IO_TR: begin\n"
			result += "                    tagin <= EthRx[ETHERNET_LENGTH-1-8-48*2-16-8:ETHERNET_LENGTH-8-48*2-16-8-32];\n"
			result += "                    clusterid <= EthRx[ETHERNET_LENGTH-1-8-48*2-16-8-32:ETHERNET_LENGTH-8-48*2-16-8-32-32];\n"
			result += "                    nodeid <= EthRx[ETHERNET_LENGTH-1-8-48*2-16-8-32-32:ETHERNET_LENGTH-8-48*2-16-8-32-32-32];\n"
			result += "                    case (EthRx[ETHERNET_LENGTH-1-8-48*2-16-8-32-32-32:ETHERNET_LENGTH-8-48*2-16-8-32-32-32-32])\n"
			result += "                        id_resource_i0: begin\n"
			result += "                            i0_copy <= EthRx[ETHERNET_LENGTH-1-8-48*2-16-8-32-32-32-32:ETHERNET_LENGTH-8-48*2-16-8-32-32-32-32-8];\n"
			result += "                        end\n"
			result += "                    endcase\n"
			result += "                    end\n"
			result += "                endcase\n"
			result += "            end\n"
			result += "        endcase\n"
			result += "        source_addess <= EthRx[ETHERNET_LENGTH-1-8-48:ETHERNET_LENGTH-8-48*2];\n"
			result += "        destination_addess <= EthRx[ETHERNET_LENGTH-1-8:ETHERNET_LENGTH-8-48];\n"
			result += "        ethernet_type <= EthRx[ETHERNET_LENGTH-1-8-48*2:ETHERNET_LENGTH-8-48*2-16];\n"
			result += "    end\n"
			result += "  end\n"

			result += "\n"
			result += "\n"
			result += "\n"
			result += "  TopModuleSPI TopModuleSPI_inst\n"
			result += "  (\n"
			result += "    .CLOCK(clk),\n"
			result += "    .SCK(sck),\n"
			result += "    .MOSI(mosi),\n"
			result += "    .CS_n(cs_n),\n"
			result += "    .MISO(miso),\n"
			result += "    .reset_n(RESET_n),\n"
			result += "    .start_conf(start_conf),\n"
			result += "    .end_conf(end_configuration),\n"
			result += "    .EthTx(EthTx),\n"
			result += "    .start_write(start_write),\n"
			result += "    .done_write(done_write),\n"
			result += "    .start_rx(start_rx),\n"
			result += "    .done_rx(done_rx),\n"
			result += "    .EthRx(EthRx)\n"
			result += "   );\n"

			//loopcheck := false
			for _, iname := range inames {
				if iname != "" {
					result += "\treg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] " + iname + "_copy;\n"
					result += "\tassign " + iname + "=" + iname + "_copy;\n"
					//loopcheck = true
				}
			}

			//if loopcheck {
			//	result += "\talways @ (posedge clk) begin\n"
			//	for _, iname := range inames {
			//		if iname != "" {
			//			result += "\t\t" + iname + " <= " + iname + "_copy;\n"
			//		}
			//	}
			//	result += "\tend\n"
			//}

			result += "endmodule\n\n\n"
			return result, nil

		}
	}
	return "", errors.New("No etherbond module found")
}

func (bmach *Bondmachine) Write_verilog_udpbond(module_name string, flavor string, iomaps *IOmap, extramods []ExtraModule) (string, error) {
	for _, mod := range extramods {
		if mod.Get_Name() == "udpbond" {
			//			udpbond_params := mod.Get_Params().Params
			//			fmt.Println(udpbond_params)
			//
			//			//			cluster_id, _ := strconv.Atoi(udpbond_params["cluster_id"])
			//			//			peer_id, _ := strconv.Atoi(udpbond_params["peer_id"])
			//
			result := ""
			//			result += "module " + module_name + "_main(\n"
			//			result += "\n"
			//
			//			clk_name := "clk"
			//			rst_name := "reset"
			//
			//			result += "\tinput " + clk_name + ",\n"
			//			result += "\tinput " + rst_name + ",\n"
			//			result += "\toutput wifi_enable,\n"
			//			result += "\tinput wifi_rx,\n"
			//			result += "\toutput wifi_tx,\n"
			//
			//			inames := strings.Split(udpbond_params["inputs"], ",")
			//			//			iids := strings.Split(udpbond_params["input_ids"], ",")
			//
			//			onames := strings.Split(udpbond_params["outputs"], ",")
			//			//			oids := strings.Split(udpbond_params["output_ids"], ",")
			//			//			odests := strings.Split(udpbond_params["destinations"], ",")
			//
			//			for _, iname := range inames {
			//				if iname != "" {
			//					result += "\toutput [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] " + iname + ",\n"
			//				}
			//			}
			//
			//			for _, oname := range onames {
			//				if oname != "" {
			//					result += "\tinput [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] " + oname + ",\n"
			//				}
			//			}
			//
			//			result = result[0:len(result)-2] + "\n);\n\n"
			//
			//			result += "\n"
			//			result += "endmodule\n\n\n"

			return result, nil

		}
	}
	return "", errors.New("No udpbond module found")
}

func (bmach *Bondmachine) Write_verilog_board(conf *Config, module_name string, flavor string, iomaps *IOmap, extramods []ExtraModule) string {

	resolved_io := make(map[string]string)

	slow_module := false
	var slow_params map[string]string

	// etherbond and udpbond are exclusive

	etherbond_module := false
	var etherbond_params map[string]string

	udpbond_module := false
	var udpbond_params map[string]string

	var inames []string
	var onames []string

	basys3_7segment_module := false
	var basys3_7segment_params map[string]string
	basys3_7segment_mapped := ""

	icebreakerLedsModule := false
	var icebreakerLedsModuleParams map[string]string
	icebreakerLedsModuleMapped := ""

	ps2KeyboardModule := false
	var ps2KeyboardParams map[string]string
	ps2KeyboardMapped := ""

	uartModule := false
	var uartParams map[string]string

	vgatext800x600Module := false
	var vgatext800x600Params map[string]string

	bmapiModule := false
	var bmapiParams map[string]string

	result := ""
	result_headers := ""

	// Extra modules preparation

	for _, mod := range extramods {

		result_headers += mod.Verilog_headers()

		result += mod.StaticVerilog()

		if mod.Get_Name() == "slow" {
			slow_module = true
			slow_params = mod.Get_Params().Params
		}
		if mod.Get_Name() == "etherbond" {
			etherbond_module = true
			etherbond_params = mod.Get_Params().Params
			//fmt.Println(etherbond_params)
			if subresult, ok := bmach.Write_verilog_etherbond("etherbond", flavor, iomaps, extramods); ok == nil {
				result += subresult
			} else {
				// TODO proper error handling
			}
			inames = strings.Split(etherbond_params["inputs"], ",")
			onames = strings.Split(etherbond_params["outputs"], ",")
		}
		if mod.Get_Name() == "udpbond" {
			udpbond_module = true
			udpbond_params = mod.Get_Params().Params
			//fmt.Println(udpbond_params)
			if subresult, ok := bmach.Write_verilog_udpbond("udpbond", flavor, iomaps, extramods); ok == nil {
				result += subresult
			} else {
				// TODO proper error handling
			}
			inames = strings.Split(udpbond_params["inputs"], ",")
			onames = strings.Split(udpbond_params["outputs"], ",")
		}
		if mod.Get_Name() == "basys3_7segment" {
			basys3_7segment_module = true
			basys3_7segment_params = mod.Get_Params().Params
			if subresult, ok := bmach.Write_verilog_basys3_7segment("basys3_7segment", flavor, iomaps, extramods); ok == nil {
				result += subresult
			} else {
				// TODO proper error handling
			}
		}
		if mod.Get_Name() == "icebreaker_leds" {
			icebreakerLedsModule = true
			icebreakerLedsModuleParams = mod.Get_Params().Params
			//if subresult, ok := bmach.Write_verilog_basys3_7segment("basys3_7segment", flavor, iomaps, extramods); ok == nil {
			//result += subresult
			//} else {
			//// TODO proper error handling
			//}
		}
		if mod.Get_Name() == "ps2keyboard" {
			ps2KeyboardModule = true
			ps2KeyboardParams = mod.Get_Params().Params
			if subresult, ok := bmach.WriteVerilogPs2Keyboard("ps2keyboard", flavor, iomaps, extramods); ok == nil {
				result += subresult
			} else {
				// TODO proper error handling
			}
		}
		if mod.Get_Name() == "uart" {
			uartModule = true
			uartParams = mod.Get_Params().Params
		}
		if mod.Get_Name() == "vga800x600" {
			vgatext800x600Module = true
			vgatext800x600Params = mod.Get_Params().Params
			if subresult, ok := bmach.WriteVerilogVgaText800x600("vga800x600", flavor, iomaps, extramods); ok == nil {
				result += subresult
			} else {
				// TODO proper error handling
			}
		}
		if mod.Get_Name() == "bmapi" {
			bmapiModule = true
			bmapiParams = mod.Get_Params().Params
			if subresult, ok := bmach.WriteVerilogBMAPI("bmapi", flavor, iomaps, extramods); ok == nil {
				result += subresult
			} else {
				// TODO proper error handling
			}
			inames = strings.Split(bmapiParams["inputs"], ",")
			onames = strings.Split(bmapiParams["outputs"], ",")
		}
	}

	// Main module building

	// External ports of the BM
	result += "module " + module_name + "_main(\n"
	result += "\n"

	clk_name := "clk"

	if cname, ok := iomaps.Assoc["clk"]; ok {
		clk_name = cname
	}

	rst_name := "reset"

	if rname, ok := iomaps.Assoc["reset"]; ok {
		rst_name = rname
	}

	result += "\tinput " + clk_name + ",\n"
	result += "\tinput " + rst_name + ",\n"

	for i := 0; i < bmach.Inputs; i++ {
		iname := Get_input_name(i)
		resolved_io[iname] = "none"
		subresult := ""
		if rname, ok := iomaps.Assoc[iname]; ok {
			resolved_io[iname] = "board"
			subresult = "\tinput " + rname + ",\n"
		}
		result += subresult
	}

	for i := 0; i < bmach.Outputs; i++ {
		oname := Get_output_name(i)
		resolved_io[oname] = "none"
		subresult := ""
		if rname, ok := iomaps.Assoc[oname]; ok {
			resolved_io[oname] = "board"
			subresult = "\toutput reg " + rname + ",\n"
		}
		result += subresult
	}

	// Processing the extra module external ports

	if basys3_7segment_module {
		result += "\toutput [6:0] segment,\n"
		result += "\toutput enable_D1,\n"
		result += "\toutput enable_D2,\n"
		result += "\toutput enable_D3,\n"
		result += "\toutput enable_D4,\n"
		result += "\toutput dp,\n"
	}

	if icebreakerLedsModule {
		result += "\toutput led1,\n"
		result += "\toutput led2,\n"
		result += "\toutput led3,\n"
		result += "\toutput led4,\n"
		result += "\toutput led5,\n"
	}

	if ps2KeyboardModule {
		result += "\tinput PS2Data,\n"
		result += "\tinput PS2Clk,\n"
	}

	if vgatext800x600Module {
		result += "\toutput wire VGA_HS_O,\n"
		result += "\toutput wire VGA_VS_O,\n"
		result += "\toutput reg [3:0] VGA_R,\n"
		result += "\toutput reg [3:0] VGA_G,\n"
		result += "\toutput reg [3:0] VGA_B,\n"
	}

	if etherbond_module {
		result += "\toutput sck,\n"
		result += "\toutput mosi,\n"
		result += "\toutput cs_n,\n"
		result += "\tinput miso,\n"
		result += "\tinput int_n,\n"
	}

	if udpbond_module {
		result += "\toutput wifi_enable,\n"
		result += "\tinput wifi_rx,\n"
		result += "\toutput wifi_tx,\n"
	}

	if uartModule {
		for name, value := range uartParams {
			fmt.Println(name[len(name)-3:])
			if name[len(name)-3:] == "_rx" {
				result += "\tinput " + value + ",\n"
			}
			if name[len(name)-3:] == "_tx" {
				result += "\toutput " + value + ",\n"
			}
		}
	}

	if bmapiModule {
		// TODO different flavors and transivers + include error handling
		if bmapiFlavor, ok := bmapiParams["bmapi_flavor"]; ok {
			switch bmapiFlavor {
			case "uartusb":
				result += "\toutput TxD,\n"
				result += "\tinput RxD,\n"
			case "aximm":
				result += "\tinput [31:0] A_DVDR_PS2PL,\n"
				result += "\toutput [31:0] A_DVDR_PL2PS,\n"
				result += "\toutput [31:0] A_changes,\n"
				result += "\tinput [31:0] A_states,\n"

				sort.Strings(onames)

				for _, oname := range onames {
					for i := 0; i < bmach.Outputs; i++ {
						ooname := Get_output_name(i)
						if ooname == oname {
							result += "\toutput [31:0] A_port_o" + strconv.Itoa(i) + ",\n"
							result += "\toutput A_port_o" + strconv.Itoa(i) + "_valid,\n"
							result += "\tinput A_port_o" + strconv.Itoa(i) + "_recv,\n"
						}
					}
				}

				sort.Strings(inames)

				for _, iname := range inames {
					for i := 0; i < bmach.Inputs; i++ {
						iiname := Get_input_name(i)
						if iiname == iname {
							result += "\tinput [31:0] A_port_i" + strconv.Itoa(i) + ",\n"
							result += "\tinput A_port_i" + strconv.Itoa(i) + "_valid,\n"
							result += "\toutput A_port_i" + strconv.Itoa(i) + "_recv,\n"
						}
					}
				}
				result += "\toutput interrupt,\n"
			default:
				// Include errors
			}
		}
	}

	result = result[0:len(result)-2] + "\n);\n\n"
	// External ports creation ended
	if conf.CommentedVerilog {
		result += "\t// External ports creation ended\n"
	}

	if rst_name != "reset" {
		if logic, ok := iomaps.Assoc["logic"]; ok {
			switch logic {
			case "negative":
				result += "\tassign reset = ~" + rst_name + ";\n"
			default:
				result += "\tassign reset = " + rst_name + ";\n"
			}
		} else {
			result += "\tassign reset = " + rst_name + ";\n"
		}
	}

	// The External_inputs connected are defined as input port
	for i := 0; i < bmach.Inputs; i++ {
		iname := Get_input_name(i)
		if _, ok := iomaps.Assoc[iname]; ok {
			result += "\treg [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] Input" + strconv.Itoa(i) + ";\n"
			continue
		}

		if etherbond_module {
			for _, ethiname := range inames {
				if iname == ethiname {
					resolved_io[iname] = "etherbond"
					result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] Input" + strconv.Itoa(i) + ";\n"
					result += "\twire Input" + strconv.Itoa(i) + "_valid;\n"
					result += "\twire Input" + strconv.Itoa(i) + "_received;\n"
					break
				}
			}
		}

		if udpbond_module {
			for _, ethiname := range inames {
				if iname == ethiname {
					resolved_io[iname] = "udpbond"
					result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] Input" + strconv.Itoa(i) + ";\n"
					result += "\twire Input" + strconv.Itoa(i) + "_valid;\n"
					result += "\twire Input" + strconv.Itoa(i) + "_received;\n"
					break
				}
			}
		}

		if ps2KeyboardModule {
			if iname == ps2KeyboardParams["mapped_input"] {
				resolved_io[iname] = "ps2keyboard"
				result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] Input" + strconv.Itoa(i) + ";\n"
				result += "\twire Input" + strconv.Itoa(i) + "_valid;\n"
				result += "\twire Input" + strconv.Itoa(i) + "_received;\n"
				ps2KeyboardMapped = "Input" + strconv.Itoa(i)
			}
		}

		if bmapiModule {
			for _, bminame := range inames {
				if iname == bminame {
					resolved_io[iname] = "bmapi"
					result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] Input" + strconv.Itoa(i) + ";\n"
					result += "\twire Input" + strconv.Itoa(i) + "_valid;\n"
					result += "\twire Input" + strconv.Itoa(i) + "_received;\n"
					break
				}
			}
		}

	}

	// The External_inputs connected are defined as input port
	for i := 0; i < bmach.Outputs; i++ {
		oname := Get_output_name(i)
		if _, ok := iomaps.Assoc[oname]; ok {
			result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] Output" + strconv.Itoa(i) + ";\n"
			continue
		}

		if etherbond_module {
			for _, ethoname := range onames {
				if oname == ethoname {
					resolved_io[oname] = "etherbond"
					result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] Output" + strconv.Itoa(i) + ";\n"
					result += "\twire Onput" + strconv.Itoa(i) + "_valid;\n"
					result += "\twire Onput" + strconv.Itoa(i) + "_received;\n"
					break
				}
			}
		}

		if udpbond_module {
			for _, ethoname := range onames {
				if oname == ethoname {
					resolved_io[oname] = "udpbond"
					result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] Output" + strconv.Itoa(i) + ";\n"
					result += "\twire Onput" + strconv.Itoa(i) + "_valid;\n"
					result += "\twire Onput" + strconv.Itoa(i) + "_received;\n"
					break
				}
			}
		}

		if bmapiModule {
			for _, bmoname := range onames {
				if oname == bmoname {
					resolved_io[oname] = "bmapi"
					result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] Output" + strconv.Itoa(i) + ";\n"
					result += "\twire Output" + strconv.Itoa(i) + "_valid;\n"
					result += "\twire Output" + strconv.Itoa(i) + "_received;\n"
					break
				}
			}
		}

		if basys3_7segment_module {
			if oname == basys3_7segment_params["mapped_output"] {
				resolved_io[oname] = "basys3_7segment"
				result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] Output" + strconv.Itoa(i) + ";\n"
				basys3_7segment_mapped = "Output" + strconv.Itoa(i)
			}
		}
		if icebreakerLedsModule {
			if oname == icebreakerLedsModuleParams["mapped_output"] {
				resolved_io[oname] = "icebreaker_module"
				result += "\twire [" + strconv.Itoa(int(bmach.Rsize)-1) + ":0] Output" + strconv.Itoa(i) + ";\n"
				icebreakerLedsModuleMapped = "Output" + strconv.Itoa(i)
			}
		}
	}

	result += "\n"
	clockString := clk_name

	// Processing per-extramodule initializazions
	if conf.CommentedVerilog {
		result += "\t// Processing per-extramodule initializazions\n"
	}

	if slow_module {
		result += "\treg [31:0] divider;\n\n"
		clockString = "divider[" + slow_params["slow_factor"] + "]"
	}

	if vgatext800x600Module {
		result += `
	reg [7:0]   hcount;
	reg [7:0]   header [0:99]; 

	initial begin
		$display("Loading memory init file header into array.");
		$readmemh("` + vgatext800x600Params["header"] + `", header);
		hcount=99;
	end
`
	}

	if bmapiModule {
		result += "\twire transconnectd;\n"
	}

	// Processing Headers from External ports of SOs
	if len(bmach.Shared_objects) > 0 {
		for proc_id, solist := range bmach.Shared_links {
			for _, so_id := range solist {
				result += strings.ReplaceAll(strings.ReplaceAll(bmach.Shared_objects[so_id].GetExternalPortsWires(bmach, proc_id, so_id, flavor), "output", "wire"), "input", "wire")
			}
		}
	}

	if uartModule {
		for name, value := range uartParams {
			fmt.Println(name[len(name)-3:])
			if name[len(name)-3:] == "_rx" {
				result += "\tassign " + name + "=" + value + ";\n"
			}
			if name[len(name)-3:] == "_tx" {
				result += "\tassign " + value + "=" + name + ";\n"
			}
		}
	}

	// Processing BM ports originated from external modules (not IO)
	if conf.CommentedVerilog {
		result += "\t// Processing BM ports originated from external modules (not IO)\n"
	}

	result += "\t" + module_name + " " + module_name + "_inst " + "(" + clockString + ", reset"

	// The External_inputs connected are defined as input port
	for i := 0; i < bmach.Inputs; i++ {
		result += ", Input" + strconv.Itoa(i)
		result += ", Input" + strconv.Itoa(i) + "_valid"
		result += ", Input" + strconv.Itoa(i) + "_received"
	}

	// The External_inputs connected are defined as input port
	for i := 0; i < bmach.Outputs; i++ {
		result += ", Output" + strconv.Itoa(i)
		result += ", Output" + strconv.Itoa(i) + "_valid"
		result += ", Output" + strconv.Itoa(i) + "_received"
	}

	// Processing Headers from External ports of SOs
	if len(bmach.Shared_objects) > 0 {
		for proc_id, solist := range bmach.Shared_links {
			for _, so_id := range solist {
				result += bmach.Shared_objects[so_id].GetExternalPortsHeader(bmach, proc_id, so_id, flavor)
			}
		}
	}

	result += ");\n\n"

	// Processing BM connected firwmwares
	if conf.CommentedVerilog {
		result += "\t// Processing BM connected firmwares\n"
	}

	if etherbond_module {
		result += "\tetherbond_main etherbond_main_inst " + "(" + clk_name + ", reset"
		result += ", sck"
		result += ", mosi"
		result += ", cs_n"
		result += ", miso"
		result += ", int_n"
		for _, iname := range inames {
			for i := 0; i < bmach.Inputs; i++ {
				iiname := Get_input_name(i)
				if iiname == iname {
					result += ", Input" + strconv.Itoa(i)
				}
			}
		}

		for _, oname := range onames {
			for i := 0; i < bmach.Outputs; i++ {
				ooname := Get_output_name(i)
				if ooname == oname {
					result += ", Output" + strconv.Itoa(i)
				}
			}
		}

		result += ");\n"
	}

	if udpbond_module {
		result += "\tudpbond_main udpbond_main (.clk100(" + clk_name + "), .reset(reset), .wifi_enable(wifi_enable), .wifi_rx(wifi_rx), .wifi_tx(wifi_tx)"
		for _, iname := range inames {
			for i := 0; i < bmach.Inputs; i++ {
				iiname := Get_input_name(i)
				if iiname == iname {
					result += ", .input_" + iname + "(Input" + strconv.Itoa(i) + ")"
				}
			}
		}

		for _, oname := range onames {
			for i := 0; i < bmach.Outputs; i++ {
				ooname := Get_output_name(i)
				if ooname == oname {
					result += ", .output_" + oname + "(Output" + strconv.Itoa(i) + ")"
				}
			}
		}

		result += ");\n"
	}

	if basys3_7segment_module {
		result += "\tbond2seg bond2seg_inst(" + clk_name + ", reset, " + basys3_7segment_mapped + ", segment ,enable_D1, enable_D2, enable_D3, enable_D4, dp);\n"
	}

	if vgatext800x600Module {
		result += "\t// Eventually create here a module\n"
	}

	if ps2KeyboardModule {
		result += "\tbondkeydrv bondkeydrv_inst(" + clk_name + ", PS2Data, PS2Clk, " + ps2KeyboardMapped + ", " + ps2KeyboardMapped + "_valid, " + ps2KeyboardMapped + "_received);\n"
	}

	if bmapiModule {
		switch bmapiParams["bmapi_flavor"] {
		case "uartusb":
			result += "\tbmapiuarttransceiver bmapiuarttransceiver_inst(" + clk_name + ", reset, TxD, RxD"

			for _, oname := range onames {
				for i := 0; i < bmach.Outputs; i++ {
					ooname := Get_output_name(i)
					if ooname == oname {
						result += ", Output" + strconv.Itoa(i) + "[" + strconv.Itoa(int(bmach.Rsize)-1) + ":0], Output" + strconv.Itoa(i) + "_valid, Output" + strconv.Itoa(i) + "_received"
					}
				}
			}

			for _, iname := range inames {
				for i := 0; i < bmach.Inputs; i++ {
					iiname := Get_input_name(i)
					if iiname == iname {
						result += ", Input" + strconv.Itoa(i) + "[" + strconv.Itoa(int(bmach.Rsize)-1) + ":0], Input" + strconv.Itoa(i) + "_valid, Input" + strconv.Itoa(i) + "_received"
					}
				}
			}

			result += ", transconnected);\n"
		case "aximm":
			result += "\tbmapiaximmtransceiver bmapiaximmtransceiver_inst(\n"
			result += "\t\t.clk(" + clk_name + "),\n"
			result += "\t\t.reset(reset),\n"
			result += "\t\t.A_DVDR_PL2PS(A_DVDR_PL2PS[31:0]),\n"
			result += "\t\t.A_DVDR_PS2PL(A_DVDR_PS2PL[31:0]),\n"
			result += "\t\t.A_changes(A_changes[31:0]),\n"
			result += "\t\t.A_states(A_states[31:0]),\n"
			for _, oname := range onames {
				for i := 0; i < bmach.Outputs; i++ {
					ooname := Get_output_name(i)
					if ooname == oname {
						result += "\t\t.A_port_o" + strconv.Itoa(i) + "(A_port_o" + strconv.Itoa(i) + "[31:0]),\n"
						result += "\t\t.port_o" + strconv.Itoa(i) + "(Output" + strconv.Itoa(i) + "[" + strconv.Itoa(int(bmach.Rsize)-1) + ":0]),\n"
						result += "\t\t.port_o" + strconv.Itoa(i) + "_valid(Output" + strconv.Itoa(i) + "_valid),\n"
						result += "\t\t.port_o" + strconv.Itoa(i) + "_recv(Output" + strconv.Itoa(i) + "_received),\n"
					}
				}
			}

			for _, iname := range inames {
				for i := 0; i < bmach.Inputs; i++ {
					iiname := Get_input_name(i)
					if iiname == iname {
						result += "\t\t.A_port_i" + strconv.Itoa(i) + "(A_port_i" + strconv.Itoa(i) + "[31:0]),\n"
						result += "\t\t.port_i" + strconv.Itoa(i) + "(Input" + strconv.Itoa(i) + "[" + strconv.Itoa(int(bmach.Rsize)-1) + ":0]),\n"
						result += "\t\t.port_i" + strconv.Itoa(i) + "_valid(Input" + strconv.Itoa(i) + "_valid),\n"
						result += "\t\t.port_i" + strconv.Itoa(i) + "_recv(Input" + strconv.Itoa(i) + "_received),\n"
					}
				}
			}
			result += "\t\t.interrupt(interrupt)\n"
			result += "\t);\n"
		}
	}

	result += "\n"

	// Processinng BM IO
	if conf.CommentedVerilog {
		result += "\t// Processing BM IO\n"
	}

	if slow_module {
		result += "\talways @ (posedge clk) begin\n"
		result += "\t\tdivider <= divider + 1;\n"
		result += "\tend\n"
	}

	switch flavor {

	case "alveou50", "basys3", "kc705", "zedboard", "ebaz4205", "icebreaker", "icefun", "icesugarnano", "max1000", "de10nano", "zc702", "ice40lp1k":

		for _, iores := range resolved_io {

			if iores == "board" {

				result += "\talways @ (posedge clk) begin\n"

				for i := 0; i < bmach.Inputs; i++ {
					iname := Get_input_name(i)
					tpname := "Input" + strconv.Itoa(i)
					aname := iname
					if rname, ok := iomaps.Assoc[iname]; ok {
						aname = rname
					}
					if aname != iname {
						for j := 0; j < int(bmach.Rsize); j++ {
							result += "\t\t" + tpname + "[" + strconv.Itoa(j) + "] <= " + nth_assoc(aname, j) + ";\n"
						}
					}
				}

				for i := 0; i < bmach.Outputs; i++ {
					oname := Get_output_name(i)
					tpname := "Output" + strconv.Itoa(i)
					aname := oname
					if rname, ok := iomaps.Assoc[oname]; ok {
						aname = rname
					}
					if aname != oname {
						for j := 0; j < int(bmach.Rsize); j++ {
							result += "\t\t" + nth_assoc(aname, j) + " <= " + tpname + "[" + strconv.Itoa(j) + "]" + ";\n"
						}
					}
				}

				result += "\tend\n"

				break
			}
		}
	}

	// Processinng Extramodules processes
	if conf.CommentedVerilog {
		result += "\t// Processing Extra modules processes\n"
	}

	if icebreakerLedsModule {
		result += "\tassign {led1, led2, led3, led4, led5} = " + icebreakerLedsModuleMapped + "[4:0];\n"
	}

	if vgatext800x600Module {

		if len(bmach.Shared_objects) > 0 {

			countso := 0
			var vgaso Vtextmem_instance
			for _, so := range bmach.Shared_objects {
				sname := so.Shortname()
				if sname == "vtm" {
					countso++
					vgaso = so.(Vtextmem_instance)
				}
			}

			if countso != 1 {
				// TODO Some fail
			}
			boxes := vgaso.Boxes

			result += `

	// generate a 40 MHz pixel strobe
	reg [15:0] cnt;
	reg pix_stb;
	always @(posedge clk)
		{pix_stb, cnt} <= cnt + 16'h6666;  // divide by 2.5: (2^16)/2.5 = 0x6666

	wire [10:0] x;  // current pixel x position: 11-bit value: 0-2047
	wire  [9:0] y;  // current pixel y position: 10-bit value: 0-1023

	// Connect the VGA display
	vga800x600 display (
		.i_clk(clk),
		.i_pix_stb(pix_stb),
		.i_rst(reset),
		.o_hs(VGA_HS_O), 
		.o_vs(VGA_VS_O), 
		.o_x(x), 
		.o_y(y)
	);

	// Font ROM initialization
	localparam FONTROM_DEPTH = 1024; 
	localparam FONTROM_A_WIDTH = 10;
	localparam FONTROM_D_WIDTH = 8;

	reg [FONTROM_A_WIDTH-1:0] fontaddress;
	wire [FONTROM_D_WIDTH-1:0] dataout;

	romfonts #(
		.ADDR_WIDTH(FONTROM_A_WIDTH), 
		.DATA_WIDTH(FONTROM_D_WIDTH), 
		.DEPTH(FONTROM_DEPTH), 
		.FONTSFILE("` + vgatext800x600Params["fonts"] + `"))
		fonts (
		.addr(fontaddress[FONTROM_A_WIDTH-1:0]), 
		.data(dataout)
	);

	// HEADER Video RAM
	localparam HEAD_VIDEORAM_LEFT = 0;
	localparam HEAD_VIDEORAM_TOP = 1;
	localparam HEAD_VIDEORAM_ROWS = 1;
	localparam HEAD_VIDEORAM_COLS = 100;
	localparam HEAD_VIDEORAM_DEPTH = 100; 
	localparam HEAD_VIDEORAM_A_WIDTH = 8;
	localparam HEAD_VIDEORAM_D_WIDTH = 8;

	reg [HEAD_VIDEORAM_A_WIDTH-1:0] headaddress;
	wire [HEAD_VIDEORAM_D_WIDTH-1:0] headdataout;
	reg headwantwrite;
	reg [7:0] headdata;

	textvideoram #(
		.ADDR_WIDTH(HEAD_VIDEORAM_A_WIDTH), 
		.DATA_WIDTH(HEAD_VIDEORAM_D_WIDTH), 
		.DEPTH(HEAD_VIDEORAM_DEPTH))
		headvideoram (
		.addr(headaddress[HEAD_VIDEORAM_A_WIDTH-1:0]), 
		.o_data(headdataout[HEAD_VIDEORAM_D_WIDTH-1:0]),
		.clk(clk),
		.wen(headwantwrite),
		.i_data(headdata[HEAD_VIDEORAM_D_WIDTH-1:0])
	);

	// Pixel processing
	reg [11:0] colour;
	wire [7:0] pix;
    
	wire within_head;
	assign within_head = ((x >= HEAD_VIDEORAM_LEFT*8) & (y >=  HEAD_VIDEORAM_TOP*8) & (x < HEAD_VIDEORAM_LEFT*8+HEAD_VIDEORAM_COLS*8) & (y < HEAD_VIDEORAM_TOP*8+HEAD_VIDEORAM_ROWS*8)) ? 1 : 0;
`

			for _, box := range boxes {
				cpS := strconv.Itoa(box.CP)
				leftS := strconv.Itoa(box.Left * 8)
				topS := strconv.Itoa(box.Top * 8)
				widthS := strconv.Itoa(box.Width * 8)
				heightS := strconv.Itoa(box.Height * 8)
				result += "\n\tlocalparam CP" + cpS + "_VIDEORAM_LEFT = " + strconv.Itoa(box.Left) + ";\n"
				result += "\tlocalparam CP" + cpS + "_VIDEORAM_TOP = " + strconv.Itoa(box.Top) + ";\n"
				result += "\tlocalparam CP" + cpS + "_VIDEORAM_COLS = " + strconv.Itoa(box.Width) + ";\n"
				result += "\tlocalparam CP" + cpS + "_VIDEORAM_ROWS = " + strconv.Itoa(box.Height) + ";\n\n"
				result += "\twire border_t_cp" + cpS + ";\n"
				result += "\twire border_b_cp" + cpS + ";\n"
				result += "\twire border_l_cp" + cpS + ";\n"
				result += "\twire border_r_cp" + cpS + ";\n"
				result += "\twire border_cp" + cpS + ";\n"
				result += "\twire within_cp" + cpS + ";\n"
				result += "\tassign within_cp" + cpS + " = ((x >= CP" + cpS + "_VIDEORAM_LEFT*8) & (y >=  CP" + cpS + "_VIDEORAM_TOP*8) & (x < CP" + cpS + "_VIDEORAM_LEFT*8+CP" + cpS + "_VIDEORAM_COLS*8) & (y < CP" + cpS + "_VIDEORAM_TOP*8+CP" + cpS + "_VIDEORAM_ROWS*8)) ? 1 : 0; ;\n"
				result += "\tassign within_cp" + cpS + " = ((x >= CP" + cpS + "_VIDEORAM_LEFT*8) & (y >=  CP" + cpS + "_VIDEORAM_TOP*8) & (x < CP" + cpS + "_VIDEORAM_LEFT*8+CP" + cpS + "_VIDEORAM_COLS*8) & (y < CP" + cpS + "_VIDEORAM_TOP*8+CP" + cpS + "_VIDEORAM_ROWS*8)) ? 1 : 0; ;\n"
				result += "\tassign border_t_cp" + cpS + " = ((x >= " + leftS + " - 2) & (x <= " + leftS + "+" + widthS + " + 2) & ( y == " + topS + " - 2)) ? 1 : 0; ;\n"
				result += "\tassign border_b_cp" + cpS + " = ((x >= " + leftS + " - 2) & (x <= " + leftS + "+" + widthS + " + 2) & ( y == " + topS + "+" + heightS + " + 2)) ? 1 : 0; ;\n"
				result += "\tassign border_l_cp" + cpS + " = ((y >= " + topS + " - 2) & (y <= " + topS + "+" + heightS + " + 2) & ( x == " + leftS + " - 2)) ? 1 : 0; ;\n"
				result += "\tassign border_r_cp" + cpS + " = ((y >= " + topS + " - 2) & (y <= " + topS + "+" + heightS + " + 2) & ( x == " + leftS + "+" + widthS + " + 2)) ? 1 : 0; ;\n"
				result += "\tassign border_cp" + cpS + " = (( border_t_cp" + cpS + " ) | ( border_b_cp" + cpS + " ) | ( border_l_cp" + cpS + " ) | ( border_r_cp" + cpS + " )) ? 1 : 0; ;\n"
			}

			result += `
//	assign pix = dataout >> ((x-1) % 8) ;
	assign pix = dataout >> ((x-2) % 8) ;

	// Font address Process
	always @ (posedge clk)
	begin
		if (within_head)
			fontaddress[FONTROM_A_WIDTH-1:0] <= headdataout * 8 + (y % 8);
`
			withinS := ""
			for _, box := range boxes {
				cpS := strconv.Itoa(box.CP)
				withinS += " || within_cp" + cpS
				result += "\t\telse if (within_cp" + cpS + ")\n"
				result += "\t\t\tfontaddress[FONTROM_A_WIDTH-1:0] <= p" + cpS + "vtm0dout * 8 + (y % 8);\n"
			}
			result += `		else
			fontaddress[FONTROM_A_WIDTH-1:0] <= 0;
	end

	// HEADER Processing
	always @ (posedge clk)
	begin
		if (hcount > 0) begin
			hcount <= hcount - 1;
			headaddress[HEAD_VIDEORAM_A_WIDTH-1:0] <= hcount;
			headdata[HEAD_VIDEORAM_D_WIDTH-1:0] <= header[hcount];
			headwantwrite <= 1;
		end 
		else
		begin
			headaddress[HEAD_VIDEORAM_A_WIDTH-1:0] <= ((y / 8) - HEAD_VIDEORAM_TOP) * HEAD_VIDEORAM_COLS +  ((x / 8) - HEAD_VIDEORAM_LEFT);
			headwantwrite <= 0;
		end
	end
`
			for _, box := range boxes {
				cpS := strconv.Itoa(box.CP)
				depth := box.Width * box.Height
				depthBits := Needed_bits(depth)
				result += `
	// CP` + cpS + ` Processing

	reg [` + strconv.Itoa(depthBits-1) + `:0] cp` + cpS + `address;

	always @ (posedge clk)
	begin
		if (within_cp` + cpS + `)
			cp` + cpS + `address[` + strconv.Itoa(depthBits-1) + `:0] <= ((y / 8) - CP` + cpS + `_VIDEORAM_TOP) * CP` + cpS + `_VIDEORAM_COLS +  ((x / 8) - CP` + cpS + `_VIDEORAM_LEFT);
	end

	assign p` + cpS + `vtm0addrfromext[` + strconv.Itoa(depthBits-1) + `:0] = cp` + cpS + `address[` + strconv.Itoa(depthBits-1) + `:0];
`
			}
			result += `
	// Channels processing
	always @ (posedge clk)
	begin
		if (( within_head ) && pix[0] == 1'b1)
			colour <= 12'b111111111111;
		else if (( ` + withinS[3:] + ` ) && pix[0] == 1'b1)
			colour <= 12'b110111011101;
		else if (( ` + withinS[3:] + ` ) && pix[0] == 1'b0)
			colour <= 12'b001100110011;
`
			colors := []string{"h01c", "h075", "ha3d", "h603"}
			for j, box := range boxes {
				cpS := strconv.Itoa(box.CP)
				result += "\t\telse if ( border_cp" + cpS + " )\n"
				result += "\t\t\tcolour <= 12'" + colors[j%len(colors)] + ";\n"
			}
			result += `
		else
			colour <= 0;
    
		VGA_R <= colour[11:8];
		VGA_G <= colour[7:4];
		VGA_B <= colour[3:0];
	end

`
		}
	}
	result += "endmodule\n"
	return result_headers + result
}
