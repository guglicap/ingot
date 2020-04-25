package net

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/ingotmc/ingot/net/protocol"
	"github.com/ingotmc/ingot/net/protocol/packet"
	"github.com/ingotmc/ingot/net/protocol/packet/login"
	"net"
)

const (
	defaultAddr = "localhost"
	defaultPort = 25565
)

// PacketHandler represents the ability to handle a Packet p coming from a Conn c
type PacketHandler interface {
	HandlePacket(c *Conn, p packet.Packet)
}

// Server listens on a given address and port, defaulting to "localhost:25565"
// Upon receiving packets from the clients it calls a given PacketHandler.
//
// By default it implements PacketHandler itself and generates a PacketReceived event which
// it then forwards to the state.
type Server struct {
	addr          string
	port          int
	l             net.Listener
	shutdown      bool
	clients       map[ConnID]*Conn
	packetHandler PacketHandler
}

// NewServer returns a new Server with the given Options
func NewServer(opts ...ServerOption) (*Server, error) {
	s := &Server{
		addr:    defaultAddr,
		port:    defaultPort,
		clients: make(map[ConnID]*Conn),
	}
	s.packetHandler = s
	for _, o := range opts {
		o(s)
	}
	return s, nil
}

// Listen starts listening on the given address and port.
// It's blocking and will return nil when Shutdown is called.
func (s *Server) Listen() error {
	listenOn := fmt.Sprintf("%s:%d", s.addr, s.port)
	l, err := net.Listen("tcp", listenOn)
	if err != nil {
		return err
	}
	s.l = l
	for {
		var c net.Conn
		c, err = s.l.Accept()
		if err != nil {
			break
		}
		go s.handleConn(c)
	}
	if s.shutdown {
		err = nil
	}
	return err
}

func (s *Server) handleConn(c net.Conn) {
	id := ConnID(uuid.New().String())
	conn := NewConnection(id, c)
	s.clients[conn.ID] = conn
	conn.Start()
	go func(conn *Conn) {
		for p := range conn.sBound {
			s.packetHandler.HandlePacket(conn, p)
		}
	}(conn)
}

// Shutdown will stop the server.
func (s *Server) Shutdown() {
	s.shutdown = true
	for k, c := range s.clients {
		err := c.Close()
		close(c.cBound)
		if err != nil {
			// TODO: logging
		}
		delete(s.clients, k)
	}
	err := s.l.Close()
	if err != nil {
		// TODO: logging
	}
}

func (s *Server) HandlePacket(c *Conn, p packet.Packet) {
	if p.ID == login.LoginStartID && p.State == protocol.Login {
		l := new(login.LoginSuccess)
		l.Name = "guglicap"
		l.UUID = "11111111-2222-3333-4444-555555555555"
		c.cBound <- packet.Packet{State: protocol.Login, ID: 0x02, Data: l}
	}
}
