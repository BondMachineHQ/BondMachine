package bondmachine

type Shared_element interface {
	Shr_get_name() string // The name
	Shr_get_desc() string // A description
	Shortname() string
	GV_config(uint8) string
	Instantiate(string) (Shared_instance, bool)
}

type Shared_instance_list []int

type Shared_instance interface {
	String() string
	Shortname() string
	Shr_get_name() string
	GV_config(uint8) string
	Write_verilog(*Bondmachine, int, string, string) string
	GetPerProcPortsWires(*Bondmachine, int, int, string) string
	GetPerProcPortsHeader(*Bondmachine, int, int, string) string
	GetExternalPortsHeader(*Bondmachine, int, int, string) string
	GetExternalPortsWires(*Bondmachine, int, int, string) string
	GetCPSharedPortsHeader(*Bondmachine, int, string) string // The ports on the Shared element shared by all the cores (header)
	GetCPSharedPortsWires(*Bondmachine, int, string) string  // The ports on the Shared element shared by all the cores (module)
}
