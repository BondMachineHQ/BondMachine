package etherbond

import (
	"github.com/mdlayher/raw"
)

// 8 bit

type Iuint8 struct {
	Value uint8
	Id    uint32
}

func (i *Iuint8) ReadValue(c *Config, p *raw.Conn, cl *Cluster, peer_id uint32, run *Cluster_runinfo, ops *Peer_op) interface{} {
	return i.Value
}

func (i *Iuint8) setValue(c *Config, p *raw.Conn, cl *Cluster, peer_id uint32, run *Cluster_runinfo, ops *Peer_op, val interface{}) {
	i.Value = val.(uint8)
}

func (i *Iuint8) SetId(id uint32) {
	i.Id = id
}

func (i *Iuint8) GetId() uint32 {
	return (i.Id)
}

type Ouint8 struct {
	Value uint8
	Id    uint32
}

func (o *Ouint8) SetId(id uint32) {
	o.Id = id
}

func (o *Ouint8) GetId() uint32 {
	return (o.Id)
}

func (o *Ouint8) WriteValue(c *Config, p *raw.Conn, cl *Cluster, peer_id uint32, run *Cluster_runinfo, ops *Peer_op, val interface{}) {
	o.Value = val.(uint8)
	send_io(c, p, cl, peer_id, run, ops, o.Id, o)
}

func (o *Ouint8) getValue(c *Config, p *raw.Conn, cl *Cluster, peer_id uint32, run *Cluster_runinfo, ops *Peer_op) interface{} {
	return o.Value
}
