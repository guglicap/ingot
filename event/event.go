package event

import "context"

type (
	Event interface {
		ID() string
	}

	Handler interface {
		Push(Event)
		Subscribe(string, func(context.Context, Event))
	}
)
