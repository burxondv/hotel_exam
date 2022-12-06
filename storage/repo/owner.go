package repo

import "time"

type Owner struct {
	ID        int64
	HotelID   int64
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type GetAllOwnersParams struct {
	Limit  int32
	Page   int32
	Search string
}

type GetAllOwnersResult struct {
	Owners []*Owner
	Count  int32
}

type UpdateOwnerPassword struct {
	OwnerID  int64
	Password string
}

type OwnerStorageI interface {
	Create(ow *Owner) (*Owner, error)
	Get(id int64) (*Owner, error)
	GetByEmail(email string) (*Owner, error)
	GetAll(params *GetAllOwnersParams) (*GetAllOwnersResult, error)
	UpdateOwnerPassword(req *UpdateOwnerPassword) error
	Update(ow *Owner) (*Owner, error)
	Delete(id int64) error
}
