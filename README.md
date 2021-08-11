# vprotocol

## Overview

This project implements the V protocol, which is really a dummy protocol
based on TCP, that features:

- ECDH key negotiation (handshake)
- Encrypted TCP connected with AES encryption
- A server that accepts concurrent & persistent TCP connections
- A client that communicates over encrypted stream with the server

### Server

- `ListenAndServe(settings Settings) (*VServer, error)`
- `Stop() error`

For `VServer` usage check out [`cmd/server/main.go`](./cmd/server/main.go)

### Client

- `Connect(addr string) error`
- `Ping() (transport.Message, error)`
- `Disconnect() error`

For `VClient` usage check out [`cmd/client/main.go`](./cmd/client/main.go)

## How to use

```shell script
# install all required dependencies
go get ./...
```

```shell script
# run the server
go run cmd/server/main.go
```

```shell script
# run the client
go run cmd/client/main.go
```

## Technologies Used

### Libraries

- viper (yaml configs)
- zap (logger)

### Encoding/Encryption

- P256 ECDH - for key negotiation
- AES - for encrypting the connection after handshake
