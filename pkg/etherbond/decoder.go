package etherbond

import (
	"fmt"
	"github.com/mdlayher/raw"
)

func frame_decode(c *Config, p *raw.Conn, cl *Cluster, peerid uint32, run *Cluster_runinfo, ops *Peer_op, n int, frame []byte) {

	dest, frame := gcat(6, frame)
	source, frame := gcat(6, frame)
	ethertype, frame := gint16(frame)
	command, frame := gint8(frame)

	if c.Debug {
		fmt.Printf("%04x %02x", ethertype, source[0])
		for i := 1; i < 6; i++ {
			fmt.Printf(":%02x", source[i])
		}
		fmt.Printf(" -> %02x", dest[0])
		for i := 1; i < 6; i++ {
			fmt.Printf(":%02x", dest[i])
		}
		fmt.Printf(" ! %02x\n", command)
	}

	switch command {
	case ADV_CLU_CM:
		clusterid, frame := gint32(frame)
		peerid, _ := gint32(frame)
		if clusterid == cl.ClusterId {
			if info, ok := run.Peers[peerid]; ok {
				info.HAddr = make([]byte, 6)
				copy(info.HAddr, source)
			} else {
				var newinfo Peer_runinfo
				newinfo.HAddr = source
				newinfo.Channels = make(map[uint32]bool)
				newinfo.Inputs = make(map[uint32]bool)
				newinfo.Outputs = make(map[uint32]bool)
				run.Peers[peerid] = newinfo
			}
		}
	case ADV_CH_CM:
		clusterid, frame := gint32(frame)
		peerid, frame := gint32(frame)
		resid, _ := gint32(frame)
		if clusterid == cl.ClusterId {
			if info, ok := run.Peers[peerid]; ok {
				list := info.Channels
				chk := false
				for ch, _ := range list {
					if ch == resid {
						chk = true
						break
					}
				}
				if !chk {
					info.Channels[resid] = true
				}
			}
		}
	case ADV_IN_CM:
		clusterid, frame := gint32(frame)
		peerid, frame := gint32(frame)
		resid, _ := gint32(frame)
		if clusterid == cl.ClusterId {
			if info, ok := run.Peers[peerid]; ok {
				list := info.Inputs
				chk := false
				for ch, _ := range list {
					if ch == resid {
						chk = true
						break
					}
				}
				if !chk {
					info.Inputs[resid] = true
				}
			}
		}
	case ADV_OUT_CM:
		clusterid, frame := gint32(frame)
		peerid, frame := gint32(frame)
		resid, _ := gint32(frame)
		if clusterid == cl.ClusterId {
			if info, ok := run.Peers[peerid]; ok {
				list := info.Outputs
				chk := false
				for ch, _ := range list {
					if ch == resid {
						chk = true
						break
					}
				}
				if !chk {
					info.Outputs[resid] = true
				}
			}
		}
	case IO_TR_CM:
		tag, frame := gint32(frame)
		clusterid, frame := gint32(frame)
		peerid, frame := gint32(frame)
		resid, frame := gint32(frame)
		data, _ := gint8(frame)
		if clusterid == cl.ClusterId {
			for _, inp := range ops.Inputs {
//fmt.Println("tag", tag,"clu",clusterid,"peer", peerid,"res", resid, "Data", data,inp.GetId())
				if inp.GetId() == resid {
					inp.setValue(c, p, cl, peerid, run, ops, data)
					send_ack(c, p, cl, peerid, run, ops, source, tag)
					break
				}
			}
		}
	case ACK_CM:
		tag, _ := gint32(frame)
		c.transaction_chan <- Transaction{TRANSDONE, tag, ""}
	}
}
