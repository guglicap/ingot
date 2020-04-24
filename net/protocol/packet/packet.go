package packet

import (
	"github.com/ingotmc/ingot/net/protocol"
	"github.com/ingotmc/ingot/net/protocol/internal/read"
	"io"
)

type Data interface {
	Encode(w io.Writer) error
	Decode(r io.Reader) error
}

type Packet struct {
	ID    int64
	State protocol.State
	Data  Data
}

type Raw struct {
	Length int64
	ID     int64
	Data   []byte
}

func ReadRaw(r io.Reader) (Raw, error) {
	p := Raw{}
	length, _, err := read.VarInt(r)
	p.Length = length
	if err != nil {
		return Raw{}, err
	}
	id, n, err := read.VarInt(r)
	if err != nil {
		return Raw{}, err
	}
	p.ID = id
	dataLen := int(length) - n
	buf := make([]byte, dataLen)
	l, err := io.ReadFull(r, buf)
	if err != nil {
		return Raw{}, err
	}
	if l != dataLen {
		// TODO: return error
	}
	p.Data = buf
	return p, nil
}
