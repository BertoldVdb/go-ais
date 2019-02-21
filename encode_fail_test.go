package ais

import "testing"

func tryEncodeTooLong(len int) bool {
	x := CodecNew(false, false)

	packet := SingleSlotBinaryMessage{
		Valid:              true,
		DestinationIDValid: false,
		ApplicationIDValid: false,
	}

	packet.Header = Header{
		MessageID: 25,
		UserID:    1337}

	packet.Payload = make([]byte, len)

	return x.EncodePacket(packet) != nil
}

func TestEncodeFailTooLong(t *testing.T) {
	if tryEncodeTooLong(129) {
		t.Error("Could encode overlength packet")
	}

	if !tryEncodeTooLong(128) {
		t.Error("Acceptable length packet rejected")
	}
}

func tryEncodeWithString(s string) bool {
	x := CodecNew(false, false)

	packet := SafetyBroadcastMessage{
		Valid: true,
		Text:  s,
	}

	packet.Header = Header{
		MessageID: 14,
		UserID:    1337}

	return x.EncodePacket(packet) != nil
}

func TestEncodeFailIllegalChar(t *testing.T) {
	if tryEncodeWithString("ILLeGAL") {
		t.Error("Could encode string with illegal char")
	}

	if !tryEncodeWithString("LEGAL") {
		t.Error("Could not encode valid string")
	}
}
func TestEncodeNumberTooLarge(t *testing.T) {
	_, ok := encodeNumber([]byte{}, false, 8, 256)
	if ok {
		t.Error("Could encode 256 in 8 bits")
	}

	_, ok = encodeNumber([]byte{}, false, 8, 255)
	if !ok {
		t.Error("Could not encode 255 in 8 bits")
	}

	_, ok = encodeNumber([]byte{}, true, 8, 128)
	if ok {
		t.Error("Could encode 128 in 8 bits signed")
	}

	_, ok = encodeNumber([]byte{}, true, 8, 127)
	if !ok {
		t.Error("Could not encode 127 in 8 bits signed")
	}

	_, ok = encodeNumber([]byte{}, true, 8, -128)
	if !ok {
		t.Error("Could not encode -128 in 8 bits signed")
	}

	_, ok = encodeNumber([]byte{}, true, 8, -129)
	if ok {
		t.Error("Could encode -129 in 8 bits signed")
	}

}

func TestEncodeFailNoUsefulData(t *testing.T) {
	x := CodecNew(false, false)

	packet := BinaryAcknowledge{
		Valid: true,
	}
	packet.Header = Header{
		MessageID: 7,
		UserID:    1337}

	if x.EncodePacket(packet) != nil {
		t.Error("Could encode packet that contained no data")
	}

	/* Add some data */
	packet.Destinations[0].Valid = true

	if x.EncodePacket(packet) == nil {
		t.Error("Could not encode packet although data was added")
	}
}

func tryEncodeWithMessageId(mID uint8) bool {
	x := CodecNew(false, false)

	packet := PositionReport{
		Valid: true,
	}

	packet.Header = Header{
		MessageID: mID,
		UserID:    1337}

	return x.EncodePacket(packet) != nil
}

func TestEncodeFailWrongMsgID(t *testing.T) {
	if !tryEncodeWithMessageId(1) {
		t.Error("Could not encode position report with msgID==1")
	}

	if !tryEncodeWithMessageId(2) {
		t.Error("Could not encode position report with msgID==2")
	}

	if !tryEncodeWithMessageId(3) {
		t.Error("Could not encode position report with msgID==3")
	}

	if tryEncodeWithMessageId(4) {
		t.Error("Could encode position report with msgID==4")
	}

	if tryEncodeWithMessageId(0) {
		t.Error("Could encode position report with msgID==0")
	}

	if tryEncodeWithMessageId(28) {
		t.Error("Could encode position report with msgID==28")
	}
}
