package etherbond

import (
	"github.com/mdlayher/raw"
)

func send_io(c *Config, p *raw.Conn, cl *Cluster, peerid uint32, run *Cluster_runinfo, ops *Peer_op, resid uint32, val interface{}) {
	ifi := c.ifi

	frame_send_chan := c.frame_send_chan
	transaction_chan := c.transaction_chan
	tag_chan := c.tag_chan

	for _, peer := range run.Peers {
		for inp, _ := range peer.Inputs {
			if inp == resid {
				bufstart, _ := packCommon(c.Rsize, IO_TR_FR)
				buf := bufstart
				dest := peer.HAddr
				buf = pcat(dest, buf)
				buf = pcat(ifi.HardwareAddr, buf)
				buf = pint16(ETHERTYPE, buf)
				buf = pint8(IO_TR_CM, buf)
				tag := <-tag_chan
				buf = pint32(tag, buf)
				buf = pint32(cl.ClusterId, buf)
				buf = pint32(peerid, buf)
				buf = pint32(resid, buf)
				switch v := val.(type) {
				case *Ouint8:
					buf = pint8(v.Value, buf)
				}

				transaction_chan <- Transaction{TRANSNEW, tag, string(bufstart)}
				frame_send_chan <- string(bufstart)
			}
		}
	}

}
