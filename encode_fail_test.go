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

func tryEncodeWithMessageID(mID uint8) bool {
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
	if !tryEncodeWithMessageID(1) {
		t.Error("Could not encode position report with msgID==1")
	}

	if !tryEncodeWithMessageID(2) {
		t.Error("Could not encode position report with msgID==2")
	}

	if !tryEncodeWithMessageID(3) {
		t.Error("Could not encode position report with msgID==3")
	}

	if tryEncodeWithMessageID(4) {
		t.Error("Could encode position report with msgID==4")
	}

	if tryEncodeWithMessageID(0) {
		t.Error("Could encode position report with msgID==0")
	}

	if tryEncodeWithMessageID(28) {
		t.Error("Could encode position report with msgID==28")
	}
}

func tryEncodeTooLargeNumber(number uint16) bool {
	x := CodecNew(false, false)

	packet := ShipStaticData{
		Valid: true,
	}
	packet.Header = Header{
		MessageID: 5,
		UserID:    1337}
	packet.Dimension.A = number

	return x.EncodePacket(packet) != nil
}

func TestEncodeTooLargeNumber(t *testing.T) {
	if !tryEncodeTooLargeNumber(1) {
		t.Error("Could not encode ShipStaticData with DimensionA=1")
	}

	if tryEncodeTooLargeNumber(65535) {
		t.Error("Could encode ShipStaticData with DimensionA=65535")
	}
}

func TestEncodeFloatOutOfRange(t *testing.T) {
	staticData := ShipStaticData{
		Valid:                true,
		MaximumStaticDraught: 1000,
	}

	staticData.Header = Header{
		MessageID: 5,
		UserID:    1337,
	}

	x := CodecNew(false, false)
	if x.EncodePacket(staticData) != nil {
		t.Error("Encoded oversized float")
	}

}
