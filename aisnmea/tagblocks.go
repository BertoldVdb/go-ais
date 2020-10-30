package aisnmea

import (
	"fmt"
	"strings"

	nmea "github.com/adrianmo/go-nmea"
)

// BlankTagBlock is a pseudo-constant for convenient comparisons to a blank value
var BlankTagBlock nmea.TagBlock

// addTagBlockChecksum performs the NMEA checksum from byte 0 (unlike addChecksum from nmea.go, which starts at 1)
func addTagBlockChecksum(sentence string) string {
	checksum := byte(0)
	for i := 0; i < len(sentence); i++ {
		checksum ^= sentence[i]
	}

	return fmt.Sprintf("%s*%02X", sentence, checksum)
}

/* 	mergeTagBlocks tries to summarise NMEA 4.10 TAG Blocks relating to a multi-sentence
	VDMVDO message (src) into one (dst) upon decoding. It is neccessarily lossy.

	Here are the assumptions:

	- Certain tags should be never have more than one value (time, source, etc), so we pick
      the first value, and ignore subsequent ones. If for some weird reason the sentences are
      timestamped instead of the message, we'd want to use the earliest timestamp.
	- The grouping tag should be different for every sentence, so there is no sense in including it.
	- Most examples in the wild seem to send the same tag block for every sentence of a message,
	  in those cases this function does no useful work, but no harm either.
	- If each sentence contributes a different subset of tags, we'll get a complete set at the end.
*/
func mergeTagBlocks(dst *nmea.TagBlock, src *nmea.TagBlock) {
	if dst == nil || src == nil {
		return
	}

	if src.Time != 0 && dst.Time == 0 {
		dst.Time = src.Time
	}

	if src.Text != "" && dst.Text == "" {
		dst.Text = src.Text
	}

	if src.Destination != "" && dst.Destination == "" {
		dst.Destination = src.Destination
	}

	if src.Source != "" && dst.Source == "" {
		dst.Source = src.Source
	}

	if src.RelativeTime != 0 && dst.RelativeTime == 0 {
		dst.RelativeTime = src.RelativeTime
	}

	if src.LineCount != 0 && dst.LineCount == 0 {
		dst.LineCount = src.LineCount
	}
}

// encodeTagBlock encodes the fields of tagBlock into a NMEA 4.10 TAG Block string
func encodeTagBlock(tagBlock *nmea.TagBlock, msgIndex, msgNum, seqNo int, addLineCount bool) string {
	if tagBlock == nil || *tagBlock == BlankTagBlock {
		return ""
	}

	tags := []string{}

	if msgNum > 1 {
		// We add the group tag for all sentences in a multi-sentence message
		tags = append(tags, fmt.Sprintf("g:%d-%d-%d", msgIndex, msgNum, seqNo))

		// We don't duplicate the rest of the tags for subsequent sentences of multi-sentence message
		if msgIndex > 1 {
			return fmt.Sprintf("\\%s\\", addTagBlockChecksum(tags[0]))
		}

		// If we got this far, it is a multi-sentence message AND this is the first sentence
		if addLineCount {
			tags = append(tags, fmt.Sprintf("n:%d", msgNum))
		}
	}

	if tagBlock.Source != "" {
		tags = append(tags, fmt.Sprintf("s:%s", tagBlock.Source))
	}

	if tagBlock.Time != 0 {
		tags = append(tags, fmt.Sprintf("c:%d", tagBlock.Time))
	}

	if tagBlock.RelativeTime != 0 {
		tags = append(tags, fmt.Sprintf("r:%d", tagBlock.RelativeTime))
	}

	if tagBlock.Destination != "" {
		tags = append(tags, fmt.Sprintf("d:%s", tagBlock.Destination))
	}

	if tagBlock.Text != "" {
		tags = append(tags, fmt.Sprintf("t:%s", tagBlock.Text))
	}

	return fmt.Sprintf("\\%s\\", addTagBlockChecksum(strings.Join(tags, ",")))
}
