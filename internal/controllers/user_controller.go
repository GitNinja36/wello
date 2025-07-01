package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/GitNinja36/wello-backend/config"
	"github.com/GitNinja36/wello-backend/internal/middleware"
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
	adminID := middleware.GetUserIDFromContext(r)
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

// patient-onboarding
func CompletePatientOnboarding(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	userID := r.Context().Value("userID").(string)

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Update user
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"name":  req.Name,
			"email": req.Email,
		}).Error; err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Patient profile completed",
	})
}

// Get all Users By Role
func GetUsersByRole(w http.ResponseWriter, r *http.Request) {
	role := strings.ToUpper(r.URL.Query().Get("role"))

	if role != "PATIENT" && role != "DOCTOR" {
		http.Error(w, "Invalid or missing role query param", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var users []models.User
	var total int64

	config.DB.Model(&models.User{}).
		Where("role = ?", role).
		Count(&total)

	db := config.DB.Preload("Appointments").Preload("Orders")

	if role == "DOCTOR" {
		db = db.Preload("DoctorProfile")
	} else {
		db = db.Preload("AdminProfile")
	}

	if err := db.
		Where("role = ?", role).
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"role":  role,
		"page":  page,
		"limit": limit,
		"total": total,
		"users": users,
	})
}

// Update the Patient Profile
func UpdatePatientProfile(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Age     int    `json:"age"`
		Gender  string `json:"gender"`
		Bio     string `json:"bio"`
		Address string `json:"address"`
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age
	user.Gender = req.Gender
	user.Bio = req.Bio
	user.Address = req.Address

	if err := config.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Patient profile updated successfully",
	})
}

// Get user info
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	UserId := middleware.GetUserIDFromContext(r)
	if UserId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", UserId).Error; err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"phone":      user.Phone,
		"role":       user.Role,
		"photoUrl":   user.PhotoURL,
		"age":        user.Age,
		"gender":     user.Gender,
		"bio":        user.Bio,
		"address":    user.Address,
		"isApproved": user.IsApproved,
	})
}

// update profile photo
func UpdateProfilePhoto(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		PhotoURL string `json:"photoUrl"`
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.PhotoURL) == "" {
		http.Error(w, "Invalid or missing photoUrl", http.StatusBadRequest)
		return
	}

	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).
		Update("photo_url", req.PhotoURL).Error; err != nil {
		http.Error(w, "Failed to update photo", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Profile photo updated successfully",
		"photoUrl": req.PhotoURL,
	})
}
