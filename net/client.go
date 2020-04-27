package net

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ingotmc/ingot/net/protocol"
	"github.com/ingotmc/ingot/net/protocol/packet"
	"github.com/ingotmc/ingot/net/protocol/packet/handshake"
	"go.uber.org/zap"
	"io"
	"net"
)

// ConnID represents a connection ID, assigned by the server upon accepting a new connection.
// Unique for every client.
type ConnID uuid.UUID

func (c ConnID) String() string {
	return uuid.UUID(c).String()
}

// Conn represents an I/O stream with a client, it handles packet encoding/decoding.
type Conn struct {
	net.Conn
	ID     ConnID
	state  protocol.State
	log    *zap.Logger
	stop   chan struct{}
	sBound chan packet.Serverbound
	cBound chan packet.Clientbound
}

// NewConnection returns a connection with the given ID and using conn as the underlying transport
func NewConnection(ID ConnID, conn net.Conn) *Conn {
	c := &Conn{
		ID:     ID,
		Conn:   conn,
		log:    zap.L().Named("net/conn").With(zap.String("id", ID.String())),
		state:  protocol.Handshaking,
		stop:   make(chan struct{}, 1),
		sBound: make(chan packet.Serverbound),
		cBound: make(chan packet.Clientbound),
	}
	return c
}

// Start initiates I/O with client.
func (c *Conn) Start() {
	go c.receive()
	go c.send()

	<-c.stop
	close(c.sBound)
	close(c.stop)
	err := c.Conn.Close()
	if err != nil {
		c.log.Error("error closing listener", zap.Error(err))
	}
}

func (c *Conn) Close() {
	c.stop <- struct{}{}
}

// blocking, reads a raw packet from the wire and decodes it.
func (c *Conn) readPacket() (packet.Serverbound, error) {
	raw, err := packet.ReadRaw(c)
	if err != nil {
		return packet.Serverbound{}, err
	}
	p, err := packet.ServerboundByIDAndState(raw.ID, c.state)
	if err != nil {
		return packet.Serverbound{}, err
	}
	err = p.Data.Decode(bytes.NewReader(raw.Data))
	fmt.Printf("decoded packet, id: %d state: %d\n", raw.ID, c.state)
	return p, err
}

func (c *Conn) writePacket(p packet.Clientbound) error {
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
func (c *Conn) updateState(info packet.Info) {
	if info.State == protocol.Handshaking && info.ID == handshake.SetProtocolID {
		c.state = protocol.Login
	}
}

// receive forwards decoded packets into the sBound channel.
func (c *Conn) receive() {
loop:
	for {
		select {
		case <-c.stop:
			break loop
		default:
			p, err := c.readPacket()
			if err != nil {
				if errors.Is(err, io.EOF) {
					c.Close()
				} else {
					c.log.Error("error receiving", zap.Error(err))
				}
				continue
			}
			p.State = c.state
			c.updateState(p.Info)
			c.sBound <- p
		}
	}
}

func (c *Conn) send() {
	for p := range c.cBound {
		err := c.writePacket(p)
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.Close()
			} else {
				c.log.Error("error sending", zap.Error(err))
			}
			continue
		}
		c.updateState(packet.Info{ID: p.ID, State: c.state})
	}
	c.log.Debug("stopping send")
}
