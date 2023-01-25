package procbuilder

import (
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Opcode interface {
	Op_get_name() string
	Op_get_desc() string
	Op_show_assembler(*Arch) string
	Op_get_instruction_len(*Arch) int
	OpInstructionVerilogHeader(*Config, *Arch, string, string) string
	Op_instruction_verilog_reset(*Arch, string) string
	Op_instruction_verilog_internal_state(*Arch, string) string
	Op_instruction_verilog_default_state(*Arch, string) string
	Op_instruction_verilog_state_machine(*Arch, string) string
	Op_instruction_verilog_footer(*Arch, string) string
	Op_instruction_verilog_extra_modules(*Arch, string) ([]string, []string)
	Op_instruction_verilog_extra_block(*Arch, string, uint8, string, []string) string
	AbstractAssembler(*Arch, []string) ([]UsageNotify, error)
	Assembler(*Arch, []string) (string, error)
	HLAssemblerMatch(*Arch) []string
	HLAssemblerNormalize(*Arch, *bmreqs.ReqRoot, string, *bmline.BasmLine) (*bmline.BasmLine, error)
	Disassembler(*Arch, string) (string, error)
	Simulate(*VM, string) error
	Generate(*Arch) string
	Required_shared() (bool, []string)
	Required_modes() (bool, []string)
	Forbidden_modes() (bool, []string)
	ExtraFiles(arch *Arch) ([]string, []string)
}

type Sharedel interface {
	Shr_get_name() string
	Shortname() string
	GetArchHeader(*Arch, string, int) string // returns the architecture header for the shared element
	GetArchParams(*Arch, string, int) string // returns the architecture module parameters for the shared element
	GetCPParams(*Arch, string, int) string   // returns the processor (CP) module internal parameters for the shared element
}

type Prerror struct {
	string
}

func (e Prerror) Error() string {
	return e.string
}

type ByName []Opcode

func (op ByName) Len() int           { return len(op) }
func (op ByName) Swap(i, j int)      { op[i], op[j] = op[j], op[i] }
func (op ByName) Less(i, j int) bool { return op[i].Op_get_name() < op[j].Op_get_name() }
