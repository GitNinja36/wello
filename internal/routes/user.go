package routes

import (
	"github.com/GitNinja36/wello-backend/internal/controllers"
	"github.com/GitNinja36/wello-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router) {
	//register Patient
	r.Post("/register/patient", controllers.RegisterPatient)

	//register Doctor
	r.Post("/register/doctor", controllers.RegisterDoctor)

	//
	r.With(middleware.JWTAuthMiddleware).Post("/onboard/patient", controllers.CompletePatientOnboarding)

	//Update Patient Profile
	r.With(middleware.JWTAuthMiddleware).Put("/edit/patient/profile", controllers.UpdatePatientProfile)

	//Get Current User
	r.With(middleware.JWTAuthMiddleware).Get("/me", controllers.GetCurrentUser)

	//update photo
	r.With(middleware.JWTAuthMiddleware).Put("/me/photo", controllers.UpdateProfilePhoto)

	//Get all user
	r.Get("/all", controllers.GetUsersByRole)
}
