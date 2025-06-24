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
	r.With(middleware.Authenticate).Post("/onboard/patient", controllers.CompletePatientOnboarding)
}
