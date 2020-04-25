package packet

import (
	"github.com/ingotmc/ingot/net/protocol"
	"github.com/ingotmc/ingot/net/protocol/packet/handshake"
	"github.com/ingotmc/ingot/net/protocol/packet/login"
)

// DataByIDAndState tries to match the given id and state to a packet type
// and returns an error if it fails to do so.
func DataByIDAndState(id int32, state protocol.State) (Data, error) {
	switch state {
	case protocol.Handshaking:
		return handshakingData(id)
	case protocol.Login:
		return loginData(id)
	}
	return nil, ErrUnknownPacket{id, state}
}

func handshakingData(id int32) (Data, error) {
	switch id {
	case handshake.SetProtocolID:
		return new(handshake.SetProtocol), nil
	}
	return nil, ErrUnknownPacket{id, protocol.Handshaking}
}

func loginData(id int32) (Data, error) {
	switch id {
	case login.LoginStartID:
		return new(login.LoginStart), nil
	}
	return nil, ErrUnknownPacket{id, protocol.Login}
}
