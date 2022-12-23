package bondgo

type BondgoConfig struct {
	Debug          bool
	Verbose        bool
	Mpm            bool
	MaxRegs        int
	Rsize          uint8
	Basic_type     string
	Basic_chantype string
	Cascading_io   bool
}

func (db *BondgoConfig) In_debug() bool {
	return db.Debug
}
