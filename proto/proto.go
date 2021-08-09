package proto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
)

func newKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Println("could not generate private key")
		return nil, err
	}
	return privateKey, nil
}
