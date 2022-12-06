package models

import "time"

type Room struct {
	ID        int64  `json:"id"`
	HotelID   int64  `json:"hotel_id"`
	Status    bool   `json:"status"`
	ImageUrl  string `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRoomRequest struct {
	HotelID   int64  `json:"hotel_id"`
    Status    bool   `json:"status" binding:"required" default:"false"`
    ImageUrl  string `json:"image_url"`
}

type GetAllRoomResponse struct {
	Rooms []*Room `json:"rooms"`
	Count int32   `json:"count"`
}
