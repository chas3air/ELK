package usersservice

import (
	"app/internal/domain/models"
	serviceerrors "app/internal/services"
	"app/pkg/lib/logger/sl"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type UserStorage interface {
	Get(context.Context) ([]models.User, error)
	GetById(context.Context, int) (models.User, error)
	Insert(context.Context, models.User) (models.User, error)
	Update(context.Context, int, models.User) (models.User, error)
	Delete(context.Context, int) (models.User, error)
}

type UserService struct {
	log     *slog.Logger
	storage UserStorage
}

func New(log *slog.Logger, storage UserStorage) *UserService {
	return &UserService{
		log:     log,
		storage: storage,
	}
}

// Get implements UserService.
func (u *UserService) Get(ctx context.Context) ([]models.User, error) {
	const op = "userService.Get"
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	users, err := u.storage.Get(ctx)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrNotFound) {
			u.log.Warn("users not found", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, serviceerrors.ErrNotFound)
		}

		u.log.Error("cannot fetch users", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

// GetById implements UserService.
func (u *UserService) GetById(ctx context.Context, id int) (models.User, error) {
	const op = "userService.GetById"
	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := u.storage.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrNotFound) {
			u.log.Warn("user not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, serviceerrors.ErrNotFound)
		}

		u.log.Error("cannot fetch user by id", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Insert implements UserService.
func (u *UserService) Insert(ctx context.Context, user models.User) (models.User, error) {
	const op = "userService.Insert"
	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := u.storage.Insert(ctx, user)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrAlreadyExists) {
			u.log.Warn("user already exists", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, serviceerrors.ErrAlreadyExists)
		}

		u.log.Error("cannot insert user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Update implements UserService.
func (u *UserService) Update(ctx context.Context, id int, user models.User) (models.User, error) {
	const op = "userService.Update"
	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := u.storage.Update(ctx, id, user)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrNotFound) {
			u.log.Warn("user not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, serviceerrors.ErrNotFound)
		}

		u.log.Error("cannot update users", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Delete implements UserService.
func (u *UserService) Delete(ctx context.Context, id int) (models.User, error) {
	const op = "userService.Delete"
	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := u.storage.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrNotFound) {
			u.log.Warn("user not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, serviceerrors.ErrNotFound)
		}

		u.log.Error("cannot delete users", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
