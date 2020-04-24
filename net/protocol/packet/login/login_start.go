package login

import (
	"github.com/ingotmc/ingot/net/protocol/internal/read"
	"io"
)

const LoginStartID = 0x00

type LoginStart struct {
	Name string
}

func (l LoginStart) Encode(w io.Writer) error {
	// might throw an error if we decide that calling encode on a serverbound packet isn't a no-op
	return nil
}

func (l *LoginStart) Decode(r io.Reader) error {
	name, err := read.String(r)
	if err != nil {
		return err
	}
	l.Name = name
	return nil
}
