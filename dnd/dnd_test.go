package dnd

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		want    Throw
		wantErr bool
	}{
		{"2d2", Throw{2, 2}, false},
		{"3d6", Throw{3, 6}, false},
		{"d6", Throw{1, 6}, false},
		{"-1d24", Throw{}, true},
		{"1d-1", Throw{}, true},
		{"dddddd", Throw{}, true},
		{"", Throw{}, true},
		{"bogus", Throw{}, true},
		{"1d1d1", Throw{}, true},
		{"1d1 d1", Throw{}, true},
		{"1d1 1", Throw{}, true},
	}
	for _, tt := range tests {
		result, err := Parse(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse(%q) = %v, want error: %t\n", tt.input, err, tt.wantErr)
		}
		if result != tt.want {
			t.Errorf("Parse(%q) = %v, want: %v\n", tt.input, result, tt.want)
		}
	}
}
