package net

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/ingotmc/ingot/net/protocol/packet"
	"net"
)

const (
	defaultAddr = "localhost"
	defaultPort = 25565
)

type PacketHandler interface {
	HandlePacket(c Conn, p packet.Raw)
}

type Server struct {
	addr           string
	port           int
	l              net.Listener
	clients        map[ConnID]*Conn
	inboundPackets chan packet.Raw
}

// NewServer returns a new Server with the given Options
func NewServer(opts ...ServerOption) (*Server, error) {
	s := &Server{
		addr: defaultAddr,
		port: defaultPort,
	}
	for _, o := range opts {
		o(s)
	}
	return s, nil
}

func (s *Server) Listen() error {
	listenOn := fmt.Sprintf("%s:%d", s.addr, s.port)
	l, err := net.Listen("tcp", listenOn)
	if err != nil {
		return err
	}
	s.l = l
	for c, err := s.l.Accept(); err == nil; {
		go s.handleConn(c)
	}
	return err
}

func (s *Server) handleConn(c net.Conn) {
	id := ConnID(uuid.New().String())
	conn, err := NewConnection(id, c)
	if err != nil {
		// TODO: log
		return
	}
	s.clients[conn.ID] = conn
	conn.Start()
}
