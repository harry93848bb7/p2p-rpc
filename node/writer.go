package node

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"time"
)

func WriteConnection(conn net.Conn, aead cipher.AEAD, pipe Pipline) error {
	defer func() {
		log.Println("closin writer...")
	}()
	var (
		n     int
		err   error
		sizei int
	)
	for {
		select {
		case event := <-pipe:
			nonce := make([]byte, 24, len(event.Data)+40)
			n, err = rand.Read(nonce)
			if err != nil {
				return err
			}
			if n != 24 {
				return errors.New("bad read nonce size")
			}

			messageType := make([]byte, 2)
			binary.BigEndian.PutUint16(messageType, uint16(event.Type))

			encryptedMsg := aead.Seal(nonce, nonce, append(messageType, event.Data...), nil)

			sizei = len(encryptedMsg)
			size := make([]byte, 4)
			binary.BigEndian.PutUint32(size, uint32(sizei))

			deadline := time.Now().Add(3 * time.Second)
			if err = conn.SetWriteDeadline(deadline); err != nil {
				return err
			}
			n, err = conn.Write(size)
			if err != nil {
				return err
			}
			if n != 4 {
				return errors.New("bad write packet size")
			}
			if err = conn.SetWriteDeadline(deadline); err != nil {
				return err
			}
			n, err = conn.Write(encryptedMsg)
			if err != nil {
				return err
			}
			if n != sizei {
				return err
			}
			if err = conn.SetWriteDeadline(time.Time{}); err != nil {
				return err
			}
		}
	}
}
