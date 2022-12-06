package repo

import "time"

type Hotel struct {
	ID           int64
	Name         string
	Address      string
	StarsPopular int16
	ImageUrl     string
	CreatedAt    time.Time
}

type GetAllHotelsParams struct {
	Limit  int32
	Page   int32
	Search string
}

type GetAllHotelsResult struct {
	Hotels []*Hotel
	Count  int32
}

type HotelStorageI interface {
	Create(ht *Hotel) (*Hotel, error)
	Get(id int64) (*Hotel, error)
	GetAll(params *GetAllHotelsParams) (*GetAllHotelsResult, error)
	Update(ht *Hotel) (*Hotel, error)
	Delete(id int64) error
}
