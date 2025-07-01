package routes

import (
	"github.com/GitNinja36/wello-backend/internal/controllers"
	"github.com/GitNinja36/wello-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func DoctorRoutes(r chi.Router) {
	//update profile
	r.With(middleware.JWTAuthMiddleware).Put("/profile", controllers.UpdateDoctorProfile)

	//update doctor slot
	r.With(middleware.JWTAuthMiddleware).Post("/slots", controllers.UpdateDoctorAvailability)

	//update doctor fee
	r.With(middleware.JWTAuthMiddleware).Put("/fee", controllers.UpdateDoctorFee)

	//get All upcoming appointments
	r.With(middleware.JWTAuthMiddleware).Get("/dashboard/appointments", controllers.GetDoctorAppointments)

	//accept or reject appointment
	r.With(middleware.JWTAuthMiddleware).Put("/dashboard/appointment/{id}", controllers.RespondToAppointment)

	//get Doctor Earnings
	r.With(middleware.JWTAuthMiddleware).Get("/earnings", controllers.GetDoctorEarnings)

	//get doctor reviews
	r.With(middleware.JWTAuthMiddleware).Get("/reviews", controllers.GetDoctorReviews)

	// to created Dummy appointment --> for test route
	r.With(middleware.JWTAuthMiddleware).Post("/debug/seed-appointment", controllers.SeedDummyAppointment)

	//to Reschedule Appointment
	r.With(middleware.JWTAuthMiddleware).Put("/appointment/reschedule/{id}", controllers.RescheduleAppointment)

	// view reschedule requests
	r.With(middleware.JWTAuthMiddleware).Get("/reschedule-requests", controllers.GetDoctorRescheduleRequests)

	// View all upcoming appointments for doctor
	r.With(middleware.JWTAuthMiddleware).Get("/upcoming-appointments", controllers.GetUpcomingAppointmentsForDoctor)

	// mark appointment as completed
	r.With(middleware.JWTAuthMiddleware).Put("/appointments/{id}/complete", controllers.CompleteAppointment)

	//Add Appointment Summary
	r.With(middleware.JWTAuthMiddleware).Put("/appointments/{id}/summary", controllers.AddAppointmentSummary)

	//Get All Unique Patients of a Doctor
	r.With(middleware.JWTAuthMiddleware).Get("/patients", controllers.GetAllPatientsForDoctor)

	// Download/Print Summary as PDF
	r.With(middleware.JWTAuthMiddleware).Get("/appointments/{id}/summary-pdf", controllers.GenerateSummaryPDF)

	// Add Test from Doctor Side
	r.With(middleware.JWTAuthMiddleware).Post("/tests/add", controllers.CreateMedicalCheck)

	//Uploading Test Reports
	r.With(middleware.JWTAuthMiddleware).Put("/tests/{id}/upload-report", controllers.UploadTestReport)
}
