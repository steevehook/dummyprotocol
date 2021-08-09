package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
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
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		fmt.Println("line", string(scanner.Bytes()))
	}
	var pub publicKey
	err := json.NewDecoder(bytes.NewReader(scanner.Bytes())).Decode(&pub)
	if err != nil {
		log.Printf("could not decode json: %v", err)
	}

	fmt.Println("from server", pub.X)
	fmt.Println("from server", pub.Y)
}

func (s server) Close() error {
	return nil
}
