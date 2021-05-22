package hash

import "testing"

func TestCityHash(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name     string
		args     args
		wantHash uint64
	}{
		{
			name:     "two",
			args:     args{
				val: "B7A7F69D-0B21-408F-870F-8C99C280797D",
			},
			wantHash: 363213291030809414,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHash := CityHashFromString(tt.args.val); gotHash != tt.wantHash {
				t.Errorf("CityHash() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}

