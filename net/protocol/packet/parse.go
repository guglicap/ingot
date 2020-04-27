package packet

import (
	"github.com/ingotmc/ingot/net/protocol"
	"github.com/ingotmc/ingot/net/protocol/packet/handshake"
	"github.com/ingotmc/ingot/net/protocol/packet/login"
)

// ServerboundByIDAndState tries to match the given id and state to a packet type
// and returns an error if it fails to do so.
func ServerboundByIDAndState(id ID, state protocol.State) (Serverbound, error) {
	p := Serverbound{}
	p.ID = id
	data, err := sBoundDataByIDAndState(id, state)
	p.Data = data
	return p, err
}

func sBoundDataByIDAndState(id ID, state protocol.State) (sBoundData, error) {
	switch state {
	case protocol.Handshaking:
		return handshakingData(id)
	case protocol.Login:
		return loginData(id)
	}
	return nil, ErrUnknownPacket{int32(id), state}
}

func handshakingData(id ID) (sBoundData, error) {
	switch id {
	case handshake.SetProtocolID:
		return new(handshake.SetProtocol), nil
	}
	return nil, ErrUnknownPacket{int32(id), protocol.Handshaking}
}

func loginData(id ID) (sBoundData, error) {
	switch id {
	case login.LoginStartID:
		return new(login.LoginStart), nil
	}
	return nil, ErrUnknownPacket{int32(id), protocol.Login}
}
