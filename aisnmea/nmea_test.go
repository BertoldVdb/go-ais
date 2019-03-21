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
	nm := NMEACodecNew(ais.CodecNew(false, false))
	err := nm.ParseSentence("$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*6A")

	if err.Error() != SentenceNotVDMVDO {
		t.Error("Wrong sentence returned invalid error")
	}
}

func TestInvalid(t *testing.T) {
	nm := NMEACodecNew(ais.CodecNew(false, false))
	err := nm.ParseSentence("$Not a NMEA sentence*bb")

	if err == nil {
		t.Error("Invalid sentence did not return error")
	}
}

func TestTooManyFragments(t *testing.T) {
	nm := NMEACodecNew(ais.CodecNew(false, false))
	err := nm.ParseSentence("!AIVDM,30,1,,A,13u08p0000QDeLNO=PvHU3M>0>`<,0*32")

	if err != nil {
		t.Error("Error returned for valid message", err)
	}

	if nm.BufferedMessages() > 0 {
		t.Error("Invalid message was added to buffer")
	}
}

func TestFailedEncode(t *testing.T) {
	nm := NMEACodecNew(ais.CodecNew(false, false))

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
	nm := NMEACodecNew(ais.CodecNew(false, false))
	nm2 := NMEACodecNew(ais.CodecNew(false, false))

	nm.DecodeCallback = func(decoded VdmPacket) {
		encoded := nm.EncodeSentence(decoded)

		nm2.DecodeCallback = func(decoded2 VdmPacket) {
			if !bytes.Equal(decoded.Payload, decoded2.Payload) {
				t.Error("Payload not identical")
			}
		}

		for _, l := range encoded {
			if nm2.ParseSentence(l) != nil {
				t.Error("Could not decode sentence we just encoded")
			}
		}
	}

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
			nm.ParseVDMVDO(&x)
		}

		if nm.BufferedMessages() > 5 {
			t.Error("Cleanup in assembler is not working", nm.BufferedMessages())
		}

	}

}
