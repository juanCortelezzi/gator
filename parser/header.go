package gatorparser

import (
	"encoding"
	"fmt"
)

const HeaderSize uint8 = 3

type Header struct {
	Version uint8
	Type    uint8
	Length  uint8
}

var _ encoding.BinaryMarshaler = &Header{}
var _ encoding.BinaryUnmarshaler = &Header{}

func (h *Header) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0, HeaderSize)
	data = append(data, h.Version)
	data = append(data, h.Type)
	data = append(data, h.Length)

	length := len(data)

	if length != int(HeaderSize) {
		panic(fmt.Sprintf(
			"Invalid header length marshaling header: len=%d but expected to be %d",
			length,
			HeaderSize,
		))
	}

	return data, nil
}

func (h *Header) UnmarshalBinary(data []byte) error {
	length := len(data)
	if length < int(HeaderSize) {
		return fmt.Errorf(
			"Not enough data to unmarshal header: got %d but expected at least %d - %w",
			length,
			HeaderSize,
			ErrNotEnoughData,
		)
	}

	if data[0] != Version {
		return fmt.Errorf(
			"Version mismatch unmarshaling header: got %d but expected %d - %w",
			data[0],
			Version,
			ErrVersionMismatch,
		)
	}

	h.Version = data[0]
	h.Type = data[1]
	h.Length = data[2]

	return nil
}
