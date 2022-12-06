package postgres_test

import (
	"testing"

	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createRoom(t *testing.T) *repo.Room {
	room, err := strg.Room().Create(&repo.Room{
		HotelID:  1,
		Status:   false,
		ImageUrl: faker.URL(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, room)

	return room
}

func deleteRoom(id int64, t *testing.T) {
	err := strg.Room().Delete(id)
	require.NoError(t, err)
}

func TestCreateRoom(t *testing.T) {
	createRoom(t)
}

func TestGetRoom(t *testing.T) {
	c := createRoom(t)

	room, err := strg.Room().Get(c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, room)
}

func TestGetAllRoom(t *testing.T) {
	room := createRoom(t)

	rooms, err := strg.Room().GetAll(&repo.GetAllRoomsParams{
		Limit: 10,
		Page:  1,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(rooms.Rooms), 1)

	deleteRoom(room.ID, t)
}

func TestUpdaterRoom(t *testing.T) {
	room := createRoom(t)
	room.HotelID = 1
	room.Status = false
    room.ImageUrl = faker.URL()
}

func TestDeleteRoom(t *testing.T) {
	room := createRoom(t)

	deleteRoom(room.ID, t)
}
