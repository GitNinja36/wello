package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GitNinja36/wello-backend/config"
	"github.com/GitNinja36/wello-backend/internal/middleware"
	"github.com/GitNinja36/wello-backend/internal/models"
	"github.com/GitNinja36/wello-backend/internal/utils"
	"github.com/go-chi/chi/v5"
)

// Patient response to reschedule request
func PatientRespondReschedule(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appointmentID := chi.URLParam(r, "id")
	if appointmentID == "" {
		http.Error(w, "Missing appointment ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Accept bool `json:"accept"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var appointment models.Appointment
	if err := config.DB.Where("id = ? AND patient_id = ?", appointmentID, userID).
		First(&appointment).Error; err != nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	if req.Accept {
		appointment.Status = models.RESCHEDULED_CONFIRMED
	} else {
		appointment.Status = models.RESCHEDULE_REJECTED
	}

	if err := config.DB.Save(&appointment).Error; err != nil {
		http.Error(w, "Failed to update appointment", http.StatusInternalServerError)
		return
	}

	message := "Reschedule rejected"
	if req.Accept {
		message = "Reschedule accepted"
	}

	var doctorProfile models.DoctorProfile
	if err := config.DB.Preload("User").Where("id = ?", appointment.DoctorProfileID).First(&doctorProfile).Error; err == nil {
		go utils.SendEmail(
			doctorProfile.User.Email,
			"Patient Response to Reschedule",
			fmt.Sprintf("The patient has %s your reschedule request.",
				map[bool]string{true: "accepted", false: "rejected"}[req.Accept]),
		)
		go utils.SendSMS(
			doctorProfile.User.Phone,
			fmt.Sprintf("Patient %s your reschedule request.",
				map[bool]string{true: "accepted", false: "rejected"}[req.Accept]),
		)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

// returns all future confirmed appointments
func GetUpcomingAppointmentsForPatient(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var appointments []models.Appointment
	if err := config.DB.
		Preload("DoctorProfile.User").
		Preload("Patient").
		Where("patient_id = ? AND scheduled_at > NOW() AND status IN ?", userID,
			[]models.AppointmentStatus{models.PENDING, models.ACCEPTED, models.RESCHEDULED_CONFIRMED}).
		Order("scheduled_at ASC").
		Find(&appointments).Error; err != nil {
		http.Error(w, "Failed to fetch upcoming appointments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"upcomingAppointments": appointments,
	})
}

// Get Patient History for a Doctor
func GetPatientHistoryForDoctor(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	patientID := chi.URLParam(r, "patientId")
	if patientID == "" {
		http.Error(w, "Missing patient ID", http.StatusBadRequest)
		return
	}

	var doctorProfile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&doctorProfile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	var appointments []models.Appointment
	if err := config.DB.
		Preload("Patient").
		Preload("DoctorProfile").
		Where("doctor_profile_id = ? AND patient_id = ? AND status = ?",
			doctorProfile.ID, patientID, models.COMPLETED).
		Order("scheduled_at DESC").
		Find(&appointments).Error; err != nil {
		http.Error(w, "Failed to fetch patient history", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"patientHistory": appointments,
	})
}

// View Past Appointment History
func GetPatientAppointmentHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var appointments []models.Appointment
	if err := config.DB.
		Preload("DoctorProfile.User").
		Preload("Patient").
		Where("patient_id = ? AND status = ?", userID, models.COMPLETED).
		Order("scheduled_at DESC").
		Find(&appointments).Error; err != nil {
		http.Error(w, "Failed to fetch history", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"history": appointments,
	})
}

// Cancel upcoming appointment
func CancelAppointmentByPatient(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appointmentID := chi.URLParam(r, "id")
	if appointmentID == "" {
		http.Error(w, "Missing appointment ID", http.StatusBadRequest)
		return
	}

	var appointment models.Appointment
	if err := config.DB.
		Preload("DoctorProfile.User").
		Where("id = ? AND patient_id = ?", appointmentID, userID).
		First(&appointment).Error; err != nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	if appointment.Status != models.PENDING &&
		appointment.Status != models.ACCEPTED &&
		appointment.Status != models.RESCHEDULED_CONFIRMED {
		http.Error(w, "Only upcoming appointments can be cancelled", http.StatusBadRequest)
		return
	}

	if appointment.ScheduledAt.Before(utils.CurrentTime()) {
		http.Error(w, "Cannot cancel past appointments", http.StatusBadRequest)
		return
	}

	appointment.Status = models.CANCELLED_BY_PATIENT
	if err := config.DB.Save(&appointment).Error; err != nil {
		http.Error(w, "Failed to cancel appointment", http.StatusInternalServerError)
		return
	}

	// Notify doctor
	go utils.SendEmail(
		appointment.DoctorProfile.User.Email,
		"Appointment Cancelled",
		fmt.Sprintf("Your appointment with the patient on %s has been cancelled by the patient.",
			appointment.ScheduledAt.Format("Jan 2, 2006 at 3:04 PM")),
	)

	go utils.SendSMS(
		appointment.DoctorProfile.User.Phone,
		fmt.Sprintf("Patient cancelled the appointment scheduled for %s.",
			appointment.ScheduledAt.Format("Jan 2, 2006 3:04 PM")),
	)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Appointment cancelled successfully",
	})
}

// Submit rating and review for a completed appointment
func SubmitReviewForAppointment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appointmentID := chi.URLParam(r, "id")
	if appointmentID == "" {
		http.Error(w, "Missing appointment ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Rating int    `json:"rating"`
		Review string `json:"review,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	var appointment models.Appointment
	if err := config.DB.
		Preload("DoctorProfile").
		Where("id = ? AND patient_id = ?", appointmentID, userID).
		First(&appointment).Error; err != nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	if appointment.Status != models.COMPLETED {
		http.Error(w, "Only completed appointments can be reviewed", http.StatusBadRequest)
		return
	}

	appointment.Rating = &req.Rating
	appointment.Review = &req.Review

	if err := config.DB.Save(&appointment).Error; err != nil {
		http.Error(w, "Failed to submit review", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Review submitted successfully",
	})
}

// Get All Appointments of a Patient
func GetAllAppointmentsForPatient(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var appointments []models.Appointment
	if err := config.DB.
		Preload("DoctorProfile.User").
		Preload("Patient").
		Where("patient_id = ?", userID).
		Order("scheduled_at DESC").
		Find(&appointments).Error; err != nil {
		http.Error(w, "Failed to fetch appointments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"appointments": appointments,
	})
}

// Get Patient Profile
func GetPatientProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var patient models.User
	if err := config.DB.Where("id = ? AND role = ?", userID, models.PATIENT).First(&patient).Error; err != nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(patient)
}

// Get Patient Test History
func GetPatientTestHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var tests []models.MedicalCheck
	if err := config.DB.Where("appointment_id IN (SELECT id FROM appointments WHERE patient_id = ?)", userID).
		Order("created_at DESC").
		Find(&tests).Error; err != nil {
		http.Error(w, "Failed to fetch test history", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"tests": tests,
	})
}
