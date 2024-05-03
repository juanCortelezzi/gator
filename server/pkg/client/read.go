package client

import (
	"bytes"
	"context"
	"net"

	"github.com/juancortelezzi/gatorparser"
	"github.com/juancortelezzi/gatorserver/pkg/clog"
)

const BufferSize = 1024

func ReadPump(ctx context.Context, logger clog.Logger, conn *net.TCPConn) {
	frameReader := NewFrameReader(conn)
	buffer := make([]byte, BufferSize)

	logger.InfoContext(ctx, "ReadPump started")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			logger.DebugContext(ctx, "ReadPump waiting for packet")
			pointer, err := frameReader.Read(buffer)
			logger.DebugContext(ctx, "ReadPump got packet")
			if err != nil {
				logger.ErrorContext(ctx, "ReadPump fail when reading packet", "err", err)
				return
			}

			rawPacket := buffer[:pointer]
			if err := executePacket(ctx, logger, conn, rawPacket); err != nil {
				logger.ErrorContext(ctx, "ReadPump fail when executing packet", "err", err)
				return
			}
		}
	}
}

func executePacket(ctx context.Context, logger clog.Logger, conn *net.TCPConn, rawPacket []byte) error {
	var packetHeader gatorparser.Header
	if err := packetHeader.UnmarshalBinary(rawPacket); err != nil {
		return err
	}

	switch packetHeader.Type {
	case gatorparser.PacketTypeLocation:
		var packetPayload gatorparser.PayloadLocation
		if err := packetPayload.UnmarshalBinary(rawPacket[gatorparser.HeaderSize:]); err != nil {
			return err
		}

		logger.InfoContext(
			ctx, "Received location",
			"latitude", packetPayload.Latitude,
			"longitude", packetPayload.Longitude,
		)

		headerBytes, err := packetHeader.MarshalBinary()
		if err != nil {
			return err
		}

		payloadBytes, err := packetPayload.MarshalBinary()
		if err != nil {
			return err
		}

		messageBytes := bytes.Join([][]byte{headerBytes, payloadBytes}, nil)
		if _, err := conn.Write(messageBytes); err != nil {
			return err
		}

		return nil
	default:
		logger.ErrorContext(ctx, "ReadPump unhandled packet due to unknown type")
		return nil
	}
}
