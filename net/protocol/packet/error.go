package packet

import "fmt"

type ErrUnknownPacket int64

func (e ErrUnknownPacket) Error() string {
	return fmt.Sprintf("unknown packet with id: %d", e)
}
