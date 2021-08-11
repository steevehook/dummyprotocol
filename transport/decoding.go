package transport

import (
	"bytes"
	"encoding/gob"

	"github.com/steevehook/vprotocol/crypto"
)

func Decode(data, secret []byte) (Message, error) {
	decrypted, err := crypto.DecryptAES(data, secret)
	if err != nil {
		return Message{}, err
	}

	decodeBuff := bytes.NewBuffer(decrypted)
	var msg Message
	err = gob.NewDecoder(decodeBuff).Decode(&msg)
	if err != nil {
		return Message{}, err
	}

	return msg, nil
}
