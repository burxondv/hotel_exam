package models

import "time"

type Hotel struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	StarsPopular int16     `json:"stars_popular"`
	ImageUrl     string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateHotelRequest struct {
	Name         string `json:"name" binding:"required, min=2, max=40"`
	Address      string `json:"address"`
	StarsPopular int16  `json:"stars_popular" binding:"required, min=1, max=5"`
	ImageUrl     string `json:"image_url"`
}

type GetAllHotelResponse struct {
	Hotels []*Hotel `json:"hotels"`
	Count  int32    `json:"count"`
}
