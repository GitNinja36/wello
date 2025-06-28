package routes

import (
	"github.com/GitNinja36/wello-backend/internal/controllers"
	"github.com/GitNinja36/wello-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func PatientRoutes(r chi.Router) {

	// Patient response to reschedule request
	r.With(middleware.JWTAuthMiddleware).Put("/appointment/respond-reschedule/{id}", controllers.PatientRespondReschedule)

	// View upcoming appointments for patient
	r.With(middleware.JWTAuthMiddleware).Get("/appointments/upcoming", controllers.GetUpcomingAppointmentsForPatient)

	//Get Patient History for a Doctor
	r.With(middleware.JWTAuthMiddleware).Get("/patients/{patientId}/history", controllers.GetPatientHistoryForDoctor)

}
