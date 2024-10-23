package aisnmeafast

import (
	"sync"

	"github.com/BertoldVdb/go-ais"
)

type stateType int

const (
	stateIdle stateType = iota
	stateTagBlock
	stateData

	stateChecksum0
	stateChecksum1
)

type recombState struct {
	MsgTotal int
	TagBlock []byte
	Data     [10][]byte
}

type Decoder struct {
	cfg DecoderConfig

	state stateType

	dataBlock []byte

	tagBlock      []byte
	tagBlockValid []byte

	aisRecombine map[[8]byte]*recombState
}

type NMEAParsed struct {
	Tagblock []byte
	Sentence []byte
}

type AISParsed struct {
	Talker  []byte
	Channel byte
	Packet  ais.Packet
}

type DecoderConfig struct {
	AIS *ais.Codec

	NMEAFunc             func(nmea NMEAParsed) error
	NMEAEncapsulatedFunc func(nmea NMEAParsed) error
	AISDecodedFunc       func(nmea NMEAParsed, ais AISParsed) error

	IgnoreChecksum bool
}

var fillDecodeTableOnce sync.Once

func New(cfg DecoderConfig) *Decoder {
	fillDecodeTableOnce.Do(fillDecodeTable)

	a := &Decoder{
		cfg: cfg,

		aisRecombine: make(map[[8]byte]*recombState),
	}

	return a
}

func (c *Decoder) nmeaCopyMany(i int, in []byte) int {
	/* Optimization to copy many bytes at once */
	for k := i; k < len(in); k++ {
		if in[k] == '*' {
			c.dataBlock = append(c.dataBlock, in[i:k+1]...)
			i = k
			c.state = stateChecksum0
			break
		}
	}

	/* Simply add all remaining bytes... */
	if c.state == stateData {
		c.dataBlock = append(c.dataBlock, in[i:]...)
		i = len(in)
		return i
	}

	return i
}

func (c *Decoder) Write(in []byte) (int, error) {
	if len(c.dataBlock) >= 8192 || len(c.tagBlock) >= 8192 {
		c.state = stateIdle
	}

	for i := 0; i < len(in); i++ {
		m := in[i]

		if m == '!' || m == '$' {
			c.dataBlock = c.dataBlock[:0]

			c.tagBlockValid = nil
			if c.state == stateTagBlock && len(c.tagBlock) >= 2 && c.tagBlock[len(c.tagBlock)-1] == '\\' {
				c.tagBlockValid = c.tagBlock[:len(c.tagBlock)-1]
			}

			c.state = stateData

			i = c.nmeaCopyMany(i, in)

		} else if c.state == stateIdle {
			if m == '\\' {
				c.state = stateTagBlock
				c.tagBlock = c.tagBlock[:0]
			}

		} else if c.state == stateTagBlock {
			c.tagBlock = append(c.tagBlock, m)

		} else if c.state == stateData {
			i = c.nmeaCopyMany(i, in)

		} else if c.state == stateChecksum0 {
			c.dataBlock = append(c.dataBlock, m)
			c.state = stateChecksum1

		} else if c.state == stateChecksum1 {
			c.dataBlock = append(c.dataBlock, m)
			c.state = stateIdle

			if err := c.handleMessage(); err != nil {
				return i, err
			}
		}
	}

	return len(in), nil
}

func nmeaReadNibble(nib byte) (uint8, bool) {
	if nib >= '0' && nib <= '9' {
		return nib - '0', true
	} else if nib >= 'A' && nib <= 'F' {
		return nib - 'A' + 10, true
	} else if nib >= 'a' && nib <= 'f' {
		return nib - 'a' + 10, true
	}

	return 0, false
}

func nmeaChecksumRead(cs []byte) (uint8, bool) {
	out0, ok := nmeaReadNibble(cs[0])
	if !ok {
		return 0, false
	}
	out1, ok := nmeaReadNibble(cs[1])
	if !ok {
		return 0, false
	}

	return out0<<4 | out1, true
}

func nmeaChecksumVerify(sentence []byte, ignoreFirst bool) bool {
	if len(sentence) < 3 || (ignoreFirst && len(sentence) < 4) {
		return false
	}

	cs, ok := nmeaChecksumRead(sentence[len(sentence)-2:])
	if !ok {
		return false
	}

	for i := 0; i < len(sentence)-3; i++ {
		cs ^= sentence[i]
	}

	if ignoreFirst {
		cs ^= sentence[0]
	}

	return cs == 0
}

var decodeTable [256]uint8

func fillDecodeTable() {
	for i := range decodeTable {
		decodeTable[i] = 0xff
	}
	for i := 0; i < 40; i++ {
		decodeTable[i+48] = uint8(i)
	}
	for i := 40; i < 64; i++ {
		decodeTable[i-40+96] = uint8(i)
	}
}

func encapsulatedToUint64(in []byte, out []uint64) int {
	cnt := 64
	ind64 := 0

	for _, v := range in {
		d64 := uint64(decodeTable[v])
		if d64 == 0xff {
			return 0
		}

		if cnt -= 6; cnt < 0 {
			out[ind64] |= d64 >> -cnt
			cnt += 64
			ind64++
			if ind64 >= len(out) {
				return 0
			}
		}

		out[ind64] |= d64 << cnt
	}

	return len(in) * 6
}

func (c *Decoder) handleMessage() error {
	if !c.cfg.IgnoreChecksum {
		if !nmeaChecksumVerify(c.dataBlock, true) {
			return nil
		}

		if c.tagBlockValid != nil {
			if !nmeaChecksumVerify(c.tagBlockValid, false) {
				return nil
			}
		}
	}

	nmeaParsed := NMEAParsed{
		Tagblock: c.tagBlockValid,
		Sentence: c.dataBlock,
	}

	/* Let other code handle normal NMEA messages */
	if c.dataBlock[0] == '$' {
		if c.cfg.NMEAFunc != nil {
			if err := c.cfg.NMEAFunc(nmeaParsed); err != nil {
				return err
			}
		}
		return nil
	}

	if c.cfg.NMEAEncapsulatedFunc != nil {
		if err := c.cfg.NMEAEncapsulatedFunc(nmeaParsed); err != nil {
			return err
		}
	}

	return c.decodeAISSentence(nmeaParsed)
}

func (c *Decoder) decodeAISSentence(nmeaParsed NMEAParsed) error {
	if c.cfg.AIS == nil || c.cfg.AISDecodedFunc == nil || len(nmeaParsed.Sentence) < 14+1+3 {
		return nil
	}

	msg := nmeaParsed.Sentence[1 : len(nmeaParsed.Sentence)-3]

	//msgType, rest := splitComma(nmeaParsed.Sentence[1:])
	if msg[5] != ',' || msg[2] != 'V' || msg[3] != 'D' || (msg[4] != 'M' && msg[4] != 'O') {
		return nil
	}
	talker := msg[:2]

	msgTotal, ok := nmeaReadNibble(msg[6])
	if !ok || msgTotal == 0 {
		return nil
	}

	msgIndex, ok := nmeaReadNibble(msg[8])
	if !ok || msgIndex > msgTotal {
		return nil
	}

	/* Many sources send things besides just 0-9 here */
	msgIDValid := msg[10] != ','
	msgID := msg[10]
	if !msgIDValid && msgTotal != 1 {
		return nil
	}

	msg = msg[11:]
	if msgIDValid {
		msg = msg[1:]
	}

	channel := msg[0]
	if channel == ',' {
		msg = msg[1:]
		channel = 0
	} else {
		msg = msg[2:]
		if channel == '2' || channel == 'b' || channel == 'B' || channel == '+' || channel == 'H' || channel == 'h' {
			channel = 2
		} else {
			channel = 1
		}
	}

	if len(msg) < 2 {
		return nil
	}

	data := msg[:len(msg)-2]

	/* Need to combine messages? */
	if msgTotal > 1 {
		data, nmeaParsed.Tagblock = c.recombineMessages(msgID, int(msgTotal), int(msgIndex), nmeaParsed.Tagblock, data)
		if data == nil {
			return nil
		}
	}

	padding, ok := nmeaReadNibble(msg[len(msg)-1])
	if !ok {
		return nil
	}

	/* This is enough to fit any valid AIS message with some extra (20 should be enough) */
	var out [32]uint64
	numBits := encapsulatedToUint64(data, out[:])
	numBits -= int(padding)

	if numBits <= 0 {
		return nil
	}

	aisParsed := AISParsed{
		Talker:  talker,
		Channel: channel,
		Packet:  c.cfg.AIS.DecodePacket64(out[:], numBits),
	}

	if aisParsed.Packet == nil {
		return nil
	}

	return c.cfg.AISDecodedFunc(nmeaParsed, aisParsed)
}

func (s *recombState) reset() {
	s.MsgTotal = 0
	s.TagBlock = s.TagBlock[:0]
	for i := range s.Data {
		s.Data[i] = s.Data[i][:0]
	}
}

func splitComma(in []byte) ([]byte, []byte) {
	for i, m := range in {
		if m == ',' {
			return in[:i], in[i+1:]
		}
	}

	return in, nil
}

func tagblockGetBlockID(key []byte, tagblock []byte) int {
	for len(tagblock) > 0 {
		var info []byte
		info, tagblock = splitComma(tagblock)

		if len(info) < 2 || (info[0] != 'g' && info[0] != 'G') || info[1] != ':' {
			continue
		}

		if len(info) >= 3 && info[len(info)-3] == '*' {
			info = info[:len(info)-3]
		}

		if len(info) == 0 {
			continue
		}

		for i := len(info) - 1; i >= 0; i-- {
			if info[i] == '-' {
				info = info[i+1:]
				break
			}
		}

		return copy(key[:], info)
	}

	return -1
}

func (c *Decoder) recombineMessages(msgID byte, msgTotal int, msgIndex int, tagBlock []byte, data []byte) ([]byte, []byte) {
	msgIndex--

	var key [8]byte
	key[0] = msgID
	if tagblockGetBlockID(key[1:6], tagBlock) < 0 {
		key[7] = 1
	}

	state, ok := c.aisRecombine[key]
	if !ok {
		state = &recombState{}
		c.aisRecombine[key] = state
	}

	if msgTotal >= len(state.Data) || msgIndex >= len(state.Data) || msgIndex < 0 {
		return nil, nil
	}

	if state.MsgTotal != msgTotal {
		state.reset()
	}

	state.MsgTotal = msgTotal

	if len(tagBlock) > len(state.TagBlock) {
		state.TagBlock = append(state.TagBlock[:0], tagBlock...)
	}

	state.Data[msgIndex] = append(state.Data[msgIndex][:0], data...)

	/* Do we have everything? */
	totalLen := 0
	for i := 0; i < msgTotal; i++ {
		if l := len(state.Data[i]); l == 0 {
			return nil, nil
		} else {
			totalLen += l
		}
	}

	out := make([]byte, 0, totalLen)
	for i := 0; i < msgTotal; i++ {
		out = append(out, state.Data[i]...)
	}

	var tb []byte
	if len(state.TagBlock) > 0 {
		tb = append([]byte{}, state.TagBlock...)
	}
	state.reset()

	return out, tb
}
