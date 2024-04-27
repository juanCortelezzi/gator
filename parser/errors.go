package gatorparser

import "errors"

var (
	ErrVersionMismatch      = errors.New("Version mismatch!")
	ErrNotEnoughData        = errors.New("Not enough data")
	ErrInvalidPacketType    = errors.New("Invalid packet type")
	ErrInvalidPayloadLength = errors.New("Invalid payload length")
	ErrInvalidField         = errors.New("Invalid field")
)
