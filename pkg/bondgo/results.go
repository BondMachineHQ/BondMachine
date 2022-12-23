package bondgo

import (
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

type BondgoResults struct {
	Config  *BondgoConfig
	Proc    *procbuilder.Arch
	Bmach   *bondmachine.Bondmachine
	Program map[int]*BondgoRoutine // The generated assembly
}

func (rs *BondgoResults) Init_Results(cfg *BondgoConfig) {
	rs.Config = cfg
	rs.Proc = nil                             // Here will be stored the eventually generated processor
	rs.Bmach = nil                            // Here will be stored the eventually generated bondmachine
	rs.Program = make(map[int]*BondgoRoutine) // Map processor Lines of assembly code
	proccode := new(BondgoRoutine)
	proccode.Lines = make([]string, 0)
	rs.Program[0] = proccode
}

func (rs *BondgoResults) WriteLine(proc_id int, line string) {
	if proccode, ok := rs.Program[proc_id]; ok {
		proccode.Append(line)
	} else {
		proccode := new(BondgoRoutine)
		proccode.Lines = make([]string, 0)
		proccode.Append(line)
		rs.Program[proc_id] = proccode
	}
}

func (rs *BondgoResults) CountLines(proc_id int) int {
	if proccode, ok := rs.Program[proc_id]; ok {
		return len(proccode.Lines)
	}
	return 0
}

func (rs *BondgoResults) GetProgram(proc_id int) []string {
	if proccode, ok := rs.Program[proc_id]; ok {
		return proccode.Lines
	}
	return []string{}
}

func (rs *BondgoResults) Replacer(proc_id int, from string, to string) {
	if _, ok := rs.Program[proc_id]; ok {
		rs.Program[proc_id].Replacer(from, to)
	}
}

func (rs *BondgoResults) Checker(proc_id int, ck string) bool {
	if _, ok := rs.Program[proc_id]; ok {
		return rs.Program[proc_id].Checker(ck)
	}
	return false
}

func (rs *BondgoResults) Shift_program_location(proc_id int, n int) {
	rs.Program[proc_id].Shift_program_location(n)
}

func (rs *BondgoResults) Write_assembly(proc_id int) string {
	result := ""
	rs.Program[proc_id].Remove_program_location()
	for _, line := range rs.Program[proc_id].Lines {
		result += line + "\n"
	}
	return result
}
