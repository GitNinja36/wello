package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/GitNinja36/wello-backend/config"
	"github.com/GitNinja36/wello-backend/internal/models"
	"github.com/GitNinja36/wello-backend/internal/utils"
)

func CreateAdminAccount(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Name       string `json:"name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Password   string `json:"password"`
		Position   string `json:"position"`
		Department string `json:"department"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		http.Error(w, "Email and Password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	adminUser := models.User{
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Password:   hashedPassword,
		Role:       models.ADMIN,
		Verified:   true,
		IsApproved: true,
	}

	if err := config.DB.Create(&adminUser).Error; err != nil {
		http.Error(w, "Failed to create admin user", http.StatusInternalServerError)
		return
	}

	adminProfile := models.AdminProfile{
		UserID:     adminUser.ID,
		Position:   req.Position,
		Department: req.Department,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := config.DB.Create(&adminProfile).Error; err != nil {
		http.Error(w, "Failed to create admin profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Admin account created successfully",
		"userId":  adminUser.ID,
	})
}
