# goais

[![cover.run](https://cover.run/go/github.com/bertoldvdb/goais.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2Fbertoldvdb%2Fgoais)
[![Go Report Card](https://goreportcard.com/badge/github.com/bertoldvdb/goais)](https://goreportcard.com/report/github.com/bertoldvdb/goais)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Summary**
This is a library written in [Go](https://golang.org/) to encode and decode Automatic Identification System packets. This is a VHF data system that is used for dependent tracking and identification of marine vessels. It is specified by the [ITU-R M.1371-5](https://www.itu.int/rec/R-REC-M.1371-5-201402-I/en) standard.

**Rationale**
There are a few other libraries offered on the internet, but all of them implement only a few messages (usually the position and static reports). In addition, very few libraries support encoding.

Many other libraries tightly integrate a NMEA (VDM/VDO) decoder with the AIS decoder. This is handy if you use a complete receiver that delivers NMEA sentences, but can be limiting. For example, you may have a low level RF receiver connected directly to an embedded system. I feel it is best to split NMEA and AIS coding as they are really two different functions.

**How to use it**
Start by getting an AIS packet from somewhere. It could for example come from [AISHub](http://www.aishub.net/) or a local receiver. Extract the payload and convert it to a byte slice containing one bit per byte. Then call the DecodePacket function on it. It will return an object containing the decoded message. For example:

    package main
    import (
        "fmt"
        "github.com/BertoldVdb/go-ais"
    )
    
    func  main() {
        msg  := []byte{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, ... , 0, 0, 0, 0, 0, 1, 1, 0, 0}
        parser  := ais.CodecNew(false, false)
        result  := parser.DecodePacket(msg)
        fmt.Printf("%T: %+v\n", result, result)
    }
    
The output of this program could be as follows:
> ais.PositionReport: {Valid:true MessageID:1 RepeatIndicator:0 UserID:235117222 NavigationalStatus:8 RateOfTurn:-128 Sog:44 PositionAccuracy:true Longitude:-0.249303 Latitude:53.7109 Cog:-54.4 TrueHeading:511 Timestamp:1 SpecialManoeuvreIndicator:0 Spare:0 Raim:true CommunicationState:59916}

To encode a packet, call the EncodePacket function. It works exactly in the opposite way of DecodePacket.

