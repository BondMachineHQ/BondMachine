package etherbond

import (
	"github.com/mdlayher/raw"
)

func send_ack(c *Config, p *raw.Conn, cl *Cluster, peerid uint32, run *Cluster_runinfo, ops *Peer_op, dest []byte, tag uint32) {
	ifi := c.ifi

	frame_send_chan := c.frame_send_chan

	bufstart, _ := packCommon(c.Rsize, ACK_FR)
	buf := bufstart
	buf = pcat(dest, buf)
	buf = pcat(ifi.HardwareAddr, buf)
	buf = pint16(ETHERTYPE, buf)
	buf = pint8(ACK_CM, buf)
	buf = pint32(tag, buf)

	frame_send_chan <- string(bufstart)

}
