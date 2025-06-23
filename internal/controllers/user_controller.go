package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/GitNinja36/wello-backend/config"
	"github.com/GitNinja36/wello-backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Register Patient
func RegisterPatient(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:       uuid.NewString(),
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     models.PATIENT,
		Verified: true,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, "User creation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Patient registered successfully",
		"userId":  user.ID,
	})
}

// Register Doctor
func RegisterDoctor(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Name             string  `json:"name"`
		Email            string  `json:"email"`
		Phone            string  `json:"phone"`
		Specialization   string  `json:"specialization"`
		LicenseNumber    string  `json:"licenseNumber"`
		ConsultationFees float64 `json:"consultationFees"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:                uuid.NewString(),
		Name:              req.Name,
		Email:             req.Email,
		Phone:             req.Phone,
		Role:              models.DOCTOR,
		RequestedAsDoctor: true,
		IsApproved:        false,
		Verified:          true,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, "Doctor user creation failed", http.StatusInternalServerError)
		return
	}

	profile := models.DoctorProfile{
		ID:                uuid.NewString(),
		UserID:            user.ID,
		Specialization:    req.Specialization,
		LicenseNumber:     req.LicenseNumber,
		ConsultationFees:  req.ConsultationFees,
		AvailabilitySlots: "[]",
		IsPending:         true,
	}

	if err := config.DB.Create(&profile).Error; err != nil {
		http.Error(w, "Doctor profile creation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Doctor registration requested. Awaiting admin approval.",
		"userId":  user.ID,
	})
}

// Admin Approve Doctor
func ApproveDoctor(w http.ResponseWriter, r *http.Request) {
	doctorId := chi.URLParam(r, "id")

	var profile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", doctorId).First(&profile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}
	profile.IsPending = false
	adminID := "admin-id-from-auth"
	profile.ApprovedBy = &adminID
	if err := config.DB.Save(&profile).Error; err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	config.DB.Model(&models.User{}).Where("id = ?", profile.UserID).
		Update("is_approved", true)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Doctor approved successfully",
	})
}
