package etherbond

import (
	"github.com/mdlayher/raw"
)

func frame_sender(c *Config, p *raw.Conn, cl *Cluster, peerid uint32, run *Cluster_runinfo, ops *Peer_op) {
	ifi := c.ifi
	frame_send_chan := c.frame_send_chan
	for {
		frame := <-frame_send_chan
		p.WriteTo([]byte(frame), &raw.Addr{HardwareAddr: ifi.HardwareAddr})
	}
}

func frame_receiver(c *Config, p *raw.Conn, cl *Cluster, peerid uint32, run *Cluster_runinfo, ops *Peer_op) {

	buf := make([]byte, 1500)
	for {

		n, _, err := p.ReadFrom(buf)
		if err != nil {
			// TODO
		}

		if n > 0 {
			newbuf := make([]byte, len(buf))
			copy(newbuf, buf)
			go frame_decode(c, p, cl, peerid, run, ops, n, newbuf)
		}
	}

}
