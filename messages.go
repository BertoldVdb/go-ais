//go:generate go run ./parser_generator --input=messages.go
package ais

import (
	"reflect"
)

type processFunc func(packet Packet) Packet

type msgMapType struct {
	rType reflect.Type
}

var msgMap [28]msgMapType

// Packet is an interface describing coded and decoded ais packets
type Packet interface {
	GetHeader() *Header
}

// Header contains the header prepended to each packet
type Header struct {
	MessageID       uint8  `aisWidth:"6"`
	RepeatIndicator uint8  `aisWidth:"2"`
	UserID          uint32 `aisWidth:"30"`
}

// GetHeader returns the header of the packet
func (h Header) GetHeader() *Header {
	return &h
}

// PositionReport should be output periodically by mobile stations. The message ID is 1, 2 or 3
// depending on the system mode.
type PositionReport struct {
	Header                    `aisWidth:"38"`
	Valid                     bool            `aisEncodeMaxLen:"168"`
	NavigationalStatus        uint8           `aisWidth:"4"`
	RateOfTurn                int16           `aisWidth:"8"`
	Sog                       Field10         `aisWidth:"10"`
	PositionAccuracy          bool            `aisWidth:"1"`
	Longitude                 FieldLatLonFine `aisWidth:"28"`
	Latitude                  FieldLatLonFine `aisWidth:"27"`
	Cog                       Field10         `aisWidth:"12"`
	TrueHeading               uint16          `aisWidth:"9"`
	Timestamp                 uint8           `aisWidth:"6"`
	SpecialManoeuvreIndicator uint8           `aisWidth:"2"`
	Spare                     uint8           `aisWidth:"3" aisEncodeAs:"0"`
	Raim                      bool            `aisWidth:"1"`
	CommunicationStateNoItdma `aisWidth:"19"`
}

// BaseStationReport should be used for reporting UTC time and date and, at the same time, position.
// A base station should use Message 4 in its periodical transmissions. Message 4 is used by AIS stations
// for determining if it is within 120 NM for response to Messages 20 and 23. A mobile station should output
// Message 11 only in response to interrogation by Message 10. Message 11 is only transmitted as a result
// of a UTC request message (Message 10). The UTC and date response should be transmitted on the channel,
// where the UTC request message was received. */
type BaseStationReport struct {
	Header                    `aisWidth:"38"`
	Valid                     bool            `aisEncodeMaxLen:"168"`
	UtcYear                   uint16          `aisWidth:"14"`
	UtcMonth                  uint8           `aisWidth:"4"`
	UtcDay                    uint8           `aisWidth:"5"`
	UtcHour                   uint8           `aisWidth:"5"`
	UtcMinute                 uint8           `aisWidth:"6"`
	UtcSecond                 uint8           `aisWidth:"6"`
	PositionAccuracy          bool            `aisWidth:"1"`
	Longitude                 FieldLatLonFine `aisWidth:"28"`
	Latitude                  FieldLatLonFine `aisWidth:"27"`
	FixType                   uint8           `aisWidth:"4"`
	LongRangeEnable           bool            `aisWidth:"1"`
	Spare                     uint16          `aisWidth:"9" aisEncodeAs:"0"`
	Raim                      bool            `aisWidth:"1"`
	CommunicationStateNoItdma `aisWidth:"19"`
}

// ShipStaticData should only be used by Class A shipborne and SAR aircraft AIS stations when reporting static
// or voyage related data.
type ShipStaticData struct {
	Header               `aisWidth:"38"`
	Valid                bool           `aisEncodeMaxLen:"424"`
	AisVersion           uint8          `aisWidth:"2"`
	ImoNumber            uint32         `aisWidth:"30"`
	CallSign             string         `aisWidth:"42"`
	Name                 string         `aisWidth:"120"`
	Type                 uint8          `aisWidth:"8"`
	Dimension            FieldDimension `aisWidth:"30"`
	FixType              uint8          `aisWidth:"4"`
	Eta                  FieldETA       `aisWidth:"20"`
	MaximumStaticDraught Field10        `aisWidth:"8"`
	Destination          string         `aisWidth:"120"`
	Dte                  bool           `aisWidth:"1"`
	Spare                bool           `aisWidth:"1" aisEncodeAs:"0"`
}

// AddressedBinaryMessage should be variable in length, based on the amount of binary data.
// The length should vary between 1 and 5 slots. See application identifiers in § 2.1, Annex 5.
type AddressedBinaryMessage struct {
	Header         `aisWidth:"38"`
	Valid          bool                       `aisEncodeMaxLen:"1008"`
	SequenceNumber uint8                      `aisWidth:"2"`
	DestinationID  uint32                     `aisWidth:"30"`
	Retransmission bool                       `aisWidth:"1"`
	Spare          bool                       `aisWidth:"1" aisEncodeAs:"0"`
	ApplicationID  FieldApplicationIdentifier `aisWidth:"16"`
	BinaryData     []byte                     `aisWidth:"-1"`
}

// BinaryBroadcastMessage will be variable in length, based on the amount of binary data.
// The length should vary between 1 and 5 slots.
type BinaryBroadcastMessage struct {
	Header        `aisWidth:"38"`
	Valid         bool                       `aisEncodeMaxLen:"1008"`
	Spare         uint8                      `aisWidth:"2" aisEncodeAs:"0"`
	ApplicationID FieldApplicationIdentifier `aisWidth:"16"`
	BinaryData    []byte                     `aisWidth:"-1"`
}

// StandardSearchAndRescueAircraftReport should be used as a standard position report for
// aircraft involved in SAR operations. Stations other than aircraft involved in SAR operations
// should not transmit this message. The default reporting interval for this message should be 10 s.
type StandardSearchAndRescueAircraftReport struct {
	Header                  `aisWidth:"38"`
	Valid                   bool            `aisEncodeMaxLen:"168"`
	Altitude                uint16          `aisWidth:"12"`
	Sog                     uint16          `aisWidth:"10"`
	PositionAccuracy        bool            `aisWidth:"1"`
	Longitude               FieldLatLonFine `aisWidth:"28"`
	Latitude                FieldLatLonFine `aisWidth:"27"`
	Cog                     Field10         `aisWidth:"12"`
	Timestamp               uint8           `aisWidth:"6"`
	AltFromBaro             bool            `aisWidth:"1"`
	Spare1                  uint8           `aisWidth:"7" aisEncodeAs:"0"`
	Dte                     bool            `aisWidth:"1"`
	Spare2                  uint8           `aisWidth:"3" aisEncodeAs:"0"`
	AssignedMode            bool            `aisWidth:"1"`
	Raim                    bool            `aisWidth:"1"`
	CommunicationStateItdma `aisWidth:"20"`
}

// CoordinatedUTCInquiry should be used when a station is requesting UTC and date from another
// station.
type CoordinatedUTCInquiry struct {
	Header        `aisWidth:"38"`
	Valid         bool   `aisEncodeMaxLen:"72"`
	Spare1        uint8  `aisWidth:"2" aisEncodeAs:"0"`
	DestinationID uint32 `aisWidth:"30"`
	Spare2        uint8  `aisWidth:"2" aisEncodeAs:"0"`
}

// AddessedSafetyMessage could be variable in length, based on the amount of safety related text.
// The length should vary between 1 and 5 slots.
type AddessedSafetyMessage struct {
	Header         `aisWidth:"38"`
	Valid          bool   `aisEncodeMaxLen:"1008"`
	SequenceNumber uint8  `aisWidth:"2"`
	DestinationID  uint32 `aisWidth:"30"`
	Retransmission bool   `aisWidth:"1"`
	Spare          bool   `aisWidth:"1" aisEncodeAs:"0"`
	Text           string `aisWidth:"-1"`
}

// SafetyBroadcastMessage could be variable in length, based on the amount of safety related
// text. The length should vary between 1 and 5 slots.
type SafetyBroadcastMessage struct {
	Header `aisWidth:"38"`
	Valid  bool   `aisEncodeMaxLen:"1008"`
	Spare  uint8  `aisWidth:"2" aisEncodeAs:"0"`
	Text   string `aisWidth:"-1"`
}

// GnssBroadcastBinaryMessage should be transmitted by a base station, which is connected to a DGNSS
// reference source, and configured to provide DGNSS data to receiving stations. The contents of the
// data should be in accordance with Recommendation ITU-R M.823, excluding preamble and parity formatting.
type GnssBroadcastBinaryMessage struct {
	Header    `aisWidth:"38"`
	Valid     bool              `aisEncodeMaxLen:"816"`
	Spare1    uint8             `aisWidth:"2" aisEncodeAs:"0"`
	Longitude FieldLatLonCoarse `aisWidth:"18"`
	Latitude  FieldLatLonCoarse `aisWidth:"17"`
	Spare2    uint8             `aisWidth:"5" aisEncodeAs:"0"`
	Data      []byte            `aisWidth:"-1"`
}

// StandardClassBPositionReport should be output periodically and autonomously instead of
// Messages 1, 2, or 3 by Class B shipborne mobile equipment, only. The reporting interval should default to
// the values given in Table 2, Annex 1, unless otherwise specified by reception of a Message 16 or 23; and
// depending on the current SOG and navigational status flag setting.
type StandardClassBPositionReport struct {
	Header                  `aisWidth:"38"`
	Valid                   bool            `aisEncodeMaxLen:"168"`
	Spare1                  uint8           `aisWidth:"8" aisEncodeAs:"0"`
	Sog                     Field10         `aisWidth:"10"`
	PositionAccuracy        bool            `aisWidth:"1"`
	Longitude               FieldLatLonFine `aisWidth:"28"`
	Latitude                FieldLatLonFine `aisWidth:"27"`
	Cog                     Field10         `aisWidth:"12"`
	TrueHeading             uint16          `aisWidth:"9"`
	Timestamp               uint8           `aisWidth:"6"`
	Spare2                  uint8           `aisWidth:"2" aisEncodeAs:"0"`
	ClassBUnit              bool            `aisWidth:"1"`
	ClassBDisplay           bool            `aisWidth:"1"`
	ClassBDsc               bool            `aisWidth:"1"`
	ClassBBand              bool            `aisWidth:"1"`
	ClassBMsg22             bool            `aisWidth:"1"`
	AssignedMode            bool            `aisWidth:"1"`
	Raim                    bool            `aisWidth:"1"`
	CommunicationStateItdma `aisWidth:"20"`
}

// ExtendedClassBPositionReport should be transmitted once every 6 min in two slots allocated by the
// use of Message 18 in the ITDMA communication state. This message should be transmitted immediately
// after the following parameter values change: dimension of ship/reference for position or type of
// electronic position fixing device.
//
//	For future equipment: this message is not needed and should not be used. All content is covered by
//	Message 18, Message 24A and 24B.
//	For legacy equipment: this message should be used by Class B shipborne mobile equipment.
type ExtendedClassBPositionReport struct {
	Header           `aisWidth:"38"`
	Valid            bool            `aisEncodeMaxLen:"312"`
	Spare1           uint8           `aisWidth:"8" aisEncodeAs:"0"`
	Sog              Field10         `aisWidth:"10"`
	PositionAccuracy bool            `aisWidth:"1"`
	Longitude        FieldLatLonFine `aisWidth:"28"`
	Latitude         FieldLatLonFine `aisWidth:"27"`
	Cog              Field10         `aisWidth:"12"`
	TrueHeading      uint16          `aisWidth:"9"`
	Timestamp        uint8           `aisWidth:"6"`
	Spare2           uint8           `aisWidth:"4" aisEncodeAs:"0"`
	Name             string          `aisWidth:"120"`
	Type             uint8           `aisWidth:"8"`
	Dimension        FieldDimension  `aisWidth:"30"`
	FixType          uint8           `aisWidth:"4"`
	Raim             bool            `aisWidth:"1"`
	Dte              bool            `aisWidth:"1"`
	AssignedMode     bool            `aisWidth:"1"`
	Spare3           uint8           `aisWidth:"4" aisEncodeAs:"0"`
}

// AidsToNavigationReport should be used by an Aids to navigation (AtoN) AIS station. This station
// may be mounted on an aid-to-navigation or this message may be transmitted by a fixed station when
// the functionality of an AtoN station is integrated into the fixed station. This message should be
// transmitted autonomously at a Rr of once every three (3) min or it may be assigned by an assigned
// mode command (Message 16) via the VHF data link, or by an external command. This message
// should not occupy more than two slots.
type AidsToNavigationReport struct {
	Header           `aisWidth:"38"`
	Valid            bool            `aisEncodeMaxLen:"356"`
	Type             uint8           `aisWidth:"5"`
	Name             string          `aisWidth:"120"`
	PositionAccuracy bool            `aisWidth:"1"`
	Longitude        FieldLatLonFine `aisWidth:"28"`
	Latitude         FieldLatLonFine `aisWidth:"27"`
	Dimension        FieldDimension  `aisWidth:"30"`
	Fixtype          uint8           `aisWidth:"4"`
	Timestamp        uint8           `aisWidth:"6"`
	OffPosition      bool            `aisWidth:"1"`
	AtoN             uint8           `aisWidth:"8"`
	Raim             bool            `aisWidth:"1"`
	VirtualAtoN      bool            `aisWidth:"1"`
	AssignedMode     bool            `aisWidth:"1"`
	Spare            bool            `aisWidth:"1" aisEncodeAs:"0"`
	NameExtension    string          `aisWidth:"-1"`
}

// GroupAssignmentCommand is transmitted by a base station when operating as a controlling
// entity(see § 4.3.3.3.2 Annex 7 and § 3.20). This message should be applied to a mobile station within
// the defined region and as selected by “Ship and Cargo Type” or “Station type”. The receiving station
// should consider all selector fields concurrently. It controls the following operating parameters of a
// mobile station:
// * transmit/receive mode;
// * reporting interval;
// * the duration of a quiet time.
// Station type 10 should be used to define the base station coverage area for control of Message 27
// transmissions by Class A and Class B “SO” mobile stations. When station type is 10 only the fields
// latitude, longitude are used, all other fields should be ignored. This information will be relevant until
// three minutes after the last reception of controlling Message 4 from the same base station (same MMSI).
type GroupAssignmentCommand struct {
	Header            `aisWidth:"38"`
	Valid             bool              `aisEncodeMaxLen:"160"`
	Spare1            uint8             `aisWidth:"2" aisEncodeAs:"0"`
	Longitude1        FieldLatLonCoarse `aisWidth:"18"`
	Latitude1         FieldLatLonCoarse `aisWidth:"17"`
	Longitude2        FieldLatLonCoarse `aisWidth:"18"`
	Latitude2         FieldLatLonCoarse `aisWidth:"17"`
	StationType       uint8             `aisWidth:"4"`
	ShipType          uint8             `aisWidth:"8"`
	Spare2            uint32            `aisWidth:"22" aisEncodeAs:"0"`
	TxRxMode          uint8             `aisWidth:"2"`
	ReportingInterval uint8             `aisWidth:"4"`
	QuietTime         uint8             `aisWidth:"4"`
	Spare3            uint8             `aisWidth:"6" aisEncodeAs:"0"`
}

// StaticDataReportA is the A part of message 24
type StaticDataReportA struct {
	Valid bool
	Name  string `aisWidth:"120"`
}

// StaticDataReportB is the B part of message 24
type StaticDataReportB struct {
	Valid          bool
	ShipType       uint8          `aisWidth:"8"`
	VendorIDName   string         `aisWidth:"18"`
	VenderIDModel  uint8          `aisWidth:"4"`
	VenderIDSerial uint32         `aisWidth:"20"`
	CallSign       string         `aisWidth:"42"`
	Dimension      FieldDimension `aisWidth:"30"`
	FixType        uint8          `aisWidth:"4"`
	Spare          uint8          `aisWidth:"2" aisEncodeAs:"0"`
}

// StaticDataReport part A shall transmit once every 6 min alternating between
// channels.
// Message 24 Part A may be used by any AIS station to associate a MMSI with a name.
// Message 24 Part A and Part B should be transmitted once every 6 min by Class B “CS” and Class B
// “SO” shipborne mobile equipment. The message consists of two parts. Message 24B should be
// transmitted within 1 min following Message 24A.
type StaticDataReport struct {
	Header     `aisWidth:"38"`
	Valid      bool              `aisEncodeMaxLen:"168"`
	Reserved   uint8             `aisWidth:"1" aisEncodeAs:"0" aisCheckValue:"0"`
	PartNumber bool              `aisWidth:"1"`
	ReportA    StaticDataReportA `aisWidth:"120" aisDependsBit:"~39" aisDependsField:"~PartNumber"`
	ReportB    StaticDataReportB `aisWidth:"120" aisDependsBit:"39" aisDependsField:"PartNumber"`
}

// LongRangeAisBroadcastMessage is primarily intended for long-range detection of AIS Class A
// and Class B “SO” equipped vessels (typically by satellite). This message has a similar content
// to Messages 1, 2 and 3, but the total number of bits has been compressed to allow for increased
// propagation delays associated with long-range detection. Refer to Annex 4 for details on
// Long-Range applications.
type LongRangeAisBroadcastMessage struct {
	Header             `aisWidth:"38"`
	Valid              bool              `aisEncodeMaxLen:"96"`
	PositionAccuracy   bool              `aisWidth:"1"`
	Raim               bool              `aisWidth:"1"`
	NavigationalStatus uint8             `aisWidth:"4"`
	Longitude          FieldLatLonCoarse `aisWidth:"18"`
	Latitude           FieldLatLonCoarse `aisWidth:"17"`
	Sog                uint8             `aisWidth:"6"`
	Cog                uint16            `aisWidth:"9"`
	PositionLatency    bool              `aisWidth:"1"`
	Spare              bool              `aisWidth:"1" aisEncodeAs:"0"`
}

// BinaryAcknowledgeData is the data part of BinaryAcknowledge
type BinaryAcknowledgeData struct {
	Valid          bool
	DestinationID  uint32 `aisWidth:"30"`
	SequenceNumber uint8  `aisWidth:"2"`
}

// BinaryAcknowledge should be used as an acknowledgement of up to four Message 6 messages received
// (see § 5.3.1, Annex 2) and should be transmitted on the channel, where the addressed message to be
// acknowledged was received.
type BinaryAcknowledge struct {
	Header       `aisWidth:"38"`
	Valid        bool                     `aisEncodeMaxLen:"168"`
	Spare        uint8                    `aisWidth:"2" aisEncodeAs:"0"`
	Destinations [4]BinaryAcknowledgeData `aisWidth:"0"`
}

// InterrogationStation1Message1 is the station 1 part of Interrogation
type InterrogationStation1Message1 struct {
	Valid      bool
	StationID  uint32 `aisWidth:"30"`
	MessageID  uint8  `aisWidth:"6"`
	SlotOffset uint16 `aisWidth:"12"`
}

// InterrogationStation1Message2 is the second station 1 part of interrogation
type InterrogationStation1Message2 struct {
	Valid      bool   `aisOptional:"1"`
	Spare      uint8  `aisWidth:"2" aisEncodeAs:"0"`
	MessageID  uint8  `aisWidth:"6"`
	SlotOffset uint16 `aisWidth:"12"`
}

// InterrogationStation2 is the station 2 part of Interrogation
type InterrogationStation2 struct {
	Valid      bool   `aisOptional:"1"`
	Spare1     uint8  `aisWidth:"2" aisEncodeAs:"0"`
	StationID  uint32 `aisWidth:"30"`
	MessageID  uint8  `aisWidth:"6"`
	SlotOffset uint16 `aisWidth:"12"`
	Spare2     uint8  `aisWidth:"2" aisEncodeAs:"0"`
}

// Interrogation should be used for interrogations via the TDMA (not DSC) VHF data link except for
// requests for UTC and date. The response should be transmitted on the channel where the interrogation
// was received.
type Interrogation struct {
	Header       `aisWidth:"38"`
	Valid        bool                          `aisEncodeMaxLen:"160"`
	Spare        uint8                         `aisWidth:"2" aisEncodeAs:"0"`
	Station1Msg1 InterrogationStation1Message1 `aisWidth:"48"`
	Station1Msg2 InterrogationStation1Message2 `aisWidth:"0"`
	Station2     InterrogationStation2         `aisWidth:"0"`
}

// AssignedModeCommandData is the data part of AssignedModeCommand
type AssignedModeCommandData struct {
	Valid         bool
	DestinationID uint32 `aisWidth:"30"`
	Offset        uint16 `aisWidth:"12"`
	Increment     uint16 `aisWidth:"10"`
}

// AssignedModeCommand be transmitted by a base station when operating as a controlling entity. Other
// stations can be assigned a transmission schedule, other than the currently used one. If a station is
// assigned a schedule, it will also enter assigned mode.
type AssignedModeCommand struct {
	Header   `aisWidth:"38"`
	Valid    bool                       `aisEncodeMaxLen:"144"`
	Spare    uint8                      `aisWidth:"2" aisEncodeAs:"0"`
	Commands [2]AssignedModeCommandData `aisWidth:"0"`
}

// DataLinkManagementMessageData is the data part of DataLinkManagementMessage
type DataLinkManagementMessageData struct {
	Valid         bool
	Offset        uint16 `aisWidth:"12"`
	NumberOfSlots uint8  `aisWidth:"4"`
	TimeOut       uint8  `aisWidth:"3"`
	Increment     uint16 `aisWidth:"11"`
}

// DataLinkManagementMessage should be used by base station(s) to pre-announce the fixed allocation
// schedule (FATDMA) for one or more base station(s) and it should be repeated as often as required. This
// way the system can provide a high level of integrity for base station(s). This is especially important in
// regions where several base stations are located adjacent to each other and mobile station(s) move
// between these different regions. These reserved slots cannot be autonomously allocated by mobile
// stations.
type DataLinkManagementMessage struct {
	Header `aisWidth:"38"`
	Valid  bool                             `aisEncodeMaxLen:"160"`
	Spare  uint8                            `aisWidth:"2" aisEncodeAs:"0"`
	Data   [4]DataLinkManagementMessageData `aisWidth:"0"`
}

// ChannelManagementBroadcastData contains the boundrary for a broadcasted channel mangement packet.
type ChannelManagementBroadcastData struct {
	Longitude1 FieldLatLonCoarse `aisWidth:"18"`
	Latitude1  FieldLatLonCoarse `aisWidth:"17"`
	Longitude2 FieldLatLonCoarse `aisWidth:"18"`
	Latitude2  FieldLatLonCoarse `aisWidth:"17"`
}

// ChannelManagementUnicastData contains the destination addresses for a unicast channel mangement packet.
type ChannelManagementUnicastData struct {
	AddressStation1 uint32 `aisWidth:"30"`
	Spare2          uint8  `aisWidth:"5" aisEncodeAs:"0"`
	AddressStation2 uint32 `aisWidth:"30"`
	Spare3          uint8  `aisWidth:"5" aisEncodeAs:"0"`
}

// ChannelManagement should be transmitted by a base station (as a broadcast message) to command the VHF
// data link parameters for the geographical area designated in this message and should be accompanied
// by a Message 4 transmission for evaluation of the message within 120 NM. The geographical area
// designated by this message should be as defined in § 4.1, Annex 2. Alternatively, this message may
// be used by a base station (as an addressed message) to command individual AIS mobile stations to
// adopt the specified VHF data link parameters. When interrogated and no channel management
// performed by the interrogated base station, the not available and/or international default settings
// should be transmitted (see § 4.1, Annex 2).
type ChannelManagement struct {
	Header               `aisWidth:"38"`
	Valid                bool                           `aisEncodeMaxLen:"168"`
	Spare1               uint8                          `aisWidth:"2" aisEncodeAs:"0"`
	ChannelA             uint16                         `aisWidth:"12"`
	ChannelB             uint16                         `aisWidth:"12"`
	TxRxMode             uint8                          `aisWidth:"4"`
	LowPower             bool                           `aisWidth:"1"`
	Area                 ChannelManagementBroadcastData `aisWidth:"70" aisDependsBit:"~139" aisDependsField:"~IsAddressed"`
	Unicast              ChannelManagementUnicastData   `aisWidth:"70" aisDependsBit:"139" aisDependsField:"IsAddressed"`
	IsAddressed          bool                           `aisWidth:"1"`
	BwA                  bool                           `aisWidth:"1"`
	BwB                  bool                           `aisWidth:"1"`
	TransitionalZoneSize uint8                          `aisWidth:"3"`
	Spare4               uint32                         `aisWidth:"23" aisEncodeAs:"0"`
}

// SingleSlotBinaryMessage is primarily intended short infrequent data transmissions. The single slot
// binary message can contain up to 128 data-bits depending on the coding method used for the contents,
// and the destination indication of broadcast or addressed. The length should not exceed one slot. See
// application identifiers in § 2.1, Annex 5.
type SingleSlotBinaryMessage struct {
	Header             `aisWidth:"38"`
	Valid              bool                       `aisEncodeMaxLen:"168"`
	DestinationIDValid bool                       `aisWidth:"1"`
	ApplicationIDValid bool                       `aisWidth:"1"`
	DestinationID      uint32                     `aisWidth:"30" aisDependsBit:"38" aisDependsField:"DestinationIDValid"`
	Spare              uint8                      `aisWidth:"2" aisDependsBit:"38" aisDependsField:"DestinationIDValid" aisEncodeAs:"0"`
	ApplicationID      FieldApplicationIdentifier `aisWidth:"16" aisDependsBit:"39" aisDependsField:"ApplicationIDValid"`
	Payload            []byte                     `aisWidth:"-1"`
}

// MultiSlotBinaryMessage is primarily intended for scheduled binary data transmissions by applying either
// the SOTDMA or ITDMA access scheme. This multiple slot binary message can contain up to 1 004 data-
// bits (using 5 slots) depending on the coding method used for the contents, and the destination
// indication of broadcast or addressed. See application identifiers in § 2.1, Annex 5.
type MultiSlotBinaryMessage struct {
	Header                  `aisWidth:"38"`
	Valid                   bool                       `aisEncodeMaxLen:"1064"`
	DestinationIDValid      bool                       `aisWidth:"1"`
	ApplicationIDValid      bool                       `aisWidth:"1"`
	DestinationID           uint32                     `aisWidth:"30" aisDependsBit:"38" aisDependsField:"DestinationIDValid"`
	Spare1                  uint8                      `aisWidth:"2" aisDependsBit:"38" aisDependsField:"DestinationIDValid" aisEncodeAs:"0"`
	ApplicationID           FieldApplicationIdentifier `aisWidth:"16" aisDependsBit:"39" aisDependsField:"ApplicationIDValid"`
	Payload                 []byte                     `aisWidth:"-1"`
	Spare2                  uint8                      `aisWidth:"4"` /* Quite a few encoders seem to put data in the Spare2 bits, so we don't force encode it as zero */
	CommunicationStateItdma `aisWidth:"20"`
}

// FieldETA represents the encoding of the estimated time of arrival
type FieldETA struct {
	Month  uint8 `aisWidth:"4"`
	Day    uint8 `aisWidth:"5"`
	Hour   uint8 `aisWidth:"5"`
	Minute uint8 `aisWidth:"6"`
}

// FieldDimension represents the encoding of the dimension
type FieldDimension struct {
	A uint16 `aisWidth:"9"`
	B uint16 `aisWidth:"9"`
	C uint8  `aisWidth:"6"`
	D uint8  `aisWidth:"6"`
}

// FieldApplicationIdentifier represents the encoding of the application identifier
type FieldApplicationIdentifier struct {
	Valid              bool
	DesignatedAreaCode uint16 `aisWidth:"10"`
	FunctionIdentifier uint8  `aisWidth:"6"`
}

// HasCommunicationState indicates that a message contains commuication state
type HasCommunicationState interface {
	IsItdma() int
	GetState() uint32
}

// CommunicationStateItdma represents the encoding of the communication state if the
// ITDMA type is included in the message
type CommunicationStateItdma struct {
	CommunicationStateIsItdma bool   `aisWidth:"1"`
	CommunicationState        uint32 `aisWidth:"19"`
}

// IsItdma indicates if Itdma is used.
func (c CommunicationStateItdma) IsItdma() int {
	if c.CommunicationStateIsItdma {
		return 1
	}
	return 0
}

// GetState will return the communication state
func (c CommunicationStateItdma) GetState() uint32 {
	return c.CommunicationState
}

// CommunicationStateNoItdma represents the encoding of the communication state if the
// type is fixed
type CommunicationStateNoItdma struct {
	CommunicationState uint32 `aisWidth:"19"`
}

// IsItdma indicates if Itdma is used. Return -1 for unknown as it depends
// on the message
func (c CommunicationStateNoItdma) IsItdma() int {
	return -1
}

// GetState will return the communication state
func (c CommunicationStateNoItdma) GetState() uint32 {
	return c.CommunicationState
}

// Field10 represents an unsigned value multiplied by 10
type Field10 float64

// FieldLatLonCoarse represents a 1/10' position
type FieldLatLonCoarse float64

// FieldLatLonFine represents a 1/10000' position
type FieldLatLonFine float64

func init() {
	msgMap[1].rType = reflect.TypeOf(PositionReport{})
	msgMap[2].rType = reflect.TypeOf(PositionReport{})
	msgMap[3].rType = reflect.TypeOf(PositionReport{})
	msgMap[4].rType = reflect.TypeOf(BaseStationReport{})
	msgMap[5].rType = reflect.TypeOf(ShipStaticData{})
	msgMap[6].rType = reflect.TypeOf(AddressedBinaryMessage{})
	msgMap[7].rType = reflect.TypeOf(BinaryAcknowledge{})
	msgMap[8].rType = reflect.TypeOf(BinaryBroadcastMessage{})
	msgMap[9].rType = reflect.TypeOf(StandardSearchAndRescueAircraftReport{})
	msgMap[10].rType = reflect.TypeOf(CoordinatedUTCInquiry{})
	msgMap[11].rType = reflect.TypeOf(BaseStationReport{})
	msgMap[12].rType = reflect.TypeOf(AddessedSafetyMessage{})
	msgMap[13].rType = reflect.TypeOf(BinaryAcknowledge{})
	msgMap[14].rType = reflect.TypeOf(SafetyBroadcastMessage{})
	msgMap[15].rType = reflect.TypeOf(Interrogation{})
	msgMap[16].rType = reflect.TypeOf(AssignedModeCommand{})
	msgMap[17].rType = reflect.TypeOf(GnssBroadcastBinaryMessage{})
	msgMap[18].rType = reflect.TypeOf(StandardClassBPositionReport{})
	msgMap[19].rType = reflect.TypeOf(ExtendedClassBPositionReport{})
	msgMap[20].rType = reflect.TypeOf(DataLinkManagementMessage{})
	msgMap[21].rType = reflect.TypeOf(AidsToNavigationReport{})
	msgMap[22].rType = reflect.TypeOf(ChannelManagement{})
	msgMap[23].rType = reflect.TypeOf(GroupAssignmentCommand{})
	msgMap[24].rType = reflect.TypeOf(StaticDataReport{})
	msgMap[25].rType = reflect.TypeOf(SingleSlotBinaryMessage{})
	msgMap[26].rType = reflect.TypeOf(MultiSlotBinaryMessage{})
	msgMap[27].rType = reflect.TypeOf(LongRangeAisBroadcastMessage{})
}

func encodeHelper(p Packet) Packet {
	switch x := p.(type) {
	case Interrogation:
		if x.Station2.Valid && !x.Station1Msg2.Valid {
			/* All values should be set to zero and the field must be encoded! */
			x.Station1Msg2 = InterrogationStation1Message2{Valid: true}
		}
		return x
	}

	return p
}

func decodeHelper(p Packet) Packet {
	switch x := p.(type) {
	case Interrogation:
		if x.Station2.Valid {
			/* If Station1Msg2 is all zeros it actually is not valid */
			if x.Station1Msg2.MessageID == 0 && x.Station1Msg2.SlotOffset == 0 {
				x.Station1Msg2.Valid = false
			}
		}
		return x
	}

	return p
}
