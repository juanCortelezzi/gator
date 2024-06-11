package gatorparser

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/google/uuid"
)

const PayloadEchoSize uint8 = 10

type PayloadEcho struct {
	Uuid          [16]byte
	UnixTimestamp int64
	Latitude      float64
	Longitude     float64
}

func (p *PayloadEcho) MarshalBinary() ([]byte, error) {
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

func (h *PayloadEcho) UnmarshalBinary(data []byte) error {
	if len(data) < int(PayloadLenLocation) {
		return fmt.Errorf(
			"Not enough data to unmarshal PayloadEcho: got %d but expected at least %d - %w",
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
