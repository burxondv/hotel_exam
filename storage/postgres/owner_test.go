package postgres_test

import (
	"testing"

	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createOwner(t *testing.T) *repo.Owner {
	owner, err := strg.Owner().Create(&repo.Owner{
		HotelID: 1,
		FirstName: faker.Name(),
		LastName: faker.Name(),
		Email: faker.Email(),
		Password: faker.Password(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, owner)

	return owner
}

func deleteOwner(id int64, t *testing.T) {
	err := strg.Owner().Delete(id)
	require.NoError(t, err)
}

func TestCreateOwner(t *testing.T) {
	createOwner(t)
}

func TestGetOwner(t *testing.T) {
	c := createOwner(t)

	owner, err := strg.Owner().Get(c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, owner)
}

func TestGetAllOwner(t *testing.T) {
	user := createOwner(t)

	users, err := strg.Owner().GetAll(&repo.GetAllOwnersParams{
		Limit: 10,
		Page: 1,
	})

    require.NoError(t, err)
    require.GreaterOrEqual(t, len(users.Owners), 1)

	deleteOwner(user.ID, t)
}

func TestUpdateOwner(t *testing.T) {
	owner := createOwner(t)
	owner.HotelID = 1
	owner.FirstName = faker.Name()
	owner.LastName = faker.Name()
	owner.Email = faker.Email()
	owner.Password = faker.Password()
}

func TestDeleteOwner(t *testing.T) {
	owner := createOwner(t)

	deleteOwner(owner.ID, t)
}
