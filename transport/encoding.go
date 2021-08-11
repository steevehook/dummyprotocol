package transport

import (
	"bytes"
	"encoding/gob"
	"net"

	"github.com/steevehook/vprotocol/crypto"
)

type Message struct {
	Operation string
	Body      interface{}
}

func Encode(conn net.Conn, secret []byte, operation string, body interface{}) error {
	data, err := msg(operation, body)
	if err != nil {
		return err
	}

	encrypted, err := crypto.EncryptAES(data, secret)
	if err != nil {
		return err
	}

	encrypted = append(encrypted, sep...)
	_, err = conn.Write(encrypted)
	if err != nil {
		return err
	}

	return nil
}

func msg(operation string, body interface{}) ([]byte, error) {
	msg := Message{
		Operation: operation,
		Body:      body,
	}
	var buff bytes.Buffer
	err := gob.NewEncoder(&buff).Encode(msg)
	if err != nil {
		return []byte{}, err
	}
	return buff.Bytes(), nil
}
