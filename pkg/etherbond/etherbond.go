package etherbond

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mdlayher/raw"
	"net"
	"strconv"
)

const (
	ETHERTYPE = 0x8888
)

const (
	ADV_CLU_FR = 0 + iota
	ADV_CH_FR
	ADV_IN_FR
	ADV_OUT_FR
	IO_TR_FR
	ACK_FR
)

const (
	ADV_CLU_CM = 0x01
	ADV_CH_CM  = 0x02
	ADV_IN_CM  = 0x03
	ADV_OUT_CM = 0x04
	IO_TR_CM   = 0x05
	ACK_CM     = 0xff
)

const (
	TRANSNEW = uint8(0) + iota
	TRANSDONE
)

// Config struct

type Config struct {
	Rsize            uint8
	ifi              *net.Interface
	Debug            bool
	Done             chan bool
	kill_sender      chan bool
	kill_receiver    chan bool
	kill_advertiser  chan bool
	frame_send_chan  chan string
	tag_chan         chan uint32
	transaction_chan chan Transaction
}

// Peers description

type Peer struct {
	PeerId   uint32
	Channels []uint32
	Inputs   []uint32
	Outputs  []uint32
}

type Cluster struct {
	ClusterId uint32
	Peers     []Peer
}

//

type peerlist map[int]int

// Peers status

type Peer_runinfo struct {
	HAddr    []byte
	Channels map[uint32]bool
	Inputs   map[uint32]bool
	Outputs  map[uint32]bool
}

type Cluster_runinfo struct {
	ClusterId uint32
	Peers     map[uint32]Peer_runinfo
	Quorate   bool
	Degraded  bool
}

//

type Macs struct {
	Assoc map[string]string
}

//

type ChOp interface {
}

type InOp interface {
	ReadValue(*Config, *raw.Conn, *Cluster, uint32, *Cluster_runinfo, *Peer_op) interface{}
	SetId(id uint32)
	GetId() uint32
	setValue(*Config, *raw.Conn, *Cluster, uint32, *Cluster_runinfo, *Peer_op, interface{})
}

type OutOp interface {
	WriteValue(*Config, *raw.Conn, *Cluster, uint32, *Cluster_runinfo, *Peer_op, interface{})
	SetId(id uint32)
	GetId() uint32
	getValue(*Config, *raw.Conn, *Cluster, uint32, *Cluster_runinfo, *Peer_op) interface{}
}

type Peer_op struct {
	PeerId   uint32
	Channels []ChOp
	Inputs   []InOp
	Outputs  []OutOp
}

func (c *Cluster) String() string {
	result, _ := json.MarshalIndent(&c, "", "\t")
	return string(result)
}

func (c *Cluster_runinfo) String() string {
	result := "ClusterId: " + strconv.Itoa(int(c.ClusterId)) + "\n"
	if c.Quorate {
		result += "Quorum: yes\n"
	} else {
		result += "Quorum: no\n"
	}
	if c.Degraded {
		result += "Degraded: yes\n"
	} else {
		result += "Degraded: no\n"
	}
	for pid, peer := range c.Peers {
		result += "\tPeer: " + strconv.Itoa(int(pid)) + "\n"
		result += fmt.Sprintf("%02x %02x %02x %02x %02x %02x", peer.HAddr[0], peer.HAddr[1], peer.HAddr[2], peer.HAddr[3], peer.HAddr[4], peer.HAddr[5])
	}
	return string(result)
}

type Transaction struct {
	Ttype uint8
	Tag   uint32
	Data  string
}

func InitializePeer(c *Config, p *raw.Conn, cl *Cluster, peer_id uint32, run *Cluster_runinfo, ops *Peer_op) error {
	pck := false
	for _, peer := range cl.Peers {
		if peer.PeerId == peer_id {
			pck = true
			ops.PeerId = peer_id
			ops.Channels = make([]ChOp, len(peer.Channels))
			ops.Inputs = make([]InOp, len(peer.Inputs))
			ops.Outputs = make([]OutOp, len(peer.Outputs))

			for i, _ := range ops.Inputs {
				switch c.Rsize {
				case 8:
					newval8 := new(Iuint8)
					newval8.SetId(peer.Inputs[i])
					ops.Inputs[i] = newval8
				}
				// TODO others
			}
			for i, _ := range ops.Outputs {
				switch c.Rsize {
				case 8:
					newval8 := new(Ouint8)
					newval8.SetId(peer.Outputs[i])
					ops.Outputs[i] = newval8
				}
				// TODO others
			}
			// TODO Channels
		}
	}
	if !pck {
		return errors.New("Missing peer")
	}

	run.Peers = make(map[uint32]Peer_runinfo)

	return nil
}

func PartecipateCluster(c *Config, cl *Cluster, peer_id uint32, ifi *net.Interface) (*raw.Conn, *Cluster_runinfo, *Peer_op, chan bool, error) {

	p, err := raw.ListenPacket(ifi, 0x8888, nil)
	if err != nil {
		// TODO
	}

	run := new(Cluster_runinfo)
	ops := new(Peer_op)
	err = InitializePeer(c, p, cl, peer_id, run, ops)
	if err != nil {
		// TODO
	}

	done := make(chan bool)
	kill_sender := make(chan bool)
	kill_receiver := make(chan bool)
	kill_advertiser := make(chan bool)

	c.Done = done
	c.kill_sender = kill_sender
	c.kill_receiver = kill_receiver
	c.kill_advertiser = kill_advertiser

	tag_chan := make(chan uint32)
	transaction_chan := make(chan Transaction)

	c.tag_chan = tag_chan
	c.transaction_chan = transaction_chan

	c.ifi = ifi

	frame_send_chan := make(chan string)
	c.frame_send_chan = frame_send_chan

	// We need a consistency check of the peer_id

	//	if p, err := raw.ListenPacket(ifi, 0x0800); err != nil {
	//		return nil, done, err
	//	}

	go transaction_manager(c, p, cl, peer_id, run, ops)

	go frame_receiver(c, p, cl, peer_id, run, ops)
	go frame_sender(c, p, cl, peer_id, run, ops)

	go advertiser(c, p, cl, peer_id, run, ops)

	return p, run, ops, done, nil
}
