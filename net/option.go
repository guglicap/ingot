package net

type ServerOption func(s *Server)

func WithAddr(addr string) ServerOption {
	return func(s *Server) {
		s.addr = addr
	}
}

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

func WithPacketHandler(ph PacketHandler) ServerOption {
	return func(s *Server) {
		s.packetHandler = ph
	}
}
