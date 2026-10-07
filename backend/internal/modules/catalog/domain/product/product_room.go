package product

import "errors"

var ErrInvalidProductID = errors.New("CATALOG_INVALID_PRODUCT_ID")

type RoomAssignment struct {
	ProductID int64
	RoomIDs   []int64
}

func NewRoomAssignment(productID int64, roomIDs []int64) (RoomAssignment, error) {
	assignment := RoomAssignment{ProductID: productID, RoomIDs: append([]int64(nil), roomIDs...)}
	if err := assignment.Validate(); err != nil {
		return RoomAssignment{}, err
	}
	return assignment, nil
}

func (a RoomAssignment) Validate() error {
	if a.ProductID <= 0 {
		return ErrInvalidProductID
	}
	seen := make(map[int64]struct{}, len(a.RoomIDs))
	for _, roomID := range a.RoomIDs {
		if roomID <= 0 {
			return ErrInvalidRoomID
		}
		if _, exists := seen[roomID]; exists {
			return ErrInvalidRoomID
		}
		seen[roomID] = struct{}{}
	}
	return nil
}
