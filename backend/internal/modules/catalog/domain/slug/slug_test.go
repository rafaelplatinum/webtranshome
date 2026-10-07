package slug

import "testing"

func TestIsValid(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		maxLength int
		want      bool
	}{
		{name: "lowercase slug", value: "living-room", maxLength: 20, want: true},
		{name: "rejects empty", maxLength: 20},
		{name: "rejects uppercase", value: "Living", maxLength: 20},
		{name: "rejects repeated hyphen", value: "living--room", maxLength: 20},
		{name: "rejects trailing hyphen", value: "living-", maxLength: 20},
		{name: "rejects length over limit", value: "living-room", maxLength: 10},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := IsValid(test.value, test.maxLength); got != test.want {
				t.Fatalf("IsValid() = %t, want %t", got, test.want)
			}
		})
	}
}
