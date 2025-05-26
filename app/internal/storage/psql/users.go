package psql

import (
	"app/internal/domain/models"
	"context"
	"database/sql"
)

type UserStorage struct {
	DB *sql.DB
}

func New(connStr string) *UserStorage {
	return &UserStorage{
		DB: nil,
	}
}

// Get implements UserStorage1.
func (u *UserStorage) Get(context.Context) ([]models.User, error) {
	return nil, nil
}

// GetById implements UserStorage1.
func (u *UserStorage) GetById(context.Context, int) (models.User, error) {
	panic("unimplemented")
}

// Insert implements UserStorage1.
func (u *UserStorage) Insert(context.Context, models.User) (models.User, error) {
	panic("unimplemented")
}

// Update implements UserStorage1.
func (u *UserStorage) Update(context.Context, int, models.User) (models.User, error) {
	panic("unimplemented")
}

// Delete implements UserStorage1.
func (u *UserStorage) Delete(context.Context, int) (models.User, error) {
	panic("unimplemented")
}
