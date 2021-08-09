package proto

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
	li net.Listener
}

func Listen(port int) error {
	li, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Printf("could not listen on port: %d", port)
		return err
	}
	srv := &server{
		li: li,
	}
	defer func() {
		_ = srv.Close()
	}()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Printf("could not accept connection: %v", err)
			continue
		}

		go srv.serve(conn)
	}
}

func (s server) serve(conn net.Conn) {
	bs, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Println("could not read from connection")
	}
	fmt.Println("public key from client", string(bs))
	time.Sleep(5*time.Second)
	_ = conn.Close()
}

func (s server) Close() error {
	return nil
}
