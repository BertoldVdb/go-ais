
# go-ais

[![Build Status](https://travis-ci.org/BertoldVdb/go-ais.svg?branch=master)](https://travis-ci.org/BertoldVdb/go-ais)
[![Coverage Status](https://coveralls.io/repos/github/BertoldVdb/go-ais/badge.svg?branch=master)](https://coveralls.io/github/BertoldVdb/go-ais?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bertoldvdb/goais)](https://goreportcard.com/report/github.com/bertoldvdb/goais)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Summary**

This is a library written in [Go](https://golang.org/) to encode and decode Automatic Identification System packets. This is a VHF data system that is used for dependent tracking and identification of marine vessels. It is specified by the [ITU-R M.1371-5](https://www.itu.int/rec/R-REC-M.1371-5-201402-I/en) standard.

**Rationale**

There are a few other libraries offered on the internet, but all of them implement only a few messages (usually the position and static reports). In addition, very few libraries support encoding.

Many other libraries tightly integrate a NMEA (VDM/VDO) decoder with the AIS decoder. This is handy if you use a complete receiver that delivers NMEA sentences, but can be limiting. For example, you may have a low level RF receiver connected directly to an embedded system. I feel it is best to split NMEA and AIS coding as they are really two different functions. However, a convenience library is provided in this repository that combines this library with a golang NMEA decoder and a self-developed NMEA encoder.

**How to use it**

Start by getting an AIS packet from somewhere. It could for example come from [AISHub](http://www.aishub.net/) or a local receiver. Extract the payload and convert it to a byte slice containing one bit per byte. Then call the DecodePacket function on it. It will return an object containing the decoded message. For example:
```go
package main
     
import (
    "fmt"
    "encoding/json"
 
    ais "github.com/BertoldVdb/go-ais"
)
    
func main() {
    msg := []byte{0, 0, ... 1, 0, 0, 0, 0, 0, 0, 0}
    
    c := ais.CodecNew(false, false, false)
    c.DropSpace = true
       
    result := c.DecodePacket(msg)
    out, _ := json.MarshalIndent(result, "", "  ")
    fmt.Printf("%T: %s\n", result, out)
}   
 ```
 
The output of this program could be as follows:
> ais.ShipStaticData: {  
> "MessageID": 5,  
> "RepeatIndicator": 0,  
> "UserID": 203999421,  
> "Valid": true,  
> "AisVersion": 0,  
> "ImoNumber": 0,  
> "CallSign": "OED3018",  
> "Name": "PRIMADONNA",  
> "Type": 69,  
> "Dimension": {  
> "A": 20,  
> "B": 93,  
> "C": 7,  
> "D": 10  
> },  
> "FixType": 1,  
> "Eta": {  
> "Month": 1,  
> "Day": 4,  
> "Hour": 18,  
> "Minute": 30  
> },  
> "MaximumStaticDraught": 1.7,  
> "Destination": "LINZ",  
> "Dte": false,  
> "Spare": false  
> }

To encode a packet, call the EncodePacket function. It works exactly in the opposite way of DecodePacket.

If you want to work with NMEA sentences you can use the following example to decode a packet:
```go
package main

import (
    "fmt"
    "github.com/BertoldVdb/go-ais"
    "github.com/BertoldVdb/go-ais/aisnmea"
)

func main() {
    nm := aisnmea.NMEACodecNew(ais.CodecNew(false, false, false))
    
    decoded, _ := nm.ParseSentence("!AIVDM,1,1,,B,33aEP2hP00PBLRFMfCp;OOw<R>`<,0*49")
    if decoded != nil {
        fmt.Printf("%+v\n", decoded.Packet)
    }
}
```

