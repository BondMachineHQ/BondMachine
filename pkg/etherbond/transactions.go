package etherbond

import (
	"github.com/mdlayher/raw"
)

func transaction_manager(c *Config, p *raw.Conn, cl *Cluster, peerid uint32, run *Cluster_runinfo, ops *Peer_op) {
	transaction_chan := c.transaction_chan
	tag_chan := c.tag_chan

	senttags := make(map[uint32]string)

	nexttag := uint32(0)

	for {
		select {
		case tag_chan <- nexttag:
			nexttag += 1
		case sent := <-transaction_chan:
			switch sent.Ttype {
			case TRANSNEW:
				senttags[sent.Tag] = sent.Data
			case TRANSDONE:
				if _, ok := senttags[sent.Tag]; ok {
					delete(senttags, sent.Tag)
				}
			}
		}
	}
}
