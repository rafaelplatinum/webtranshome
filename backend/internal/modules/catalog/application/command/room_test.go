package command

import (
	"context"
	"errors"
	"testing"

	"webtranshome/internal/modules/catalog/domain/room"
)

type roomRepositoryStub struct {
	item        room.Room
	findErr     error
	createErr   error
	updateErr   error
	createCalls int
	updateCalls int
}

func (r *roomRepositoryStub) ListRooms(context.Context) ([]room.Room, error) {
	return nil, nil
}

func (r *roomRepositoryStub) FindRoom(context.Context, int64) (room.Room, error) {
	return r.item, r.findErr
}

func (r *roomRepositoryStub) CreateRoom(_ context.Context, item room.Room) (room.Room, error) {
	r.createCalls++
	item.ID = 11
	return item, r.createErr
}

func (r *roomRepositoryStub) UpdateRoom(_ context.Context, item room.Room) error {
	r.updateCalls++
	r.item = item
	return r.updateErr
}

func TestCreateRoom(t *testing.T) {
	tests := []struct {
		name      string
		input     CreateRoomInput
		createErr error
		wantErr   error
		wantCalls int
	}{
		{name: "creates room", input: CreateRoomInput{Name: " Living Room ", Slug: "living-room", ImageCover: "/cover.png"}, wantCalls: 1},
		{name: "rejects invalid room", input: CreateRoomInput{Name: "Living Room", Slug: "living-room"}, wantErr: room.ErrInvalidRoom},
		{name: "propagates repository error", input: CreateRoomInput{Name: "Living Room", Slug: "living-room", ImageCover: "/cover.png"}, createErr: errors.New("database unavailable"), wantCalls: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &roomRepositoryStub{createErr: test.createErr}
			result, err := NewCreateRoom(repository).Execute(context.Background(), test.input)
			wantErr := test.wantErr
			if wantErr == nil {
				wantErr = test.createErr
			}
			if !errors.Is(err, wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, wantErr)
			}
			if repository.createCalls != test.wantCalls {
				t.Fatalf("CreateRoom() calls = %d, want %d", repository.createCalls, test.wantCalls)
			}
			if wantErr == nil && (result.ID != 11 || result.Name != "Living Room") {
				t.Fatalf("created room = %+v", result)
			}
		})
	}
}

func TestUpdateRoom(t *testing.T) {
	inactive := false
	emptyCover := ""
	tests := []struct {
		name    string
		id      int64
		changes UpdateRoomInput
		findErr error
		wantErr error
	}{
		{name: "deactivates room", id: 11, changes: UpdateRoomInput{IsActive: &inactive}},
		{name: "rejects invalid id", wantErr: ErrInvalidRoomID},
		{name: "propagates missing room", id: 11, findErr: room.ErrRoomNotFound, wantErr: room.ErrRoomNotFound},
		{name: "rejects clearing required cover", id: 11, changes: UpdateRoomInput{ImageCover: &emptyCover}, wantErr: room.ErrInvalidRoom},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &roomRepositoryStub{
				item:    room.Room{ID: 11, Name: "Living Room", Slug: "living-room", ImageCover: "/cover.png", IsActive: true},
				findErr: test.findErr,
			}
			test.changes.ID = test.id
			_, err := NewUpdateRoom(repository, repository).Execute(context.Background(), test.changes)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, test.wantErr)
			}
			if test.wantErr == nil && repository.updateCalls != 1 {
				t.Fatalf("UpdateRoom() calls = %d, want 1", repository.updateCalls)
			}
			if test.name == "deactivates room" && repository.item.IsActive {
				t.Fatal("UpdateRoom() did not deactivate room")
			}
		})
	}
}

var (
	_ room.Reader = (*roomRepositoryStub)(nil)
	_ room.Writer = (*roomRepositoryStub)(nil)
)
