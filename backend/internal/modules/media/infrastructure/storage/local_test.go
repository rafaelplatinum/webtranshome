package storage

import "testing"

func TestWebpDimensionsSupportsAllWebpBitstreams(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		wantW int
		wantH int
	}{
		{
			name:  "extended",
			data:  []byte{'R', 'I', 'F', 'F', 0, 0, 0, 0, 'W', 'E', 'B', 'P', 'V', 'P', '8', 'X', 0, 0, 0, 0, 0, 0, 0, 0, 99, 0, 0, 49, 0, 0},
			wantW: 100,
			wantH: 50,
		},
		{
			name:  "lossless",
			data:  []byte{'R', 'I', 'F', 'F', 0, 0, 0, 0, 'W', 'E', 'B', 'P', 'V', 'P', '8', 'L', 0, 0, 0, 0, 0x2f, 99, 64, 12, 0},
			wantW: 100,
			wantH: 50,
		},
		{
			name:  "lossy",
			data:  []byte{'R', 'I', 'F', 'F', 0, 0, 0, 0, 'W', 'E', 'B', 'P', 'V', 'P', '8', ' ', 0, 0, 0, 0, 0, 0, 0, 0x9d, 0x01, 0x2a, 100, 0, 50, 0},
			wantW: 100,
			wantH: 50,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			width, height, ok := webpDimensions(test.data)
			if !ok || width != test.wantW || height != test.wantH {
				t.Fatalf("webpDimensions() = (%d, %d, %v), want (%d, %d, true)", width, height, ok, test.wantW, test.wantH)
			}
		})
	}
}
