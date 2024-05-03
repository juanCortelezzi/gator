package main_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/juancortelezzi/gatorparser"
	"github.com/juancortelezzi/gatorserver/pkg/clog"
	"github.com/juancortelezzi/gatorserver/pkg/server"
)

const waitForReadyTimeout = 3 * time.Second

func testLookupEnv(key string) (string, bool) {
	if key == "PORT" {
		return "3000", true
	}

	return "", false
}

func getBaseUrl() string {
	port, found := testLookupEnv("PORT")
	if !found {
		panic("PORT environment variable not found")
	}

	return net.JoinHostPort("127.0.0.1", port)
}

func connectToServer(ctx context.Context, addr string) (net.Conn, error) {
	startTime := time.Now()
	for {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				if time.Since(startTime) > waitForReadyTimeout {
					return nil, fmt.Errorf("timeout reached while waiting for server: %e", err)
				}
				time.Sleep(time.Millisecond * 250)
				continue
			}
		}
		return conn, nil
	}
}

func TestMain(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	{
		logger := clog.NewLogger(os.Stdout, slog.LevelDebug)
		go server.Run(ctx, logger, testLookupEnv)
	}

	conn, err := connectToServer(context.Background(), getBaseUrl())
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	log.Println("writing to server")

	header := gatorparser.NewHeader(gatorparser.PacketTypeLocation)
	payload := gatorparser.PayloadLocation{
		Uuid:          uuid.New(),
		UnixTimestamp: time.Now().UTC().Unix(),
		Latitude:      100.22,
		Longitude:     44.65,
	}

	headerBytes, err := header.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal header: %v", err)
	}

	payloadBytes, err := payload.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	packetToSend := bytes.Join([][]byte{headerBytes, payloadBytes}, nil)
	_, err = conn.Write(packetToSend)
	if err != nil {
		t.Fatalf("Failed to write to server: %v", err)
	}

	buffer := make([]byte, 1024)
	pointer, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read from server: %v", err)
	}

	packetReceived := buffer[:pointer]
	if !bytes.Equal(packetReceived, packetToSend) {
		t.Fatalf("Received packet does not match sent packet")
	}
}
