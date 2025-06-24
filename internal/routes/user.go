package routes

import (
	"github.com/GitNinja36/wello-backend/internal/controllers"
	"github.com/GitNinja36/wello-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router) {
	r.Post("/register/patient", controllers.RegisterPatient)
	r.Post("/register/doctor", controllers.RegisterDoctor)
	r.Post("/admin/approved/doctor/{id}", controllers.ApproveDoctor)
	r.With(middleware.JWTAuthMiddleware).Post("/onboard/patient", controllers.CompletePatientOnboarding)
	r.Get("/all", controllers.GetUsersByRole)

	r.With(middleware.JWTAuthMiddleware).Put("/edit/doctor/profile", controllers.UpdateDoctorProfile)
	r.With(middleware.JWTAuthMiddleware).Put("/edit/patient/profile", controllers.UpdatePatientProfile)
	r.With(middleware.JWTAuthMiddleware).Get("/me", controllers.GetCurrentUser)
}
