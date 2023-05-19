package ais_test

import (
	"bufio"
	"compress/gzip"
	"github.com/BertoldVdb/go-ais"
	"github.com/BertoldVdb/go-ais/aisnmea"
	"io"
	"os"
	"strings"
	"testing"
)

func runTest(b *testing.B, parseFast bool) {
	x := ais.CodecNew(false, false, parseFast)
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

func readFileTest(b *testing.B, benchFile string, parseFast bool) {
	c := ais.CodecNew(false, false, parseFast)
	nm := aisnmea.NMEACodecNew(c)
	var reader io.Reader
	fp, err := os.Open(benchFile)
	if err != nil {
		b.Fatal("could not open ", benchFile, "\nerr:\n", err)
	}
	stat, err := fp.Stat()
	b.SetBytes(stat.Size())
	defer fp.Close()

	if strings.Contains(benchFile, ".gz") {
		reader, err = gzip.NewReader(fp)
		if err != nil {
			b.Fatal("could not create gzip reader", err)
		}
	} else {
		reader = fp
	}

	// create line by line scanner, read each line and decode it
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Bytes()
		_, err := nm.ParseSentence(string(line))
		if err != nil {
			return
		}
	}
	if err := scanner.Err(); err != nil {
		b.Fatal("error reading file", err)
	}

}
func BenchmarkDecodeFile(b *testing.B) {
	testFile := os.Getenv("NMEA_BENCH_FILE")
	if testFile == "" {
		b.Skip("NMEA_BENCH_FILE not set")
	}
	b.Run("prefer reflection", func(b *testing.B) {
		readFileTest(b, testFile, false)
	})
	b.Run("do not prefer reflection", func(b *testing.B) {
		readFileTest(b, testFile, true)
	})
}
