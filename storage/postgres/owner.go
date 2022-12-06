package postgres

import (
	"database/sql"
	"fmt"

	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/jmoiron/sqlx"
)

type ownerRepo struct {
	db *sqlx.DB
}

func NewOwner(db *sqlx.DB) repo.OwnerStorageI {
	return &ownerRepo{
		db: db,
	}
}

func (ur *ownerRepo) Create(owner *repo.Owner) (*repo.Owner, error) {
	query := `
		INSERT INTO owners(
			hotel_id,
			first_name,
            last_name,
            email,
			password
		) VALUES($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	row := ur.db.QueryRow(
		query,
		owner.HotelID,
		owner.FirstName,
		owner.LastName,
		owner.Email,
		owner.Password,
	)

	err := row.Scan(
		&owner.ID,
		&owner.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return owner, nil
}

func (ur *ownerRepo) Get(id int64) (*repo.Owner, error) {
	var result repo.Owner

	query := `
		SELECT
			id,
			hotel_id,
			first_name,
            last_name,
            email,
			password,
			created_at
		FROM owners
		WHERE id=$1
	`

	row := ur.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID,
		&result.HotelID,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.Password,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (gr *ownerRepo) GetByEmail(email string) (*repo.Owner, error) {
	var result repo.Owner

	query := `
		SELECT
			id,
			hotel_id,
			first_name,
			last_name,
			email,
			password,
			created_at
		FROM owners
		WHERE email=$1
	`

	row := gr.db.QueryRow(query, email)
	err := row.Scan(
		&result.ID,
		&result.HotelID,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.Password,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *ownerRepo) GetAll(params *repo.GetAllOwnersParams) (*repo.GetAllOwnersResult, error) {
	result := repo.GetAllOwnersResult{
		Owners: make([]*repo.Owner, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	filter := ""
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
			WHERE first_name ILIKE '%s' OR last_name ILIKE '%s' OR email ILIKE '%s'`,
			str, str, str,
		)
	}

	query := `
		SELECT
			id,
			hotel_id,
			first_name,
			last_name,
			email,
			password,
			created_at
		FROM owners
		` + filter + `
		ORDER BY created_at desc
		` + limit

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u repo.Owner

		err := rows.Scan(
			&u.ID,
			&u.HotelID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Owners = append(result.Owners, &u)
	}

	queryCount := `SELECT count(1) FROM owners ` + filter
	err = ur.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (gr *ownerRepo) UpdateOwnerPassword(req *repo.UpdateOwnerPassword) error {
	query := `UPDATE owners SET password=$1 WHERE id=$2`

	_, err := gr.db.Exec(query, req.Password, req.OwnerID)
	if err != nil {
		return err
	}

	return nil
}

func (gr *ownerRepo) Update(guest *repo.Owner) (*repo.Owner, error) {
	query := `
		UPDATE owners SET
			hotel_id=$1,
            first_name=$2,
			last_name=$3,
			email=$4,
			password=$5,
		WHERE id=$6
		RETURNING id, hotel_id, first_name, last_name, email, password
	`
	row := gr.db.QueryRow(
		query,
		guest.HotelID,
		guest.FirstName,
		guest.LastName,
		guest.Email,
		guest.Password,
		guest.ID,
	)

	var result *repo.Guest
	err := row.Scan(
		&result.ID,
		&result.HotelID,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.Password,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return guest, nil
}

func (gr *ownerRepo) Delete(id int64) error {
	query := `DELETE FROM owners WHERE id=$1`

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
