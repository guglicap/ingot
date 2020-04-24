package net

import (
	"github.com/ingotmc/ingot/net/protocol"
	"github.com/ingotmc/ingot/net/protocol/packet"
	"net"
)

type ConnID string

type Conn struct {
	ID ConnID
	net.Conn
	state    protocol.State
	inbound  chan packet.Packet
	outbound chan packet.Packet
}

func NewConnection(ID ConnID, conn net.Conn) (*Conn, error) {
	c := &Conn{
		ID:   ID,
		Conn: conn,
	}
	return c, nil
}

func (c Conn) readPacket() (p packet.Packet, err error) {
	rp, err := packet.ReadRaw(c)
	if err != nil {
		return
	}
	p.ID = rp.ID
	p.State = c.state
	data := packet.DataByIDAndState(p.ID, c.state)
	err = data.Decode(c)
	p.Data = data
	return
}

func (c Conn) receive() {
	var err error
	for err == nil {
		var p packet.Packet
		p, err = c.readPacket()
		if err != nil {
			// TODO: logging
			continue
		}
		c.inbound <- p
	}
}

func (c Conn) send() {

}

func (c Conn) Start() {
	go c.receive()
}
