package repo

import "time"

type Guest struct {
	ID          int64
	HotelID     int64
	RoomID      int64
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Password    string
	CreatedAt   time.Time
}

type GetAllGuestsParams struct {
	Limit  int32
	Page   int32
	Search string
}

type GetAllGuestsResult struct {
	Guests []*Guest
	Count  int32
}

type UpdateGuestPassword struct {
	GuestID  int64
	Password string
}

type GuestStorageI interface {
	Create(gt *Guest) (*Guest, error)
	Get(id int64) (*Guest, error)
	GetByEmail(email string) (*Guest, error)
	GetAll(params *GetAllGuestsParams) (*GetAllGuestsResult, error)
	UpdateGuestPassword(req *UpdateGuestPassword) error
	Update(gt *Guest) (*Guest, error)
	Delete(id int64) error
}
