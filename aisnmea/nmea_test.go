package aisnmea

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/BertoldVdb/go-ais"
	"github.com/adrianmo/go-nmea"
)

func TestWrongType(t *testing.T) {
	nm := NMEACodecNew(ais.CodecNew(false, false, false))
	_, err := nm.ParseSentence("$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*6A")

	if err.Error() != SentenceNotVDMVDO {
		t.Error("Wrong sentence returned invalid error")
	}
}

func TestInvalid(t *testing.T) {
	nm := NMEACodecNew(ais.CodecNew(false, false, false))
	_, err := nm.ParseSentence("$Not a NMEA sentence*bb")

	if err == nil {
		t.Error("Invalid sentence did not return error")
	}
}

func TestTooManyFragments(t *testing.T) {
	nm := NMEACodecNew(ais.CodecNew(false, false, false))
	_, err := nm.ParseSentence("!AIVDM,30,1,,A,13u08p0000QDeLNO=PvHU3M>0>`<,0*32")

	if err != nil {
		t.Error("Error returned for valid message", err)
	}

	if nm.BufferedMessages() > 0 {
		t.Error("Invalid message was added to buffer")
	}
}

func TestFailedEncode(t *testing.T) {
	nm := NMEACodecNew(ais.CodecNew(false, false, false))

	p := VdmPacket{}

	if nm.EncodeSentence(p) != nil {
		t.Error("Output was produced although neither a valid packet or payload was provided")
	}

	p = VdmPacket{
		Packet: ais.PositionReport{Valid: false},
	}

	if nm.EncodeSentence(p) != nil {
		t.Error("Output was produced although the packet was not valid")
	}
}

func TestNMEAReencode(t *testing.T) {
	nm := NMEACodecNew(ais.CodecNew(false, false, false))
	nm2 := NMEACodecNew(ais.CodecNew(false, false, false))

	file, err := os.Open("testdata/aistest.nmea")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sentence := scanner.Text()
		indexSpace := strings.Index(sentence, " ")
		sentence = sentence[indexSpace+1:]

		result, err := nmea.Parse(sentence)
		if err != nil {
			t.Error("Decoding failed", err)
		}

		switch x := result.(type) {
		case nmea.VDMVDO:
			decoded, err := nm.ParseVDMVDO(&x)
			if err != nil {
				t.Error("Decoding failed", err)
				continue
			}

			if decoded != nil {
				encoded := nm.EncodeSentence(*decoded)

				for _, l := range encoded {
					decoded2, err := nm2.ParseSentence(l)
					if err != nil {
						t.Error("Could not decode sentence we just encoded")
						continue
					}

					if decoded2 != nil {
						if !bytes.Equal(decoded.Payload, decoded2.Payload) {
							t.Error("Payload not identical", decoded.Payload, decoded2.Payload)
						}
					}
				}
			}
		}

		if nm.BufferedMessages() > 5 {
			t.Error("Cleanup in assembler is not working", nm.BufferedMessages())
		}

	}

}

func TestNMEATagBlockDecodeSingleSentence(t *testing.T) {
	nm := NMEACodecNew(ais.CodecNew(false, false, false))
	msg, err := nm.ParseSentence("\\s:2156,c:1560234814*36\\!AIVDM,1,1,,B,23aDqDOP0S0:mk2Kv3Ip=wvpR>`<,0*3D")

	if err != nil {
		t.Error("Error returned for valid message", err)
	}

	if msg == nil {
		t.Error("No error, but no message for single-sentence message")
	}

	if msg.TagBlock.Source != "2156" {
		t.Error("TAG block Source not parsed")
	}

	if msg.TagBlock.Time != 1560234814 {
		t.Error("TAG block Time not parsed")
	}
}

func TestNMEATagBlockDecodeMultiSentence(t *testing.T) {
	nm := NMEACodecNew(ais.CodecNew(false, false, false))
	msg, err := nm.ParseSentence(
		"\\g:1-2-2449555,s:2251,c:1560234814*7E\\!AIVDM,2,1,7,A," +
			"8h3OwjQKP@5UUEPPP121IoCol54cd0Wws7wwjp:@`P1UUFD9e2B94oCPH54M`3kw,0*7A")

	if err != nil {
		t.Error("Error returned for valid message", err)
	}

	if msg != nil {
		t.Error("Premature return of message")
	}

	msg, err = nm.ParseSentence("\\g:2-2-2449555*63\\!AIVDM,2,2,7,A,sUwwjt;HvP1,2*4F")

	if err != nil {
		t.Error("Error returned for valid message", err)
	}

	if msg == nil {
		t.Error("No error, but no message for multi-sentence message")
	}

	if msg.TagBlock.Source != "2251" {
		t.Error("TAG block Source not parsed")
	}

	if msg.TagBlock.Time != 1560234814 {
		t.Error("TAG block Time not parsed")
	}

	if msg.TagBlock.Grouping != "" {
		t.Error("TAG block Grouping parsed (should be ignored)")
	}
}
