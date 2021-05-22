package hash

import (
	"testing"
)

func TestJavaHash(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name     string
		args     args
		wantHash int32
	}{
		{
			name:     "empty string",
			args:     args{val: ""},
			wantHash: 0,
		},
		{
			name:     "hello world",
			args:     args{val: "Hello, world!"},
			wantHash: -1880044555,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHash := JavaHashFromString(tt.args.val); gotHash != tt.wantHash {
				t.Errorf("JavaHash() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}

func TestInt32Abs(t *testing.T) {
	type args struct {
		val int32
	}
	tests := []struct {
		name    string
		args    args
		wantAbs int64
	}{
		{
			name:    "zero",
			args:    args{val: 0},
			wantAbs: 0,
		},
		{
			name:    "one",
			args:    args{val: -1},
			wantAbs: 1,
		},
		{
			name:    "min int32",
			args:    args{val: -2147483648},
			wantAbs: 2147483648,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAbs := Int32Abs(tt.args.val); gotAbs != tt.wantAbs {
				t.Errorf("Int32Abs() = %v, want %v", gotAbs, tt.wantAbs)
			}
		})
	}
}

