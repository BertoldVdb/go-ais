package ais

import "testing"

func TestDecodeInvalidType24PartNumber(t *testing.T) {
	// The allowed range of type 24 PartNumber is 0-1, but the field
	// occupies two bits. The extra bit is reserved, and must be zero.

	// Below is a type 24B message from this real-world NMEA:
	// !AIVDM,1,1,,A,HTRBMh>T<wwFj443@24?pJKt00p0,0*7a
	// The reserved bit of PartNumber is set, and message used to be rejected.
	data := []byte{
		0, 1, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1,
		0, 0, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0,
		0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 1, 1,
		1, 1, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	x := CodecNew(false, false)
	msg := x.DecodePacket(data)
	if msg == nil {
		t.Fatal("Failed to decode valid packet")
	}
	report, ok := msg.(StaticDataReport)
	if !ok {
		t.Fatal("Unexpected dynamic type")
	}

	// Check that reserved and part number are properly decoded.
	if report.Reserved != 1 {
		t.Error("Unexpected Reserved:", report.Reserved)
	}
	if !report.PartNumber {
		t.Error("Unexpected PartNumber:", report.PartNumber)
	}
}
