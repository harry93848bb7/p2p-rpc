package node

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

func ReadConnection(conn net.Conn, aead cipher.AEAD, pipe Pipline) error {
	defer func() {
		log.Println("closin reader...")
	}()
	var (
		n    int
		size int
		err  error
	)
	for {

		header := make([]byte, 4)
		n, err = io.ReadAtLeast(conn, header, 4)
		if err != nil {
			return err
		}
		if n != 4 {
			return errors.New("bad read packet size")
		}

		// minimum size = nonce + overhead + messagetype
		size = int(binary.BigEndian.Uint32(header))
		if size < 42 || size > 32000000 {
			return errors.New("too large / small packet size")
		}
		packet := make([]byte, size)

		if err = conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
			return err
		}
		n, err = io.ReadAtLeast(conn, packet, size)
		if err != nil {
			return err
		}
		if n != size {
			return errors.New("bad read packet size")
		}
		if err = conn.SetReadDeadline(time.Time{}); err != nil {
			log.Println(err)
			return err
		}
		plaintext, err := aead.Open(nil, packet[:24], packet[24:], nil)
		if err != nil {
			log.Println(err)
			return err
		}
		if len(plaintext) < 2 {
			return errors.New("bad plaintext packet size")
		}
		pipe <- &Event{
			Type: MessageType(binary.BigEndian.Uint16(plaintext[:2])),
			Data: plaintext[2:],
		}
	}
}
