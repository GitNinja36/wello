package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

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
	_ = json.NewDecoder(r.Body).Decode(&req)

	otp := generateOTP()
	utils.SaveOTP(req.Phone, otp, 5*time.Minute)

	err := utils.SendSMS(req.Phone, fmt.Sprintf("Your Wello OTP is: %s", otp))
	if err != nil {
		http.Error(w, "Failed to send OTP via SMS", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
	_ = json.NewDecoder(r.Body).Decode(&req)

	if utils.VerifyOTP(req.Phone, req.OTP) {
		token := utils.GenerateJWT(req.Phone)
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
		return
	}
	http.Error(w, "Invalid OTP", http.StatusUnauthorized)
}

func SendOTPEmail(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email string `json:"email"`
	}
	var req Request
	_ = json.NewDecoder(r.Body).Decode(&req)

	otp := generateOTP()
	utils.SaveOTP(req.Email, otp, 5*time.Minute)

	err := utils.SendEmailOTP(req.Email, otp)
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
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
	_ = json.NewDecoder(r.Body).Decode(&req)

	if utils.VerifyOTP(req.Email, req.OTP) {
		token := utils.GenerateJWT(req.Email)
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
		return
	}
	http.Error(w, "Invalid OTP", http.StatusUnauthorized)
}
