package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Base route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server setup"))
	})

	// All routes
	r.Route("/auth", AuthRoutes)
	r.Route("/admin", AdminRoutes)
	r.Route("/user", UserRoutes)
	r.Route("/doctor", DoctorRoutes)
	r.Route("/patient", PatientRoutes)
	r.Route("/appointment", AppointmentRoutes)
	r.Route("/medical-check", MedicalCheckRoutes)
	r.Route("/order", OrderRoutes)

	return r
}
