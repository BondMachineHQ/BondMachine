package procbuilder

const (
	C_OPCODE = uint8(0) + iota
	C_REGSIZE
	C_INPUT
	C_OUTPUT
	C_ROMSIZE
	C_RAMSIZE
	C_SHAREDOBJECT
	C_CONNECTED
)

const (
	I_NIL = 0
)
const (
	S_NIL = ""
)

type UsageNotify struct {
	ComponentType uint8
	Components    string
	Componenti    int
}
