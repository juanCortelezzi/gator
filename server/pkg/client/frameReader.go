package client

import (
	"errors"
	"fmt"
	"io"

	"github.com/juancortelezzi/gatorparser"
)

type FrameReader struct {
	reader   io.Reader
	previous []byte
	scratch  []byte
}

var _ io.Reader = &FrameReader{}

func NewFrameReader(reader io.Reader) *FrameReader {
	return &FrameReader{
		reader:   reader,
		previous: make([]byte, 0, BufferSize),
		scratch:  make([]byte, BufferSize),
	}
}

func (r *FrameReader) getPacketLength(data []byte) (int, error) {
	var header gatorparser.Header
	if err := header.UnmarshalBinary(data); err != nil {
		if errors.Is(err, gatorparser.ErrNotEnoughData) {
			return 0, nil
		}

		return 0, fmt.Errorf("Received unexpected packet: %w", err)
	}

	packetLength := int(gatorparser.HeaderSize + header.Length)
	if len(data) < packetLength {
		return 0, nil
	}

	return packetLength, nil
}

func (r *FrameReader) Read(p []byte) (n int, err error) {
	for {
		packetLength, err := r.getPacketLength(r.previous)
		if err != nil {
			return 0, err
		}

		if packetLength > 0 {
			copy(p, r.previous[:packetLength])
			r.previous = r.previous[packetLength:]
			return packetLength, nil
		}

		pointer, err := r.reader.Read(r.scratch)
		if err != nil {
			return 0, nil
		}

		r.previous = append(r.previous, r.scratch[:pointer]...)
	}
}
