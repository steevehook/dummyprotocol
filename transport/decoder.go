package transport

import (
	"encoding/json"
	"io"
)

// Decode decodes a certain reader into the specified type
func Decode(r io.Reader, v interface{}) error {
	err := json.NewDecoder(r).Decode(v)
	if err != nil {
		return err
	}
	return nil
}
