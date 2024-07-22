package ais

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
)

func fuzzToBinary(pkt []byte) []byte {
	out := make([]byte, 1)
	outIndex := 1

	for i, m := range pkt {
		if outIndex >= len(out) {
			out = append(out, 0)
		}

		if m == 1 {
			out[outIndex] |= (1 << (i % 8))
		}

		if i%8 == 7 {
			outIndex++
		}
	}

	/* Store the amount of unused bits */
	out[0] = byte((len(out)-1)*8 - len(pkt))

	return out
}

func fuzzFromBinary(bin []byte) []byte {
	if len(bin) < 1 || bin[0] >= 8 {
		return nil
	}

	out := make([]byte, (len(bin)-1)*8-int(bin[0]))

	binIndex := 1
	for i := range out {
		out[i] = (bin[binIndex] >> (i % 8)) & 1

		if i%8 == 7 {
			binIndex++
		}
	}

	return out
}

func TestFuzzConvert(t *testing.T) {
	for i := 0; i < 1024; i++ {
		in := make([]byte, i)

		for j := range in {
			in[j] = byte(rand.Intn(2))
		}

		pkt := fuzzToBinary(in)
		out := fuzzFromBinary(pkt)

		if !bytes.Equal(in, out) {
			t.Error("Value different for", i)
		}
	}
}

func FuzzDecode(f *testing.F) {
	for msgID := 1; msgID <= 27; msgID++ {
		msgFile := fmt.Sprintf("testmsg/%d.msg", msgID)
		file, err := os.Open(msgFile)
		if err != nil {
			f.Error("Failed to open file", msgID)
			return
		}

		r := bufio.NewReader(file)
		for index := 0; index < 10; index++ {

			line, err := r.ReadString('\n')
			if err != nil {
				break
			}

			line = line[:len(line)-2]
			source := []byte(line)

			for i := range source {
				source[i] -= '0'
			}

			if index == 0 {
				for i := range source {
					fuzzBin := fuzzToBinary(source[:i])
					f.Add(fuzzBin)
				}

				for i := range source {
					source[i] = 1 - source[i]
					fuzzBin := fuzzToBinary(source)
					source[i] = 1 - source[i]
					f.Add(fuzzBin)
				}
			} else {
				fuzzBin := fuzzToBinary(source)
				f.Add(fuzzBin)
			}
		}

		file.Close()
	}

	fastParse := CodecNewFast(false, false, true)
	fastParse.FloatWithoutConversion = true

	slowParse := CodecNewFast(false, false, false)
	slowParse.FloatWithoutConversion = true

	f.Fuzz(func(t *testing.T, msg []byte) {
		decoded := fastParse.DecodePacket(msg)
		decodedSlow := slowParse.DecodePacket(msg)

		if decoded == nil && decodedSlow != nil {
			t.Error("Only slow decoder decoded packet")
			t.Logf("%T", decodedSlow)
			t.Logf("%+v", decodedSlow)
		} else if decoded != nil && decodedSlow == nil {
			t.Error("Only fast decoder decoded packet")
			t.Logf("%T", decoded)
			t.Logf("%+v", decoded)
		} else if !reflect.DeepEqual(decoded, decodedSlow) {
			t.Logf("%T %T", decodedSlow, decoded)
			t.Logf("%+v", decodedSlow)
			t.Logf("%+v", decoded)
			t.Error("Decoded messages are not equal")
		}
	})
}
