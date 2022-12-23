package bondgo

import (
	"strconv"
)

const (
	TR_PROC = uint8(0) + iota
	TR_CHAN
	TR_SHAREDMEM
	TR_EXIT
)

const (
	C_OPCODE = uint8(0) + iota
	C_REGSIZE
	C_INPUT
	C_OUTPUT
	C_ROMSIZE
	C_RAMSIZE
	C_SHAREDOBJECT
	C_CONNECTED
	C_DEVICE
)

const (
	I_NIL = 0
)
const (
	S_NIL = ""
)

type ProcRequirements struct {
	Opcodes       []string
	Registersize  int
	Inputs        int
	Outputs       int
	Romsize       int
	Ramsize       int
	SharedObjects []string
	Device        string
}

type IORequirements struct {
	Inputs_ids  []int
	Outputs_ids []int
}

// Maybe a rewrite in term of SO (with interfaces) is need.
type ChanRequirements struct {
	Connected []int
}

type SharedMemRequirements struct {
}

type BondgoRequirements struct {
	Config *BondgoConfig // Global compiler config
	Procr  map[int]*ProcRequirements
	IOr    map[int]*IORequirements
	Chanr  map[int]*ChanRequirements
	Shrdr  map[int]*SharedMemRequirements
}

type UsageNotify struct {
	TargetType    uint8
	TargetId      int
	ComponentType uint8
	Components    string
	Componenti    int
}

func (reqmnt *ProcRequirements) Dump_Requirements() string {
	result := "Opcodes: "
	for i, opcode := range reqmnt.Opcodes {
		result += opcode
		if i != len(reqmnt.Opcodes)-1 {
			result += ","
		} else {
			result += "\n"
		}
	}
	result += "Registersize: " + strconv.Itoa(reqmnt.Registersize) + "\n"
	result += "Inputs: " + strconv.Itoa(reqmnt.Inputs) + "\n"
	result += "Outputs: " + strconv.Itoa(reqmnt.Outputs) + "\n"
	result += "Romsize: " + strconv.Itoa(reqmnt.Romsize) + "\n"
	result += "Ramsize: " + strconv.Itoa(reqmnt.Ramsize) + "\n"
	result += "Device: " + reqmnt.Device + "\n"
	result += "Shared Objects: "
	for i, so := range reqmnt.SharedObjects {
		result += so
		if i != len(reqmnt.SharedObjects)-1 {
			result += ","
		}
	}
	result += "\n"

	return result
}

func (reqmnt *IORequirements) Dump_Requirements() string {
	result := "Inputs: "
	for i, con := range reqmnt.Inputs_ids {
		result += strconv.Itoa(con)
		if i != len(reqmnt.Inputs_ids)-1 {
			result += ","
		} else {
			result += "\n"
		}
	}

	result += "Outputs: "

	for i, con := range reqmnt.Outputs_ids {
		result += strconv.Itoa(con)
		if i != len(reqmnt.Outputs_ids)-1 {
			result += ","
		} else {
			result += "\n"
		}
	}

	return result
}

func (reqmnt *ChanRequirements) Dump_Requirements() string {
	result := "Connected: "
	for i, con := range reqmnt.Connected {
		result += strconv.Itoa(con)
		if i != len(reqmnt.Connected)-1 {
			result += ","
		} else {
			result += "\n"
		}
	}
	result += "\n"

	return result
}

func (reqmnt *BondgoRequirements) Init_Requirements(cfg *BondgoConfig) {
	reqmnt.Config = cfg
	reqmnt.Procr = make(map[int]*ProcRequirements)
	reqmnt.IOr = make(map[int]*IORequirements)
	reqmnt.Chanr = make(map[int]*ChanRequirements)
	reqmnt.Shrdr = make(map[int]*SharedMemRequirements)
}

func (reqmnt *BondgoRequirements) Dump_Requirements() string {
	result := ""
	result += "--- Processors ---\n"
	for proc_id, procreq := range reqmnt.Procr {
		result += "p" + strconv.Itoa(proc_id) + "\n"
		result += procreq.Dump_Requirements()
	}
	result += "--- IO ---\n"
	for proc_id, ioreq := range reqmnt.IOr {
		result += "proc " + strconv.Itoa(proc_id) + "\n"
		result += ioreq.Dump_Requirements()
	}
	result += "--- Channels ---\n"
	for chan_id, chanreq := range reqmnt.Chanr {
		result += "ch" + strconv.Itoa(chan_id) + "\n"
		result += chanreq.Dump_Requirements()
	}
	result += "--- Shared Memory ---\n"
	// TODO
	return result
}

func (reqmnt *BondgoRequirements) Usage_Monitor(useditem chan UsageNotify, usagedone chan bool) {
	//debug := reqmnt.Config.Debug
UB:
	for {
		notif := <-useditem

		targettype := notif.TargetType
		targetid := notif.TargetId
		componenttype := notif.ComponentType
		components := notif.Components
		componenti := notif.Componenti

		switch targettype {
		case TR_PROC:
			var proc *ProcRequirements
			if exists, ok := reqmnt.Procr[targetid]; ok {
				proc = exists
			} else {
				proc = new(ProcRequirements)
				reqmnt.Procr[targetid] = proc
			}

			switch componenttype {
			case C_OPCODE:
				present := false
				for _, op := range proc.Opcodes {
					if op == components {
						present = true
						break
					}
				}
				if !present {
					proc.Opcodes = append(proc.Opcodes, components)
				}
			case C_REGSIZE:
				if componenti > proc.Registersize {
					proc.Registersize = componenti
				}
			case C_INPUT:

				var ior *IORequirements
				if exists, ok := reqmnt.IOr[targetid]; ok {
					ior = exists
				} else {
					ior = new(IORequirements)
					ior.Inputs_ids = make([]int, 0)
					ior.Outputs_ids = make([]int, 0)
					reqmnt.IOr[targetid] = ior
				}

				present := false
				for _, inp := range ior.Inputs_ids {
					if inp == componenti {
						present = true
						break
					}
				}

				if !present {
					ior.Inputs_ids = append(ior.Inputs_ids, componenti)
					proc.Inputs = len(ior.Inputs_ids)
				}

			case C_OUTPUT:

				var ior *IORequirements
				if exists, ok := reqmnt.IOr[targetid]; ok {
					ior = exists
				} else {
					ior = new(IORequirements)
					ior.Inputs_ids = make([]int, 0)
					ior.Outputs_ids = make([]int, 0)
					reqmnt.IOr[targetid] = ior
				}

				present := false
				for _, outp := range ior.Outputs_ids {
					if outp == componenti {
						present = true
						break
					}
				}

				if !present {
					ior.Outputs_ids = append(ior.Outputs_ids, componenti)
					proc.Outputs = len(ior.Outputs_ids)
				}

			case C_ROMSIZE:
				if componenti > proc.Romsize {
					proc.Romsize = componenti
				}
			case C_RAMSIZE:
				if componenti > proc.Ramsize {
					proc.Ramsize = componenti
				}
			case C_SHAREDOBJECT:
				proc.SharedObjects = append(proc.SharedObjects, components)
			case C_DEVICE:
				proc.Device = components
			}

			// TODO Other cases
		case TR_CHAN:
			var cchan *ChanRequirements
			if exists, ok := reqmnt.Chanr[targetid]; ok {
				cchan = exists
			} else {
				cchan = new(ChanRequirements)
				reqmnt.Chanr[targetid] = cchan
			}
			switch componenttype {
			case C_CONNECTED:
				present := false
				for _, op := range cchan.Connected {
					if op == componenti {
						present = true
						break
					}
				}
				if !present {
					cchan.Connected = append(cchan.Connected, componenti)
				}
			}

		case TR_EXIT:
			break UB
		}
	}
	usagedone <- true
}
