package postgres

import (
	"database/sql"
	"fmt"

	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/jmoiron/sqlx"
)

type roomRepo struct {
	db *sqlx.DB
}

func NewRoom(db *sqlx.DB) repo.RoomStorageI {
	return &roomRepo{
		db: db,
	}
}

func (ur *roomRepo) Create(room *repo.Room) (*repo.Room, error) {
	query := `
		INSERT INTO rooms(
			hotel_id,
			status,
			image_url
		) VALUES($1, $2, $3)
		RETURNING id, created_at
	`

	row := ur.db.QueryRow(
		query,
		room.HotelID,
		room.Status,
		room.ImageUrl,
	)

	err := row.Scan(
		&room.ID,
		&room.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (ur *roomRepo) Get(id int64) (*repo.Room, error) {
	var result repo.Room

	query := `
		SELECT
			id,
			hotel_id,
            status,
            image_url,
			created_at
		FROM rooms
		WHERE id=$1
	`

	row := ur.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID,
		&result.HotelID,
		&result.Status,
		&result.ImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *roomRepo) GetAll(params *repo.GetAllRoomsParams) (*repo.GetAllRoomsResult, error) {
	result := repo.GetAllRoomsResult{
		Rooms: make([]*repo.Room, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, offset)

	filter := ""
	if params.HotelID != 0 {
		filter += fmt.Sprintf(" WHERE hotel_id='%d'", params.HotelID)
	}

	query := `
		SELECT
			id,
			hotel_id,
            status,
			image_url,
			created_at
		FROM rooms
		` + filter + `
		ORDER BY created_at desc
        ` + limit

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u repo.Room

		err := rows.Scan(
			&u.ID,
			&u.HotelID,
			&u.Status,
			&u.ImageUrl,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Rooms = append(result.Rooms, &u)
	}

	queryCount := `SELECT count(1) FROM rooms ` + filter
	err = ur.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *roomRepo) Update(room *repo.Room) (*repo.Room, error) {
	query := `
		UPDATE rooms SET
		    hotel_id=$1,
            status=$2,
            image_url=$3
		WHERE id=$4
		RETURNING id, created_at
	`

	row := ur.db.QueryRow(
		query,
		room.HotelID,
		room.Status,
		room.ImageUrl,
		room.ID,
	)

	var result repo.Room
	err := row.Scan(
		&result.ID,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (ur *roomRepo) Delete(id int64) error {
	query := "DELETE FROM rooms WHERE id=$1"

	result, err := ur.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsCount == 0 {
		return sql.ErrNoRows
	}

	return nil
}
