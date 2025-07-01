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
	r.With(middleware.JWTAuthMiddleware).Get("/{patientId}/history", controllers.GetPatientHistoryForDoctor)

	//View Past Appointment History
	r.With(middleware.JWTAuthMiddleware).Get("/appointments/history", controllers.GetPatientAppointmentHistory)

	// Cancel upcoming appointment
	r.With(middleware.JWTAuthMiddleware).Put("/appointments/{id}/cancel", controllers.CancelAppointmentByPatient)

	//to give review
	r.With(middleware.JWTAuthMiddleware).Post("/appointments/{id}/review", controllers.SubmitReviewForAppointment)

	//Get All Appointments of a Patient
	r.With(middleware.JWTAuthMiddleware).Get("/appointments/all", controllers.GetAllAppointmentsForPatient)

	//Get Patient Profile
	r.With(middleware.JWTAuthMiddleware).Get("/profile", controllers.GetPatientProfile)

	//Get Patient Test History
	r.With(middleware.JWTAuthMiddleware).Get("/tests/history", controllers.GetPatientTestHistory)
}
