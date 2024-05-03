package gatorparser

import (
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

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

var (
	ErrVersionMismatch      = errors.New("Version mismatch!")
	ErrNotEnoughData        = errors.New("Not enough data")
	ErrInvalidPacketType    = errors.New("Invalid packet type")
	ErrInvalidPayloadLength = errors.New("Invalid payload length")
	ErrInvalidField         = errors.New("Invalid field")
)

type Header struct {
	Version uint8
	Type    uint8
	Length  uint8
}

var _ encoding.BinaryMarshaler = &Header{}
var _ encoding.BinaryUnmarshaler = &Header{}

func NewHeader(packetType uint8) *Header {
	// FIX: I hate this so much: there is no way that every time we
	// create a new payload we have to change 17 lines of
	// boilerplate code
	var length uint8
	switch packetType {
	case PacketTypeLocation:
		length = PayloadLenLocation
	default:
		panic(fmt.Sprintf("Invalid packet type: %d", packetType))
	}

	return &Header{
		Version: Version,
		Type:    packetType,
		Length:  length,
	}
}

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

type PayloadLocation struct {
	Uuid          [16]byte
	UnixTimestamp int64
	Latitude      float64
	Longitude     float64
}

func (p *PayloadLocation) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0, PayloadLenLocation)

	data = append(data, p.Uuid[:]...)
	data = binary.BigEndian.AppendUint64(data, uint64(p.UnixTimestamp))
	data = binary.BigEndian.AppendUint64(data, math.Float64bits(p.Latitude))
	data = binary.BigEndian.AppendUint64(data, math.Float64bits(p.Longitude))

	length := len(data)

	if length != int(PayloadLenLocation) {
		panic(fmt.Sprintf(
			"Invalid packet length: got len=%d but expected to be %d",
			length,
			HeaderSize,
		))
	}

	return data, nil
}

func (h *PayloadLocation) UnmarshalBinary(data []byte) error {
	if len(data) < int(PayloadLenLocation) {
		return fmt.Errorf(
			"Not enough data to unmarshal PayloadLocation: got %d but expected at least %d - %w",
			len(data),
			PayloadLenLocation,
			ErrNotEnoughData,
		)
	}

	rawId, data := data[:16], data[16:]
	id, err := uuid.FromBytes(rawId)
	if err != nil {
		return fmt.Errorf("Invalid field Uuid unmarshaling PacketLocation: %w", err)
	}

	rawUnixTimestamp, data := data[:8], data[8:]
	unixTimestamp := int64(binary.BigEndian.Uint64(rawUnixTimestamp))

	rawLat, rawLng := data[:8], data[8:]
	lat := math.Float64frombits(binary.BigEndian.Uint64(rawLat))
	lng := math.Float64frombits(binary.BigEndian.Uint64(rawLng))

	h.Uuid = id
	h.UnixTimestamp = unixTimestamp
	h.Latitude = lat
	h.Longitude = lng

	return nil
}
