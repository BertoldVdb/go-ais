package ais

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
)

func tryFile(t *testing.T, x *Codec, msgID int) {
	f, err := os.Open(fmt.Sprintf("testmsg/%d.msg", msgID))
	if err != nil {
		t.Error("Failed to open file", msgID)
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for index := 0; true; index++ {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}

		line = line[:len(line)-2]
		source := []byte(line)

		/* Convert ascii '0' and '1' to real 0 and 1 */
		for i := 0; i < len(source); i++ {
			source[i] -= '0'
		}

		/* Decode the packet */
		decoded := x.DecodePacket(source)
		if decoded == nil {
			/* Failed to decode... */
			t.Error("Could not decode", msgID, index)
			continue
		}

		encoded := x.EncodePacket(decoded)

		/* Check if the bitstream is identical */
		if len(encoded) < len(source) {
			/* Output is too short */
			t.Error("Output too short", msgID, index)
			continue
		}

		if !bytes.Equal(encoded[:len(source)], source) {
			t.Error("Bitstream does not match", msgID, index)
		}
	}
}

func TestReEncode(t *testing.T) {
	x := CodecNew(false, false)

	/* Convenience conversion disabled to avoid float inaccuracies */
	x.FloatWithoutConversion = true

	/* Test all message types */
	for i := 1; i <= 27; i++ {
		log.Println("Working on message type", i)
		tryFile(t, x, i)
	}
}

func checkFloat(a float64, b float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}

	return diff < 0.001
}

func TestFloatReEncoder(t *testing.T) {
	x := CodecNew(true, true)

	packet := PositionReport{
		Valid:     true,
		Longitude: 12.345,
		Latitude:  1.2345,
		Cog:       123.2,
	}
	packet.Header = Header{
		MessageID: 2,
		UserID:    1337}

	encoded := x.EncodePacket(packet)
	if encoded == nil {
		t.Error("Failed to encode position report")
		return
	}

	switch newPacket := x.DecodePacket(encoded).(type) {
	case PositionReport:
		if !checkFloat(float64(packet.Latitude), float64(newPacket.Latitude)) ||
			!checkFloat(float64(packet.Longitude), float64(newPacket.Longitude)) ||
			!checkFloat(float64(packet.Cog), float64(newPacket.Cog)) {
			t.Error("Packet has a floating point error")
		}
	default:
		t.Error("Packet was returned as a different type")
	}
}

func testInterfacecInternal(p Packet, t *testing.T) (bool, int, uint32) {
	x := CodecNew(false, false)

	encoded := x.EncodePacket(p)
	if encoded == nil {
		t.Error("Failed to encode packet")
	}

	decoded := x.DecodePacket(encoded)
	if decoded.GetHeader().MessageID != p.GetHeader().MessageID {
		t.Error("GetHeader failed")
	}

	switch newPacket := decoded.(type) {
	case HasCommunicationState:
		return true, newPacket.IsItdma(), newPacket.GetState()
	default:
		return false, 0, 0
	}
}

func TestInterfaceAccess(t *testing.T) {

	packet := MultiSlotBinaryMessage{
		Valid: true,
	}
	packet.Header = Header{MessageID: 26}
	packet.CommunicationStateIsItdma = true
	packet.CommunicationState = 0x1234

	a, b, c := testInterfacecInternal(packet, t)
	if !a {
		t.Error("Comm state not found")
	}
	if b != 1 || c != 0x1234 {
		t.Error("Comm state decoding error", b, c)
	}

	packet.CommunicationStateIsItdma = false
	packet.CommunicationState = 0xCAFE

	a, b, c = testInterfacecInternal(packet, t)
	if !a {
		t.Error("Comm state not found")
	}
	if b != 0 || c != 0xCAFE {
		t.Error("Comm state decoding error", b, c)
	}

	packet2 := PositionReport{
		Valid: true,
	}
	packet2.Header = Header{MessageID: 1}
	packet2.CommunicationState = 0x4321
	a, b, c = testInterfacecInternal(packet2, t)
	if !a {
		t.Error("Comm state not found")
	}
	if b != -1 || c != 0x4321 {
		t.Error("Comm state decoding error", b, c)
	}

}
