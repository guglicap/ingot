// package handshake implements packet data for the connection state Handshaking
package handshake

import (
	"github.com/ingotmc/ingot/net/protocol/internal/read"
	"io"
)

const SetProtocolID = 0x00

type SetProtocol struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	NextState       int32
}

func (s SetProtocol) Encode(w io.Writer) error {
	return nil
}

func (s *SetProtocol) Decode(r io.Reader) (err error) {
	s.ProtocolVersion, _, err = read.VarInt(r)
	if err != nil {
		return
	}
	s.ServerAddress, err = read.String(r)
	if err != nil {
		return
	}
	s.ServerPort, err = read.UShort(r)
	if err != nil {
		return
	}
	s.NextState, _, err = read.VarInt(r)
	return
}
