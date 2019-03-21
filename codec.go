// Package ais is an Automatic Identification System (ITU-R M.1371-5) packet decoder and encoder
package ais

import (
	"reflect"
	"strconv"
)

// Codec encodes and decodes AIS messages (ITU-R M.1371-5)
type Codec struct {
	StrictByteAlignment bool
	minValidMap         map[string]uint

	// DecoderCheckFixedValues will validate that spare values are zero. This is only for debugging
	// and should not be used as newer protocl versions may use the spares for something.
	DecoderCheckFixedValues bool

	// FloatWithoutConversion will disable float conversion, making re-encoding the decoded output yield
	// an identical bitstream.
	FloatWithoutConversion bool

	// DropSpace will drop spaces at the end of strings. Many encoders pad strings
	// with spaces instead of @.
	DropSpace bool
}

func assert(condition bool, err string) {
	if !condition {
		panic("BUG: " + err)
	}
}

// CodecNew creates and initializes the AIS parser. The two parameters allow accepting
// messages that a few types of existing encoders seem to transmit with invalid length.
// This is very rare, passing false to both should be fine.
func CodecNew(acceptShortAck bool, acceptShortShipStaticData bool) *Codec {
	t := &Codec{}

	t.minValidMap = make(map[string]uint)

	if acceptShortAck {
		t.minValidMap["BinaryAcknowledgeData"] = 1
	}

	if acceptShortShipStaticData {
		t.minValidMap["ShipStaticData"] = 420
	}

	return t
}

// ChannelToFrequency converts an AIS channel number into its frequency in Hz
func (t *Codec) ChannelToFrequency(channel uint16) uint {
	/* https://www.itu.int/dms_pubrec/itu-r/rec/m/R-REC-M.1084-5-201203-I!!PDF-E.pdf */

	c := uint(channel % 1000)
	duplexType := channel / 1000

	shipFrequency := uint(0)
	if 1 <= c && c <= 28 {
		shipFrequency = 156050000 + (c-1)*50000
	} else if 60 <= c && c <= 88 {
		shipFrequency = 156025000 + (c-60)*50000
	} else if 260 <= c && c <= 287 {
		shipFrequency = 156037500 + (c-260)*50000
	} else if 201 <= c && c <= 228 {
		shipFrequency = 156062500 + (c-201)*50000
	} else {
		return 0
	}

	/* Is 287 a special case or a typo in the standard?
	   It is 1MHz higher than expected */

	/* In this range only simplex operation is defined */
	if 156375000 <= shipFrequency && shipFrequency <= 156887500 {
		return shipFrequency
	}

	switch duplexType {
	case 1:
		return shipFrequency
	case 2:
		return shipFrequency + 4600000
	default:
		return 0
	}
}

func makeSigned(input uint64, length uint) int64 {
	result := int64(input)
	maxValue := int64(1) << length

	if result >= maxValue/2 {
		result = result - maxValue
	}

	return result
}

func extractNumber(payload []byte, isSigned bool, offset *uint, width uint) int64 {
	result := uint64(0)

	for i := *offset; i < *offset+width; i++ {
		result <<= 1
		if i < uint(len(payload)) {
			result |= uint64(payload[i])
		}
	}

	*offset += width

	if isSigned {
		return makeSigned(result, width)
	}

	return int64(result)
}

func extractString(payload []byte, offset *uint, width uint, dropSpace bool) string {
	numChars := width / 6

	result := make([]byte, numChars)

	for i := uint(0); i < numChars; i++ {
		number := extractNumber(payload, false, offset, 6)

		if number < 32 {
			number = number + 64
		}

		result[i] = byte(number)
	}

	/* The string is closed by @ */
	stripSpace := len(result)
	for i := len(result) - 1; i >= 0; i-- {
		if result[i] != '@' {
			if !dropSpace || result[i] != ' ' {
				break
			}
		}
		stripSpace--
	}

	result = result[:stripSpace]

	return string(result)
}

func aisFindFieldLength(sf reflect.StructField, payload []byte) (valid bool, skip bool, fixedLength bool, length uint) {
	depends, ok := sf.Tag.Lookup("aisDependsBit")
	if ok {
		target := byte(1)
		var dependsI int
		if depends[0] == '~' {
			target = 0
			dependsI, _ = strconv.Atoi(depends[1:])
		} else {
			dependsI, _ = strconv.Atoi(depends)
		}

		if dependsI >= len(payload) {
			return false, false, false, 0
		}

		if payload[dependsI] != target {
			return true, true, false, 0
		}
	}

	vi, _ := strconv.Atoi(sf.Tag.Get("aisWidth"))
	if vi < 0 {
		return true, false, false, 0
	}
	return true, false, true, uint(vi)
}

func isBasicValue(val reflect.Value) (bool, int64) {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true, val.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true, int64(val.Uint())
	case reflect.Bool:
		if val.Bool() {
			return true, 1
		}
		return true, 0
	default:
		return false, 0
	}
}

func isSigned(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func (t *Codec) aisFillMessage(val reflect.Value, payload []byte, offset *uint) int {
	/* Return value of -2 propagates all the way up and fails the decode */

	optional := false

	/* Look up the struct */
	st := val.Type()
	strType := st.Name()

	validField := val.FieldByName("Valid")
	if validField.IsValid() {
		f, _ := st.FieldByName("Valid")
		_, optional = f.Tag.Lookup("aisOptional")
		validField.SetBool(false)
	}

	/* Calculate minimum length */
	minLength := uint(0)
	for i := 0; i < val.NumField(); i++ {
		valid, skip, fixedLength, length := aisFindFieldLength(st.Field(i), payload)
		if !valid {
			return -1
		}
		if fixedLength && !skip {
			minLength += length
		}
	}

	minBitsForValid, ok := t.minValidMap[strType]
	if !ok {
		minBitsForValid = minLength
	}

	/* Is the message long enough? */
	if len(payload)-int(*offset) < int(minBitsForValid) {
		if optional {
			return 0
		}
		return -1
	}

	for i := 0; i < val.NumField(); i++ {
		if st.Field(i).Name == "Valid" {
			continue
		}

		field := val.Field(i)

		_, skip, fixedLength, v := aisFindFieldLength(st.Field(i), payload)

		if skip {
			continue
		}

		/* Some fields take up the entire remainder of the message */
		if !fixedLength {
			varLen := len(payload) - int(minLength)
			assert(varLen >= 0, "Variable length payload length is negative.")

			v = uint(varLen)
		}

		var basicValue int64
		if b, _ := isBasicValue(field); b {
			checkValue := false
			correctValue := int64(0)

			encodeAsStr, encodeAsFound := st.Field(i).Tag.Lookup("aisCheckValue")
			if !encodeAsFound && t.DecoderCheckFixedValues {
				encodeAsStr, encodeAsFound = st.Field(i).Tag.Lookup("aisEncodeAs")
			}

			if encodeAsFound {
				correctValue, _ = strconv.ParseInt(encodeAsStr, 10, 64)
				checkValue = true
			}

			basicValue = extractNumber(payload, isSigned(field), offset, v)
			if checkValue && (basicValue != correctValue) {
				return -2
			}

		}

		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(basicValue)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetUint(uint64(basicValue))
		case reflect.Bool:
			field.SetBool(basicValue == 1)
		case reflect.String:
			field.SetString(extractString(payload, offset, v, t.DropSpace))
		case reflect.Slice:
			field.SetBytes(payload[*offset : *offset+v])
			*offset += v
		case reflect.Array:
			for k := 0; k < field.Len(); k++ {
				subField := field.Index(k)
				ok := t.aisFillMessage(subField, payload, offset)
				assert(ok != -2, "Could not encode array subfield. This should not happen with the current message definition")
				if ok == -1 {
					if k == 0 {
						return -2
					}
					break
				}
			}
		case reflect.Struct:
			ok := t.aisFillMessage(field, payload, offset)
			if ok < 0 {
				return ok
			}
		case reflect.Float64, reflect.Float32:
			value := extractNumber(payload, true, offset, v)
			field.SetFloat(float64(value))

			if !t.FloatWithoutConversion {
				switch field.Type().Name() {
				case "FieldLatLonFine":
					field.SetFloat(float64(value) / 10000.0 / 60.0)
				case "FieldLatLonCoarse":
					field.SetFloat(float64(value) / 10.0 / 60.0)
				case "Field10":
					field.SetFloat(float64(value) / 10.0)
				}
			}
		}
	}

	if validField.IsValid() {
		validField.SetBool(true)
	}

	return 0
}

// DecodePacket will convert a []byte containing 0 and 1 to an object containing the decoded packet.
// It will return nil if decoding failed.
func (t *Codec) DecodePacket(payload []byte) Packet {
	if len(payload)%8 != 0 {
		/* AIS messages should be a multiple of 8-bits:
		 *  [Order AIS message bits into 8-bit bytes for assembly of transmission packet, see § 3.3.7.]
		 * Most AIS messages are already a multiple of 8-bits and not all transmitters seem to implement
		 * this padding requirement, so you may get some messages that are not a multiple of 8 bits long.
		 * Also, some receivers seem to return invalid fillBits values (off by 1 or 2), therefore it is
		 * recommended to not treat bad padding as an error condition */
		if t.StrictByteAlignment {
			return nil
		}
	}

	if len(payload) < 6 {
		return nil
	}

	offset := uint(0)
	msgID := extractNumber(payload, false, &offset, 6)

	offset = 0

	/* Use default decoder */
	if msgID >= 1 && msgID <= 27 {
		msgType := msgMap[msgID]
		msgPtr := reflect.New(msgType.rType)
		if t.aisFillMessage(msgPtr.Elem(), payload, &offset) == 0 {
			switch out := msgPtr.Elem().Interface().(type) {
			case Packet:
				return decodeHelper(out)
			}
		}
	}

	return nil
}

func encodeNumber(packet []byte, isSigned bool, width uint, number int64) ([]byte, bool) {
	if !isSigned {
		maxVal := (int64(1) << width) - 1
		if number > maxVal {
			return packet, false
		}
	} else {
		minVal := -(int64(1) << width) / 2
		maxVal := (int64(1)<<width)/2 - 1

		if number < minVal || number > maxVal {
			return packet, false
		}
	}

	numUnsigned := uint64(number)

	for i := int(width - 1); i >= 0; i-- {
		packet = append(packet, byte((numUnsigned>>uint(i))&1))
	}

	return packet, true
}

func encodeString(packet []byte, width uint, fixedWidth bool, str string) ([]byte, bool) {
	var i uint
	for i = 0; i < uint(len(str)) && (i < width/6 || !fixedWidth); i++ {
		char := byte(str[i])

		if 64 <= char && char <= 95 {
			char -= 64
		} else if 32 <= char && char <= 63 {
			/* No translation needed */
		} else {
			/* This character is not valid */
			return packet, false
		}

		packet, _ = encodeNumber(packet, false, 6, int64(char))
	}

	/* Pad fixed width strings */
	if fixedWidth {
		for ; i < width/6; i++ {
			packet, _ = encodeNumber(packet, false, 6, 0)
		}
	}

	return packet, true
}

func aisEncodedLength(val reflect.Value, i int) (skip bool, fixedLength bool, length uint) {
	st := val.Type()

	sf := st.Field(i)
	depends, dependsFound := sf.Tag.Lookup("aisDependsField")
	if dependsFound {
		value := true
		if depends[0] == '~' {
			value = false
			depends = depends[1:]
		}
		df := val.FieldByName(depends)
		if df.Bool() != value {
			return true, true, 0
		}
	}
	vi, _ := strconv.Atoi(sf.Tag.Get("aisWidth"))

	if vi < 0 {
		return false, false, 0
	}

	return false, true, uint(vi)
}

func (t *Codec) aisEncodeMessage(val reflect.Value, packet []byte) ([]byte, bool) {
	vf := val.FieldByName("Valid")
	if vf.IsValid() && !vf.Bool() {

		/* Is it optional? */
		tf, _ := val.Type().FieldByName("Valid")
		_, opt := tf.Tag.Lookup("aisOptional")
		return packet, opt
	}

	st := val.Type()
	var ok bool

	for i := 0; i < val.NumField(); i++ {
		if st.Field(i).Name == "Valid" {
			continue
		}

		field := val.Field(i)
		skip, fixedLength, v := aisEncodedLength(val, i)

		if skip {
			continue
		}

		if b, k := isBasicValue(field); b {
			encodeAsStr, encodeAsFound := st.Field(i).Tag.Lookup("aisEncodeAs")
			if encodeAsFound {
				k, _ = strconv.ParseInt(encodeAsStr, 10, 64)
			}

			packet, ok = encodeNumber(packet, isSigned(field), v, k)
			if !ok {
				return packet, false
			}
		}

		switch field.Kind() {
		case reflect.String:
			packet, ok = encodeString(packet, v, fixedLength, field.String())
			if !ok {
				return packet, false
			}
		case reflect.Slice:
			tmp := field.Bytes()
			assert(!fixedLength, "Fixed length slices are not supported since they do not occur in the current spec")
			packet = append(packet, tmp...)

		case reflect.Array:
			for k := 0; k < field.Len(); k++ {
				subField := field.Index(k)
				packet, ok = t.aisEncodeMessage(subField, packet)
				if !ok {
					if k == 0 {
						return packet, false
					}
					break
				}
			}
		case reflect.Struct:
			packet, ok = t.aisEncodeMessage(field, packet)
			if !ok {
				return packet, false
			}
		case reflect.Float64, reflect.Float32:
			value := int64(field.Float())

			if !t.FloatWithoutConversion {
				switch field.Type().Name() {
				case "FieldLatLonFine":
					value = int64(field.Float() * 10000.0 * 60.0)
				case "FieldLatLonCoarse":
					value = int64(field.Float() * 10.0 * 60.0)
				case "Field10":
					value = int64(field.Float() * 10.0)
				}
			}

			packet, _ = encodeNumber(packet, true, v, value)
		}

	}

	return packet, true
}

// EncodePacket encodes a valid AIS object to a binary []byte.
// nil is returned if encoding failed.
func (t *Codec) EncodePacket(message Packet) []byte {
	/* Check if the type is correct */
	mID := message.GetHeader().MessageID
	if mID < 1 || mID > 27 {
		return nil
	}
	expectedType := msgMap[mID].rType
	if reflect.TypeOf(message) != expectedType {
		return nil
	}

	val := reflect.ValueOf(encodeHelper(message))

	vt, _ := val.Type().FieldByName("Valid")

	encodeString, ok := vt.Tag.Lookup("aisEncodeMaxLen")
	assert(ok, "aisEncodeMaxLen not found")
	encodeLen, _ := strconv.Atoi(encodeString)

	/* AIS packets need to be a multiple of 8 bits */
	if encodeLen%8 != 0 {
		encodeLen += 7 - (encodeLen % 8)
	}

	packet := make([]byte, 0, encodeLen)
	packet, ok = t.aisEncodeMessage(val, packet)
	if !ok {
		return nil
	}

	if encodeLen > 0 && len(packet) > encodeLen {
		return nil
	}

	/* Pad packet to 8-bit boundary:
	 * From 3.3.7: Unused bits in the last byte should be set to zero in order to preserve byte boundary.
	 * This means that you cannot send variable length binary messages where the payload is not a multiple of 8-bits.
	 * Otherwise, the padding will damage the attached communication state.
	 *
	 * Some encoders do seem to produce these messages (total length is not a multiple of 8), we can receive that
	 * but not encode it. These are likely proprietary systems anyway, all publically documented messages I have seen
	 * are a multiple of 8-bits long.
	 */
	for len(packet)%8 != 0 {
		packet = append(packet, 0)
	}

	return packet
}
