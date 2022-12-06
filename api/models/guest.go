package models

import "time"

type Guest struct {
	ID          int64  `json:"id"`
	HotelID     int64  `json:"hotel_id"`
	RoomID      int64  `json:"room_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CreateGuestRequest struct {
	HotelID     int64 `json:"hotel_id"`
	RoomID      int64 `json:"room_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type GetAllGuestsResponse struct {
	Guests []*Guest `json:"guests"`
	Count  int32    `json:"count"`
}
