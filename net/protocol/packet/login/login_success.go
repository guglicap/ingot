package login

import (
	"github.com/ingotmc/ingot/net/protocol/internal/write"
	"io"
)

type LoginSuccess struct {
	UUID string
	Name string
}

func (l LoginSuccess) Encode(w io.Writer) error {
	err := write.String(l.UUID, w)
	if err != nil {
		return err
	}
	err = write.String(l.Name, w)
	return err
}
