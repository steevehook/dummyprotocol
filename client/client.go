package client

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/steevehook/vprotocol/crypto"
	"github.com/steevehook/vprotocol/transport"
)

type request struct {
	Operation string      `json:"operation"`
	Body      interface{} `json:"body"`
}

type response struct {
	Body struct{
		Key crypto.PublicKey `json:"key"`
	} `json:"body"`
}

func New() Client {
	return Client{}
}

type Client struct {

}

func (c *Client) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("could not dial connection")
		return err
	}

	k, err := crypto.NewECDHKey()
	if err != nil {
		return err
	}

	pub := crypto.PublicKey{
		Curve: k.PublicKey.Curve.Params(),
		X:     k.PublicKey.X,
		Y:     k.PublicKey.Y,
	}

	r := request{
		Operation: "connect",
		Body: map[string]interface{}{
			"key": pub,
		},
	}
	err = json.NewEncoder(conn).Encode(r)
	if err != nil {
		log.Println("could not encode json")
		return err
	}

	var res response
	err = transport.Decode(conn, &res)
	if err != nil {
		log.Println("could not decode json")
		return err
	}
	fmt.Println("res", res.Body.Key.X)
	fmt.Println("res", res.Body.Key.Y)

	kk, _ := res.Body.Key.Curve.ScalarMult(res.Body.Key.X, res.Body.Key.Y, k.D.Bytes())
	shared := sha256.Sum256(kk.Bytes())
	fmt.Printf("%x\n", shared)
	return nil
}

func (c *Client) Disconnect() error {
	return nil
}

func (c *Client) Ping() ([]byte, error) {
	return []byte{}, nil
}
