package repository

import (
	"context"

	"inventory-tracker/models"
)

type UserRepository interface {
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	Create(ctx context.Context, name, categoryGroup, username, passwordHash string) (models.Users, error)
	GetByUsername(ctx context.Context, username string) (models.Users, error)
	GetByID(ctx context.Context, id int) (models.Users, error)
}

type userRepo struct {
	DB DBTX
}

func NewUserRepo(db DBTX) UserRepository {
	return &userRepo{DB: db}
}

func (r *userRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, username).Scan(&exists)
	if err != nil {
		return false, translate(err)
	}
	return exists, nil
}

func (r *userRepo) Create(ctx context.Context, name, categoryGroup, username, passwordHash string) (models.Users, error) {
	var u models.Users
	err := r.DB.QueryRow(ctx,
		`INSERT INTO users (name, category_group, username, password)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, category_group, username`,
		name, categoryGroup, username, passwordHash).
		Scan(&u.Id, &u.Name, &u.CategoryGroup, &u.Username)
	if err != nil {
		return models.Users{}, translate(err)
	}
	return u, nil
}

// GetByUsername returns the password hash along with the user, because the
// service needs it to compare. Everything above the service drops it.
func (r *userRepo) GetByUsername(ctx context.Context, username string) (models.Users, error) {
	var u models.Users
	err := r.DB.QueryRow(ctx,
		`SELECT id, name, category_group, username, password FROM users WHERE username = $1`, username).
		Scan(&u.Id, &u.Name, &u.CategoryGroup, &u.Username, &u.Password)
	if err != nil {
		return models.Users{}, translate(err)
	}
	return u, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int) (models.Users, error) {
	var u models.Users
	err := r.DB.QueryRow(ctx,
		`SELECT id, name, category_group, username FROM users WHERE id = $1`, id).
		Scan(&u.Id, &u.Name, &u.CategoryGroup, &u.Username)
	if err != nil {
		return models.Users{}, translate(err)
	}
	return u, nil
}
