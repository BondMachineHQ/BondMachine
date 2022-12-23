package bondmachine

const (
	bmapiEncoder = `package bmapiusbuart

import (
	"context"
	"log"
	"time"
)

func (ba *BMAPI) encoder(ctx context.Context) {
	if ba.debug {
		log.Println("encoder: starting")
	}
	for {
		select {
		case <-ctx.Done():
			if ba.debug {
				log.Println("encoder: exiting")
			}
			return
		default:
			if ctx.Err() != nil {
				if ba.debug {
					log.Println("encoder: exiting")
				}
				return
			}
			// All the actions depends on the internal state of the library (ba.state)
			switch ba.stateGet() {
			case stateWAIT:
				if ba.debug {
					log.Println("encoder: BMAPI in state WAIT")
				}
				select {
				case <-ctx.Done():
					if ba.debug {
						log.Println("encoder: exiting")
					}
					return
				case ba.sendChan <- cmdHANDSH | ba.hs:
				}
				if ba.debug {
					log.Printf("encoder: HandShake sent %08b - new state HSSENT", cmdHANDSH|ba.hs)
				}
				ba.stateSet(stateHSSENT)
				ba.notresetChan <- struct{}{}
			case stateMASKRECV:
				if ba.debug {
					log.Println("encoder: BMAPI in state MASKRECV")
				}
				select {
				case <-ctx.Done():
					if ba.debug {
						log.Println("encoder: exiting")
					}
					return
				case ba.sendChan <- ba.hsmask & (cmdHANDSH | ba.hs):
				}
				if ba.debug {
					log.Printf("encoder: Ack sent %08b - new state ACK", ba.hsmask&(cmdHANDSH|ba.hs))
				}
				ba.stateSet(stateACK)
				ba.notresetChan <- struct{}{}
			case stateCONNECT:
				if ba.debug {
					log.Println("encoder: BMAPI in state CONNECT")
				}
				timeout := time.NewTimer(timeoutDelay * time.Second)
				select {
				case <-ctx.Done():
					if ba.debug {
						log.Println("encoder: exiting")
					}
					return
				case <-ba.keep:
					cmd := cmdKEEP
					ba.sendChan <- cmd
					if ba.debug {
						log.Printf("encoder: KEEP sent\n")
					}
				case o0Data := <-ba.o0Send:
					if o0Data != ba.o0 {
						if ba.monitor {
							ba.monitorChan <- struct{}{}
						}
						ba.o0 = o0Data
						reg := uint8(0)
						cmd := cmdNEWVAL | reg
						ba.sendChan <- cmd
						// BMAPI Customize
						mask8 := uint8(255)
						for i := 0; i < regsizeB; i++ {
							value := uint8((o0Data << (8 * i)) & mask8)
							ba.sendChan <- value
							if ba.debug {
								log.Printf("encoder: NEWVAL sent - Reg:%d - Value:%d\n", reg, value)
							}
						}
					}
				case o0Valid := <-ba.o0validSend:
					if o0Valid != ba.o0valid {
						if ba.monitor {
							ba.monitorChan <- struct{}{}
						}
						ba.o0valid = o0Valid
						if ba.debug {
							log.Println("encoder: o0valid changed to", o0Valid)
						}
						reg := uint8(0)
						if o0Valid {
							if ba.debug {
								log.Printf("encoder: DVALIDH sent - Reg: %d\n", reg)
							}
							cmd := cmdDVALIDH | reg
							ba.sendChan <- cmd
						} else {
							if ba.debug {
								log.Printf("encoder: DVALIDL sent - Reg: %d\n", reg)
							}
							cmd := cmdDVALIDL | reg
							ba.sendChan <- cmd
						}
					}
				case i0Recv := <-ba.i0recvSend:
					if i0Recv != ba.i0recv {
						if ba.monitor {
							ba.monitorChan <- struct{}{}
						}
						ba.i0recv = i0Recv
						if ba.debug {
							log.Println("encoder: i0recv changed to", i0Recv)
						}
						reg := uint8(0)
						if i0Recv {
							if ba.debug {
								log.Printf("encoder: RECVH sent - Reg: %d\n", reg)
							}
							cmd := cmdDRECVH | reg
							ba.sendChan <- cmd
						} else {
							if ba.debug {
								log.Printf("encoder: RECVL sent - Reg: %d\n", reg)
							}
							cmd := cmdDRECVL | reg
							ba.sendChan <- cmd
						}
					}
				case <-timeout.C:
				}
				timeout.Stop()
			}
		}
	}
}

func (ba *BMAPI) sendAll() {

	reg := uint8(0)
	cmd := cmdNEWVAL | reg
	ba.sendChan <- cmd
	// FakeAccel Customize
	mask8 := uint8(255)
	for i := 0; i < regsizeB; i++ {
		value := uint8((ba.o0 << (8 * i)) & mask8)
		ba.sendChan <- value
	}

	reg = uint8(0)
	if ba.o0valid {
		cmd := cmdDVALIDH | reg
		ba.sendChan <- cmd
	} else {
		cmd := cmdDVALIDL | reg
		ba.sendChan <- cmd
	}

	reg = uint8(0)
	if ba.i0recv {
		cmd := cmdDRECVH | reg
		ba.sendChan <- cmd
	} else {
		cmd := cmdDRECVL | reg
		ba.sendChan <- cmd
	}

}
`
)
