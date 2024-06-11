package gatorparser

import (
	"errors"
	"fmt"
)

const (
	Version uint8 = 1

	PacketTypeLocation uint8 = iota
	PacketTypeEcho
)

var (
	ErrVersionMismatch      = errors.New("Version mismatch!")
	ErrNotEnoughData        = errors.New("Not enough data")
	ErrInvalidPacketType    = errors.New("Invalid packet type")
	ErrInvalidPayloadLength = errors.New("Invalid payload length")
	ErrInvalidField         = errors.New("Invalid field")
)

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
