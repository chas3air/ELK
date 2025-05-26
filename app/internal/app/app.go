package app

import (
	usershandlers "app/internal/handlers/users"
	usersservice "app/internal/services/users"
	"app/internal/storage/mock"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	log  *slog.Logger
	port int
}

func New(log *slog.Logger, port int) *App {
	return &App{
		log:  log,
		port: port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	r := mux.NewRouter()

	usersStorage := mock.New(a.log)
	usersService := usersservice.New(a.log, usersStorage)
	usersHandler := usershandlers.New(a.log, usersService)

	r.HandleFunc("/api/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("200 OK"))
	}).Methods(http.MethodGet)

	r.HandleFunc("/api/users", usersHandler.GetUsersHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/users/{id}", usersHandler.GetUserByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/users", usersHandler.InsertHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/users/{id}", usersHandler.UpdateHandler).Methods(http.MethodPut)
	r.HandleFunc("/api/users/{id}", usersHandler.DeleteHandler).Methods(http.MethodDelete)

	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", a.port),
		r,
	); err != nil {
		panic(err)
	}

	return nil
}
