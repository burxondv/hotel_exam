package postgres

import (
	"database/sql"
	"fmt"

	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/jmoiron/sqlx"
)

type hotelRepo struct {
	db *sqlx.DB
}

func NewHotel(db *sqlx.DB) repo.HotelStorageI {
	return &hotelRepo{
		db: db,
	}
}

func (ur *hotelRepo) Create(hotel *repo.Hotel) (*repo.Hotel, error) {
	query := `
		INSERT INTO hotels(
			name,
			address,
            stars_popular,
			image_url
		) VALUES($1, $2, $3, $4)
		RETURNING id, created_at
	`

	row := ur.db.QueryRow(
		query,
		hotel.Name,
		hotel.Address,
		hotel.StarsPopular,
		hotel.ImageUrl,
	)

	err := row.Scan(
		&hotel.ID,
		&hotel.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return hotel, nil
}

func (ur *hotelRepo) Get(id int64) (*repo.Hotel, error) {
	var result repo.Hotel

	query := `
		SELECT
			id,
			name,
			address,
            stars_popular,
            image_url,
			created_at
		FROM hotels
		WHERE id=$1
	`

	row := ur.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID,
		&result.Name,
		&result.Address,
		&result.StarsPopular,
		&result.ImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *hotelRepo) GetAll(params *repo.GetAllHotelsParams) (*repo.GetAllHotelsResult, error) {
	result := repo.GetAllHotelsResult{
		Hotels: make([]*repo.Hotel, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, offset)

	filter := ""
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
			WHERE name ilike '%s' OR address ilike '%s'`,
			str, str,
		)
	}

	query := `
		SELECT
			id,
			name,
			address,
			stars_popular,
			image_url,
			created_at
		FROM hotels
		` + filter + `
		ORDER BY stars_popular desc
        ` + limit

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u repo.Hotel

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Address,
			&u.StarsPopular,
			&u.ImageUrl,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Hotels = append(result.Hotels, &u)
	}

	queryCount := `SELECT count(1) FROM hotels ` + filter
	err = ur.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *hotelRepo) Update(hotel *repo.Hotel) (*repo.Hotel, error) {
	query := `
		UPDATE hotels SET
		    name=$1,
            address=$2,
			stars_popular=$3,
            image_url=$4
		WHERE id=$5
		RETURNING id, name, address, stars_popular, image_url, created_at
	`

	row := ur.db.QueryRow(
		query,
		hotel.Name,
		hotel.Address,
		hotel.StarsPopular,
		hotel.ImageUrl,
		hotel.ID,
	)

	var result repo.Hotel
	err := row.Scan(
		&result.ID,
		&result.Name,
		&result.Address,
		&result.StarsPopular,
		&result.ImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return hotel, nil
}

func (ur *hotelRepo) Delete(id int64) error {
	query := "DELETE FROM hotels WHERE id=$1"

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
