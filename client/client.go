package client

import (
	"crypto/ecdsa"
	"errors"
	"log"
	"net"

	"github.com/steevehook/vprotocol/crypto"
	"github.com/steevehook/vprotocol/transport"
)

type VClient struct {
	conn   net.Conn
	secret []byte
}

func (c *VClient) Connect(addr string) error {
	if c.conn != nil {
		return errors.New("already connected")
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("could not dial connection")
		return err
	}
	c.conn = conn

	privateKey, err := crypto.NewECDHKey()
	if err != nil {
		return err
	}

	err = crypto.EncodeECDHPublicKey(conn, privateKey.PublicKey)
	if err != nil {
		return err
	}

	var serverPublicKey *ecdsa.PublicKey
	err = crypto.DecodeECDHPublicKey(conn, &serverPublicKey)
	if err != nil {
		return err
	}

	c.secret = crypto.ECDHSecret(serverPublicKey, privateKey)
	return nil
}

func (c *VClient) Ping() (transport.Message, error) {
	if c.conn == nil {
		return transport.Message{}, errors.New("client has not yet connected")
	}

	err := transport.Encode(c.conn, c.secret, "ping", nil)
	if err != nil {
		return transport.Message{}, err
	}

	scanner := transport.NewVScanner(c.conn)
	if !scanner.Scan() {
		return transport.Message{}, errors.New("could not scan connection")
	}

	msg, err := transport.Decode(scanner.Bytes(), c.secret)
	if err != nil {
		return transport.Message{}, err
	}

	return msg, nil
}

func (c *VClient) Disconnect() error {
	if c.conn == nil {
		return errors.New("already disconnected")
	}

	return nil
}
