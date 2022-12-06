package models

import "time"

type RegisterOwnerRequest struct {
	HotelID   int64  `json:"hotel_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type AuthOwnerResponse struct {
	ID          int64     `json:"id"`
	HotelID     int64     `json:"hotel_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	AccessToken string    `json:"access_token"`
}