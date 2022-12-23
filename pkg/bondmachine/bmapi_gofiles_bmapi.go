package bondmachine

const (
	bmapi = `package {{ $.PackageName }}

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	// BMAPI Customize
	regsizeB = 1
)

const (
	stateWAIT = uint8(0) + iota
	stateHSSENT
	stateMASKRECV
	stateACK
	stateCONNECT
)

const (
	resetTime    = 5  // seconds
	timeoutDelay = 1  // seconds
	waitDelay    = 10 // milliseconds
)

const (
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	normal = "\033[0m"
)

type BMAPI struct {
	// BMAPI Customize
	state        uint8
	stateMutex   sync.RWMutex
	hsmask       uint8
	hs           uint8
	keep         chan struct{}
	debug        bool
	monitor      bool
	monitorChan  chan struct{}
	i0           uint8
	i0Mutex      sync.RWMutex
	i0valid      bool
	i0validMutex sync.RWMutex
	i0validWait  chan struct{}
	i0recv       bool
	i0recvSend   chan bool
	o0           uint8
	o0Send       chan uint8
	o0valid      bool
	o0validSend  chan bool
	o0recv       bool
	o0recvMutex  sync.RWMutex
	o0recvWait   chan struct{}
	recvChan     <-chan uint8
	sendChan     chan<- uint8
	endedChan    <-chan struct{}
	notresetChan chan struct{}
	cancel       context.CancelFunc
}

func (ba *BMAPI) stateGet() uint8 {
	ba.stateMutex.RLock()
	defer ba.stateMutex.RUnlock()
	return ba.state
}

func (ba *BMAPI) stateCheck(s uint8) bool {
	ba.stateMutex.RLock()
	defer ba.stateMutex.RUnlock()
	return ba.state == s
}

func (ba *BMAPI) stateSet(s uint8) {
	ba.stateMutex.Lock()
	defer ba.stateMutex.Unlock()
	ba.state = s
}

func (ba *BMAPI) sendkeep(ctx context.Context) {
	if ba.debug {
		log.Println("sendkeep: starting")
	}
	for {
		select {
		case <-ctx.Done():
			if ba.debug {
				log.Println("sendkeep: exiting")
			}
			return
		case <-time.After(200 * time.Millisecond):
			if ba.stateCheck(stateCONNECT) {
				ba.keep <- struct{}{}
			}
		}
	}
}

func (ba *BMAPI) reset(ctx context.Context) {
	if ba.debug {
		log.Println("reset: starting")
	}
	for {
		t := time.NewTimer(resetTime * time.Second)
		select {
		case <-ctx.Done():
			if ba.debug {
				log.Println("reset: exiting")
			}
			return
		case <-ba.notresetChan:
		case <-t.C:
			ba.stateSet(stateWAIT)
			if ba.debug {
				log.Println("reset: state resetted to WAIT")
			}
		}
		t.Stop()
	}
}
func bmapiInit(device string, tr func(context.Context, string, bool) (chan<- uint8, <-chan uint8, <-chan struct{})) (*BMAPI, error) {
	result := new(BMAPI)

	ctx, cancel := context.WithCancel(context.Background())

	result.debug = false
	result.monitor = false
	result.monitorChan = make(chan struct{})

	result.cancel = cancel
	result.sendChan, result.recvChan, result.endedChan = tr(ctx, device, result.debug)

	result.notresetChan = make(chan struct{})

	result.hs = uint8(28)

	result.state = stateWAIT

	result.keep = make(chan struct{})

	// TODO SendChan

	result.o0Send = make(chan uint8)
	result.o0validSend = make(chan bool)
	result.i0recvSend = make(chan bool)
	result.i0validWait = make(chan struct{})
	result.o0recvWait = make(chan struct{})

	// Initialize bm struct
	result.i0 = 0
	result.i0valid = false
	result.i0recv = false
	result.o0 = 0
	result.o0valid = false
	result.o0recv = false

	go result.decoder(ctx)
	go result.encoder(ctx)
	go result.sendkeep(ctx)
	go result.reset(ctx)
	if result.monitor {
		go result.monitorLoop(ctx)
	}

	return result, nil
}

func AcceleratorInit(device string, tr func(context.Context, string, bool) (chan<- uint8, <-chan uint8, <-chan struct{})) (*BMAPI, error) {
	return bmapiInit(device, tr)
}

func (ba *BMAPI) AcceleratorStop() {
	ba.cancel()
	time.Sleep(time.Second)
	<-ba.endedChan
	time.Sleep(time.Second)
}

func (ba *BMAPI) WaitConnection() {
	for {
		if ba.stateCheck(stateCONNECT) {
			break
		}
		time.Sleep(waitDelay * time.Millisecond)
	}
}

func (ba *BMAPI) monitorLoop(ctx context.Context) {

	// Store the ba values
	oldo0 := ba.o0
	oldo0valid := ba.o0valid
	oldi0 := ba.i0
	oldi0valid := ba.i0valid
	oldi0recv := ba.i0recv
	oldo0recv := ba.o0recv

	for {
		select {
		case <-ba.monitorChan:
			result := ""
			if oldi0 != ba.i0 {
				result += red
				oldi0 = ba.i0
			}
			result += "i0: " + fmt.Sprintf("%08d", ba.i0)
			result += normal

			if oldi0valid != ba.i0valid {
				result += red
				oldi0valid = ba.i0valid
			}
			if ba.i0valid {
				result += " i0valid: high "
			} else {
				result += " i0valid: low  "
			}
			result += normal

			if oldi0recv != ba.i0recv {
				result += red
				oldi0recv = ba.i0recv
			}
			if ba.i0recv {
				result += " i0recv: high "
			} else {
				result += " i0recv: low  "
			}
			result += normal

			if oldo0 != ba.o0 {
				result += red
				oldo0 = ba.o0
			}
			result += " o0: " + fmt.Sprintf("%08d", ba.o0)

			if oldo0valid != ba.o0valid {
				result += red
				oldo0valid = ba.o0valid
			}
			if ba.o0valid {
				result += " o0valid: high "
			} else {
				result += " o0valid: low  "
			}
			result += normal

			if oldo0recv != ba.o0recv {
				result += red
				oldo0recv = ba.o0recv
			}
			if ba.o0recv {
				result += " o0recv: high "
			} else {
				result += " o0recv: low  "
			}
			result += normal

			log.Println(result)
		case <-ctx.Done():
			return
		}
	}
}
`
)
