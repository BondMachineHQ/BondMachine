package bondgo

const (
	MAX_REGS     = 65536
	MAX_MEMORY   = 65536
	MAX_CHANNELS = 65536
	MAX_INPUTS   = 16
	MAX_OUTPUTS  = 16
	MAXPEERID    = uint32(10000)
)

type BondgoCheck struct {
	// The results
	*BondgoResults // The generated objects

	// Globally
	*BondgoConfig // Global compiler config

	*BondgoRequirements // The pointer to the requirements struct
	*BondgoRuninfo      // Running data

	*BondgoMessages // Compiler messages and errors

	*BondgoFunctions // The map of declared functions in all scopes

	Used    chan UsageNotify // Used to notify the used resource
	Reqs    chan VarReq      // Variable request
	Answers chan VarAns      // Variable response

	// Scope wide
	Outer *BondgoCheck // The previous scope
	Clean *BondgoCheck // The scope that need to be clean by its next brother aka when the scope ends

	Vars    map[string]VarCell // The map of variable in the current scope
	Returns []VarCell          // Returns for functions call

	CurrentLoop   string
	CurrentSwitch string
	CurrentDevice string

	CurrentRoutine int // Aka current processor
}

type Abs_assembly struct {
	ProcProgs []string
	Bonds     []string
}
