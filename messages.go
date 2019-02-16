package goais

import "reflect"

var msgMap [28]reflect.Type

// AISPositionReport should be output periodically by mobile stations. The message ID is 1, 2 or 3
// depending on the system mode.
type AISPositionReport struct {
	Valid                     bool               `aisEncodeMaxLen:"168"`
	MessageID                 uint8              `aisWidth:"6"`
	RepeatIndicator           uint8              `aisWidth:"2"`
	UserID                    uint32             `aisWidth:"30"`
	NavigationalStatus        uint8              `aisWidth:"4"`
	RateOfTurn                int8               `aisWidth:"8"`
	Sog                       uint16             `aisWidth:"10"`
	PositionAccuracy          bool               `aisWidth:"1"`
	Longitude                 AISFieldLatLonFine `aisWidth:"28"`
	Latitude                  AISFieldLatLonFine `aisWidth:"27"`
	Cog                       AISField10         `aisWidth:"12"`
	TrueHeading               uint16             `aisWidth:"9"`
	Timestamp                 uint8              `aisWidth:"6"`
	SpecialManoeuvreIndicator uint8              `aisWidth:"2"`
	Spare                     uint8              `aisWidth:"3" aisEncodeAs:"0"`
	Raim                      bool               `aisWidth:"1"`
	CommunicationState        uint32             `aisWidth:"19"`
}

// AISBaseStationReport should be used for reporting UTC time and date and, at the same time, position.
// A base station should use Message 4 in its periodical transmissions. Message 4 is used by AIS stations
// for determining if it is within 120 NM for response to Messages 20 and 23. A mobile station should output
// Message 11 only in response to interrogation by Message 10. Message 11 is only transmitted as a result
// of a UTC request message (Message 10). The UTC and date response should be transmitted on the channel,
// where the UTC request message was received. */
type AISBaseStationReport struct {
	Valid              bool               `aisEncodeMaxLen:"168"`
	MessageID          uint8              `aisWidth:"6"`
	RepeatIndicator    uint8              `aisWidth:"2"`
	UserID             uint32             `aisWidth:"30"`
	UtcYear            uint16             `aisWidth:"14"`
	UtcMonth           uint8              `aisWidth:"4"`
	UtcDay             uint8              `aisWidth:"5"`
	UtcHour            uint8              `aisWidth:"5"`
	UtcMinute          uint8              `aisWidth:"6"`
	UtcSecond          uint8              `aisWidth:"6"`
	PositionAccuracy   bool               `aisWidth:"1"`
	Longitude          AISFieldLatLonFine `aisWidth:"28"`
	Latitude           AISFieldLatLonFine `aisWidth:"27"`
	FixType            uint8              `aisWidth:"4"`
	LongRangeEnable    bool               `aisWidth:"1"`
	Spare              uint16             `aisWidth:"9" aisEncodeAs:"0"`
	Raim               bool               `aisWidth:"1"`
	CommunicationState uint32             `aisWidth:"19"`
}

// AISShipStaticData should only be used by Class A shipborne and SAR aircraft AIS stations when reporting static
// or voyage related data.
type AISShipStaticData struct {
	Valid                bool              `aisEncodeMaxLen:"424"`
	MessageID            uint8             `aisWidth:"6"`
	RepeatIndicator      uint8             `aisWidth:"2"`
	UserID               uint32            `aisWidth:"30"`
	AisVersion           uint8             `aisWidth:"2"`
	ImoNumber            uint32            `aisWidth:"30"`
	CallSign             string            `aisWidth:"42"`
	Name                 string            `aisWidth:"120"`
	Type                 uint8             `aisWidth:"8"`
	Dimension            AISFieldDimension `aisWidth:"30"`
	FixType              uint8             `aisWidth:"4"`
	Eta                  AISFieldETA       `aisWidth:"20"`
	MaximumStaticDraught uint8             `aisWidth:"8"`
	Destination          string            `aisWidth:"120"`
	Dte                  bool              `aisWidth:"1"`
	Spare                bool              `aisWidth:"1" aisEncodeAs:"0"`
}

// AISAddressedBinaryMessage should be variable in length, based on the amount of binary data.
// The length should vary between 1 and 5 slots. See application identifiers in § 2.1, Annex 5.
type AISAddressedBinaryMessage struct {
	Valid           bool                          `aisEncodeMaxLen:"1008"`
	MessageID       uint8                         `aisWidth:"6"`
	RepeatIndicator uint8                         `aisWidth:"2"`
	SourceID        uint32                        `aisWidth:"30"`
	SequenceNumber  uint8                         `aisWidth:"2"`
	DestinationID   uint32                        `aisWidth:"30"`
	Retransmission  bool                          `aisWidth:"1"`
	Spare           bool                          `aisWidth:"1" aisEncodeAs:"0"`
	ApplicationID   AISFieldApplicationIdentifier `aisWidth:"16"`
	BinaryData      []byte                        `aisWidth:"-1"`
}

// AISBinaryBroadcastMessage will be variable in length, based on the amount of binary data.
// The length should vary between 1 and 5 slots.
type AISBinaryBroadcastMessage struct {
	Valid           bool                          `aisEncodeMaxLen:"1008"`
	MessageID       uint8                         `aisWidth:"6"`
	RepeatIndicator uint8                         `aisWidth:"2"`
	SourceID        uint32                        `aisWidth:"30"`
	Spare           uint8                         `aisWidth:"2" aisEncodeAs:"0"`
	ApplicationID   AISFieldApplicationIdentifier `aisWidth:"16"`
	BinaryData      []byte                        `aisWidth:"-1"`
}

// AISStandardSearchAndRescueAircraftReport should be used as a standard position report for
// aircraft involved in SAR operations. Stations other than aircraft involved in SAR operations
// should not transmit this message. The default reporting interval for this message should be 10 s.
type AISStandardSearchAndRescueAircraftReport struct {
	Valid              bool                       `aisEncodeMaxLen:"168"`
	MessageID          uint8                      `aisWidth:"6"`
	RepeatIndicator    uint8                      `aisWidth:"2"`
	UserID             uint32                     `aisWidth:"30"`
	Altitude           uint16                     `aisWidth:"12"`
	Sog                uint16                     `aisWidth:"10"`
	PositionAccuracy   bool                       `aisWidth:"1"`
	Longitude          AISFieldLatLonFine         `aisWidth:"28"`
	Latitude           AISFieldLatLonFine         `aisWidth:"27"`
	Cog                AISField10                 `aisWidth:"12"`
	Timestamp          uint8                      `aisWidth:"6"`
	AltFromBaro        bool                       `aisWidth:"1"`
	Spare1             uint8                      `aisWidth:"7" aisEncodeAs:"0"`
	Dte                bool                       `aisWidth:"1"`
	Spare2             uint8                      `aisWidth:"3" aisEncodeAs:"0"`
	AssignedMode       bool                       `aisWidth:"1"`
	Raim               bool                       `aisWidth:"1"`
	CommunicationState AISFieldCommunicationState `aisWidth:"20"`
}

// AISCoordinatedUTCInquiry should be used when a station is requesting UTC and date from another
// station.
type AISCoordinatedUTCInquiry struct {
	Valid           bool   `aisEncodeMaxLen:"72"`
	MessageID       uint8  `aisWidth:"6"`
	RepeatIndicator uint8  `aisWidth:"2"`
	SourceID        uint32 `aisWidth:"30"`
	Spare1          uint8  `aisWidth:"2" aisEncodeAs:"0"`
	DestinationID   uint32 `aisWidth:"30"`
	Spare2          uint8  `aisWidth:"2" aisEncodeAs:"0"`
}

// AISAddessedSafetyMessage could be variable in length, based on the amount of safety related text.
// The length should vary between 1 and 5 slots.
type AISAddessedSafetyMessage struct {
	Valid           bool   `aisEncodeMaxLen:"1008"`
	MessageID       uint8  `aisWidth:"6"`
	RepeatIndicator uint8  `aisWidth:"2"`
	SourceID        uint32 `aisWidth:"30"`
	SequenceNumber  uint8  `aisWidth:"2"`
	DestinationID   uint32 `aisWidth:"30"`
	Retransmission  bool   `aisWidth:"1"`
	Spare           bool   `aisWidth:"1" aisEncodeAs:"0"`
	Text            string `aisWidth:"-1"`
}

// AISSafetyBroadcastMessage could be variable in length, based on the amount of safety related
// text. The length should vary between 1 and 5 slots.
type AISSafetyBroadcastMessage struct {
	Valid           bool   `aisEncodeMaxLen:"1008"`
	MessageID       uint8  `aisWidth:"6"`
	RepeatIndicator uint8  `aisWidth:"2"`
	SourceID        uint32 `aisWidth:"30"`
	Spare           uint8  `aisWidth:"2" aisEncodeAs:"0"`
	Text            string `aisWidth:"-1"`
}

// AISGnssBroadcastBinaryMessage should be transmitted by a base station, which is connected to a DGNSS
// reference source, and configured to provide DGNSS data to receiving stations. The contents of the
// data should be in accordance with Recommendation ITU-R M.823, excluding preamble and parity formatting.
type AISGnssBroadcastBinaryMessage struct {
	Valid           bool                 `aisEncodeMaxLen:"816"`
	MessageID       uint8                `aisWidth:"6"`
	RepeatIndicator uint8                `aisWidth:"2"`
	SourceID        uint32               `aisWidth:"30"`
	Spare1          uint8                `aisWidth:"2" aisEncodeAs:"0"`
	Longitude       AISFieldLatLonCoarse `aisWidth:"18"`
	Latitude        AISFieldLatLonCoarse `aisWidth:"17"`
	Spare2          uint8                `aisWidth:"5" aisEncodeAs:"0"`
	Data            []byte               `aisWidth:"-1"`
}

// AISStandardClassBPositionReport should be output periodically and autonomously instead of
// Messages 1, 2, or 3 by Class B shipborne mobile equipment, only. The reporting interval should default to
// the values given in Table 2, Annex 1, unless otherwise specified by reception of a Message 16 or 23; and
// depending on the current SOG and navigational status flag setting.
type AISStandardClassBPositionReport struct {
	Valid              bool                       `aisEncodeMaxLen:"168"`
	MessageID          uint8                      `aisWidth:"6"`
	RepeatIndicator    uint8                      `aisWidth:"2"`
	UserID             uint32                     `aisWidth:"30"`
	Spare1             uint8                      `aisWidth:"8" aisEncodeAs:"0"`
	Sog                uint16                     `aisWidth:"10"`
	PositionAccuracy   bool                       `aisWidth:"1"`
	Longitude          AISFieldLatLonFine         `aisWidth:"28"`
	Latitude           AISFieldLatLonFine         `aisWidth:"27"`
	Cog                AISField10                 `aisWidth:"12"`
	TrueHeading        uint16                     `aisWidth:"9"`
	Timestamp          uint8                      `aisWidth:"6"`
	Spare2             uint8                      `aisWidth:"2" aisEncodeAs:"0"`
	ClassBUnit         bool                       `aisWidth:"1"`
	ClassBDisplay      bool                       `aisWidth:"1"`
	ClassBDsc          bool                       `aisWidth:"1"`
	ClassBBand         bool                       `aisWidth:"1"`
	ClassBMsg22        bool                       `aisWidth:"1"`
	AssignedMode       bool                       `aisWidth:"1"`
	Raim               bool                       `aisWidth:"1"`
	CommunicationState AISFieldCommunicationState `aisWidth:"20"`
}

// AISExtendedClassBPositionReport should be transmitted once every 6 min in two slots allocated by the
// use of Message 18 in the ITDMA communication state. This message should be transmitted immediately
// after the following parameter values change: dimension of ship/reference for position or type of
// electronic position fixing device.
//   For future equipment: this message is not needed and should not be used. All content is covered by
//   Message 18, Message 24A and 24B.
//   For legacy equipment: this message should be used by Class B shipborne mobile equipment.
type AISExtendedClassBPositionReport struct {
	Valid            bool               `aisEncodeMaxLen:"312"`
	MessageID        uint8              `aisWidth:"6"`
	RepeatIndicator  uint8              `aisWidth:"2"`
	UserID           uint32             `aisWidth:"30"`
	Spare1           uint8              `aisWidth:"8" aisEncodeAs:"0"`
	Sog              uint16             `aisWidth:"10"`
	PositionAccuracy bool               `aisWidth:"1"`
	Longitude        AISFieldLatLonFine `aisWidth:"28"`
	Latitude         AISFieldLatLonFine `aisWidth:"27"`
	Cog              AISField10         `aisWidth:"12"`
	TrueHeading      uint16             `aisWidth:"9"`
	Timestamp        uint8              `aisWidth:"6"`
	Spare2           uint8              `aisWidth:"4" aisEncodeAs:"0"`
	Name             string             `aisWidth:"120"`
	Type             uint8              `aisWidth:"8"`
	Dimension        AISFieldDimension  `aisWidth:"30"`
	FixType          uint8              `aisWidth:"4"`
	Raim             bool               `aisWidth:"1"`
	Dte              bool               `aisWidth:"1"`
	AssignedMode     bool               `aisWidth:"1"`
	Spare3           uint8              `aisWidth:"4" aisEncodeAs:"0"`
}

// AISAidsToNavigationReport should be used by an Aids to navigation (AtoN) AIS station. This station
// may be mounted on an aid-to-navigation or this message may be transmitted by a fixed station when
// the functionality of an AtoN station is integrated into the fixed station. This message should be
// transmitted autonomously at a Rr of once every three (3) min or it may be assigned by an assigned
// mode command (Message 16) via the VHF data link, or by an external command. This message
// should not occupy more than two slots.
type AISAidsToNavigationReport struct {
	Valid            bool               `aisEncodeMaxLen:"356"`
	MessageID        uint8              `aisWidth:"6"`
	RepeatIndicator  uint8              `aisWidth:"2"`
	ID               uint32             `aisWidth:"30"`
	Type             uint8              `aisWidth:"5"`
	Name             string             `aisWidth:"120"`
	PositionAccuracy bool               `aisWidth:"1"`
	Longitude        AISFieldLatLonFine `aisWidth:"28"`
	Latitude         AISFieldLatLonFine `aisWidth:"27"`
	Dimension        AISFieldDimension  `aisWidth:"30"`
	Fixtype          uint8              `aisWidth:"4"`
	Timestamp        uint8              `aisWidth:"6"`
	OffPosition      bool               `aisWidth:"1"`
	AtoN             uint8              `aisWidth:"8"`
	Raim             bool               `aisWidth:"1"`
	VirtualAtoN      bool               `aisWidth:"1"`
	AssignedMode     bool               `aisWidth:"1"`
	Spare            bool               `aisWidth:"1" aisEncodeAs:"0"`
	NameExtension    string             `aisWidth:"-1"`
}

// AISGroupAssignmentCommand is transmitted by a base station when operating as a controlling
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
type AISGroupAssignmentCommand struct {
	Valid             bool                 `aisEncodeMaxLen:"160"`
	MessageID         uint8                `aisWidth:"6"`
	RepeatIndicator   uint8                `aisWidth:"2"`
	SourceID          uint32               `aisWidth:"30"`
	Spare1            uint8                `aisWidth:"2" aisEncodeAs:"0"`
	Longitude1        AISFieldLatLonCoarse `aisWidth:"18"`
	Latitude1         AISFieldLatLonCoarse `aisWidth:"17"`
	Longitude2        AISFieldLatLonCoarse `aisWidth:"18"`
	Latitude2         AISFieldLatLonCoarse `aisWidth:"17"`
	StationType       uint8                `aisWidth:"4"`
	ShipType          uint8                `aisWidth:"8"`
	Spare2            uint32               `aisWidth:"22" aisEncodeAs:"0"`
	TxRxMode          uint8                `aisWidth:"2"`
	ReportingInterval uint8                `aisWidth:"4"`
	QuietTime         uint8                `aisWidth:"4"`
	Spare3            uint8                `aisWidth:"6" aisEncodeAs:"0"`
}

// AISStaticDataReportA is the A part of message 24
type AISStaticDataReportA struct {
	Valid           bool   `aisEncodeMaxLen:"160"`
	MessageID       uint8  `aisWidth:"6"`
	RepeatIndicator uint8  `aisWidth:"2"`
	UserID          uint32 `aisWidth:"30"`
	PartNumber      uint8  `aisWidth:"2" aisEncodeAs:"0"`
	Name            string `aisWidth:"120"`
}

// AISStaticDataReportB is the B part of message 24
type AISStaticDataReportB struct {
	Valid           bool              `aisEncodeMaxLen:"168"`
	MessageID       uint8             `aisWidth:"6"`
	RepeatIndicator uint8             `aisWidth:"2"`
	UserID          uint32            `aisWidth:"30"`
	PartNumber      uint8             `aisWidth:"2" aisEncodeAs:"1"`
	ShipType        uint8             `aisWidth:"8"`
	VendorIDName    string            `aisWidth:"18"`
	VenderIDModel   uint8             `aisWidth:"4"`
	VenderIDSerial  uint32            `aisWidth:"20"`
	CallSign        string            `aisWidth:"42"`
	Dimension       AISFieldDimension `aisWidth:"30"`
	FixType         uint8             `aisWidth:"4"`
	Spare           uint8             `aisWidth:"2" aisEncodeAs:"0"`
}

// AISLongRangeAisBroadcastMessage is primarily intended for long-range detection of AIS Class A
// and Class B “SO” equipped vessels (typically by satellite). This message has a similar content
// to Messages 1, 2 and 3, but the total number of bits has been compressed to allow for increased
// propagation delays associated with long-range detection. Refer to Annex 4 for details on
//Long-Range applications.
type AISLongRangeAisBroadcastMessage struct {
	Valid              bool                 `aisEncodeMaxLen:"96"`
	MessageID          uint8                `aisWidth:"6"`
	RepeatIndicator    uint8                `aisWidth:"2"`
	UserID             uint32               `aisWidth:"30"`
	PositionAccuracy   bool                 `aisWidth:"1"`
	Raim               bool                 `aisWidth:"1"`
	NavigationalStatus uint8                `aisWidth:"4"`
	Longitude          AISFieldLatLonCoarse `aisWidth:"18"`
	Latitude           AISFieldLatLonCoarse `aisWidth:"17"`
	Sog                uint8                `aisWidth:"6"`
	Cos                uint16               `aisWidth:"9"`
	PositionLatency    bool                 `aisWidth:"1"`
	Spare              bool                 `aisWidth:"1" aisEncodeAs:"0"`
}

// AISBinaryAcknowledgeData is the data part of AISBinaryAcknowledge
type AISBinaryAcknowledgeData struct {
	Valid          bool
	DestinationID  uint32 `aisWidth:"30"`
	SequenceNumber uint8  `aisWidth:"2"`
}

// AISBinaryAcknowledge should be used as an acknowledgement of up to four Message 6 messages received
// (see § 5.3.1, Annex 2) and should be transmitted on the channel, where the addressed message to be
// acknowledged was received.
type AISBinaryAcknowledge struct {
	Valid           bool                        `aisEncodeMaxLen:"168"`
	MessageID       uint8                       `aisWidth:"6"`
	RepeatIndicator uint8                       `aisWidth:"2"`
	SourceID        uint32                      `aisWidth:"30"`
	Spare           uint8                       `aisWidth:"2" aisEncodeAs:"0"`
	Destinations    [4]AISBinaryAcknowledgeData `aisWidth:"0"`
}

// AISInterrogationStation1 is the station 1 part of AISInterrogation
type AISInterrogationStation1 struct {
	Valid      bool
	MessageID  uint8  `aisWidth:"6"`
	SlotOffset uint16 `aisWidth:"12"`
	Spare      uint8  `aisWidth:"2" aisEncodeAs:"0"`
}

// AISInterrogationStation2 is the station 2 part of AISInterrogation
type AISInterrogationStation2 struct {
	Valid      bool
	Station2ID uint32 `aisWidth:"30"`
	MessageID  uint8  `aisWidth:"6"`
	SlotOffset uint16 `aisWidth:"12"`
	Spare      uint8  `aisWidth:"2" aisEncodeAs:"0"`
}

// AISInterrogation should be used for interrogations via the TDMA (not DSC) VHF data link except for
// requests for UTC and date. The response should be transmitted on the channel where the interrogation
// was received.
type AISInterrogation struct {
	Valid           bool                        `aisEncodeMaxLen:"160"`
	MessageID       uint8                       `aisWidth:"6"`
	RepeatIndicator uint8                       `aisWidth:"2"`
	SourceID        uint32                      `aisWidth:"30"`
	Spare           uint8                       `aisWidth:"2" aisEncodeAs:"0"`
	Station1ID      uint32                      `aisWidth:"30"`
	Station1        [2]AISInterrogationStation1 `aisWidth:"0"`
	Station2        [1]AISInterrogationStation2 `aisWidth:"0"`
}

// AISAssignedModeCommandData is the data part of AISAssignedModeCommand
type AISAssignedModeCommandData struct {
	Valid         bool
	DestinationID uint32 `aisWidth:"30"`
	Offset        uint16 `aisWidth:"12"`
	Increment     uint16 `aisWidth:"10"`
}

// AISAssignedModeCommand be transmitted by a base station when operating as a controlling entity. Other
// stations can be assigned a transmission schedule, other than the currently used one. If a station is
// assigned a schedule, it will also enter assigned mode.
type AISAssignedModeCommand struct {
	Valid           bool                          `aisEncodeMaxLen:"144"`
	MessageID       uint8                         `aisWidth:"6"`
	RepeatIndicator uint8                         `aisWidth:"2"`
	SourceID        uint32                        `aisWidth:"30"`
	Spare           uint8                         `aisWidth:"2" aisEncodeAs:"0"`
	Commands        [2]AISAssignedModeCommandData `aisWidth:"0"`
}

// AISDataLinkManagementMessageData is the data part of AISDataLinkManagementMessage
type AISDataLinkManagementMessageData struct {
	Valid         bool
	Offset        uint16 `aisWidth:"12"`
	NumberOfSlots uint8  `aisWidth:"4"`
	TimeOut       uint8  `aisWidth:"3"`
	Increment     uint16 `aisWidth:"11"`
}

// AISDataLinkManagementMessage should be used by base station(s) to pre-announce the fixed allocation
// schedule (FATDMA) for one or more base station(s) and it should be repeated as often as required. This
// way the system can provide a high level of integrity for base station(s). This is especially important in
// regions where several base stations are located adjacent to each other and mobile station(s) move
// between these different regions. These reserved slots cannot be autonomously allocated by mobile
// stations.
type AISDataLinkManagementMessage struct {
	Valid           bool                                `aisEncodeMaxLen:"160"`
	MessageID       uint8                               `aisWidth:"6"`
	RepeatIndicator uint8                               `aisWidth:"2"`
	SourceID        uint32                              `aisWidth:"30"`
	Spare           uint8                               `aisWidth:"2" aisEncodeAs:"0"`
	Data            [4]AISDataLinkManagementMessageData `aisWidth:"0"`
}

// AISChannelManagement should be transmitted by a base station (as a broadcast message) to command the VHF
// data link parameters for the geographical area designated in this message and should be accompanied
// by a Message 4 transmission for evaluation of the message within 120 NM. The geographical area
// designated by this message should be as defined in § 4.1, Annex 2. Alternatively, this message may
// be used by a base station (as an addressed message) to command individual AIS mobile stations to
// adopt the specified VHF data link parameters. When interrogated and no channel management
// performed by the interrogated base station, the not available and/or international default settings
// should be transmitted (see § 4.1, Annex 2).
type AISChannelManagement struct {
	Valid                bool                 `aisEncodeMaxLen:"168"`
	MessageID            uint8                `aisWidth:"6"`
	RepeatIndicator      uint8                `aisWidth:"2"`
	StationID            uint32               `aisWidth:"30"`
	Spare1               uint8                `aisWidth:"2" aisEncodeAs:"0"`
	ChannelA             uint16               `aisWidth:"12"`
	ChannelB             uint16               `aisWidth:"12"`
	TxRxMode             uint8                `aisWidth:"4"`
	LowPower             bool                 `aisWidth:"1"`
	Longitude1           AISFieldLatLonCoarse `aisWidth:"18" aisDependsBit:"~139" aisDependsField:"~IsAddressed"`
	Latitude1            AISFieldLatLonCoarse `aisWidth:"17" aisDependsBit:"~139" aisDependsField:"~IsAddressed"`
	Longitude2           AISFieldLatLonCoarse `aisWidth:"18" aisDependsBit:"~139" aisDependsField:"~IsAddressed"`
	Latitude2            AISFieldLatLonCoarse `aisWidth:"17" aisDependsBit:"~139" aisDependsField:"~IsAddressed"`
	AddressStation1      uint32               `aisWidth:"30" aisDependsBit:"139" aisDependsField:"IsAddressed"`
	Spare2               uint8                `aisWidth:"5" aisDependsBit:"139" aisDependsField:"IsAddressed" aisEncodeAs:"0"`
	AddressStation2      uint32               `aisWidth:"30" aisDependsBit:"139" aisDependsField:"IsAddressed"`
	Spare3               uint8                `aisWidth:"5" aisDependsBit:"139" aisDependsField:"IsAddressed" aisEncodeAs:"0"`
	IsAddressed          bool                 `aisWidth:"1"`
	BwA                  bool                 `aisWidth:"1"`
	BwB                  bool                 `aisWidth:"1"`
	TransitionalZoneSize uint8                `aisWidth:"3"`
	Spare4               uint32               `aisWidth:"23" aisEncodeAs:"0"`
}

// AISSingleSlotBinaryMessage is primarily intended short infrequent data transmissions. The single slot
// binary message can contain up to 128 data-bits depending on the coding method used for the contents,
// and the destination indication of broadcast or addressed. The length should not exceed one slot. See
// application identifiers in § 2.1, Annex 5.
type AISSingleSlotBinaryMessage struct {
	Valid              bool                          `aisEncodeMaxLen:"168"`
	MessageID          uint8                         `aisWidth:"6"`
	RepeatIndicator    uint8                         `aisWidth:"2"`
	SourceID           uint32                        `aisWidth:"30"`
	DestinationIDValid bool                          `aisWidth:"1"`
	ApplicationIDValid bool                          `aisWidth:"1"`
	DestinationID      uint32                        `aisWidth:"30" aisDependsBit:"38" aisDependsField:"DestinationIDValid"`
	Spare              uint8                         `aisWidth:"2" aisDependsBit:"38" aisDependsField:"DestinationIDValid" aisEncodeAs:"0"`
	ApplicationID      AISFieldApplicationIdentifier `aisWidth:"16" aisDependsBit:"39" aisDependsField:"ApplicationIDValid"`
	Payload            []byte                        `aisWidth:"-1"`
}

// AISMultiSlotBinaryMessage is primarily intended for scheduled binary data transmissions by applying either
// the SOTDMA or ITDMA access scheme. This multiple slot binary message can contain up to 1 004 data-
// bits (using 5 slots) depending on the coding method used for the contents, and the destination
// indication of broadcast or addressed. See application identifiers in § 2.1, Annex 5.
type AISMultiSlotBinaryMessage struct {
	Valid              bool                          `aisEncodeMaxLen:"1064"`
	MessageID          uint8                         `aisWidth:"6"`
	RepeatIndicator    uint8                         `aisWidth:"2"`
	SourceID           uint32                        `aisWidth:"30"`
	DestinationIDValid bool                          `aisWidth:"1"`
	ApplicationIDValid bool                          `aisWidth:"1"`
	DestinationID      uint32                        `aisWidth:"30" aisDependsBit:"38" aisDependsField:"DestinationIDValid"`
	Spare              uint8                         `aisWidth:"2" aisDependsBit:"38" aisDependsField:"DestinationIDValid" aisEncodeAs:"0"`
	ApplicationID      AISFieldApplicationIdentifier `aisWidth:"16" aisDependsBit:"39" aisDependsField:"ApplicationIDValid"`
	Payload            []byte                        `aisWidth:"-1"`
	CommunicationState AISFieldCommunicationState    `aisWidth:"20"`
}

// AISFieldETA represents the encoding of the estimated time of arrival
type AISFieldETA struct {
	Valid  bool
	Month  uint8 `aisWidth:"4"`
	Day    uint8 `aisWidth:"5"`
	Hour   uint8 `aisWidth:"5"`
	Minute uint8 `aisWidth:"6"`
}

// AISFieldDimension represents the encoding of the dimension
type AISFieldDimension struct {
	Valid bool
	A     uint16 `aisWidth:"9"`
	B     uint16 `aisWidth:"9"`
	C     uint8  `aisWidth:"6"`
	D     uint8  `aisWidth:"6"`
}

// AISFieldApplicationIdentifier represents the encoding of the application identifier
type AISFieldApplicationIdentifier struct {
	Valid              bool
	DesignatedAreaCode uint16 `aisWidth:"10"`
	FunctionIdentifier uint8  `aisWidth:"6"`
}

// AISFieldCommunicationState represents the encoding of the communication state if the
// TDMA type is included in the message
type AISFieldCommunicationState struct {
	Valid   bool
	IsItdma bool   `aisWidth:"1"`
	State   uint32 `aisWidth:"19"`
}

// AISField10 represents a value multiplied by 10
type AISField10 float64

// AISFieldLatLonCoarse representes a 1/10' position
type AISFieldLatLonCoarse float64

// AISFieldLatLonFine representes a 1/10000' position
type AISFieldLatLonFine float64

func init() {
	msgMap[1] = reflect.TypeOf(AISPositionReport{})
	msgMap[2] = reflect.TypeOf(AISPositionReport{})
	msgMap[3] = reflect.TypeOf(AISPositionReport{})
	msgMap[4] = reflect.TypeOf(AISBaseStationReport{})
	msgMap[5] = reflect.TypeOf(AISShipStaticData{})
	msgMap[6] = reflect.TypeOf(AISAddressedBinaryMessage{})
	msgMap[7] = reflect.TypeOf(AISBinaryAcknowledge{})
	msgMap[8] = reflect.TypeOf(AISBinaryBroadcastMessage{})
	msgMap[9] = reflect.TypeOf(AISStandardSearchAndRescueAircraftReport{})
	msgMap[10] = reflect.TypeOf(AISCoordinatedUTCInquiry{})
	msgMap[11] = reflect.TypeOf(AISBaseStationReport{})
	msgMap[12] = reflect.TypeOf(AISAddessedSafetyMessage{})
	msgMap[13] = reflect.TypeOf(AISBinaryAcknowledge{})
	msgMap[14] = reflect.TypeOf(AISSafetyBroadcastMessage{})
	msgMap[15] = reflect.TypeOf(AISInterrogation{})
	msgMap[16] = reflect.TypeOf(AISAssignedModeCommand{})
	msgMap[17] = reflect.TypeOf(AISGnssBroadcastBinaryMessage{})
	msgMap[18] = reflect.TypeOf(AISStandardClassBPositionReport{})
	msgMap[19] = reflect.TypeOf(AISExtendedClassBPositionReport{})
	msgMap[20] = reflect.TypeOf(AISDataLinkManagementMessage{})
	msgMap[21] = reflect.TypeOf(AISAidsToNavigationReport{})
	msgMap[22] = reflect.TypeOf(AISChannelManagement{})
	msgMap[23] = reflect.TypeOf(AISGroupAssignmentCommand{})
	msgMap[25] = reflect.TypeOf(AISSingleSlotBinaryMessage{})
	msgMap[26] = reflect.TypeOf(AISMultiSlotBinaryMessage{})
	msgMap[27] = reflect.TypeOf(AISLongRangeAisBroadcastMessage{})
}
