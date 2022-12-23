package bondmachine

const (
	bmapiDecoder = `package bmapiusbuart

import (
	"context"
	"log"
	"time"
)

func (ba *BMAPI) decoder(ctx context.Context) {
	if ba.debug {
		log.Printf("decoder: starting\n")
	}
	for {
		select {
		case <-ctx.Done():
			if ba.debug {
				log.Printf("decoder: exiting\n")
			}
			return
		case b := <-ba.recvChan:
			if ba.debug {
				log.Printf("decoder: received data: %08b\n", b)
			}
			// All the actions depends on the internal state of the library (ba.state)
			switch ba.stateGet() {
			case stateWAIT:
				if ba.debug {
					log.Printf("decoder: BMAPI in state WAIT, ignoring received data: %08b\n", b)
				}
			case stateHSSENT:
				if ba.debug {
					log.Printf("decoder: received hsdata %08b - new state MASKRECV", b)
				}
				ba.hsmask = b
				ba.stateSet(stateMASKRECV)
				ba.notresetChan <- struct{}{}
			case stateMASKRECV:
				if ba.debug {
					log.Printf("decoder: BMAPI in state MASKRECV, ignoring received data: %08b\n", b)
				}
			case stateACK:
				if b == ba.hsmask&(cmdHANDSH|ba.hs) {
					if ba.debug {
						log.Println("decoder: BMAPI in state ACK, valid handshake - new state CONNECT")
					}
					ba.stateSet(stateCONNECT)
					//ba.sendAll()
					ba.notresetChan <- struct{}{}
				} else {
					if ba.debug {
						log.Println("decoder: BMAPI in state ACK, not valid handshake - new state WAIT")
					}
					ba.stateSet(stateWAIT)
				}
			case stateCONNECT:
				if ba.debug {
					log.Println("decoder: BMAPI in state CONNECT data will be decoded normally")
				}
				switch {
				case (b & cmdMASK) == cmdKEEP:
					ba.notresetChan <- struct{}{}
					if ba.debug {
						log.Println("decoder: KEEP received")
					}
				case (b & cmdMASK) == cmdNEWVAL:
					reg := (b & ^cmdMASK)
					if ba.debug {
						log.Printf("decoder: received command: NEWVAL - Reg:%d\n", reg)
					}
					// BMAPI Customize
					value := uint8(0)
					for i := 0; i < regsizeB; i++ {
						value = (value << 8) + <-ba.recvChan
					}
					switch reg {
					// BMAPI Customize
					case uint8(0):
						if ba.i0 != value {
							if ba.monitor {
								ba.monitorChan <- struct{}{}
							}
							ba.i0Mutex.Lock()
							ba.i0 = value
							ba.i0Mutex.Unlock()
						}
					}
				case (b & cmdMASK) == cmdDVALIDH:
					reg := (b & ^cmdMASK)
					if ba.debug {
						log.Printf("decoder: received command: DVALIDH - Reg:%d\n", reg)
					}
					// FakeAccel Customize
					switch reg {
					// FakeAccel Customize
					case uint8(0):
						if !ba.i0valid {
							if ba.monitor {
								ba.monitorChan <- struct{}{}
							}
							ba.i0validMutex.Lock()
							ba.i0valid = true
							ba.i0validMutex.Unlock()
						}

						select {
						case <-ba.i0validWait:
						default:
						}
					}
				case (b & cmdMASK) == cmdDVALIDL:
					reg := (b & ^cmdMASK)
					if ba.debug {
						log.Printf("decoder: received command: DVALIDL - Reg:%d\n", reg)
					}
					// FakeAccel Customize
					switch reg {
					// FakeAccel Customize
					case uint8(0):
						if ba.i0valid {
							if ba.monitor {
								ba.monitorChan <- struct{}{}
							}
							ba.i0validMutex.Lock()
							ba.i0valid = false
							ba.i0validMutex.Unlock()
						}

						timeout := time.NewTimer(timeoutDelay * time.Second)
						select {
						case ba.i0recvSend <- false:
						case <-timeout.C:
						}
						timeout.Stop()

						select {
						case <-ba.i0validWait:
						default:
						}

					}
				case (b & cmdMASK) == cmdDRECVH:
					reg := (b & ^cmdMASK)
					if ba.debug {
						log.Printf("decoder: received command: DRECVH - Reg:%d\n", reg)
					}
					// FakeAccel Customize
					switch reg {
					// FakeAccel Customize
					case uint8(0):
						if !ba.o0recv {
							if ba.monitor {
								ba.monitorChan <- struct{}{}
							}
							ba.o0recvMutex.Lock()
							ba.o0recv = true
							ba.o0recvMutex.Unlock()
						}

						timeout := time.NewTimer(timeoutDelay * time.Second)
						select {
						case ba.o0validSend <- false:
						case <-timeout.C:
						}
						timeout.Stop()

						select {
						case <-ba.o0recvWait:
						default:
						}

					}
				case (b & cmdMASK) == cmdDRECVL:
					reg := (b & ^cmdMASK)
					if ba.debug {
						log.Printf("decoder: received command: DRECVL - Reg:%d\n", reg)
					}
					// FakeAccel Customize
					switch reg {
					// FakeAccel Customize
					case uint8(0):
						if ba.o0recv {
							if ba.monitor {
								ba.monitorChan <- struct{}{}
							}
							ba.o0recvMutex.Lock()
							ba.o0recv = false
							ba.o0recvMutex.Unlock()
						}
					}
				default:
				}
			}
		}
	}
}
`
)
