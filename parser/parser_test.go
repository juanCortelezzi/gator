package gatorparser_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/juancortelezzi/gatorparser"
)

func TestUnix(t *testing.T) {
	storeTime := time.Now().Unix()
	parsedStoreTime := time.Unix(storeTime, 0)

	if storeTime != parsedStoreTime.Unix() {
		t.Fatalf("Unix timestamp mismatch: %d != %d", storeTime, parsedStoreTime.Unix())
	}
}

func TestPacketHeaderMarshaling(t *testing.T) {
	header := &gatorparser.Header{
		Version: 1,
		Type:    5,
		Length:  40,
	}

	data, err := header.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	var parsedHeader gatorparser.Header
	if err := parsedHeader.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}

	if header.Version != parsedHeader.Version {
		t.Fatalf("Version mismatch: %d != %d", header.Version, parsedHeader.Version)
	}

	if header.Type != parsedHeader.Type {
		t.Fatalf("Type mismatch: %d != %d", header.Type, parsedHeader.Type)
	}

	if header.Length != parsedHeader.Length {
		t.Fatalf("Length mismatch: %d != %d", header.Length, parsedHeader.Length)
	}
}

func TestLocatorPacketMarshaling(t *testing.T) {
	packet := gatorparser.PacketLocation{
		Uuid:          uuid.New(),
		UnixTimestamp: time.Now().Unix(),
		Latitude:      rand.Float64(),
		Longitude:     rand.Float64(),
	}

	data, err := packet.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	var parsedPacket gatorparser.PacketLocation
	if err := parsedPacket.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}

	if packet.Uuid != parsedPacket.Uuid {
		t.Fatalf("UUID mismatch: %s != %s", packet.Uuid, parsedPacket.Uuid)
	}

	if packet.UnixTimestamp != parsedPacket.UnixTimestamp {
		t.Fatalf("Timestamp mismatch: %d != %d", packet.UnixTimestamp, parsedPacket.UnixTimestamp)
	}

	if packet.Latitude != parsedPacket.Latitude {
		t.Fatalf("Latitude mismatch: %f != %f", packet.Latitude, parsedPacket.Latitude)
	}

	if packet.Longitude != parsedPacket.Longitude {
		t.Fatalf("Longitude mismatch: %f != %f", packet.Longitude, parsedPacket.Longitude)
	}
}
