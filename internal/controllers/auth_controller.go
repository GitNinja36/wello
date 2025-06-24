package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/GitNinja36/wello-backend/config"
	"github.com/GitNinja36/wello-backend/internal/models"
	"github.com/GitNinja36/wello-backend/internal/utils"
)

func generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SendOTPPhone(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Phone string `json:"phone"`
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	otp := generateOTP()
	utils.SaveOTP(req.Phone, otp, 5*time.Minute)

	if err := utils.SendSMS(req.Phone, fmt.Sprintf("Your Wello OTP is: %s", otp)); err != nil {
		http.Error(w, "Failed to send OTP via SMS", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "OTP sent to phone",
	})
}

func VerifyOTPPhone(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if !utils.VerifyOTP(req.Phone, req.OTP) {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	var user models.User
	result := config.DB.Where("phone = ?", req.Phone).First(&user)
	if result.Error == nil {
		token := utils.GenerateJWT(user.ID, string(user.Role), user.IsApproved, false)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"token":        token,
			"role":         user.Role,
			"isApproved":   user.IsApproved,
			"newUser":      false,
			"needsProfile": false,
		})
		return
	}

	token := utils.GenerateJWT(req.Phone, "PATIENT", true, true)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":        token,
		"role":         "PATIENT",
		"isApproved":   true,
		"newUser":      true,
		"needsProfile": true,
	})
}

func SendOTPEmail(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email string `json:"email"`
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	otp := generateOTP()
	utils.SaveOTP(req.Email, otp, 5*time.Minute)

	if err := utils.SendEmailOTP(req.Email, otp); err != nil {
		http.Error(w, "Failed to send OTP via email", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "OTP sent to email",
	})
}

func VerifyOTPEmail(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if !utils.VerifyOTP(req.Email, req.OTP) {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	var user models.User
	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error == nil {
		token := utils.GenerateJWT(user.ID, string(user.Role), user.IsApproved, false)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"token":        token,
			"role":         user.Role,
			"isApproved":   user.IsApproved,
			"newUser":      false,
			"needsProfile": false,
		})
		return
	}

	token := utils.GenerateJWT(req.Email, "PATIENT", true, true)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":        token,
		"role":         "PATIENT",
		"isApproved":   true,
		"newUser":      true,
		"needsProfile": true,
	})
}
