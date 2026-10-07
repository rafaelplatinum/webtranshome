package room

import (
	"errors"
	"strings"
	"testing"
)

func TestNewRoom(t *testing.T) {
	tests := []struct {
		name       string
		roomName   string
		slug       string
		imageCover string
		sortOrder  int64
		wantErr    error
	}{
		{name: "trims values", roomName: " Living Room ", slug: "living-room", imageCover: " /cover.png "},
		{name: "rejects empty name", slug: "living-room", imageCover: "/cover.png", wantErr: ErrInvalidRoom},
		{name: "rejects missing cover", roomName: "Living Room", slug: "living-room", wantErr: ErrInvalidRoom},
		{name: "rejects invalid slug", roomName: "Living Room", slug: "Living Room", imageCover: "/cover.png", wantErr: ErrInvalidRoom},
		{name: "rejects negative order", roomName: "Living Room", slug: "living-room", imageCover: "/cover.png", sortOrder: -1, wantErr: ErrInvalidRoom},
		{name: "rejects order beyond database range", roomName: "Living Room", slug: "living-room", imageCover: "/cover.png", sortOrder: maxSortOrder + 1, wantErr: ErrInvalidRoom},
		{name: "rejects image beyond database range", roomName: "Living Room", slug: "living-room", imageCover: strings.Repeat("a", 256), wantErr: ErrInvalidRoom},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := New(test.roomName, test.slug, test.imageCover, test.sortOrder)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("New() error = %v, want %v", err, test.wantErr)
			}
			if test.wantErr == nil && (got.Name != "Living Room" || got.ImageCover != "/cover.png" || !got.IsActive) {
				t.Fatalf("New() room = %+v", got)
			}
		})
	}
}

func TestRoomApply(t *testing.T) {
	inactive := false
	emptyCover := ""
	tests := []struct {
		name    string
		changes Changes
		wantErr error
	}{
		{name: "deactivates room", changes: Changes{IsActive: &inactive}},
		{name: "rejects empty patch", wantErr: ErrInvalidRoom},
		{name: "rejects clearing required cover", changes: Changes{ImageCover: &emptyCover}, wantErr: ErrInvalidRoom},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item := Room{ID: 2, Name: "Living Room", Slug: "living-room", ImageCover: "/cover.png", IsActive: true}
			err := item.Apply(test.changes)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Apply() error = %v, want %v", err, test.wantErr)
			}
			if test.name == "deactivates room" && item.IsActive {
				t.Fatal("Apply() did not deactivate room")
			}
			if test.name == "rejects clearing required cover" && item.ImageCover != "/cover.png" {
				t.Fatal("Apply() changed room after invalid patch")
			}
		})
	}
}
