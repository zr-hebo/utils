package strings

import "testing"

func TestBuildString(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal string",
			args: args{[]string{"haha", "hehe"}},
			want: "hahahehe",
		},
		{
			name: "string with number",
			args: args{[]string{"haha", "1", "1", "1", "1", "1"}},
			want: "haha11111",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildString(tt.args.args...); got != tt.want {
				t.Errorf("BuildString() = %v, want %v", got, tt.want)
			}
		})
	}
}
