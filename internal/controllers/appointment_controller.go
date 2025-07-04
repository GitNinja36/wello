package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GitNinja36/wello-backend/config"
	"github.com/GitNinja36/wello-backend/internal/middleware"
	"github.com/GitNinja36/wello-backend/internal/models"
)

// Book Appointment
type BookAppointmentRequest struct {
	DoctorID      string `json:"doctorId"`
	ScheduledAt   string `json:"scheduledAt"`
	Mode          string `json:"mode"`
	Location      string `json:"location"`
	FeePaid       bool   `json:"feePaid"`
	PatientName   string `json:"patientName"`
	ContactNumber string `json:"contactNumber"`
	Age           int    `json:"age"`
}

func BookAppointment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req BookAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		http.Error(w, "Invalid date format. Expected RFC3339", http.StatusBadRequest)
		return
	}

	appt := models.Appointment{
		PatientID:       userID,
		DoctorProfileID: req.DoctorID,
		ScheduledAt:     scheduledTime,
		Mode:            models.AppointmentMode(req.Mode),
		Location:        &req.Location,
		FeePaid:         req.FeePaid,
		Status:          models.PENDING,
	}

	if err := config.DB.Create(&appt).Error; err != nil {
		http.Error(w, "Failed to book appointment", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Appointment booked successfully",
		"id":      appt.ID,
	})
}
