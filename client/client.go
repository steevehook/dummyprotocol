package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func NewClient() Client {
	return Client{}
}

type Client struct {
}

func (c *Client) Connect(ipv4 string, port int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ipv4, port))
	if err != nil {
		log.Printf("could not dial connection")
		return err
	}

	k, err := newECDHKey()
	if err != nil {
		return err
	}
	//bs, err := json.Marshal(k.PublicKey)
	//if err != nil {
	//	return err
	//}

	//_, err = conn.Write(bs)
	//if err != nil {
	//	log.Println("could not write to connection")
	//	return err
	//}
	//fmt.Println("client", k.PublicKey.X)
	//fmt.Println("client", k.PublicKey.Y)

	pub := publicKey{
		Curve: k.PublicKey.Curve.Params(),
		X: k.PublicKey.X,
		Y: k.PublicKey.Y,
	}
	err = json.NewEncoder(conn).Encode(pub)
	if err != nil {
		log.Println("could not encode json")
		return err
	}
	//conn.Close()
	return nil
}

func (c *Client) Disconnect() error {
	return nil
}

func (c *Client) Ping() ([]byte, error) {
	return []byte{}, nil
}
