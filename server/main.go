package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"strconv"

	"github.com/juancortelezzi/gatorparser"
	"github.com/juancortelezzi/gatorserver/pkg/client"
	"github.com/juancortelezzi/gatorserver/pkg/clog"
	"github.com/juancortelezzi/gatorserver/pkg/ka"
)

func main() {
	ctx := context.Background()
	logger := clog.NewLogger(os.Stdout, slog.LevelInfo)

	if err := Run(ctx, logger, os.LookupEnv); err != nil {
		logger.ErrorContext(ctx, "error in top level", "err", err)
		os.Exit(1)
	}
}

func Run(ctx context.Context, logger clog.Logger, lookupEnv func(string) (string, bool)) error {

	logger.DebugContext(ctx, "looking env variables")

	portString, found := lookupEnv("PORT")
	if !found {
		return ka.Boom("PORT environment variable not found")
	}

	port, err := strconv.Atoi(portString)
	if err != nil {
		return ka.Boom("PORT environment variable is not a number")
	}

	addr := &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				logger.WarnContext(ctx, "Failed to accept connection", "error", err)
				continue
			}

			tcpConn, ok := conn.(*net.TCPConn)
			if !ok {
				logger.WarnContext(ctx, "Failed to cast connection to TCPConn")
				if err := conn.Close(); err != nil {
					logger.ErrorContext(ctx, "Failed to close connection", "error", err)
				}
				continue
			}

			go handleConnection(ctx, logger, tcpConn)
		}
	}
}

func handleConnection(ctx context.Context, logger clog.Logger, conn *net.TCPConn) {
	defer conn.Close()

	logger.InfoContext(ctx, "Handling connection", "remote_addr", conn.RemoteAddr())

	frameReader := client.NewFrameReader(conn)
	buffer := make([]byte, 0, client.BufferSize)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			pointer, err := frameReader.Read(buffer)
			if err != nil {
				logger.ErrorContext(ctx, "ReadPump fail when reading packet", "err", err)
				return
			}

			rawPacket := buffer[:pointer]

			var packetHeader gatorparser.Header
			if err := packetHeader.UnmarshalBinary(rawPacket); err != nil {
				logger.ErrorContext(ctx, "ReadPump fail when reading packet header", "err", err)
				return
			}

			switch packetHeader.Type {
			case gatorparser.PacketTypeLocation:
				var packetPayload gatorparser.PacketLocation
				if err := packetPayload.UnmarshalBinary(rawPacket[gatorparser.HeaderSize:]); err != nil {
					logger.ErrorContext(ctx, "ReadPump fail when reading packet payload", "err", err)
					return
				}

			default:
				logger.ErrorContext(ctx, "ReadPump unhandled packet due to unknown type")
				return
			}
		}
	}
}
