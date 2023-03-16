package ais

import "testing"

func runTest(b *testing.B, parseFast bool) {
	x := CodecNew(false, false, parseFast)
	line := "000001000011100000000111101001010010000000000000000000000000011111110111110111001101001000100000111001110001011110001000101011000101000111011110000000001110101000001100"
	source := []byte(line)
	/* Convert ascii '0' and '1' to real 0 and 1 */
	for i := 0; i < len(source); i++ {
		source[i] -= '0'
	}
	for i := 0; i < b.N; i++ {
		x.DecodePacket(source)
	}
}
func BenchmarkDecode(b *testing.B) {
	b.Run("prefer reflection", func(b *testing.B) {
		runTest(b, false)
	})
	b.Run("do not prefer reflection", func(b *testing.B) {
		runTest(b, true)
	})
}
