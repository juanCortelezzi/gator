package gatorparser

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/juancortelezzi/uniqueid"
)

const PayloadLenLocation uint8 = 40

type PayloadLocation struct {
	ID            uniqueid.NanoID
	UnixTimestamp int64
	Latitude      float64
	Longitude     float64
}

func (p *PayloadLocation) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0, PayloadLenLocation)

	data = append(data, p.ID[:]...)
	data = binary.BigEndian.AppendUint64(data, uint64(p.UnixTimestamp))
	data = binary.BigEndian.AppendUint64(data, math.Float64bits(p.Latitude))
	data = binary.BigEndian.AppendUint64(data, math.Float64bits(p.Longitude))

	length := len(data)

	if length != int(PayloadLenLocation) {
		panic(fmt.Sprintf(
			"Invalid packet length: got len=%d but expected to be %d",
			HeaderSize,
			length,
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

	rawId, data := data[:uniqueid.Length], data[uniqueid.Length:]
	id, err := uniqueid.FromBytes(rawId)
	if err != nil {
		return fmt.Errorf("%w: could not parser ID", ErrInvalidField)
	}

	rawUnixTimestamp, data := data[:8], data[8:]
	unixTimestamp := int64(binary.BigEndian.Uint64(rawUnixTimestamp))

	rawLat, rawLng := data[:8], data[8:]
	lat := math.Float64frombits(binary.BigEndian.Uint64(rawLat))
	lng := math.Float64frombits(binary.BigEndian.Uint64(rawLng))

	h.ID = id
	h.UnixTimestamp = unixTimestamp
	h.Latitude = lat
	h.Longitude = lng

	return nil
}
