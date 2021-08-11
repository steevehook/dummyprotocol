package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"encoding/gob"
	"io"
)

func NewECDHKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func EncodeECDHPublicKey(w io.Writer, publicKey ecdsa.PublicKey) error {
	gob.Register(ecdsa.PublicKey{})
	gob.Register(elliptic.P256())
	return gob.NewEncoder(w).Encode(publicKey)
}

func DecodeECDHPublicKey(r io.Reader, v interface{}) error {
	gob.Register(ecdsa.PublicKey{})
	gob.Register(elliptic.P256())
	return gob.NewDecoder(r).Decode(v)
}

func ECDHSecret(publicKey *ecdsa.PublicKey, privateKey *ecdsa.PrivateKey) []byte {
	x, _ := publicKey.Curve.ScalarMult(publicKey.X, publicKey.Y, privateKey.D.Bytes())
	return x.Bytes()
}

func EncryptAES(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(createHash(key))
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}

	bs := gcm.Seal(nonce, nonce, data, nil)
	return bs, nil
}

func DecryptAES(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(createHash(key))
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	bs, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return []byte{}, err
	}

	return bs, nil
}

func createHash(key []byte) []byte {
	hash := md5.New()
	hash.Write(key)
	return hash.Sum(nil)
}
