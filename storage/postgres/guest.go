package postgres

import (
	"database/sql"
	"fmt"

	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/jmoiron/sqlx"
)

type guestRepo struct {
	db *sqlx.DB
}

func NewGuest(db *sqlx.DB) repo.GuestStorageI {
	return &guestRepo{
		db: db,
	}
}

func (gr *guestRepo) Create(guest *repo.Guest) (*repo.Guest, error) {
	query := `
		INSERT INTO guests(
			hotel_id,
			room_id,
			first_name,
			last_name,
			phone_number,
			email,
			password
		) VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	row := gr.db.QueryRow(
		query,
		guest.HotelID,
		guest.RoomID,
		guest.FirstName,
		guest.LastName,
		guest.PhoneNumber,
		guest.Email,
		guest.Password,
	)

	err := row.Scan(
		&guest.ID,
		&guest.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return guest, nil
}

func (gr *guestRepo) Get(id int64) (*repo.Guest, error) {
	var result repo.Guest

	query := `
		SELECT
			id,
			hotel_id,
            room_id,
			first_name,
            last_name,
			phone_number,
            email,
            password,
			created_at
		FROM guests
        WHERE id=$1
	`

	row := gr.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID,
		&result.HotelID,
		&result.RoomID,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Email,
		&result.Password,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (gr *guestRepo) GetByEmail(email string) (*repo.Guest, error) {
	var result repo.Guest

	query := `
		SELECT
			id,
			hotel_id,
			room_id,
			first_name,
			last_name,
			phone_number,
			email,
			password,
			created_at
		FROM guests
		WHERE email=$1
	`

	row := gr.db.QueryRow(query, email)
	err := row.Scan(
		&result.ID,
		&result.HotelID,
		&result.RoomID,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Email,
		&result.Password,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (gr *guestRepo) GetAll(params *repo.GetAllGuestsParams) (*repo.GetAllGuestsResult, error) {
	result := repo.GetAllGuestsResult{
		Guests: make([]*repo.Guest, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, offset)

	filter := ""
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
		WHERE first_name ilike '%s' OR last_name ilike '%s' OR email ilike '%s' OR phone_number ilike '%s'`,
			str, str, str, str,
		)
	}

	query := `
		SELECT
			id,
			hotel_id,
			room_id,
			first_name,
			last_name,
			phone_number,
			email,
			password,
			created_at
		FROM guests
		` + filter + `
		ORDER BY created_at desc
		` + limit

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u repo.Guest

		err := rows.Scan(
			&u.ID,
			&u.HotelID,
			&u.RoomID,
			&u.FirstName,
			&u.LastName,
			&u.PhoneNumber,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Guests = append(result.Guests, &u)
	}

	queryCount := `SELECT count(1) FROM guests ` + filter
	err = gr.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (gr *guestRepo) UpdateGuestPassword(req *repo.UpdateGuestPassword) error {
	query := `UPDATE guests SET password=$1 WHERE id=$2`

	_, err := gr.db.Exec(query, req.Password, req.GuestID)
	if err != nil {
		return err
	}

	return nil
}

func (gr *guestRepo) Update(guest *repo.Guest) (*repo.Guest, error) {
	query := `
		UPDATE guests SET
			hotel_id=$1,
			room_id=$2,
            first_name=$3,
			last_name=$4,
            phone_number=$5,
			email=$6,
			password=$7,
		WHERE id=$8
		RETURNING id, hotel_id, room_id, first_name, last_name, phone_number, email, password
	`
	row := gr.db.QueryRow(
		query,
		guest.HotelID,
		guest.RoomID,
		guest.FirstName,
		guest.LastName,
		guest.PhoneNumber,
		guest.Email,
		guest.Password,
		guest.ID,
	)

	var result *repo.Guest
	err := row.Scan(
		&result.ID,
		&result.HotelID,
		&result.RoomID,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Email,
		&result.Password,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return guest, nil
}

func (gr *guestRepo) Delete(id int64) error {
	query := `DELETE FROM guests WHERE id=$1`

	result, err := gr.db.Exec(query, id)
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
