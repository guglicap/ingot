package protocol

type State int

const (
	Handshaking State = iota
	Status
	Login
	Play
)
