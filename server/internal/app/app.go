package app

import (
	"chat/internal/api"
	"chat/internal/migrations"
	"chat/internal/store"
	"chat/internal/utils"
	"database/sql"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger      *log.Logger
	DB          *sql.DB
	UserHandler *api.UserHandler
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	userStore := store.NewPostgresUserStore(pgDB)

	userHandler := api.NewUserHandler(userStore, logger)

	app := &Application{
		Logger:      logger,
		DB:          pgDB,
		UserHandler: userHandler,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, _r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"status": "OK"})
	a.Logger.Println("INFO: Health check")
}
