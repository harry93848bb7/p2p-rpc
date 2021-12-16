package main

import (
	"net"
	"time"

	"github.com/harry93848bb7/p2p-rpc/node"
	"github.com/harry93848bb7/p2p-rpc/protocol"
	"google.golang.org/protobuf/proto"
)

func main() {
	conn, err := net.DialTimeout("tcp", ":8085", 3*time.Second)
	if err != nil {
		panic(err)
	}
	aead, err := node.InboundHandshake(conn, 3*time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	n := node.NewNode()

	go func() {
		err := node.ReadConnection(conn, aead, n.Reader)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := node.WriteConnection(conn, aead, n.Writer)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)

			out, err := proto.Marshal(&protocol.Ping{
				Message:   "ping!",
				Timestamp: time.Now().Unix(),
			})
			if err != nil {
				panic(err)
			}
			n.Writer <- &node.Event{
				Type: node.PingMesage,
				Data: out,
			}
		}
	}()
	n.HandleEvents()
}
