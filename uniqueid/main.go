package uniqueid

import (
	"errors"
	nanoid "github.com/matoous/go-nanoid/v2"
)

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	Length   = 18
)

var (
	ErrInvalidID = errors.New("Error Invalid ID")
)

type NanoID = [Length]byte

func New() (NanoID, error) {
	id := [Length]byte{}
	stringId, err := nanoid.Generate(alphabet, Length)
	if err != nil {
		return id, err
	}

	for i, char := range stringId {
		id[i] = byte(char)
	}

	if len(id) != Length {
		panic("Invalid ID length")
	}

	return id, nil
}

func Must() NanoID {
	stringId := nanoid.MustGenerate(alphabet, Length)
	id := [Length]byte{}

	for i, char := range stringId {
		id[i] = byte(char)
	}

	if len(id) != Length {
		panic("Invalid ID length")
	}

	return id
}

func FromBytes(b []byte) (NanoID, error) {
	id := [Length]byte{}

	if len(b) != Length {
		return id, ErrInvalidID
	}

	for i, byte := range b {
		id[i] = byte
	}

	return id, nil
}
