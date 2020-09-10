package aisnmea

import (
	"errors"
	"fmt"

	"github.com/BertoldVdb/go-ais"

	nmea "github.com/adrianmo/go-nmea"
)

const (
	//SentenceNotVDMVDO will be returned if the provided sentence is of the wrong type
	SentenceNotVDMVDO string = "Sentence was not of the VDM or VDO type"
)

// NMEACodec is a convenience code that allows easy encoding and decoding of NMEA as produced by
// an AIS receiver
type NMEACodec struct {
	assembler      *vdmAssembler
	codec          *ais.Codec
	MaxLineLength  int
	seqNo          int
	AppendChecksum bool
}

// NMEACodecNew creates a NMEACodec. You need to provide a configured ais.Codec
func NMEACodecNew(codec *ais.Codec) *NMEACodec {
	a := &NMEACodec{
		assembler:      vdmAssemblerCreate(),
		codec:          codec,
		MaxLineLength:  82,
		AppendChecksum: true,
	}

	return a
}

func (nc *NMEACodec) handleAssembledMessage(assembled *VdmPacket) {
	channel := byte(1)
	c := assembled.Channel
	if c == '2' || c == 'b' || c == 'B' || c == '+' || c == 'H' || c == 'h' {
		channel = byte(2)
	}

	assembled.Packet = nc.codec.DecodePacket(assembled.Payload)
	assembled.Channel = channel
}

// BufferedMessages return the number of messages buffered in the reassembler
func (nc *NMEACodec) BufferedMessages() int {
	return nc.assembler.bufferedMessages()
}

// ParseVDMVDO parses a message contained in a nmea.VDMVDO struct
func (nc *NMEACodec) ParseVDMVDO(m *nmea.VDMVDO) (*VdmPacket, error) {
	assembled, ok := nc.assembler.process(m)
	if ok {
		nc.handleAssembledMessage(&assembled)
		return &assembled, nil
	}

	return nil, nil
}

// ParseSentence decodes a NMEA sentence containing an AIS message
func (nc *NMEACodec) ParseSentence(sentence string) (*VdmPacket, error) {
	s, err := nmea.Parse(sentence)
	if err != nil {
		return nil, err
	}

	switch m := s.(type) {
	case nmea.VDMVDO:
		return nc.ParseVDMVDO(&m)
	}

	return nil, errors.New(SentenceNotVDMVDO)
}

func valueToChar(value byte) byte {
	result := value + 48

	if result >= 88 {
		result += 8
	}

	return result
}

func addChecksum(sentence string) string {
	checksum := byte(0)
	for i := 1; i < len(sentence); i++ {
		checksum ^= sentence[i]
	}

	return fmt.Sprintf("%s*%02X", sentence, checksum)
}

// EncodeSentence encodes the provided packet into zero or more NMEA sentences
func (nc *NMEACodec) EncodeSentence(p VdmPacket) []string {
	if p.Payload == nil && p.Packet != nil {
		p.Payload = nc.codec.EncodePacket(p.Packet)
	}

	if p.Payload == nil {
		return nil
	}

	asciiPayload := make([]byte, 0, len(p.Payload)/6*8+8)

	value := byte(0)
	bitsUsed := 0
	for i := 0; i < len(p.Payload); i++ {
		if p.Payload[i] > 1 {
			return nil
		}

		value <<= 1
		value += p.Payload[i]
		bitsUsed++
		if bitsUsed >= 6 {
			asciiPayload = append(asciiPayload, valueToChar(value))

			bitsUsed = 0
			value = 0
		}
	}

	fillBits := 0
	if bitsUsed != 0 {
		fillBits = 6 - bitsUsed
		value <<= uint(fillBits)

		asciiPayload = append(asciiPayload, valueToChar(value))
	}

	channel := 'A'
	if p.Channel == 2 {
		channel = 'B'
	}

	maxDataLength := nc.MaxLineLength - 3 - 2 - 1 - 16 //Subtract: *CRC, \r\n, !, AIVDM,x,y,z,c,,f
	if maxDataLength < 1 {
		maxDataLength = 1
	}

	var output []string

	/* Does the packet fit in one sentence? */
	if nc.MaxLineLength <= 0 || len(asciiPayload) <= maxDataLength {
		output = make([]string, 1)
		output[0] = fmt.Sprintf("!%s%s,1,1,,%c,%s,%d", p.TalkerID, p.MessageType, channel, asciiPayload, fillBits)
		output[0] = addChecksum(output[0])

		tagBlock := encodeTagBlock(&p.TagBlock, 1,1,0, false)
		if tagBlock != "" {
			output[0] = fmt.Sprintf("%s%s", tagBlock, output[0])
		}
	} else {
		dataIndex := 0
		msgIndex := 1

		msgNum := len(asciiPayload)/maxDataLength + 1

		for dataIndex < len(asciiPayload) {
			bytesUsed := maxDataLength

			if dataIndex+bytesUsed > len(asciiPayload) {
				bytesUsed = len(asciiPayload) - dataIndex
			}

			sub := asciiPayload[dataIndex : dataIndex+bytesUsed]
			dataIndex += bytesUsed

			suffix := 0
			if dataIndex >= len(asciiPayload) {
				suffix = fillBits
			}

			sentence := fmt.Sprintf("!%s%s,%d,%d,%d,%c,%s,%d", p.TalkerID, p.MessageType, msgNum, msgIndex, nc.seqNo, channel, sub, suffix)

			sentence = addChecksum(sentence)

			tagBlock := encodeTagBlock(&p.TagBlock, msgIndex, msgNum, nc.seqNo, true)
			if tagBlock != "" {
				sentence = fmt.Sprintf("%s%s", tagBlock, sentence)
			}

			output = append(output, sentence)

			msgIndex++
		}

		nc.seqNo++
		if nc.seqNo == 10 {
			nc.seqNo = 0
		}
	}

	return output
}
