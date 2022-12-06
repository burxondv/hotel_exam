package postgres_test

import (
	"testing"

	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createHotel(t *testing.T) *repo.Hotel {
	hotel, err := strg.Hotel().Create(&repo.Hotel{
		Name:         faker.Name(),
		Address:      faker.Name(),
		StarsPopular: 4,
		ImageUrl:     faker.URL(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, hotel)

	return hotel
}

func deleteHotel(id int64, t *testing.T) {
	err := strg.Hotel().Delete(id)
	require.NoError(t, err)
}

func TestCreateHotel(t *testing.T) {
	createHotel(t)
}

func TestGetHotel(t *testing.T) {
	c := createHotel(t)

	hotel, err := strg.Hotel().Get(c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, hotel)
}

func TestGetAllHotel(t *testing.T) {
	hotel := createHotel(t)

	hotels, err := strg.Hotel().GetAll(&repo.GetAllHotelsParams{
		Limit: 10,
		Page:  1,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(hotels.Hotels), 1)

	deleteHotel(hotel.ID, t)
}

func TestUpdateHotel(t *testing.T) {
	hotel := createHotel(t)
	hotel.Name = faker.Name()
	hotel.Address = faker.Name()
	hotel.StarsPopular = 4
	hotel.ImageUrl = faker.URL()
}

func TestDeleteHotel(t *testing.T) {
	hotel := createHotel(t)

	deleteHotel(hotel.ID, t)
}
