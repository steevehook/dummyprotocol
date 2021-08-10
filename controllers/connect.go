package controllers

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/steevehook/vprotocol/crypto"
	"github.com/steevehook/vprotocol/transport"
)

type connectRequest struct {
	Key crypto.PublicKey `json:"key"`
}

type connectResponse struct {
	Key crypto.PublicKey `json:"key"`
}

func (router Router) connect() func(io.Writer, request) error {
	return func(w io.Writer, r request) error {
		var body connectRequest
		err := parseReq(r, &body)
		if err != nil {
			transport.SendError(w, connectOperation, err)
			return err
		}

		key, err := crypto.NewECDHKey()
		if err != nil {
			transport.SendError(w, connectOperation, errors.New("could not connect"))
			return err
		}

		k, _ := body.Key.Curve.ScalarMult(body.Key.X, body.Key.Y, key.D.Bytes())
		shared := sha256.Sum256(k.Bytes())
		fmt.Printf("%x\n", shared)

		pub := crypto.PublicKey{
			Curve: key.PublicKey.Curve.Params(),
			X:     key.PublicKey.X,
			Y:     key.PublicKey.Y,
		}
		res := connectResponse{
			Key: pub,
		}
		transport.SendJSON(w, connectOperation, res)
		return nil
	}
}
