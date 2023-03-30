package bondmachine

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

//reorg {"name": "Init", "descr": "Initialization functions"}

type Prerror struct {
	string
}

var Allshared []Shared_element

func init() {
	Allshared = make([]Shared_element, 0)
	Allshared = append(Allshared, Sharedmem{})
	Allshared = append(Allshared, Channel{})
	Allshared = append(Allshared, Barrier{})
	Allshared = append(Allshared, Lfsr8{})
	Allshared = append(Allshared, Vtextmem{})
	Allshared = append(Allshared, Queue{})
	Allshared = append(Allshared, Stack{})
	Allshared = append(Allshared, Uart{})
}

func (e Prerror) Error() string {
	return e.string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Bond struct {
	Map_to uint8 // 0 Means a shell input, 1 a shell output, 2 a processor input, 3 a processor output
	Res_id int   // The reasource id of a processor or external IO
	Ext_id int   // The order within the resource if appliable
}

//reorg {"name": "Config typedef", "descr": "Definition og the Config structure"}

type Config struct {
	*bminfo.BMinfo
	*bmreqs.ReqRoot
	procbuilder.HwOptimizations
	Debug            bool
	Dotdetail        uint8
	CommentedVerilog bool
}

//reorg {"name": "BondMachine typedefs", "descr": "Definition of BondMachine and BondMachine JSON data structures"}

type Bondmachine struct {
	Rsize uint8

	Domains []*procbuilder.Machine

	Processors []int

	Inputs  int
	Outputs int

	Internal_inputs  []Bond // These are what internally in considered and input, i.e. Processors inputs and Shell outputs, lets say they are N
	Internal_outputs []Bond // These are what internally in considered and output, i.e. Processors outputs and Shell inputs, lets say they are M

	Links []int // Is the link matrix for Internal_inputs that may be connected only to one Internal_output (they are N!). A value of -1 means the input is unconnected, otherwise it points to the Internal_output id

	Shared_objects []Shared_instance      // This is the list of shared elements present in a bondmachine
	Shared_links   []Shared_instance_list // For every processor the list of attached shared objects

}

type Bondmachine_json struct {
	Rsize uint8

	Domains []*procbuilder.Machine_json

	Processors []int

	Inputs  int
	Outputs int

	Internal_inputs  []Bond
	Internal_outputs []Bond
	Links            []int

	Shared_objects []string
	Shared_links   []Shared_instance_list
	// TODO Missing shared elements
}

//reorg {"name": "Multi Abstract Assembly data structure definition"}

type Abs_assembly struct {
	ProcProgs []string
	Bonds     []string
}

type Residual struct {
	Map IOmap
}

// Mapping between boards ports and bondmachine IOs
type IOmap struct {
	Assoc map[string]string
}

// Parameters for extra modules
type ExtraParams struct {
	Params map[string]string
}

// Extra modules (etherbond)
type ExtraModule interface {
	Get_Name() string
	Get_Params() *ExtraParams
	Import(string) error
	Export() string
	Check(*Bondmachine) error
	Verilog_headers() string
	StaticVerilog() string
	ExtraFiles() ([]string, []string)
}

type EmuDriver interface {
	PushCommand([]byte) ([]byte, error)
	Init() error
	Run() error
}

//reorg {"name": "Configuration converter", "descr": "Funcion to extract a Config structure for procbuilder from a Config of the bondmachine"}

func (c *Config) ProcbuilderConfig() *procbuilder.Config {
	result := new(procbuilder.Config)
	result.ReqRoot = c.ReqRoot
	result.HwOptimizations = c.HwOptimizations
	result.Commented_verilog = c.CommentedVerilog
	return result
}

func (b *Bond) String() string {
	result := ""
	if b.Map_to == 0 {
		result += "i" + strconv.Itoa(b.Res_id)
	} else if b.Map_to == 1 {
		result += "o" + strconv.Itoa(b.Res_id)
	} else if b.Map_to == 2 {
		result += "p" + strconv.Itoa(b.Res_id) + "i" + strconv.Itoa(b.Ext_id)
	} else if b.Map_to == 3 {
		result += "p" + strconv.Itoa(b.Res_id) + "o" + strconv.Itoa(b.Ext_id)
	}
	return result
}

//reorg {"name": "BondMachine Init", "descr": "Initialization function for BondMachine istances"}

func (bmach *Bondmachine) Init() {
	// It contains an idempotent set of operation to ensure bondmachine consistency

	if bmach.Shared_links == nil {
		bmach.Shared_links = make([]Shared_instance_list, len(bmach.Processors))
		for i, _ := range bmach.Shared_links {
			bmach.Shared_links[i] = make([]int, 0)
		}
	}

}

func (bmach *Bondmachine) Jsoner() *Bondmachine_json {
	result := new(Bondmachine_json)
	result.Rsize = bmach.Rsize
	result.Processors = bmach.Processors
	result.Inputs = bmach.Inputs
	result.Outputs = bmach.Outputs
	result.Internal_inputs = bmach.Internal_inputs
	result.Internal_outputs = bmach.Internal_outputs
	result.Links = bmach.Links
	result.Domains = make([]*procbuilder.Machine_json, len(bmach.Domains))
	for i, mach := range bmach.Domains {
		result.Domains[i] = mach.Jsoner()
	}
	result.Shared_objects = make([]string, len(bmach.Shared_objects))
	for i, so := range bmach.Shared_objects {
		result.Shared_objects[i] = so.String()
	}
	result.Shared_links = bmach.Shared_links
	return result
}

func (bmachj *Bondmachine_json) Dejsoner() *Bondmachine {
	result := new(Bondmachine)
	result.Rsize = bmachj.Rsize
	result.Processors = bmachj.Processors
	result.Inputs = bmachj.Inputs
	result.Outputs = bmachj.Outputs
	result.Internal_inputs = bmachj.Internal_inputs
	result.Internal_outputs = bmachj.Internal_outputs
	result.Links = bmachj.Links
	result.Domains = make([]*procbuilder.Machine, len(bmachj.Domains))
	for i, machj := range bmachj.Domains {
		result.Domains[i] = machj.Dejsoner()
	}
	result.Shared_objects = make([]Shared_instance, len(bmachj.Shared_objects))
	for i, so := range bmachj.Shared_objects {
		// TODO loading checks missing
		for _, shr := range Allshared {
			if inst, ok := shr.Instantiate(so); ok {
				//              loaded = true
				result.Shared_objects[i] = inst
				break
			}
		}

	}
	result.Shared_links = bmachj.Shared_links
	return result
}

func (bmach *Bondmachine) List_domains() string {
	result := ""
	if len(bmach.Domains) != 0 {
		for i, dom := range bmach.Domains {
			result += fmt.Sprintf("----------\nDomain %03d\n", i) + dom.Descr() + fmt.Sprintf("\n")
			proccheck := false
			proclist := ""
			for proc_id, proc_dom := range bmach.Processors {
				if proc_dom == i {
					proccheck = true
					proclist = proclist + fmt.Sprintf("%03d ", proc_id)
				}
			}
			if proccheck {
				result += "Processors " + proclist + "\n"
			} else {
				result += "No processor in this domain\n"
			}
		}
	} else {
		result += "No defined domains"
	}
	return result
}
func (bmach *Bondmachine) GetMultiAssembly() (*Abs_assembly, error) {

	if len(bmach.Processors) != 0 {
		result := new(Abs_assembly)
		result.ProcProgs = make([]string, 0)
		result.Bonds = make([]string, 0)

		for _, dom_id := range bmach.Processors {
			mymachine := bmach.Domains[dom_id]
			if disass_text, err := mymachine.Disassembler(); err == nil {
				result.ProcProgs = append(result.ProcProgs, disass_text)
			} else {
				return nil, err
			}
		}

		for _, bond := range bmach.List_bonds() {
			result.Bonds = append(result.Bonds, fmt.Sprintf("%s", bond))
		}

		return result, nil
	} else {
		return nil, Prerror{"Undefined processors"}
	}
	return nil, nil
}

func (bmach *Bondmachine) List_inputs() []string {
	result := []string{}
	if bmach.Inputs != 0 {
		for i := 0; i < bmach.Inputs; i++ {
			result = append(result, Get_input_name(i))
		}
	}
	return result
}

func (bmach *Bondmachine) Add_input() (string, error) {
	newbond := Bond{0, bmach.Inputs, 0}
	bmach.Inputs = bmach.Inputs + 1
	bmach.Internal_outputs = append(bmach.Internal_outputs, newbond)
	return "Added input " + strconv.Itoa(bmach.Inputs-1) + " successfully", nil
}

func (bmach *Bondmachine) Del_input(iid int) error {
	if iid < bmach.Inputs {
		bondtoremove := Bond{0, iid, 0}
		// 1 - remove all the bonds using the input (an input is within Internal_outputs)
		for i, linked := range bmach.Links {
			if linked != -1 {
				chkbond := bmach.Internal_outputs[linked]
				if chkbond.Map_to == 0 && chkbond.Res_id == bondtoremove.Res_id {
					bmach.Links[i] = -1
				}
			}
		}

		// 2 - Remove the internal output
		// 3 - lower by 1 all ids > iid in bonds
		newinternal := make([]Bond, len(bmach.Internal_outputs)-1)
		j := 0
		keeppos := -1
		for i, bond := range bmach.Internal_outputs {
			if bond.Map_to == 0 {
				if bond.Res_id == bondtoremove.Res_id {
					j--
					// Keep the position of the bond removed from Internal_outputs
					keeppos = i
				} else if bond.Res_id > bondtoremove.Res_id {
					newinternal[j] = Bond{0, bond.Res_id - 1, 0}
				} else {
					newinternal[j] = bond
				}
			} else {
				newinternal[j] = bond
			}
			j++
		}

		bmach.Internal_outputs = newinternal

		// 4 - Swift links
		if keeppos > -1 {
			for i, linked := range bmach.Links {
				if linked > keeppos {
					bmach.Links[i] = linked - 1
				}
			}
		}

		// 5 - Lower Inputs
		bmach.Inputs--
	} else {
		return Prerror{"Input id outside limit"}
	}
	return nil
}

func (bmach *Bondmachine) List_outputs() []string {
	result := []string{}
	if bmach.Outputs != 0 {
		for i := 0; i < bmach.Outputs; i++ {
			result = append(result, Get_output_name(i))
		}
	}
	return result
}

func (bmach *Bondmachine) Add_output() (string, error) {
	newbond := Bond{1, bmach.Outputs, 0}
	bmach.Outputs = bmach.Outputs + 1
	bmach.Internal_inputs = append(bmach.Internal_inputs, newbond)
	bmach.Links = append(bmach.Links, -1)
	return "Added output " + strconv.Itoa(bmach.Outputs-1) + " successfully", nil
}

func (bmach *Bondmachine) Del_output(oid int) error {
	if oid < bmach.Outputs {
		bondtoremove := Bond{1, oid, 0}

		// 1 - Remove the internal input
		// 1 - lower by 1 all ids > oid in bonds
		newinternal := make([]Bond, len(bmach.Internal_inputs)-1)
		newlinks := make([]int, len(bmach.Links)-1)
		j := 0
		for i, bond := range bmach.Internal_inputs {
			if bond.Map_to == 1 {
				if bond.Res_id == bondtoremove.Res_id {
					j--
				} else if bond.Res_id > bondtoremove.Res_id {
					newinternal[j] = Bond{1, bond.Res_id - 1, 0}
					newlinks[j] = bmach.Links[i]
				} else {
					newinternal[j] = bond
					newlinks[j] = bmach.Links[i]
				}
			} else {
				newinternal[j] = bond
				newlinks[j] = bmach.Links[i]
			}
			j++
		}

		bmach.Internal_inputs = newinternal
		bmach.Links = newlinks

		bmach.Outputs--
	} else {
		return Prerror{"Output id outside limit"}
	}
	return nil
}

func (bmach *Bondmachine) List_processors() string {
	result := ""
	if len(bmach.Processors) != 0 {
		for i, _ := range bmach.Processors {
			result += fmt.Sprintf("----------\nProcessor %03d\nP%d", i, i) + fmt.Sprintf("\n")
		}
	} else {
		result += "No defined processor"
	}
	return result
}

func (bmach *Bondmachine) EnumProcessors() int {
	return len(bmach.Processors)
}

func (bmach *Bondmachine) EnumBonds() int {
	result := 0
	for _, p := range bmach.Links {
		if p != -1 {
			result++
		}
	}
	return result
}

func (bmach *Bondmachine) Add_processor(dom_id int) (string, error) {
	if dom_id >= len(bmach.Domains) {
		return "", Prerror{"Non existent domain"}
	}
	inps := int(bmach.Domains[dom_id].N)
	outs := int(bmach.Domains[dom_id].M)
	for i := 0; i < inps; i++ {
		newbond := Bond{2, len(bmach.Processors), i}
		bmach.Internal_inputs = append(bmach.Internal_inputs, newbond)
		bmach.Links = append(bmach.Links, -1)
	}
	for i := 0; i < outs; i++ {
		newbond := Bond{3, len(bmach.Processors), i}
		bmach.Internal_outputs = append(bmach.Internal_outputs, newbond)
	}
	bmach.Processors = append(bmach.Processors, dom_id)

	newsolist := make([]int, 0)

	bmach.Shared_links = append(bmach.Shared_links, newsolist)

	return "Processor " + strconv.Itoa(len(bmach.Processors)-1) + " successfully added", nil
}

func (bmach *Bondmachine) List_bonds() map[int]string {
	result := make(map[int]string)
	if len(bmach.Links) != 0 {
		for i, linked := range bmach.Links {
			if linked != -1 {
				result[i] = bmach.Internal_outputs[linked].String() + "," + bmach.Internal_inputs[i].String()
			}
		}
	}
	return result
}

func (bmach *Bondmachine) Add_bond(endpoints []string) {
	for i, inp := range bmach.Internal_inputs {
		if inp.String() == endpoints[0] {
			for j, outp := range bmach.Internal_outputs {
				if outp.String() == endpoints[1] {
					bmach.Links[i] = j
					break
				}
			}
			break
		}

		if inp.String() == endpoints[1] {
			for j, outp := range bmach.Internal_outputs {
				if outp.String() == endpoints[0] {
					bmach.Links[i] = j
					break
				}
			}
			break
		}
	}
}

func (bmach *Bondmachine) Del_bond(bid int) error {
	if bid < len(bmach.Links) {
		bmach.Links[bid] = -1
	} else {
		return Prerror{"Bond id outside limit"}
	}
	return nil
}

func (bmach *Bondmachine) List_internal_inputs() []string {
	result := []string{}
	if len(bmach.Internal_inputs) != 0 {
		for i, _ := range bmach.Internal_inputs {
			result = append(result, bmach.Internal_inputs[i].String())
		}
	}
	return result
}

func (bmach *Bondmachine) List_internal_outputs() []string {
	result := []string{}
	if len(bmach.Internal_outputs) != 0 {
		for i, _ := range bmach.Internal_outputs {
			result = append(result, bmach.Internal_outputs[i].String())
		}
	}
	return result
}

func (bmach *Bondmachine) List_shared_objects() string {
	result := ""
	if len(bmach.Shared_objects) != 0 {
		for i, so := range bmach.Shared_objects {
			result += fmt.Sprintf("%03d - ", i) + so.String() + "\n"
		}
	} else {
		result += "No defined shared object"
	}
	return result
}

func (bmach *Bondmachine) Add_shared_objects(sos []string) {
	for _, so := range sos {
		// TODO Error repotting
		//		loaded := false
		for _, shr := range Allshared {
			if inst, ok := shr.Instantiate(so); ok {
				bmach.Shared_objects = append(bmach.Shared_objects, inst)
				//		loaded = true
				break
			}
		}
		//		if !loaded {
		//			result += "How to make a shared object from \"" + so + "\" is unkwown"
		//		}
	}
}

func (bmach *Bondmachine) List_processor_shared_object_links() string {
	result := ""
	for proc_id, curr_links := range bmach.Shared_links {
		result += "Processor " + strconv.Itoa(proc_id) + "\n\t"
		for _, so_id := range curr_links {
			result += strconv.Itoa(so_id) + " "
		}
		result += "\n"
	}
	return result
}

func (bmach *Bondmachine) Connect_processor_shared_object(endpoints []string) {
	if proc_id, ok := strconv.Atoi(endpoints[0]); ok == nil {
		if so_id, ok := strconv.Atoi(endpoints[1]); ok == nil {
			if len(bmach.Shared_links) > proc_id {
				curr_links := bmach.Shared_links[proc_id]
				already := false
				for _, link := range curr_links {
					if link == so_id {
						already = true
						break
					}
				}
				if !already {
					curr_links = append(curr_links, so_id)
					bmach.Shared_links[proc_id] = curr_links
				}
			}
		}
	}
}

func (bmach *Bondmachine) GetProgramsAlias() ([]string, error) {

	if len(bmach.Processors) != 0 {
		result := make([]string, 0)
		// TODO Make the code better
		for _, dom_id := range bmach.Processors {
			mymachine := bmach.Domains[dom_id]
			apstr, _ := mymachine.Program_alias()
			result = append(result, apstr)
		}

		return result, nil
	} else {
		return nil, Prerror{"Undefined processors"}
	}
	return nil, nil
}

func (bmach *Bondmachine) Attach_benchmark_core(endpoints []string) error {
	var e0 string
	var e1 string

	for _, outp := range bmach.Internal_outputs {
		if outp.String() == endpoints[0] {
			e0 = endpoints[0]
		}
		if outp.String() == endpoints[1] {
			e1 = endpoints[1]
		}
	}

	if e0 != "" && e1 != "" {
		mybcore := new(procbuilder.Machine)

		myarch := &mybcore.Arch

		myarch.Rsize = uint8(bmach.Rsize)

		myarch.Modes = make([]string, 1)
		myarch.Modes[0] = "ha"

		opcodes := make([]procbuilder.Opcode, 0)

		for _, op := range procbuilder.Allopcodes {
			for _, opn := range []string{"sic", "r2o", "j"} {
				if opn == op.Op_get_name() {
					opcodes = append(opcodes, op)
					break
				}
			}
		}

		sort.Sort(procbuilder.ByName(opcodes))

		myarch.Op = opcodes

		myarch.R = uint8(1)
		myarch.L = uint8(0)
		myarch.N = uint8(2)
		myarch.M = uint8(1)
		myarch.O = uint8(2)
		myarch.Shared_constraints = ""

		prog := "sic r0 i0\nsic r0 i1\nr2o r0 o0\nj 0\n"

		if prog, err := myarch.Assembler([]byte(prog)); err == nil {
			mybcore.Program = prog
		} else {
			fmt.Println(err)
		}

		bmach.Domains = append(bmach.Domains, mybcore)
		bmach.Add_processor(len(bmach.Domains) - 1)
		newpnum := strconv.Itoa(len(bmach.Processors) - 1)
		bmach.Add_bond([]string{"p" + newpnum + "i0", e0})
		bmach.Add_bond([]string{"p" + newpnum + "i1", e1})
		bmach.Add_output()
		newonum := strconv.Itoa(bmach.Outputs - 1)
		bmach.Add_bond([]string{"p" + newpnum + "o0", "o" + newonum})
	} else {
		return Prerror{"Benchmark core endpoints has to be internal outputs"}
	}
	return nil
}

func (bmach *Bondmachine) Dot(conf *Config, prefix string, vm *VM, oldvmstate *VM) string {
	// TODO Finish the layout building
	result := ""
	strutti := 0
	if prefix == "" {
		result += "digraph callgraph {\n"
		result += "\tbgcolor=\"#a2a6a7\";\n"
		result += "\tcompound=true;\n"
		result += "\tnode [fontname=\"monospace\"];\n"
		result += "\tfontname=\"Monospace\";\n"
	}

	for i, dom_id := range bmach.Processors {
		result += "\tsubgraph cluster_" + prefix + "_p" + strconv.Itoa(i) + " {\n"
		if GV_config(GVCLUSPROC) != "" {
			result += "\t" + GV_config(GVCLUSPROC) + ";\n"
		}
		if conf.BMinfo != nil {
			if conf.BMinfo.CPNames != nil {
				result += "\tlabel=\"" + conf.BMinfo.CPNames[i] + "\";\n"
			} else {
				result += "\t\tlabel=\"Processor " + strconv.Itoa(i) + "\";\n"
			}
		} else {
			result += "\t\tlabel=\"Processor " + strconv.Itoa(i) + "\";\n"
		}

		if conf.Dotdetail > 3 {
			result += "\t\tsubgraph cluster_" + prefix + "_p" + strconv.Itoa(i) + "_specs {\n"
			result += "\t\tlabel=\"Specs\";\n"

			// Register information
			reg_num := 1 << bmach.Domains[dom_id].R
			result += "\t\tnode [label=\"Registers: " + strconv.Itoa(reg_num) + "\"" + GV_config(GVINFOPROCREGS) + "] " + prefix + "p" + strconv.Itoa(i) + "procinforegs;\n"

			// Opcodes
			if conf.Dotdetail > 3 {

				oplist := "Opcodes:"
				j := 0
				for _, op := range bmach.Domains[dom_id].Op {
					if j%3 == 0 {
						oplist += "\\n" + op.Op_get_name()
					} else {
						oplist += ", " + op.Op_get_name()
					}
					j++
				}

				result += "\t\tnode [label=\"" + oplist + "\"" + GV_config(GVINFOPROCOPCODES) + "] " + prefix + "p" + strconv.Itoa(i) + "procinfocodes;\n"
			}

			// Program
			if conf.Dotdetail > 4 {

				if disass_text, err := bmach.Domains[dom_id].Disassembler(); err == nil {
					result += "\t\tsubgraph cluster_" + prefix + "_p" + strconv.Itoa(i) + "_prog {\n"
					result += "\t\tlabel=\"Program:\";\n"
					result += "\t\tstruct" + prefix + strconv.Itoa(strutti) + " [label=<<TABLE BORDER=\"0\">\n"
					strutti++
					for j, line := range strings.Split(disass_text, "\n") {
						if line != "" {
							if vm != nil && vm.Processors[i].Pc == uint64(j) {
								result += "\t\t<TR><TD BGCOLOR=\"red\">" + line + "</TD></TR>\n"
							} else {
								result += "\t\t<TR><TD>" + line + "</TD></TR>\n"
							}
						}
					}
					result += "\t\t</TABLE>>];\n"
					result += "\t\t}\n"
				}
			}

			result += "\t\t}\n"
		}

		if vm != nil {
			result += "\t\tsubgraph cluster_" + prefix + "_p" + strconv.Itoa(i) + "_vm_pc {\n"
			result += "\t\t\tlabel=\"PC:\";\n"
			result += "\t\tnode [label=\"" + strconv.Itoa(int(vm.Processors[i].Pc)) + "\"" + GV_config(GVINFOPROCPC) + "] " + prefix + "p" + strconv.Itoa(i) + "_pc;\n"
			result += "\t\t}\n"

			result += "\t\tsubgraph cluster_" + prefix + "_p" + strconv.Itoa(i) + "_vm_regs {\n"
			result += "\t\t\tlabel=\"Registers:\";\n"

			result += "\t\tstruct" + prefix + strconv.Itoa(strutti) + " [label=<<TABLE BORDER=\"0\">\n"
			strutti++

			for j, reg := range vm.Processors[i].Registers {
				colorstring := ""
				if oldvmstate != nil && oldvmstate.Processors[i].Registers[j] != reg {
					colorstring = " BGCOLOR=\"red\""
				}
				switch vm.Bmach.Rsize {
				case 8:
					result += "\t\t<TR><TD" + colorstring + ">" + procbuilder.Get_register_name(j) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint8)))) + "</TD></TR>\n"
				case 16:
					result += "\t\t<TR><TD" + colorstring + ">" + procbuilder.Get_register_name(j) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint16)))) + "</TD></TR>\n"
				default:
					// TODO Fix
				}
			}

			result += "\t\t</TABLE>> " + GV_config(GVINFOPROCPROG) + "];\n"
			result += "\t\t}\n"

		}

		inps := int(bmach.Domains[dom_id].N)
		outs := int(bmach.Domains[dom_id].M)
		result += "\t\tsubgraph cluster_" + prefix + "_p" + strconv.Itoa(i) + "_inputs {\n"
		if GV_config(GVCLUSININPROC) != "" {
			result += "\t\t" + GV_config(GVCLUSININPROC) + ";\n"
		}
		result += "\t\t\tlabel=\"P" + strconv.Itoa(i) + " inputs\";\n"

		for j := 0; j < inps; j++ {
			if vm != nil {
				result += "\t\tsubgraph cluster_" + prefix + "_p" + strconv.Itoa(i) + "_input_" + strconv.Itoa(j) + " {\n"
				//result += "\t\t\tlabel=\"p" + strconv.Itoa(i) + "i" + strconv.Itoa(j) + ":\";\n"
				result += "\t\t\tlabel=\"\";\n"
				result += "\t\t\tnode [label=\"p" + strconv.Itoa(i) + "i" + strconv.Itoa(j) + "\" " + GV_config(GVNODEININPROC) + "] " + prefix + "p" + strconv.Itoa(i) + "i" + strconv.Itoa(j) + ";\n"

				result += "\t\tstruct" + prefix + strconv.Itoa(strutti) + " [label=<<TABLE BORDER=\"0\">\n"
				strutti++

				inp := vm.Processors[i].Inputs[j]

				colorstring := ""
				if oldvmstate != nil && oldvmstate.Processors[i].Inputs[j] != inp {
					colorstring = " BGCOLOR=\"red\""
				}
				switch vm.Bmach.Rsize {
				case 8:
					result += "\t\t<TR><TD" + colorstring + ">" + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(inp.(uint8)))) + "</TD></TR>\n"
				case 16:
					result += "\t\t<TR><TD" + colorstring + ">" + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(inp.(uint16)))) + "</TD></TR>\n"
				default:
					// TODO Fix
				}

				result += "\t\t</TABLE>> " + GV_config(GVINFOPROCPROG) + "];\n"
				result += "\t\t}\n"

			} else {
				result += "\t\t\tnode [label=\"p" + strconv.Itoa(i) + "i" + strconv.Itoa(j) + "\" " + GV_config(GVNODEININPROC) + "] " + prefix + "p" + strconv.Itoa(i) + "i" + strconv.Itoa(j) + ";\n"
			}
		}

		result += "\t\t}\n"
		result += "\t\tsubgraph cluster_" + prefix + "_p" + strconv.Itoa(i) + "_outputs {\n"
		if GV_config(GVCLUSOUTINPROC) != "" {
			result += "\t\t" + GV_config(GVCLUSOUTINPROC) + ";\n"
		}
		result += "\t\t\tlabel=\"p" + strconv.Itoa(i) + " outputs\";\n"
		for j := 0; j < outs; j++ {
			if vm != nil {
				result += "\t\tsubgraph cluster_" + prefix + "_p" + strconv.Itoa(i) + "_output_" + strconv.Itoa(j) + " {\n"
				//result += "\t\t\tlabel=\"p" + strconv.Itoa(i) + "i" + strconv.Itoa(j) + ":\";\n"
				result += "\t\t\tlabel=\"\";\n"
				result += "\t\tnode [label=\"p" + strconv.Itoa(i) + "o" + strconv.Itoa(j) + "\" " + GV_config(GVNODEOUTINPROC) + "] " + prefix + "p" + strconv.Itoa(i) + "o" + strconv.Itoa(j) + ";\n"

				result += "\t\tstruct" + prefix + strconv.Itoa(strutti) + " [label=<<TABLE BORDER=\"0\">\n"
				strutti++

				outp := vm.Processors[i].Outputs[j]

				colorstring := ""
				if oldvmstate != nil && oldvmstate.Processors[i].Outputs[j] != outp {
					colorstring = " BGCOLOR=\"red\""
				}
				switch vm.Bmach.Rsize {
				case 8:
					result += "\t\t<TR><TD" + colorstring + ">" + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(outp.(uint8)))) + "</TD></TR>\n"
				case 16:
					result += "\t\t<TR><TD" + colorstring + ">" + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(outp.(uint16)))) + "</TD></TR>\n"
				default:
					// TODO Fix
				}

				result += "\t\t</TABLE>> " + GV_config(GVINFOPROCPROG) + "];\n"
				result += "\t\t}\n"

			} else {
				result += "\t\tnode [label=\"p" + strconv.Itoa(i) + "o" + strconv.Itoa(j) + "\" " + GV_config(GVNODEOUTINPROC) + "] " + prefix + "p" + strconv.Itoa(i) + "o" + strconv.Itoa(j) + ";\n"
			}
		}
		result += "\t\t}\n"

		if bmach.Shared_links != nil && len(bmach.Shared_links[i]) > 0 {

			seq := make(map[string]int)
			subresult := make(map[string]string)

			for _, so := range bmach.Shared_links[i] {
				sname := bmach.Shared_objects[so].Shortname()
				if j, ok := seq[sname]; ok {
					subresult[sname] += "\t\tnode [label=\"p" + strconv.Itoa(i) + sname + strconv.Itoa(j) + "\" " + bmach.Shared_objects[so].GV_config(GVNODEINPROC) + "] " + prefix + "p" + strconv.Itoa(i) + sname + strconv.Itoa(j) + ";\n"
				} else {
					seq[sname] = 0
					subresult[sname] = "\t\tsubgraph cluster_" + prefix + "_p" + strconv.Itoa(i) + "_" + bmach.Shared_objects[so].Shr_get_name() + " {\n"
					if bmach.Shared_objects[so].GV_config(GVCLUSINPROC) != "" {
						subresult[sname] += "\t\t" + bmach.Shared_objects[so].GV_config(GVCLUSINPROC) + ";\n"
					}
					subresult[sname] += "\t\t\tlabel=\"P" + strconv.Itoa(i) + " " + bmach.Shared_objects[so].Shr_get_name() + "\";\n"
					subresult[sname] += "\t\t\tnode [label=\"p" + strconv.Itoa(i) + sname + strconv.Itoa(j) + "\" " + bmach.Shared_objects[so].GV_config(GVNODEINPROC) + "] " + prefix + "p" + strconv.Itoa(i) + sname + strconv.Itoa(j) + ";\n"
				}
				seq[sname]++
			}

			for _, subres := range subresult {
				result += subres
				result += "\t\t}\n"
			}

		}

		result += "\t}\n"
	}

	result += "\tsubgraph cluster_" + prefix + "_inputs {\n"
	if GV_config(GVCLUSIN) != "" {
		result += "\t" + GV_config(GVCLUSIN) + ";\n"
	}
	result += "\t\tlabel=\"Inputs\";\n"
	for i := 0; i < bmach.Inputs; i++ {
		result += "\t\tnode [label=\"i" + strconv.Itoa(i) + "\" " + GV_config(GVNODEIN) + "] " + prefix + "i" + strconv.Itoa(i) + ";\n"
	}
	result += "\t}\n"

	result += "\tsubgraph cluster_" + prefix + "_outputs {\n"
	if GV_config(GVCLUSOUT) != "" {
		result += "\t" + GV_config(GVCLUSOUT) + ";\n"
	}
	result += "\t\tlabel=\"Outputs\";\n"
	for i := 0; i < bmach.Outputs; i++ {
		result += "\t\tnode [label=\"o" + strconv.Itoa(i) + "\" " + GV_config(GVNODEOUT) + "] " + prefix + "o" + strconv.Itoa(i) + ";\n"
	}
	result += "\t}\n"

	if len(bmach.Shared_objects) > 0 {

		seqout := make([]int, len(bmach.Shared_objects))

		seq := make(map[string]int)
		subresult := make(map[string]string)

		for i, so := range bmach.Shared_objects {
			sname := so.Shortname()
			if j, ok := seq[sname]; ok {
				subresult[sname] += "\t\tnode [label=\"" + sname + strconv.Itoa(j) + "\" " + so.GV_config(GVNODE) + "] " + prefix + sname + strconv.Itoa(j) + ";\n"
			} else {
				seq[sname] = 0
				subresult[sname] = "\tsubgraph cluster_" + prefix + "_" + so.Shr_get_name() + " {\n"
				if so.GV_config(GVCLUS) != "" {
					subresult[sname] += "\t\t" + so.GV_config(GVCLUS) + ";\n"
				}
				subresult[sname] += "\t\tlabel=\"" + so.Shr_get_name() + "\";\n"
				subresult[sname] += "\t\tnode [label=\"" + sname + strconv.Itoa(j) + "\" " + so.GV_config(GVNODE) + "] " + prefix + sname + strconv.Itoa(j) + ";\n"
			}
			seqout[i] = seq[sname]
			seq[sname]++
		}

		for _, subres := range subresult {
			result += subres
			result += "\t}\n"
		}

		for i, _ := range bmach.Processors {

			if len(bmach.Shared_links[i]) > 0 {

				seq := make(map[string]int)

				for _, so := range bmach.Shared_links[i] {
					sname := bmach.Shared_objects[so].Shortname()
					if _, ok := seq[sname]; !ok {
						seq[sname] = 0
					}

					result += prefix + "p" + strconv.Itoa(i) + sname + strconv.Itoa(seq[sname]) + " -> " + prefix + sname + strconv.Itoa(seqout[so]) + "[" + bmach.Shared_objects[so].GV_config(GVEDGE) + "];\n"
					seq[sname]++
				}
			}
		}

	}

	for i, linked := range bmach.Links {
		if linked != -1 {
			result += prefix + bmach.Internal_outputs[linked].String() + " -> " + prefix + bmach.Internal_inputs[i].String() + ";\n"
		}
	}

	if prefix == "" {
		result += "}\n"
	}

	return result
}
