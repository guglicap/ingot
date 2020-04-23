package event

type (
	Event interface {
		ID() string
	}

	Handler interface {
		Emit(Event)
		Subscribe(string, func(Event))
	}
)
