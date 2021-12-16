package node

type Pipline chan *Event

type MessageType uint16

type Event struct {
	Type MessageType
	Data []byte
}

var (
	UnknownMesage MessageType = 0
	PingMesage    MessageType = 1
	PongMessage   MessageType = 2
)
