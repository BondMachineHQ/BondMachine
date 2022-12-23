package etherbond

import (
	//"fmt"
	"github.com/mdlayher/raw"
	"time"
)

func advertiser(c *Config, p *raw.Conn, cl *Cluster, peerid uint32, run *Cluster_runinfo, ops *Peer_op) {

	ifi := c.ifi
	frame_send_chan := c.frame_send_chan

	for {
		{
			bufstart, _ := packCommon(c.Rsize, ADV_CLU_FR)
			buf := bufstart
			buf = pcat([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, buf)
			buf = pcat(ifi.HardwareAddr, buf)
			buf = pint16(ETHERTYPE, buf)
			buf = pint8(ADV_CLU_CM, buf)
			buf = pint32(cl.ClusterId, buf)
			buf = pint32(peerid, buf)

			frame_send_chan <- string(bufstart)
		}

		for _, peer := range cl.Peers {
			if peer.PeerId == peerid {
				for _, inid := range peer.Channels {
					bufstart, _ := packCommon(c.Rsize, ADV_CH_FR)
					buf := bufstart
					buf = pcat([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, buf)
					buf = pcat(ifi.HardwareAddr, buf)
					buf = pint16(ETHERTYPE, buf)
					buf = pint8(ADV_CH_CM, buf)
					buf = pint32(cl.ClusterId, buf)
					buf = pint32(peerid, buf)
					buf = pint32(inid, buf)

					frame_send_chan <- string(bufstart)
				}

				for _, inid := range peer.Inputs {
					bufstart, _ := packCommon(c.Rsize, ADV_IN_FR)
					buf := bufstart
					buf = pcat([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, buf)
					buf = pcat(ifi.HardwareAddr, buf)
					buf = pint16(ETHERTYPE, buf)
					buf = pint8(ADV_IN_CM, buf)
					buf = pint32(cl.ClusterId, buf)
					buf = pint32(peerid, buf)
					buf = pint32(inid, buf)

					frame_send_chan <- string(bufstart)
				}

				for _, inid := range peer.Outputs {
					bufstart, _ := packCommon(c.Rsize, ADV_OUT_FR)
					buf := bufstart
					buf = pcat([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, buf)
					buf = pcat(ifi.HardwareAddr, buf)
					buf = pint16(ETHERTYPE, buf)
					buf = pint8(ADV_OUT_CM, buf)
					buf = pint32(cl.ClusterId, buf)
					buf = pint32(peerid, buf)
					buf = pint32(inid, buf)

					frame_send_chan <- string(bufstart)
				}
			}
		}
		time.Sleep(1 * time.Second)
	}

}
