package ais_test

import (
	"fmt"
	"log"

	"github.com/BertoldVdb/go-ais"
	"github.com/BertoldVdb/go-ais/aisnmea"
)

var sentences = []string{
	"!AIVDM,1,1,,B,10bb7q@P0lPGHlVMhbl0Qgw>2>`<,0*7F",
	"!AIVDM,1,1,,A,13u08p0000QDeLNO=PvHU3M>0>`<,0*00",
	"!AIVDM,1,1,,A,18Kk3f0vQ37K@m>0cnR9poo@0>`<,0*25",
	"!ABVDM,2,1,9,A,602=WITp2uLn01mVIj<04CH>NB0000PCEnUdK6UQKG=4HGIaI21:KnqQ,0*7E",
	"!ABVDM,2,2,9,A,M6QQKR1@JG9aI@0,2*69",
	"!AIVDM,1,1,,A,3815>kEwh00F2rfMvm<tCRk@0>`<,0*30",
	"!AIVDM,1,1,,B,E>ldCi?;Pb2a@22`:4@HrGK6P0044b3T6Jde@00003v01P,4*0F",
	"!AIVDM,2,1,6,B,53Jmvl82Bw3CTP7??K5<D5<lThF222222222221:I0oK;4e20L3FH42i,0*5C",
	"!AIVDM,2,2,6,B,p88888888888880,2*69",
	"!AIVDM,2,1,3,A,53Jmvl82Bw3CTP7??K5<D5<lThF222222222221:I0oK;4e20L3FH42ip888,0*12",
	"!AIVDM,1,1,,B,13`h7T001TwBtlvFII=>205>0>`<,0*6C",
	"!AIVDM,1,1,,B,14`V6d002kG;oFfKmQh=F:`60>`<,0*74",
	"!AIVDM,1,1,,B,16KMGt0000a>gnlD61a`KIa40>`<,0*6A",
	"!AIVDM,2,1,6,A,53aGCoD000010O?O7P1<lUB1L44hP5HDr3L0001?3H034vQh?2QDlkh00000,0*7F",
	"!AIVDM,1,1,,B,13aFe10P00PCwWJMh9U2vOw<2>`<,0*14",
	"!AIVDM,2,1,4,B,53aGFRT000010GO?KH0@F1H4hd0000000000001?98=63t@PJ0888888,0*19",
	"!AIVDM,2,2,4,B,888888888888880,2*23",
}

func Example() {
	nm := aisnmea.NMEACodecNew(ais.CodecNew(false, false))

	for _, sentence := range sentences {
		decoded, err := nm.ParseSentence(sentence)

		if err != nil {
			log.Fatal(err)
		}

		if decoded == nil {
			// packet not assembled yet
			continue
		}

		switch t := decoded.Packet.(type) {
		case ais.PositionReport:
			// do something with a position report is received
			fmt.Printf("received a position report which is of type: %T\n", t)
		case ais.ShipStaticData:
			// do something with a static data report is received
			fmt.Printf("received a static data report which is of type: %T\n", t)
		}
	}

	// Output:
	// received a position report which is of type: ais.PositionReport
	// received a position report which is of type: ais.PositionReport
	// received a position report which is of type: ais.PositionReport
	// received a position report which is of type: ais.PositionReport
	// received a static data report which is of type: ais.ShipStaticData
	// received a position report which is of type: ais.PositionReport
	// received a position report which is of type: ais.PositionReport
	// received a position report which is of type: ais.PositionReport
	// received a position report which is of type: ais.PositionReport
	// received a static data report which is of type: ais.ShipStaticData
}
