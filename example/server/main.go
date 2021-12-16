package main

import (
	"net"
	"time"

	"github.com/harry93848bb7/p2p-rpc/node"
	"github.com/harry93848bb7/p2p-rpc/protocol"
	"google.golang.org/protobuf/proto"
)

func main() {
	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		panic(err)
	}
	for {
		peer, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go HandlePeer(peer)
	}
}

func HandlePeer(peer net.Conn) {
	aead, err := node.OutboundHandshake(peer, 3*time.Second)
	if err != nil {
		panic(err)
	}
	defer peer.Close()

	n := node.NewNode()

	go func() {
		err := node.ReadConnection(peer, aead, n.Reader)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := node.WriteConnection(peer, aead, n.Writer)
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
