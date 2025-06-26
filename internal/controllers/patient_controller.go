package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/GitNinja36/wello-backend/config"
	"github.com/GitNinja36/wello-backend/internal/middleware"
	"github.com/GitNinja36/wello-backend/internal/models"
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

	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}
