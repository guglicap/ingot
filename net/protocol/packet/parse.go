package packet

import (
	"github.com/ingotmc/ingot/net/protocol"
	"github.com/ingotmc/ingot/net/protocol/packet/login"
)

func DataByIDAndState(id int64, state protocol.State) Data {
	switch state {
	case protocol.Login:
		return loginData(id)
	}
	return nil
}

func loginData(id int64) Data {
	switch id {
	case login.LoginStartID:
		return new(login.LoginStart)
	}
	return nil
}
