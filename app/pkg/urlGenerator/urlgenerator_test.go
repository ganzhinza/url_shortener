package urlgenerator

import (
	"testing"
)

func TestGenerate_LenAndSymbols(t *testing.T) {
	type args struct {
		size uint
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test 1",
			args: args{10},
			want: 10,
		},
		{
			name: "Test 2",
			args: args{7},
			want: 7,
		},
		{
			name: "Test 3",
			args: args{0},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Generate(tt.args.size)
			if len(got) != tt.want {
				t.Errorf("Generate len = %v, want %v", len(got), tt.want)
			}

			for i := range got {
				exists := false
				for j := range letters {
					if got[i] == letters[j] {
						exists = true
						break
					}
				}
				if !exists {
					t.Errorf("Unexpected symbol: %s", string(got[i]))
				}
			}
		})
	}
}
