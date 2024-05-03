package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/juancortelezzi/gatorserver/pkg/client"
	"github.com/juancortelezzi/gatorserver/pkg/clog"
)

var (
	ErrEnvVariableNotFound    = errors.New("Environment variable not found")
	ErrEnvVariableParseFailed = errors.New("Failed to parse environment variable")
)

func Run(ctx context.Context, logger clog.Logger, lookupEnv func(string) (string, bool)) error {
	logger.DebugContext(ctx, "looking env variables")

	portString, found := lookupEnv("PORT")
	if !found {
		return fmt.Errorf("%w: PORT", ErrEnvVariableNotFound)
	}

	port, err := strconv.Atoi(portString)
	if err != nil {
		return fmt.Errorf("%w: PORT is not a number", ErrEnvVariableParseFailed)
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

	logger.InfoContext(ctx, "Server started", "addr", listener.Addr())

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

	client.ReadPump(ctx, logger, conn)

	logger.InfoContext(ctx, "Connection closed", "remote_addr", conn.RemoteAddr())
}
