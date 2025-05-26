package mock

import (
	"app/internal/domain/models"
	"context"
	"errors"
	"log/slog"
	"sync"
)

type UserStorage struct {
	users  map[int]models.User
	mu     sync.Mutex
	nextId int
}

func New(log *slog.Logger) *UserStorage {
	return &UserStorage{
		users:  make(map[int]models.User),
		nextId: 1,
	}
}

// Get возвращает всех пользователей.
func (u *UserStorage) Get(ctx context.Context) ([]models.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	users := make([]models.User, 0, len(u.users))
	for _, user := range u.users {
		users = append(users, user)
	}
	return users, nil
}

// GetById возвращает пользователя по ID.
func (u *UserStorage) GetById(ctx context.Context, id int) (models.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	user, exists := u.users[id]
	if !exists {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

// Insert добавляет нового пользователя.
func (u *UserStorage) Insert(ctx context.Context, user models.User) (models.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	user.Id = u.nextId
	u.users[u.nextId] = user
	u.nextId++
	return user, nil
}

// Update обновляет информацию о пользователе.
func (u *UserStorage) Update(ctx context.Context, id int, user models.User) (models.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, exists := u.users[id]; !exists {
		return models.User{}, errors.New("user not found")
	}
	user.Id = id
	u.users[id] = user
	return user, nil
}

// Delete удаляет пользователя по ID.
func (u *UserStorage) Delete(ctx context.Context, id int) (models.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	user, exists := u.users[id]
	if !exists {
		return models.User{}, errors.New("user not found")
	}
	delete(u.users, id)
	return user, nil
}
