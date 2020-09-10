package aisnmea

import "testing"
import nmea "github.com/adrianmo/go-nmea"

func Test_encodeTagBlock(t *testing.T) {
	type args struct {
		tagBlock *nmea.TagBlock
		msgIndex int
		msgNum   int
		seqNo    int
		addLineCount bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name:"nil input tagblock", args:args{
			tagBlock: nil, msgIndex: 0, msgNum: 0, seqNo: 0, addLineCount: true}, want:""},
		{name:"blank input tagblock", args:args{
			tagBlock: &nmea.TagBlock{}, msgIndex: 0, msgNum: 0, seqNo: 0, addLineCount: true}, want:""},
		{name:"single sentence message",
			args:args{
				tagBlock: &nmea.TagBlock{
					Time:         1,
				},
				msgIndex: 1, msgNum: 1, seqNo: 0, addLineCount: true,
			},
			want:"\\c:1*68\\",
		},
		{name:"multi sentence message, sentence 1",
			args:args{
				tagBlock: &nmea.TagBlock{
					Time:         1,
				},
				msgIndex: 1, msgNum: 2, seqNo: 0, addLineCount: true,
			},
			want:"\\g:1-2-0,n:2,c:1*60\\",
		},
		{name:"multi sentence message, sentence 2",
			args:args{
				tagBlock: &nmea.TagBlock{
					Time:         1,
				},
				msgIndex: 2, msgNum: 2, seqNo: 0, addLineCount: true,
			},
			want:"\\g:2-2-0*6D\\",
		},
		{name:"correct checksum time source and grouping",
			args:args{
				tagBlock: &nmea.TagBlock{
					Time:         1560234814,
					Source:       "2251",
					Grouping:     "1-2-2449555",
				},
				msgIndex: 1, msgNum: 2, seqNo: 2449555, addLineCount: false,
			},
			want:"\\g:1-2-2449555,s:2251,c:1560234814*7E\\",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeTagBlock(
				tt.args.tagBlock, tt.args.msgIndex, tt.args.msgNum, tt.args.seqNo, tt.args.addLineCount);
			got != tt.want {
				t.Errorf("encodeTagBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeTagBlocks(t *testing.T) {
	type args struct {
		dst *nmea.TagBlock
		src *nmea.TagBlock
	}
	tests := []struct {
		name string
		args args
		ok func(t *testing.T, a args)
	}{
		{name:"null src", args:args{dst:&nmea.TagBlock{Text: "foo"}, src:nil}, ok:func(t *testing.T, a args) {
			same := nmea.TagBlock{Text: "foo"}
			if *a.dst != same {
				t.Error("null src is expected to leave dst unchanged")
			}
		}},
		{name:"null dst", args:args{src:&nmea.TagBlock{Text: "foo"}, dst:nil}, ok:func(t *testing.T, a args) {
			if a.dst != nil {
				t.Error("null dst is expected to remain unchanged")
			}
		}},
		{name:"zero src time", args:args{src:&nmea.TagBlock{Time: 0}, dst:&nmea.TagBlock{Time: 1}},
			ok:func(t *testing.T, a args) {
				if a.dst.Time != 1 {
					t.Error("zero src time should not overwrite dst time")
				}
		}},
		{name:"non-zero dst time", args:args{src:&nmea.TagBlock{Time: 2}, dst:&nmea.TagBlock{Time: 1}},
			ok:func(t *testing.T, a args) {
				if a.dst.Time != 1 {
					t.Error("non-zero dst time should remain unchanged")
				}
		}},
		{name:"orthogonal tags", args:args{src:&nmea.TagBlock{Source: "outer_space"}, dst:&nmea.TagBlock{Time: 1}},
			ok:func(t *testing.T, a args) {
				if a.dst.Time != 1 {
					t.Error("time not set")
				}
				if a.dst.Source != "outer_space" {
					t.Error("source not set")
				}
		}},
		{name:"ignore grouping tag", args:args{src:&nmea.TagBlock{Grouping: "1-2-99"}, dst:&nmea.TagBlock{Time: 1}},
			ok:func(t *testing.T, a args) {
				if a.dst.Time != 1 {
					t.Error("dst time should remain unchanged")
				}
				if a.dst.Grouping != "" {
					t.Error("dst grouping should never be set")
				}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergeTagBlocks(tt.args.dst, tt.args.src)
			tt.ok(t, tt.args)
		})
	}
}