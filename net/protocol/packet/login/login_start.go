// package login implements packet data for the connection state Login
package login

import (
	"github.com/ingotmc/ingot/net/protocol/internal/read"
	"io"
)

const LoginStartID = 0x00

type LoginStart struct {
	Name string
}

func (l *LoginStart) Decode(r io.Reader) error {
	name, err := read.String(r)
	if err != nil {
		return err
	}
	l.Name = name
	return nil
}
