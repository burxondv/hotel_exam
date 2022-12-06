package models

import "time"

type RegisterGuestRequest struct {
	HotelID   int64  `json:"hotel_id"`
	RoomID    int64  `json:"room_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type AuthGuestResponse struct {
	ID          int64     `json:"id"`
	HotelID     int64     `json:"hotel_id"`
	RoomID      int64     `json:"room_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	AccessToken string    `json:"access_token"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}
