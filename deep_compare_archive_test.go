package ais_test

import (
	"bufio"
	"compress/gzip"
	"github.com/BertoldVdb/go-ais"
	"github.com/BertoldVdb/go-ais/aisnmea"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

func readAllPackets(t *testing.T, benchFile string, parseFast bool, cb func(eachPacket *aisnmea.VdmPacket)) {
	c := ais.CodecNewFast(false, false, parseFast)
	nm := aisnmea.NMEACodecNew(c)
	var reader io.Reader
	fp, err := os.Open(benchFile)
	if err != nil {
		t.Fatal("could not open ", benchFile, "\nerr:\n", err)
	}
	defer fp.Close()
	reader, err = gzip.NewReader(fp)
	if err != nil {
		t.Fatal()
	}

	// create line by line scanner, read each line and decode it
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Bytes()
		packet, err := nm.ParseSentence(string(line))
		if err != nil {
			t.Fatal("error parsing sentence", err)
		}
		cb(packet)
	}
	if err := scanner.Err(); err != nil {
		t.Fatal("error reading file", err)
	}

}

func TestDecodeFile(t *testing.T) {
	file := os.Getenv("NMEA_BENCH_FILE")
	if file == "" {
		t.Skip("NMEA_BENCH_FILE not set")
	}
	slowChan := make(chan *aisnmea.VdmPacket, 1)
	fastChan := make(chan *aisnmea.VdmPacket, 1)
	finished := false
	go func() {
		readAllPackets(t, file, false, func(eachPacket *aisnmea.VdmPacket) {
			slowChan <- eachPacket
		})
		finished = true
	}()
	go func() {
		readAllPackets(t, file, false, func(eachPacket *aisnmea.VdmPacket) {
			fastChan <- eachPacket
		})
		finished = true
	}()
	i := 0
	for {
		if finished {
			break
		}
		cancelTimer := time.AfterFunc(10*time.Second, func() {
			t.Fatal("Timeout should never happen")
			os.Exit(1)
		})
		leftValue := <-slowChan
		rightValue := <-fastChan
		if !reflect.DeepEqual(leftValue, rightValue) {
			t.Fatalf("Packet %d not equal (slow: %v, fast: %v)", i, leftValue, rightValue)
		}
		i += 1
		cancelTimer.Stop()
	}
}
