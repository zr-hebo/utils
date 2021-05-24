package hash

import (
	"encoding/binary"
	"strings"
	"testing"
)

func TestCityHashFromString(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name     string
		args     args
		wantHash uint64
	}{
		{
			name: "one",
			args: args{
				val: "B7A7F69D-0B21-408F-870F-8C99C280797D",
			},
			wantHash: 363213291030809414,
		},
		{
			name: "two",
			args: args{
				val: "ANDROID_a88d345dfdaa1825",
			},
			wantHash: 12149212146700503755,
		},
		{
			name: "three",
			args: args{
				val: "B7A7F69D-0B21-408F-870F-8C99C280797D",
			},
			wantHash: 363213291030809414,
		},
		{
			name: "four",
			args: args{
				val: "ANDROID_0972a26a2f9e4284",
			},
			wantHash: 15982774355910090593,
		},
		{
			name: "five",
			args: args{
				val: "A0EE275A-DCA3-4F8E-BCD6-ACD1F0DC33BD",
			},
			wantHash: 12016482007678147712,
		},
		{
			name: "five",
			args: args{
				val: "A0EE275A-DCA3-4F8E-BCD6-ACD1F0DC33BD",
			},
			wantHash: 12016482007678147712,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHash := CityHashFromString(strings.ToLower(tt.args.val)); gotHash != tt.wantHash {
				t.Errorf("CityHash() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}


func TestCityHashFromInt(t *testing.T) {
	type args struct {
		val []byte
	}

	one := make([]byte, 8)
	binary.LittleEndian.PutUint64(one, 4)
	two := make([]byte, 8)
	binary.LittleEndian.PutUint64(two, 666)
	three := make([]byte, 8)
	binary.LittleEndian.PutUint64(three, 43359)
	four := make([]byte, 8)
	binary.LittleEndian.PutUint64(four, 5837457283055)
	five := make([]byte, 8)
	binary.LittleEndian.PutUint64(five, 98168135770576)
	six := make([]byte, 8)
	binary.LittleEndian.PutUint64(six, 6811194006946514881)
	tests := []struct {
		name     string
		args     args
		wantHash uint64
	}{
		{
			name: "one",
			args: args{
				val: one,
			},
			wantHash: 3255232038643208583,
		},
		{
			name: "two",
			args: args{
				val: two,
			},
			wantHash: 15703263142150285816,
		},
		{
			name: "three",
			args: args{
				val: three,
			},
			wantHash: 17303161218879695284,
		},
		{
			name: "four",
			args: args{
				val: four,
			},
			wantHash: 14762139812735696189,
		},
		{
			name: "five",
			args: args{
				val: five,
			},
			wantHash: 2767479816985530557,
		},
		{
			name: "six",
			args: args{
				val: six,
			},
			wantHash: 16014426317825635999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHash := CityHash(tt.args.val); gotHash != tt.wantHash {
				t.Errorf("CityHash() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}
