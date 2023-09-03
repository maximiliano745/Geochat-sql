package data

import (
	"context"
	"fmt"
	"time"

	"github.com/github.com/maximiliano745/Geochat-sql/pkg/user"
)

type UserRepository struct {
	Data *Data
}

// Obtener todos
func (ur *UserRepository) GetAll(ctx context.Context) ([]user.User, error) {
	q := `
    SELECT id, first_name, last_name, username, email, picture,
        created_at, updated_at
        FROM users;
    `

	rows, err := ur.Data.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Username,
			&u.Email, &u.Picture, &u.CreatedAt, &u.UpdatedAt)
		users = append(users, u)
	}

	return users, nil
}

// Obtener uno
func (ur *UserRepository) GetOne(ctx context.Context, id uint) (user.User, error) {
	q := `
    SELECT id, first_name, last_name, username, email, picture,
        created_at, updated_at
        FROM users WHERE id = $1;
    `

	row := ur.Data.DB.QueryRowContext(ctx, q, id)

	var u user.User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Username, &u.Email,
		&u.Picture, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Obtener por username
func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (user.User, error) {
	q := `
    SELECT id, first_name, last_name, username, email, picture,
        password, created_at, updated_at
        FROM users WHERE username = $1;
    `

	row := ur.Data.DB.QueryRowContext(ctx, q, username)

	var u user.User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Username,
		&u.Email, &u.Picture, &u.Hash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Obtener X Email  ----------------- MIO -----------------------
func (ur *UserRepository) GetByMail(ctx context.Context, email string) (user.User, error) {
	fmt.Println("Revisanso Email existe..???   GetByMail()")
	q := `
    SELECT  id, username, email, hash, created_at, updated_at, pasword
        FROM users WHERE email = $1;
    `

	row := ur.Data.DB.QueryRowContext(ctx, q, email)

	var u user.User
	err := row.Scan(&u.ID, &u.Username, &u.Email,
		&u.Hash, &u.CreatedAt, &u.UpdatedAt, &u.Password)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Insertar
func (ur *UserRepository) Create(ctx context.Context, u *user.User) error {
	q := `
    INSERT INTO users (username, pasword, email, created_at, updated_at, hash)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id;
    `

	if err := u.HashPassword(); err != nil {
		return err
	}

	row := ur.Data.DB.QueryRowContext(
		ctx, q, u.Username, u.Password, u.Email, time.Now(), time.Now(), u.Hash,
	)

	err := row.Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

// Actualizar
func (ur *UserRepository) Update(ctx context.Context, id uint, u user.User) error {
	q := `
    UPDATE users set first_name=$1, last_name=$2, email=$3, picture=$4, updated_at=$5
        WHERE id=$6;
    `

	stmt, err := ur.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx, u.FirstName, u.LastName, u.Email,
		u.Picture, time.Now(), id,
	)
	if err != nil {
		return err
	}

	return nil
}

// Borrar
func (ur *UserRepository) Delete(ctx context.Context, id uint) error {
	q := `DELETE FROM users WHERE id=$1;`

	stmt, err := ur.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
