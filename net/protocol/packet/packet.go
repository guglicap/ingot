package packet

import (
	"bytes"
	"github.com/ingotmc/ingot/net/protocol"
	"github.com/ingotmc/ingot/net/protocol/internal/read"
	"github.com/ingotmc/ingot/net/protocol/internal/write"
	"io"
)

// Data represents packet data.
type Data interface {
	Encode(w io.Writer) error
	Decode(r io.Reader) error
}

// Packet represents a packet whose Data has been decoded / is ready to be encoded.
type Packet struct {
	ID    int32
	State protocol.State
	Data  Data
}

// Raw represents a packet on the wire, without the length metadata.
// This data is held as a []byte.
type Raw struct {
	ID   int32
	Data []byte
}

// ReadRaw reads from the Reader r and returns a packet.Raw
func ReadRaw(r io.Reader) (Raw, error) {
	p := Raw{}
	length, _, err := read.VarInt(r)
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
		return p, ErrInvalidDataLength{
			length:   l,
			expected: dataLen,
		}
	}
	p.Data = buf
	return p, nil
}

// WriteRaw writes a Raw packet to the wire and returns in error
// if the encoding cannot take place successfully.
func WriteRaw(p Raw, w io.Writer) error {
	buf := bytes.NewBuffer(make([]byte, 0, 5))
	n, err := write.VarInt(p.ID, buf)
	if err != nil {
		return err
	}
	_, err = write.VarInt(int32(n+len(p.Data)), w)
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes()[:n])
	if err != nil {
		return err
	}
	_, err = w.Write(p.Data)
	return err
}
