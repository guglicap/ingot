package net

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/ingotmc/ingot/net/protocol"
	"github.com/ingotmc/ingot/net/protocol/packet"
	"github.com/ingotmc/ingot/net/protocol/packet/login"
	"go.uber.org/zap"
	"net"
)

const (
	defaultAddr = "localhost"
	defaultPort = 25565
)

// PacketHandler represents the ability to handle a Packet p coming from a Conn c
type PacketHandler interface {
	HandlePacket(c *Conn, p packet.Serverbound)
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
	log           *zap.SugaredLogger
	shutdown      bool
	clients       map[ConnID]*Conn
	packetHandler PacketHandler
}

// NewServer returns a new Server with the given Options
func NewServer(opts ...ServerOption) (*Server, error) {
	s := &Server{
		addr:    defaultAddr,
		port:    defaultPort,
		log:     zap.S().Named("net/server"),
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
		var c net.Conn // needed not to shadow err
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

// Shutdown will stop the server.
func (s *Server) Shutdown() {
	s.shutdown = true
	s.log.Info("shutting down...")
	for _, c := range s.clients {
		c.Close()
	}
	err := s.l.Close()
	if err != nil {
		s.log.Error("error closing listener", err)
	}
}

func (s *Server) HandlePacket(c *Conn, p packet.Serverbound) {
	if p.ID == login.LoginStartID && p.State == protocol.Login {
		l := new(login.LoginSuccess)
		l.Name = "guglicap"
		l.UUID = "11111111-2222-3333-4444-555555555555"
		c.cBound <- packet.Clientbound{ID: 0x02, Data: l}
	}
}

func (s *Server) handleConn(c net.Conn) {
	id := ConnID(uuid.New())
	conn := NewConnection(id, c)
	s.log.Debugw("new connection", "id", id.String())
	s.clients[conn.ID] = conn
	go conn.Start()
	for p := range conn.sBound {
		s.packetHandler.HandlePacket(conn, p)
	}
	s.removeConnection(conn)
}

// removeConnection removes a connection c from the server.
// calling removeConnection with an already closed connection is a noop.
func (s *Server) removeConnection(c *Conn) {
	if _, ok := s.clients[c.ID]; !ok {
		s.log.Debug("called removeConnection on an already closed connection", c.ID.String())
		return
	}
	close(c.cBound)
	delete(s.clients, c.ID)
	s.log.Debugw("client disconnected", "id", c.ID.String())
}
