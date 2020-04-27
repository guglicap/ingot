package net

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/ingotmc/ingot/net/protocol"
	"github.com/ingotmc/ingot/net/protocol/packet"
	"github.com/ingotmc/ingot/net/protocol/packet/handshake"
	"net"
)

// ConnID represents a connection ID, assigned by the server upon accepting a new connection.
// Unique for every client.
type ConnID uuid.UUID

// Conn represents an I/O stream with a client, it handles packet encoding/decoding.
type Conn struct {
	ID ConnID
	net.Conn
	state  protocol.State
	stop   chan struct{}
	sBound chan packet.Packet
	cBound chan packet.Packet
}

// NewConnection returns a connection with the given ID and using conn as the underlying transport
func NewConnection(ID ConnID, conn net.Conn) *Conn {
	c := &Conn{
		ID:     ID,
		Conn:   conn,
		state:  protocol.Handshaking,
		stop:   make(chan struct{}),
		sBound: make(chan packet.Packet),
		cBound: make(chan packet.Packet),
	}
	return c
}

// blocking, reads a raw packet from the wire and decodes it.
func (c Conn) readPacket() (packet.Packet, error) {
	raw, err := packet.ReadRaw(c)
	if err != nil {
		return packet.Packet{}, err
	}
	var p packet.Packet
	p.ID = raw.ID
	p.State = c.state
	data, err := packet.DataByIDAndState(p.ID, c.state)
	if err != nil {
		return packet.Packet{}, err
	}
	err = data.Decode(bytes.NewReader(raw.Data))
	p.Data = data
	return p, nil
}

func (c Conn) writePacket(p packet.Packet) error {
	rp := packet.Raw{}
	rp.ID = p.ID
	buf := bytes.NewBuffer([]byte{})
	err := p.Data.Encode(buf)
	if err != nil {
		return err
	}
	rp.Data = buf.Bytes()
	return packet.WriteRaw(rp, c)
}

// updateState takes care of maintaining the correct state in order to decode
// upcoming packets correctly.
func (c *Conn) updateState(p packet.Packet) {
	if p.State == protocol.Handshaking && p.ID == handshake.SetProtocolID {
		c.state = protocol.Login
	}
}

// receive forwards decoded packets into the sBound channel.
func (c Conn) receive() {
loop:
	for {
		select {
		case <-c.stop:
			break loop
		default:
			p, err := c.readPacket()
			if err != nil {
				// TODO: logging
				continue
			}
			c.updateState(p)
			c.sBound <- p
		}
	}
	close(c.sBound)
}

func (c Conn) send() {
	for p := range c.cBound {
		err := c.writePacket(p)
		if err != nil {
			// TODO: logging
			continue
		}
		c.updateState(p)
	}
}

// Start initiates I/O with client.
func (c Conn) Start() {
	go c.receive()
	go c.send()
}

func (c *Conn) Close() error {
	c.stop <- struct{}{}
	return c.Conn.Close()
}
