package routes

import (
	"github.com/GitNinja36/wello-backend/internal/controllers"
	"github.com/GitNinja36/wello-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func AppointmentRoutes(r chi.Router) {
	// Book Appointment
	r.With(middleware.JWTAuthMiddleware).Post("/book", controllers.BookAppointment)
}
