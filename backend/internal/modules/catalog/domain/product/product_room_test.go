package product

import (
	"errors"
	"testing"
)

func TestNewRoomAssignment(t *testing.T) {
	assignment, err := NewRoomAssignment(7, []int64{1, 2, 3})
	if err != nil {
		t.Fatalf("NewRoomAssignment() error = %v", err)
	}
	if assignment.ProductID != 7 || len(assignment.RoomIDs) != 3 {
		t.Fatalf("NewRoomAssignment() = %+v", assignment)
	}
}

func TestNewRoomAssignmentRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name string
		id   int64
		rows []int64
		want error
	}{
		{name: "invalid product", id: 0, rows: []int64{1}, want: ErrInvalidProductID},
		{name: "invalid room id", id: 7, rows: []int64{0}, want: ErrInvalidRoomID},
		{name: "duplicate room", id: 7, rows: []int64{1, 1}, want: ErrInvalidRoomID},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewRoomAssignment(test.id, test.rows)
			if !errors.Is(err, test.want) {
				t.Fatalf("NewRoomAssignment() error = %v, want %v", err, test.want)
			}
		})
	}
}
