package ais

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestPositionParse(t *testing.T) {
	fastParse := CodecNew(false, false, true)
	/* Convenience conversion disabled to avoid float inaccuracies */
	fastParse.FloatWithoutConversion = true

	slowParse := CodecNew(false, false, false)
	slowParse.FloatWithoutConversion = true

	for msgID := 1; msgID <= 27; msgID++ {
		msgFile := fmt.Sprintf("testmsg/%d.msg", msgID)
		t.Logf("Testing message file %s", msgFile)
		f, err := os.Open(msgFile)
		if err != nil {
			t.Error("Failed to open file", msgID)
			return
		}

		r := bufio.NewReader(f)
		for index := 0; true; index++ {

			line, err := r.ReadString('\n')
			if err != nil {
				break
			}

			line = line[:len(line)-2]
			source := []byte(line)

			/* Convert ascii '0' and '1' to real 0 and 1 */
			for i := 0; i < len(source); i++ {
				source[i] -= '0'
			}
			/* Decode the packet */
			decoded := fastParse.DecodePacket(source)
			if decoded == nil {
				/* Failed to decode... */
				t.Error("Could not decode with fastParse", msgID, index)
				return
			}

			decodedSlow := slowParse.DecodePacket(source)
			if decodedSlow == nil {
				/* Failed to decode... */
				t.Error("Could not decode with reflection", msgID, index)
				return
			}

			if !reflect.DeepEqual(decoded, decodedSlow) {
				t.Error("Decoded messages are not equal", msgID, index)
				t.Log("FastParse/slowparse:\n", decoded, "\n", decodedSlow)
				t.Logf("error on file msgID %d at index %d", msgID, index)
			}
		}
		f.Close()
		// also run another the encoding/decoding
		tryFile(t, fastParse, msgID)
	}
}
