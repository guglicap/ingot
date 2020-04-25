package packet

import (
	"fmt"
	"github.com/ingotmc/ingot/net/protocol"
)

// ErrUnknownPacket is returned when a packet can't be decoded because the combination of id and state is unknown.
type ErrUnknownPacket struct {
	id    int32
	state protocol.State
}

func (e ErrUnknownPacket) Error() string {
	return fmt.Sprintf("unknown packet for state: %d with id: %d", e.state, e.id)
}

// ErrInvalidDataLength is returned when the amount of bytes read while receiving a packet data doesn't
// match the length declared in the packet header.
type ErrInvalidDataLength struct {
	length, expected int
}

func (e ErrInvalidDataLength) Error() string {
	return fmt.Sprintf("read less bytes than expected: read %d, expected %d", e.length, e.expected)
}
