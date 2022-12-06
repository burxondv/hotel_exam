package models

import "time"

type Owner struct {
	ID        int64     `json:"id"`
	HotelID   int64     `json:"hotel_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateOwnerRequest struct {
	HotelID   int64  `json:"hotel_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type GetAllOwnersResponse struct {
	Owners []*Owner `json:"owners"`
	Count  int32    `json:"count"`
}
