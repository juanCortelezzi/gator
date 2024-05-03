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

func TestHeaderMarshaling(t *testing.T) {
	originalHeader := &gatorparser.Header{
		Version: 1,
		Type:    5,
		Length:  40,
	}

	data, err := originalHeader.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	var parsedHeader gatorparser.Header
	if err := parsedHeader.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}

	if originalHeader.Version != parsedHeader.Version {
		t.Fatalf("Version mismatch: %d != %d", originalHeader.Version, parsedHeader.Version)
	}

	if originalHeader.Type != parsedHeader.Type {
		t.Fatalf("Type mismatch: %d != %d", originalHeader.Type, parsedHeader.Type)
	}

	if originalHeader.Length != parsedHeader.Length {
		t.Fatalf("Length mismatch: %d != %d", originalHeader.Length, parsedHeader.Length)
	}
}

func TestPayloadLocationMarshaling(t *testing.T) {
	originalPacket := gatorparser.PayloadLocation{
		Uuid:          uuid.New(),
		UnixTimestamp: time.Now().Unix(),
		Latitude:      rand.Float64(),
		Longitude:     rand.Float64(),
	}

	data, err := originalPacket.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	var parsedPacket gatorparser.PayloadLocation
	if err := parsedPacket.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}

	if originalPacket.Uuid != parsedPacket.Uuid {
		t.Fatalf("UUID mismatch: %s != %s", originalPacket.Uuid, parsedPacket.Uuid)
	}

	if originalPacket.UnixTimestamp != parsedPacket.UnixTimestamp {
		t.Fatalf("Timestamp mismatch: %d != %d", originalPacket.UnixTimestamp, parsedPacket.UnixTimestamp)
	}

	if originalPacket.Latitude != parsedPacket.Latitude {
		t.Fatalf("Latitude mismatch: %f != %f", originalPacket.Latitude, parsedPacket.Latitude)
	}

	if originalPacket.Longitude != parsedPacket.Longitude {
		t.Fatalf("Longitude mismatch: %f != %f", originalPacket.Longitude, parsedPacket.Longitude)
	}
}
