package storage

import (
	"github.com/burxondv/hotel_exam/storage/postgres"

	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	Guest() repo.GuestStorageI
	Owner() repo.OwnerStorageI
	Hotel() repo.HotelStorageI
	Room() repo.RoomStorageI
}

type storagePg struct {
	guestRepo repo.GuestStorageI
	ownerRepo repo.OwnerStorageI
	hotelRepo repo.HotelStorageI
	roomRepo  repo.RoomStorageI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &storagePg{
		guestRepo: postgres.NewGuest(db),
		ownerRepo: postgres.NewOwner(db),
		hotelRepo: postgres.NewHotel(db),
		roomRepo:  postgres.NewRoom(db),
	}
}

func (s *storagePg) Guest() repo.GuestStorageI {
	return s.guestRepo
}

func (s *storagePg) Owner() repo.OwnerStorageI {
	return s.ownerRepo
}

func (s *storagePg) Hotel() repo.HotelStorageI {
	return s.hotelRepo
}

func (s *storagePg) Room() repo.RoomStorageI {
	return s.roomRepo
}