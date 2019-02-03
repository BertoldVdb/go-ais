package goais

import "testing"

func TestDecodeFailMsgIDTooShort(t *testing.T) {
	x := AISParserCreate(false, false)

	data := []byte{0, 1, 0, 1}

	if x.DecodePacket(data) != nil {
		t.Error("Could decode 4-bit packet")
	}
}

func TestDecodeStrictAlignment(t *testing.T) {
	data := []byte{0, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0}

	x := AISParserCreate(false, false)
	x.StrictByteAlignment = true

	if x.DecodePacket(data) == nil {
		t.Error("Failed to decode valid packet")
	}

	data = append(data, 0)

	if x.DecodePacket(data) != nil {
		t.Error("Could decode packet with invalid padding")
	}
}

func TestDecodeDependentBitNotReceived(t *testing.T) {
	data := []byte{0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	x := AISParserCreate(false, false)

	if x.DecodePacket(data) == nil {
		t.Error("Failed to decode valid packet")
	}

	data = data[:139]

	if x.DecodePacket(data) != nil {
		t.Error("Decoded packet although IsAddressed was not received")
	}
}
