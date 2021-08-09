package client

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"math/big"
)

type publicKey struct {
	Curve *elliptic.CurveParams `json:"Curve"`
	X     *big.Int              `json:"X"`
	Y     *big.Int              `json:"Y"`
}

func newECDHKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Println("could not generate private key")
		return nil, err
	}
	return privateKey, nil
}
