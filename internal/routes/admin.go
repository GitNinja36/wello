package routes

import (
	"github.com/GitNinja36/wello-backend/internal/controllers"
	"github.com/go-chi/chi/v5"
)

func AdminRoutes(r chi.Router) {
	r.Post("/register", controllers.CreateAdminAccount)
}
