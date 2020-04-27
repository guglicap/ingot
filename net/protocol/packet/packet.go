package packet

import (
	"bytes"
	"github.com/ingotmc/ingot/net/protocol"
	"github.com/ingotmc/ingot/net/protocol/internal/read"
	"github.com/ingotmc/ingot/net/protocol/internal/write"
	"io"
)

type sBoundData interface{ Decode(r io.Reader) error }

// PacketID represent the ID field of a packet.
type ID int32

// Info represents the combination of ID + State which uniquely identifies a packet, given its direction
type Info struct {
	ID
	protocol.State
}

// Serverbound represents a packet sent by the client to the server.
type Serverbound struct {
	Info
	Data sBoundData
}

type cBoundData interface{ Encode(w io.Writer) error }

// Clientbound represents a packet sent by the server to the client.
type Clientbound struct {
	ID   ID // we don't use Info here because State isn't needed to encode a packet.
	Data cBoundData
}

// Raw represents a packet on the wire, without the length metadata.
// This data is held as a []byte.
type Raw struct {
	ID   ID
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
	p.ID = ID(id)
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
	n, err := write.VarInt(int32(p.ID), buf)
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
