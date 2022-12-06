package postgres_test

import (
	"testing"

	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createGuest(t *testing.T) *repo.Guest {
	guest, err := strg.Guest().Create(&repo.Guest{
		HotelID:     1,
		RoomID:      1,
		FirstName:   faker.Name(),
		LastName:    faker.Name(),
		PhoneNumber: faker.Phonenumber(),
		Email:       faker.Email(),
		Password:    faker.Password(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, guest)

	return guest
}

func deleteGuest(id int64, t *testing.T) {
	err := strg.Guest().Delete(id)
	require.NoError(t, err)
}

func TestCreateGuest(t *testing.T) {
	createGuest(t)
}

func TestGetGuest(t *testing.T) {
	c := createGuest(t)

	guest, err := strg.Guest().Get(c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, guest)
}

func TestGetAllGuest(t *testing.T) {
	guest := createGuest(t)

	guests, err := strg.Guest().GetAll(&repo.GetAllGuestsParams{
		Limit: 10,
		Page:  1,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(guests.Guests), 1)

	deleteGuest(guest.ID, t)
}

func TestUpdateGuest(t *testing.T) {
	guest := createGuest(t)
	guest.HotelID = 1
	guest.RoomID = 1
	guest.FirstName = faker.Name()
	guest.LastName = faker.Name()
	guest.PhoneNumber = faker.Phonenumber()
	guest.Email = faker.Email()
	guest.Password = faker.Password()
}
func TestDeleteGuest(t *testing.T) {
	guest := createGuest(t)

	deleteGuest(guest.ID, t)
}
