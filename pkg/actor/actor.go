package actor

type Actor struct {
	ID         string
	State      string // actor state
	LogicIndex int    // index of logic in state
	AwaitEvent bool   // is we stack in await event

	Data         []byte
	DataTypeName string
}
