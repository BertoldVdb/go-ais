package ais

import "testing"

func TestChannelToFrequency(t *testing.T) {
	values := []struct {
		channel uint16
		freq    uint
	}{
		{2087, 161975000}, /* The standard AIS1 channel */
		{2088, 162025000}, /* The standard AIS2 channel */
		{1087, 157375000}, /* AIS1 should it use ship instead of shore frequencies */
		{87, 0},           /* AIS1 should it be a simplex frequency (invalid) */
		{11, 156550000},   /* Some random simplex channels */
		{211, 156562500},
		{270, 156537500},
		{71, 156575000},
		{11111, 0}, /* A channel that does not exist */
	}

	x := CodecNew(false, false)

	for _, v := range values {
		myFreq := x.ChannelToFrequency(v.channel)
		if myFreq != v.freq {
			t.Error("Failed to calculate channel frequency", v.channel, v.freq, myFreq)
		}
	}
}
