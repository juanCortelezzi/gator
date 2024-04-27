package gatorparser

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"math"
	"slices"

	"github.com/google/uuid"
)

const (
	Version            uint8 = 1
	HeaderSize         uint8 = 3
	PayloadLenLocation uint8 = 40
)

const (
	PacketTypeLocation uint8 = iota
)

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

type PacketLocation struct {
	Uuid          [16]byte
	UnixTimestamp int64
	Latitude      float64
	Longitude     float64
}

func (p *PacketLocation) MarshalBinary() ([]byte, error) {
	header := &Header{
		Version: Version,
		Type:    PacketTypeLocation,
		Length:  PayloadLenLocation,
	}

	data, err := header.MarshalBinary()
	if err != nil {
		panic(fmt.Sprintf("Error marshaling hardcoded header in PacketLocation: %v", err))
	}

	data = slices.Grow(data, int(header.Length))

	data = append(data, p.Uuid[:]...)
	data = binary.BigEndian.AppendUint64(data, uint64(p.UnixTimestamp))
	data = binary.BigEndian.AppendUint64(data, math.Float64bits(p.Latitude))
	data = binary.BigEndian.AppendUint64(data, math.Float64bits(p.Longitude))

	length := len(data)

	if length != int(HeaderSize+header.Length) {
		panic(fmt.Sprintf(
			"Invalid packet length: len=%d but expected to be %d",
			length,
			HeaderSize,
		))
	}

	return data, nil
}

func (h *PacketLocation) UnmarshalBinary(data []byte) error {
	var header Header
	if err := header.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("Invalid header unmarshaling PacketLocation: %w", err)
	}

	if header.Length != PayloadLenLocation {
		return fmt.Errorf(
			"Invalid payload length unmarshaling PacketLocation: got %d but expected %d - %w",
			header.Length,
			PayloadLenLocation,
			ErrInvalidPayloadLength,
		)
	}

	if header.Type != PacketTypeLocation {
		return fmt.Errorf(
			"Invalid packet type unmarshaling PacketLocation: got %d but expected %d - %w",
			header.Type,
			PacketTypeLocation,
			ErrInvalidPacketType,
		)
	}

	if len(data) < int(HeaderSize+header.Length) {
		return fmt.Errorf(
			"Not enough data to unmarshal PacketLocaiton: got %d but expected at least %d - %w",
			len(data),
			HeaderSize+header.Length,
			ErrNotEnoughData,
		)
	}

	payload := data[HeaderSize : HeaderSize+header.Length]

	rawId, payload := payload[:16], payload[16:]
	id, err := uuid.FromBytes(rawId)
	if err != nil {
		return fmt.Errorf("Invalid field Uuid unmarshaling PacketLocation: %w", err)
	}

	rawUnixTimestamp, payload := payload[:8], payload[8:]
	unixTimestamp := int64(binary.BigEndian.Uint64(rawUnixTimestamp))

	rawLat, rawLng := payload[:8], payload[8:]
	lat := math.Float64frombits(binary.BigEndian.Uint64(rawLat))
	lng := math.Float64frombits(binary.BigEndian.Uint64(rawLng))

	h.Uuid = id
	h.UnixTimestamp = unixTimestamp
	h.Latitude = lat
	h.Longitude = lng

	return nil
}
