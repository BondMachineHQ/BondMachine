package bondgo

import (
	"fmt"
	"strconv"
)

const (
	REQ_REMOVE = uint8(0) + iota // Removes the resoure from the local used list, if associated with a global one the global is never removed
	REQ_NEW                      // New variable request (and eventually attach to a global one)
	REQ_ATTACH                   // Attach a local var to a global one
	REQ_EXIT
)

const (
	ANS_OK = uint8(0) + iota
	ANS_FAIL
)

type VarCells []VarCell

type BondgoRuninfo struct {
	Config     *BondgoConfig // Global compiler config
	Used_cells map[int]VarCells
	IO         []IOInfo
	Channels   []ChanInfo
	SharedRAM  []SharedRAMInfo
}

// IO topology
type IOInfo struct {
	Global_id int
	Inputs    []int
	Output    int
}

// bondmachine channels
type ChanInfo struct {
	Global_id int
	Connected []int
	Readers   []int
	Writers   []int
}

// bondmachine shared RAM
type SharedRAMInfo struct {
	Global_id int
	Connected []int
}

type VarReq struct {
	ReqType      uint8
	Processor_id int
	Cell         VarCell
}

type VarAns struct {
	AnsType uint8
	Cell    VarCell
}

func (vr *VarReq) String() string {
	result := ""
	switch vr.ReqType {
	case REQ_NEW:
		result += "new "
	case REQ_REMOVE:
		result += "del "
	}
	result += "proc " + strconv.Itoa(vr.Processor_id) + " "
	result += "cell " + vr.Cell.String()
	// TODO Finish
	return result
}

func (ri *BondgoRuninfo) Init_Runinfo(cfg *BondgoConfig) {
	ri.Config = cfg
	ri.Used_cells = make(map[int]VarCells)
	ri.IO = make([]IOInfo, 0)
	ri.Channels = make([]ChanInfo, 0)
	ri.SharedRAM = make([]SharedRAMInfo, 0)
}

// This goroutine assign or frees used memory within a processor
func (ri *BondgoRuninfo) Var_assigner(req chan VarReq, resp chan VarAns, useditem chan UsageNotify, assignerdone chan bool) {
	debug := ri.Config.Debug
	busylist := ri.Used_cells
	busychan := ri.Channels
	busyio := ri.IO
VA:
	for {
		r := <-req

		rtype := r.ReqType
		rproc := r.Processor_id
		rcell := r.Cell

		if debug {
			fmt.Println("\tList of resources in use: ", busylist)
			fmt.Println("\tOperation: ", &r)
		}

		switch rtype {
		case REQ_REMOVE:
			switch rcell.Procobjtype {
			case REGISTER, MEMORY:
				if gent, _ := Type_from_string(ri.Config.Basic_type); Same_Type(rcell.Vtype, gent) {
					if _, ok := busylist[rproc]; ok {
						blist := busylist[rproc]
						if i, ok := memused(r.Cell, blist); ok {
							if i < len(blist)-1 {
								blist = append(blist[:i], blist[i+1:]...)
							} else {
								blist = blist[:i]
							}
							busylist[rproc] = blist
							resp <- VarAns{ANS_OK, r.Cell}
						} else {
							panic("Attempt to remove an unused Memory cell")
						}
					} else {
						panic("Attempt to remove an unused Memory cell")
					}
				} else if gent, _ := Type_from_string("bool"); Same_Type(rcell.Vtype, gent) {
					if _, ok := busylist[rproc]; ok {
						blist := busylist[rproc]
						if i, ok := memused(r.Cell, blist); ok {
							if i < len(blist)-1 {
								blist = append(blist[:i], blist[i+1:]...)
							} else {
								blist = blist[:i]
							}
							busylist[rproc] = blist
							resp <- VarAns{ANS_OK, r.Cell}
						} else {
							panic("Attempt to remove an unused Memory cell")
						}
					} else {
						panic("Attempt to remove an unused Memory cell")
					}
				} else {
					panic("Allocator received a wrong type, this cannot happen. A bug is here")
				}
			case INPUT, OUTPUT, CHANNEL:
				// TODO Check and make better
				resp <- VarAns{ANS_OK, r.Cell}
			}
		case REQ_NEW:
			switch rcell.Procobjtype {
			case REGISTER:
				if gent, _ := Type_from_string(ri.Config.Basic_type); Same_Type(rcell.Vtype, gent) {
					created := false
					if _, ok := busylist[rproc]; !ok {
						vcells := make([]VarCell, 0)
						busylist[rproc] = vcells
					}
					for i := 0; i < MAX_REGS; i++ {
						guessed := VarCell{gent, REGISTER, i, 0, 0, 0, 0, 0}
						present := false
						for _, assigned := range busylist[rproc] {
							if assigned.Procobjtype == guessed.Procobjtype && assigned.Id == guessed.Id {
								present = true
								break
							}
						}
						if !present {
							resp <- VarAns{ANS_OK, guessed}
							useditem <- UsageNotify{TR_PROC, rproc, C_REGSIZE, S_NIL, i + 1}
							busylist[rproc] = append(busylist[rproc], guessed)
							created = true
							break
						}
					}

					if !created {
						panic("Recursion function not allowed")
					}

				} else if gent, _ := Type_from_string("bool"); Same_Type(rcell.Vtype, gent) {
					created := false
					if _, ok := busylist[rproc]; !ok {
						vcells := make([]VarCell, 0)
						busylist[rproc] = vcells
					}
					for i := 0; i < MAX_REGS; i++ {
						guessed := VarCell{gent, REGISTER, i, 0, 0, 0, 0, 0}
						present := false
						for _, assigned := range busylist[rproc] {
							if assigned.Procobjtype == guessed.Procobjtype && assigned.Id == guessed.Id {
								present = true
								break
							}
						}
						if !present {
							resp <- VarAns{ANS_OK, guessed}
							useditem <- UsageNotify{TR_PROC, rproc, C_REGSIZE, S_NIL, i + 1}
							busylist[rproc] = append(busylist[rproc], guessed)
							created = true
							break
						}
					}

					if !created {
						panic("Recursion function not allowed")
					}

				} else {
					panic("Allocator received a wrong type, this cannot happen. A bug is here")
				}
			case MEMORY:
				if gent, _ := Type_from_string(ri.Config.Basic_type); Same_Type(rcell.Vtype, gent) {
					created := false
					if _, ok := busylist[rproc]; !ok {
						vcells := make([]VarCell, 0)
						busylist[rproc] = vcells
					}
					for i := 0; i < MAX_MEMORY; i++ {
						//  TODO This code uses only 1 memory area per variable, remember whenever it will happen that other types will be inserted that this code has to be substituted with something else
						guessed := VarCell{gent, MEMORY, i, i, i, 0, 0, 0}
						present := false
						for _, assigned := range busylist[rproc] {
							if assigned.Procobjtype == guessed.Procobjtype && assigned.Id == guessed.Id {
								present = true
								break
							}
						}
						if !present {
							resp <- VarAns{ANS_OK, guessed}
							useditem <- UsageNotify{TR_PROC, rproc, C_RAMSIZE, S_NIL, i + 1}
							busylist[rproc] = append(busylist[rproc], guessed)
							created = true
							break
						}
					}

					if !created {
						panic("Recursion function not allowed")
					}

				} else if gent, _ := Type_from_string("bool"); Same_Type(rcell.Vtype, gent) {
					created := false
					if _, ok := busylist[rproc]; !ok {
						vcells := make([]VarCell, 0)
						busylist[rproc] = vcells
					}
					for i := 0; i < MAX_MEMORY; i++ {
						//  TODO This code uses only 1 memory area per variable, remember whenever it will happen that other types will be inserted that this code has to be substituted with something else
						guessed := VarCell{gent, MEMORY, i, i, i, 0, 0, 0}
						present := false
						for _, assigned := range busylist[rproc] {
							if assigned.Procobjtype == guessed.Procobjtype && assigned.Id == guessed.Id {
								present = true
								break
							}
						}
						if !present {
							resp <- VarAns{ANS_OK, guessed}
							useditem <- UsageNotify{TR_PROC, rproc, C_RAMSIZE, S_NIL, i + 1}
							busylist[rproc] = append(busylist[rproc], guessed)
							created = true
							break
						}
					}

					if !created {
						panic("Recursion function not allowed")
					}

				} else {
					panic("Allocator received a wrong type, this cannot happen. A bug is here")
				}
			case INPUT:
				if gent, _ := Type_from_string(ri.Config.Basic_type); Same_Type(rcell.Vtype, gent) {
					if _, ok := busylist[rproc]; !ok {
						vcells := make([]VarCell, 0)
						busylist[rproc] = vcells
					}

					// Check for IO consistency (Only if the IO is inittializated) and global data update
					consistent := true
					globalpresent := false

					if rcell.Global_id != 0 {
						for _, ioinfo := range busyio {
							if ioinfo.Global_id == rcell.Global_id {
								globalpresent = true
								if ioinfo.Output == rproc {
									// A processor cannot have and input connected to an output
									consistent = false
									break
								} else {
									ioinfo.Inputs = append(ioinfo.Inputs, rproc)
								}
							}
						}
						if !globalpresent {
							newioinputs := make([]int, 1)
							newioinputs[0] = rproc
							newio := IOInfo{rcell.Global_id, newioinputs, 0}
							busyio = append(busyio, newio)
						}
					}

					if consistent {
						for i := 0; i < MAX_INPUTS; i++ {
							guessed := VarCell{gent, INPUT, i, i, i, rcell.Global_id, rcell.Start_globalid, rcell.End_globalid}
							present := false
							for _, assigned := range busylist[rproc] {
								if assigned.Procobjtype == guessed.Procobjtype && assigned.Id == guessed.Id {
									present = true
									break
								}
							}
							if !present {
								resp <- VarAns{ANS_OK, guessed}
								// Only in the IO is inittializated its use has to be notified
								if rcell.Global_id != 0 {
									useditem <- UsageNotify{TR_PROC, rproc, C_INPUT, S_NIL, rcell.Global_id}
								}
								busylist[rproc] = append(busylist[rproc], guessed)
								break
							}
						}
					} else {
						resp <- VarAns{ANS_FAIL, VarCell{gent, 0, 0, 0, 0, 0, 0, 0}}
					}
				} else {
					panic("Allocator received a wrong type, this cannot happen. A bug is here")
				}
			case OUTPUT:
				if gent, _ := Type_from_string(ri.Config.Basic_type); Same_Type(rcell.Vtype, gent) {
					if _, ok := busylist[rproc]; !ok {
						vcells := make([]VarCell, 0)
						busylist[rproc] = vcells
					}

					// Check for IO consistency (Only if the IO is inittializated)
					consistent := true
					globalpresent := false

					if rcell.Global_id != 0 {
						for _, ioinfo := range busyio {
							if ioinfo.Global_id == rcell.Global_id {
								globalpresent = true
								for _, ioimp := range ioinfo.Inputs {
									if ioimp == rproc {
										// A processor cannot have and input connected to an output
										consistent = false
										break
									}
								}
								if ioinfo.Output != 0 {
									// There can be only one output for every global id
									consistent = false
									break
								} else {
									ioinfo.Output = rproc
								}
							}
						}
						if !globalpresent {
							newioinputs := make([]int, 0)
							newio := IOInfo{rcell.Global_id, newioinputs, rproc}
							busyio = append(busyio, newio)
						}
					}

					if consistent {

						for i := 0; i < MAX_OUTPUTS; i++ {
							guessed := VarCell{gent, OUTPUT, i, i, i, rcell.Global_id, rcell.Start_globalid, rcell.End_globalid}
							present := false
							for _, assigned := range busylist[rproc] {
								if assigned.Procobjtype == guessed.Procobjtype && assigned.Id == guessed.Id {
									present = true
									break
								}
							}
							if !present {
								resp <- VarAns{ANS_OK, guessed}
								// Only in the IO is inittializated its use has to be notified
								if rcell.Global_id != 0 {
									useditem <- UsageNotify{TR_PROC, rproc, C_OUTPUT, S_NIL, rcell.Global_id}
								}
								busylist[rproc] = append(busylist[rproc], guessed)
								break
							}
						}
					} else {
						resp <- VarAns{ANS_FAIL, VarCell{gent, 0, 0, 0, 0, 0, 0, 0}}
					}
				} else {
					panic("Allocator received a wrong type, this cannot happen. A bug is here")
				}
			case CHANNEL:
				if gent, _ := Type_from_string(ri.Config.Basic_chantype); Same_Type(rcell.Vtype, gent) {
					if _, ok := busylist[rproc]; !ok {
						vcells := make([]VarCell, 0)
						busylist[rproc] = vcells
					}

					guessed_global_id := 0
					created := false
					for i := 0; i < MAX_CHANNELS; i++ {
						guessed_global_id = i
						present := false
						for _, channel := range busychan {
							if channel.Global_id == guessed_global_id {
								present = true
								break
							}
						}
						if !present {
							connected := make([]int, 1)
							connected[0] = rproc
							readers := make([]int, 0)
							writers := make([]int, 0)
							busychan = append(busychan, ChanInfo{guessed_global_id, connected, readers, writers})
							useditem <- UsageNotify{TR_CHAN, guessed_global_id, C_CONNECTED, S_NIL, rproc}
							created = true
							break
						}
					}
					if created {
						created = false
						for i := 0; i < MAX_CHANNELS; i++ {
							guessed := VarCell{gent, CHANNEL, i, i, i, guessed_global_id, guessed_global_id, guessed_global_id}
							present := false
							for _, assigned := range busylist[rproc] {
								if assigned.Procobjtype == guessed.Procobjtype && assigned.Id == guessed.Id {
									present = true
									break
								}
							}
							if !present {
								resp <- VarAns{ANS_OK, guessed}
								useditem <- UsageNotify{TR_PROC, rproc, C_SHAREDOBJECT, "channel:", I_NIL}
								busylist[rproc] = append(busylist[rproc], guessed)
								useditem <- UsageNotify{TR_CHAN, guessed_global_id, C_CONNECTED, S_NIL, rproc}
								created = true
								break
							}
						}
						if !created {
							panic("Channel creation failed")
						}
					} else {
						panic("global channel id failed")
					}
				} else if gent, _ := Type_from_string("chan bool"); Same_Type(rcell.Vtype, gent) {
					if _, ok := busylist[rproc]; !ok {
						vcells := make([]VarCell, 0)
						busylist[rproc] = vcells
					}

					guessed_global_id := 0
					created := false
					for i := 0; i < MAX_CHANNELS; i++ {
						guessed_global_id = i
						present := false
						for _, channel := range busychan {
							if channel.Global_id == guessed_global_id {
								present = true
								break
							}
						}
						if !present {
							connected := make([]int, 1)
							connected[0] = rproc
							readers := make([]int, 0)
							writers := make([]int, 0)
							busychan = append(busychan, ChanInfo{guessed_global_id, connected, readers, writers})
							useditem <- UsageNotify{TR_CHAN, guessed_global_id, C_CONNECTED, S_NIL, rproc}
							created = true
							break
						}
					}
					if created {
						created = false
						for i := 0; i < MAX_CHANNELS; i++ {
							guessed := VarCell{gent, CHANNEL, i, i, i, guessed_global_id, guessed_global_id, guessed_global_id}
							present := false
							for _, assigned := range busylist[rproc] {
								if assigned.Procobjtype == guessed.Procobjtype && assigned.Id == guessed.Id {
									present = true
									break
								}
							}
							if !present {
								resp <- VarAns{ANS_OK, guessed}
								useditem <- UsageNotify{TR_PROC, rproc, C_SHAREDOBJECT, "channel:", I_NIL}
								busylist[rproc] = append(busylist[rproc], guessed)
								useditem <- UsageNotify{TR_CHAN, guessed_global_id, C_CONNECTED, S_NIL, rproc}
								created = true
								break
							}
						}
						if !created {
							panic("Channel creation failed")
						}
					} else {
						panic("global channel id failed")
					}
				}
			}
		case REQ_ATTACH:
			switch rcell.Procobjtype {
			case CHANNEL:
				if gent, _ := Type_from_string(ri.Config.Basic_chantype); Same_Type(rcell.Vtype, gent) {
					if _, ok := busylist[rproc]; !ok {
						vcells := make([]VarCell, 0)
						busylist[rproc] = vcells
					}

					guessed_global_id := rcell.Global_id
					created := false
					for i := 0; i < MAX_CHANNELS; i++ {
						guessed := VarCell{gent, CHANNEL, i, i, i, guessed_global_id, guessed_global_id, guessed_global_id}
						present := false
						for _, assigned := range busylist[rproc] {
							if assigned.Procobjtype == guessed.Procobjtype && assigned.Id == guessed.Id {
								present = true
								break
							}
						}
						if !present {
							resp <- VarAns{ANS_OK, guessed}
							useditem <- UsageNotify{TR_PROC, rproc, C_SHAREDOBJECT, "channel:", I_NIL}
							busylist[rproc] = append(busylist[rproc], guessed)
							busychan[guessed_global_id].Connected = append(busychan[guessed_global_id].Connected, rproc)
							useditem <- UsageNotify{TR_CHAN, guessed_global_id, C_CONNECTED, S_NIL, rproc}
							created = true
							break
						}
					}
					if !created {
						panic("channel attach failed")
					}
				} else if gent, _ := Type_from_string("chan bool"); Same_Type(rcell.Vtype, gent) {
					if _, ok := busylist[rproc]; !ok {
						vcells := make([]VarCell, 0)
						busylist[rproc] = vcells
					}

					guessed_global_id := rcell.Global_id
					created := false
					for i := 0; i < MAX_CHANNELS; i++ {
						guessed := VarCell{gent, CHANNEL, i, i, i, guessed_global_id, guessed_global_id, guessed_global_id}
						present := false
						for _, assigned := range busylist[rproc] {
							if assigned.Procobjtype == guessed.Procobjtype && assigned.Id == guessed.Id {
								present = true
								break
							}
						}
						if !present {
							resp <- VarAns{ANS_OK, guessed}
							useditem <- UsageNotify{TR_PROC, rproc, C_SHAREDOBJECT, "channel:", I_NIL}
							busylist[rproc] = append(busylist[rproc], guessed)
							busychan[guessed_global_id].Connected = append(busychan[guessed_global_id].Connected, rproc)
							useditem <- UsageNotify{TR_CHAN, guessed_global_id, C_CONNECTED, S_NIL, rproc}
							created = true
							break
						}
					}
					if !created {
						panic("channel attach failed")
					}
				}
			}
		case REQ_EXIT:
			//fmt.Println(busychan)
			break VA
		}
	}
	assignerdone <- true
}
