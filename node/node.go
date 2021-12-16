package node

import (
	"log"

	"github.com/harry93848bb7/p2p-rpc/protocol"
	"google.golang.org/protobuf/proto"
)

type Node struct {
	Reader Pipline
	Writer Pipline
}

func NewNode() *Node {
	return &Node{
		Reader: make(Pipline),
		Writer: make(Pipline),
	}
}

func (n *Node) HandleEvents() {
	var err error
	for {
		select {
		case event := <-n.Reader:
			switch event.Type {
			case PingMesage:

				// handle ping message events
				message := &protocol.Ping{}
				if err = proto.Unmarshal(event.Data, message); err != nil {
					log.Println("failed to parse ping message:", err)
				} else {
					log.Println("ping message found: ", message)
				}

				// send pong reply back
				out, err := proto.Marshal(&protocol.Pong{
					Message:   "pong!",
					Timestamp: message.Timestamp,
				})
				if err != nil {
					panic(err)
				}
				n.Writer <- &Event{
					Type: PongMessage,
					Data: out,
				}

			case PongMessage:

				// handle pong reply!
				message := &protocol.Pong{}
				if err = proto.Unmarshal(event.Data, message); err != nil {
					log.Println("failed to parse pong message:", err)
				} else {
					log.Println("pong message found: ", message)
				}

			default:
				log.Println("unknown message type from peer: ", event.Type)
			}
		}
	}
}
