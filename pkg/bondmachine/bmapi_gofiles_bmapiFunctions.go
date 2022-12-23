package bondmachine

const (
	bmapiFunctions = `package bmapiusbuart

import (
	"errors"
	"time"
)

func (ba *BMAPI) BMi2r(register uint8) (uint8, error) {
	if ba.stateGet() == stateCONNECT {
		switch register {
		case uint8(0):
			// Async: read the value from the stored chache
			ba.i0Mutex.RLock()
			result := ba.i0
			ba.i0Mutex.RUnlock()

			ba.i0validMutex.RLock()
			validCheck := ba.i0valid
			ba.i0validMutex.RUnlock()

			// Send the data recv
			if validCheck {
				timeout := time.NewTimer(timeoutDelay * time.Second)
				select {
				case ba.i0recvSend <- true:
				case <-timeout.C:
				}
				timeout.Stop()
			}

			return result, nil
		default:
			return 0, errors.New("unknown register")
		}
	} else {
		return 0, errors.New("unconnected")
	}
}

func (ba *BMAPI) BMi2rw(register uint8) (uint8, error) {
	if ba.stateGet() == stateCONNECT {
		switch register {
		case uint8(0):

			ba.i0validMutex.RLock()
			validCheck := ba.i0valid
			ba.i0validMutex.RUnlock()

			if !validCheck {
				ba.i0validWait <- struct{}{}
			}

			// Sync: read the value from the stored chache
			ba.i0Mutex.RLock()
			result := ba.i0
			ba.i0Mutex.RUnlock()

			// Send the data recv
			timeout := time.NewTimer(timeoutDelay * time.Second)
			select {
			case ba.i0recvSend <- true:
			case <-timeout.C:
			}
			timeout.Stop()

			ba.i0validMutex.RLock()
			validCheck = ba.i0valid
			ba.i0validMutex.RUnlock()

			if validCheck {
				ba.i0validWait <- struct{}{}
			}

			return result, nil
		default:
			return 0, errors.New("unknown register")
		}
	} else {
		return 0, errors.New("unconnected")
	}
}
func (ba *BMAPI) BMr2o(register uint8, value uint8) error {
	if ba.stateGet() == stateCONNECT {
		switch register {
		case uint8(0):
			timeout := time.NewTimer(timeoutDelay * time.Second)
			select {
			case ba.o0Send <- value:
			case <-timeout.C:
			}
			timeout.Stop()

			timeout = time.NewTimer(timeoutDelay * time.Second)
			select {
			case ba.o0validSend <- true:
			case <-timeout.C:
			}
			timeout.Stop()

		default:
			return errors.New("unknown register")
		}
	} else {
		return errors.New("unconnected")
	}
	return nil
}
func (ba *BMAPI) BMr2ow(register uint8, value uint8) error {
	if ba.stateGet() == stateCONNECT {
		switch register {
		case uint8(0):
			timeout := time.NewTimer(timeoutDelay * time.Second)
			select {
			case ba.o0Send <- value:
			case <-timeout.C:
			}
			timeout.Stop()

			timeout = time.NewTimer(timeoutDelay * time.Second)
			select {
			case ba.o0validSend <- true:
			case <-timeout.C:
			}
			timeout.Stop()

			ba.o0recvMutex.RLock()
			validCheck := ba.o0recv
			ba.o0recvMutex.RUnlock()

			if !validCheck {
				ba.o0recvWait <- struct{}{}
			}
			return nil
		default:
			return errors.New("unknown register")
		}
	} else {
		return errors.New("unconnected")
	}
}
func (*BMAPI) BMr2owa(register uint8, value uint8) error {
	return nil
}
`
)
