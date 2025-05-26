package usershandlers

import (
	"app/internal/domain/models"
	serviceerrors "app/internal/services"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserService interface {
	Get(context.Context) ([]models.User, error)
	GetById(context.Context, int) (models.User, error)
	Insert(context.Context, models.User) (models.User, error)
	Update(context.Context, int, models.User) (models.User, error)
	Delete(context.Context, int) (models.User, error)
}

type UsersHandlers struct {
	log     *slog.Logger
	service UserService
}

func New(log *slog.Logger, service UserService) *UsersHandlers {
	return &UsersHandlers{
		log: log,
		service: service,
	}
}

func (u *UsersHandlers) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	const op = "usersHandlers.GetUsersHandler"

	users, err := u.service.Get(r.Context())
	if err != nil {
		if errors.Is(err, serviceerrors.ErrNotFound) {
			log.Println("WARN: users not found")
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte{})
			return
		}

		log.Println("ERROR: cannot fetch users")
		http.Error(w, "Cannot fetch users", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("ERROR: cannot write users to response body")
		http.Error(w, "ERROR: cannot write users to response body", http.StatusInternalServerError)
		return
	}
}

func (u *UsersHandlers) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	const op = "usersHandlers.GetUserByIdHandler"

	id_s, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("ERROR: id is a required parametr")
		http.Error(w, "ERROR: id is a required parametr", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("ERROR: id must be int")
		http.Error(w, "ERROR: id must be int", http.StatusBadRequest)
		return
	}

	user, err := u.service.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrNotFound) {
			log.Println("WARN: user not found")
			http.Error(w, "WARN: user not found", http.StatusNotFound)
			return
		}

		log.Println("ERROR: cannot fetch user by id")
		http.Error(w, "Cannot fetch user by id", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println("ERROR: cannot write user to response body")
		http.Error(w, "ERROR: cannot write user to response body", http.StatusInternalServerError)
		return
	}
}

func (u *UsersHandlers) InsertHandler(w http.ResponseWriter, r *http.Request) {
	const op = "usersHandler.Insert"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR: body is empty")
		http.Error(w, "ERROR: body is empty", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Println("ERROR: cannot parse body to user")
		http.Error(w, "ERROR: cannot parse body to user", http.StatusBadRequest)
		return
	}

	insertedUser, err := u.service.Insert(r.Context(), user)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrAlreadyExists) {
			log.Println("WARN: user already exists")
			http.Error(w, "WARN: user already exists", http.StatusConflict)
			return
		}

		log.Println("ERROR: cannot insert user")
		http.Error(w, "ERROR: cannot insert user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(insertedUser); err != nil {
		log.Println("ERROR: cannot write user to response body")
		http.Error(w, "ERROR: cannot write user to response body", http.StatusInternalServerError)
		return
	}
}

func (u *UsersHandlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	const op = "usersHandler.Insert"

	id_s, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("ERROR: id is a required parametr")
		http.Error(w, "ERROR: id is a required parametr", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("ERROR: id must be int")
		http.Error(w, "ERROR: id must be int", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR: body is empty")
		http.Error(w, "ERROR: body is empty", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Println("ERROR: cannot parse body to user")
		http.Error(w, "ERROR: cannot parse body to user", http.StatusBadRequest)
		return
	}

	updatedUser, err := u.service.Update(r.Context(), id, user)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrNotFound) {
			log.Println("WARN: user not found")
			http.Error(w, "WARN: user not found", http.StatusNotFound)
			return
		}

		log.Println("ERROR: cannot update user")
		http.Error(w, "Cannot update user", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		log.Println("ERROR: cannot write user to response body")
		http.Error(w, "ERROR: cannot write user to response body", http.StatusInternalServerError)
		return
	}
}

func (u *UsersHandlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	const op = "usersHandler.Delete"

	id_s, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("ERROR: id is a required parametr")
		http.Error(w, "ERROR: id is a required parametr", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("ERROR: id must be int")
		http.Error(w, "ERROR: id must be int", http.StatusBadRequest)
		return
	}

	user, err := u.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrNotFound) {
			log.Println("WARN: user not found")
			http.Error(w, "WARN: user not found", http.StatusNotFound)
			return
		}

		log.Println("ERROR: cannot delete user")
		http.Error(w, "Cannot fetch user", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println("ERROR: cannot write user to response body")
		http.Error(w, "ERROR: cannot write user to response body", http.StatusInternalServerError)
		return
	}
}
