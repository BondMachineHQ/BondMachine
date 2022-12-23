package bondgo

type BondgoMessages struct {
	Config       *BondgoConfig
	Error_status bool     // The compiler general status
	Message_list []string // The list of compiler messages
}

func (b *BondgoMessages) Init_Messages(cfg *BondgoConfig) {
	b.Message_list = make([]string, 0)
	b.Error_status = false
}

func (b *BondgoMessages) Log(line string) {
	b.Message_list = append(b.Message_list, line)
}

func (b *BondgoMessages) Warning(line string) {
	b.Message_list = append(b.Message_list, "Warning: "+line)
}

func (b *BondgoMessages) Set_faulty(line string) {
	b.Message_list = append(b.Message_list, "Error: "+line)
	b.Error_status = true
}

func (b *BondgoMessages) Is_faulty() bool {
	return b.Error_status
}

func (b *BondgoMessages) Dump_log() string {
	result := ""
	for _, line := range b.Message_list {
		result += line + "\n"
	}
	return result
}
