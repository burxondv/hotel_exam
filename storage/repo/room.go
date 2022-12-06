package repo

import "time"

type Room struct {
	ID        int64
	HotelID   int64
	Status    bool
	ImageUrl  string
	CreatedAt time.Time
}

type GetAllRoomsParams struct {
	Limit   int32
	Page    int32
	HotelID int32
}

type GetAllRoomsResult struct {
	Rooms []*Room
	Count int32
}

type RoomStorageI interface {
	Create(rm *Room) (*Room, error)
	Get(id int64) (*Room, error)
	GetAll(params *GetAllRoomsParams) (*GetAllRoomsResult, error)
	Update(rm *Room) (*Room, error)
	Delete(id int64) error
}
