package bondmachine

const (
	bmapiCommands = `package bmapiusbuart

const (
	cmdNEWVAL  = uint8(0)   // 000 00000
	cmdDVALIDH = uint8(32)  // 001 00000
	cmdDVALIDL = uint8(64)  // 010 00000
	cmdDRECVH  = uint8(96)  // 011 00000
	cmdDRECVL  = uint8(128) // 100 00000
	cmdHANDSH  = uint8(160) // 101 00000
	cmdKEEP    = uint8(192) // 110 00000

	cmdMASK = uint8(224) // 111 00000
)
`
)
