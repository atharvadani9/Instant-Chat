package routes

import (
	"chat/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/healthcheck", app.HealthCheck)
	r.Post("/user.register", app.UserHandler.Register)
	r.Post("/user.login", app.UserHandler.Login)
	r.Get("/user.get", app.UserHandler.GetUsers)

	return r
}
