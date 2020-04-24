package protocol

type State int

const (
	Handshaking State = iota
	Login
	Play
	Status
)
