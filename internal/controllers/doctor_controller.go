package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GitNinja36/wello-backend/config"
	"github.com/GitNinja36/wello-backend/internal/middleware"
	"github.com/GitNinja36/wello-backend/internal/models"
	"github.com/GitNinja36/wello-backend/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jung-kurt/gofpdf"
)

// update Doctor profile
func UpdateDoctorProfile(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Name           string  `json:"name"`
		Age            *int    `json:"age"`
		Gender         *string `json:"gender"`
		Address        *string `json:"address"`
		Specialization string  `json:"specialization"`
		LicenseNumber  string  `json:"licenseNumber"`
		ClinicName     string  `json:"clinicName"`
		Experience     string  `json:"experience"`
		Bio            string  `json:"bio"`
		Certifications string  `json:"certifications"`
		PhotoURL       *string `json:"photoUrl"`
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userUpdates := map[string]interface{}{
		"name": req.Name,
	}
	if req.Age != nil {
		userUpdates["age"] = *req.Age
	} else {
		userUpdates["age"] = nil
	}
	if req.Gender != nil {
		userUpdates["gender"] = *req.Gender
	} else {
		userUpdates["gender"] = nil
	}
	if req.Address != nil {
		userUpdates["address"] = *req.Address
	} else {
		userUpdates["address"] = nil
	}

	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(userUpdates).Error; err != nil {
		http.Error(w, "Failed to update user details", http.StatusInternalServerError)
		return
	}

	var profile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	profile.Specialization = req.Specialization
	profile.LicenseNumber = req.LicenseNumber
	profile.ClinicName = req.ClinicName
	profile.Experience = req.Experience
	profile.Bio = req.Bio
	profile.Certifications = req.Certifications
	if req.PhotoURL != nil {
		profile.PhotoURL = req.PhotoURL
	} else {
		profile.PhotoURL = nil
	}

	if err := config.DB.Save(&profile).Error; err != nil {
		http.Error(w, "Failed to update doctor profile", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Doctor profile updated successfully",
	})
}

// Availability slot input type
type Slot struct {
	Day   string   `json:"day"`
	Slots []string `json:"slots"`
}

// update doctor slot
func UpdateDoctorAvailability(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		AvailabilitySlots []Slot `json:"availabilitySlots"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	availabilityJSON, err := json.Marshal(req.AvailabilitySlots)
	if err != nil {
		http.Error(w, "Failed to marshal availability", http.StatusInternalServerError)
		return
	}

	var profile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	profile.AvailabilitySlots = string(availabilityJSON)
	if err := config.DB.Save(&profile).Error; err != nil {
		http.Error(w, "Failed to update availability", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Availability slots updated successfully",
	})
}

// update doctor fee
func UpdateDoctorFee(w http.ResponseWriter, r *http.Request) {
	userId := middleware.GetUserIDFromContext(r)
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ConsultationFees float64 `json:"consultationFees"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var profile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	profile.ConsultationFees = req.ConsultationFees
	if err := config.DB.Save(&profile).Error; err != nil {
		http.Error(w, "Failed to update consultation fee", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Consultation fee updated successfully",
	})
}

// get All upcoming appointments
func GetDoctorAppointments(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var profile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	var appointments []models.Appointment
	if err := config.DB.Preload("Patient").
		Preload("DoctorProfile").
		Where("doctor_profile_id = ?", profile.ID).
		Order("scheduled_at ASC").
		Find(&appointments).Error; err != nil {
		http.Error(w, "Failed to fetch appointments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(appointments)
}

// get Doctor Earnings
func GetDoctorEarnings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var profile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	var totalEarnings float64
	err := config.DB.Model(&models.Appointment{}).
		Select("COALESCE(SUM(doctor_profiles.consultation_fees), 0)").
		Joins("JOIN doctor_profiles ON appointments.doctor_profile_id = doctor_profiles.id").
		Where("appointments.doctor_profile_id = ? AND appointments.status = ?", profile.ID, models.COMPLETED).
		Scan(&totalEarnings).Error
	if err != nil {
		http.Error(w, "Failed to calculate earnings", http.StatusInternalServerError)
		return
	}

	var totalAppointments int64
	err = config.DB.Model(&models.Appointment{}).
		Where("doctor_profile_id = ? AND status = ?", profile.ID, models.COMPLETED).
		Count(&totalAppointments).Error
	if err != nil {
		http.Error(w, "Failed to count appointments", http.StatusInternalServerError)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "all"
	}

	type GroupedResult struct {
		Label string  `json:"label"`
		Total float64 `json:"total"`
		Count int64   `json:"count"`
	}

	var groupedData []GroupedResult

	timeFormat := ""
	switch period {
	case "daily":
		timeFormat = "YYYY-MM-DD"
	case "weekly":
		timeFormat = "IYYY-IW"
	case "monthly":
		timeFormat = "YYYY-MM"
	case "yearly":
		timeFormat = "YYYY"
	}

	if period != "all" {
		err = config.DB.
			Model(&models.Appointment{}).
			Select(`to_char(appointments.scheduled_at, ?) as label,
				    SUM(doctor_profiles.consultation_fees) as total,
				    COUNT(*) as count`, timeFormat).
			Joins("JOIN doctor_profiles ON appointments.doctor_profile_id = doctor_profiles.id").
			Where("appointments.doctor_profile_id = ? AND appointments.status = ?", profile.ID, models.COMPLETED).
			Group("label").
			Order("label ASC").
			Scan(&groupedData).Error
		if err != nil {
			http.Error(w, "Failed to get grouped earnings", http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"totalEarnings":     totalEarnings,
		"totalAppointments": totalAppointments,
		"groupedData":       groupedData,
		"period":            period,
	})
}

// get doctor reviews
func GetDoctorReviews(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var reviews []models.Review
	if err := config.DB.Where("doctor_id = ?", userID).Find(&reviews).Error; err != nil {
		http.Error(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reviews)
}

// accept or reject appointment
func RespondToAppointment(w http.ResponseWriter, r *http.Request) {
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
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status != "ACCEPTED" && req.Status != "REJECTED" {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	var profile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	var appointment models.Appointment
	if err := config.DB.Where("id = ? AND doctor_profile_id = ?", appointmentID, profile.ID).First(&appointment).Error; err != nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	appointment.Status = models.AppointmentStatus(req.Status)
	if err := config.DB.Save(&appointment).Error; err != nil {
		http.Error(w, "Failed to update appointment status", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Appointment status updated successfully",
	})
}

// Seed Dummy Appointment --> Test route
func SeedDummyAppointment(w http.ResponseWriter, r *http.Request) {
	doctorUserID := middleware.GetUserIDFromContext(r)
	if doctorUserID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var doctorProfile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", doctorUserID).First(&doctorProfile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	dummyAppointment := models.Appointment{
		DoctorProfileID: doctorProfile.ID,
		PatientID:       "1f3171ff-3a9a-420e-9d0d-d5d097fdb118",
		Status:          models.PENDING,
		ScheduledAt:     time.Now().Add(48 * time.Hour),
	}

	if err := config.DB.Create(&dummyAppointment).Error; err != nil {
		http.Error(w, "Failed to create dummy appointment", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":        "Dummy appointment created",
		"appointment_id": dummyAppointment.ID,
	})
}

// to Reschedule Appointment
func RescheduleAppointment(w http.ResponseWriter, r *http.Request) {
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
		NewDate time.Time `json:"newDate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var profile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	var appointment models.Appointment
	if err := config.DB.Preload("Patient").Where("id = ? AND doctor_profile_id = ?", appointmentID, profile.ID).First(&appointment).Error; err != nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	appointment.ScheduledAt = req.NewDate
	appointment.Status = models.RESCHEDULE_REQUESTED

	if err := config.DB.Save(&appointment).Error; err != nil {
		http.Error(w, "Failed to update appointment", http.StatusInternalServerError)
		return
	}

	go utils.SendEmail(
		appointment.Patient.Email,
		"Reschedule Request from Doctor",
		fmt.Sprintf("Your appointment has been rescheduled to: %s", req.NewDate.Format("Jan 2, 2006 3:04PM")),
	)

	go utils.SendSMS(
		appointment.Patient.Phone,
		fmt.Sprintf("Your appointment has been rescheduled to: %s", req.NewDate.Format("Jan 2, 2006 3:04PM")),
	)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Appointment reschedule request sent to patient",
	})
}

// get appointments where reschedule is requested
func GetDoctorRescheduleRequests(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
		Where("doctor_profile_id = ? AND status = ?", doctorProfile.ID, models.RESCHEDULE_REQUESTED).
		Order("scheduled_at ASC").
		Find(&appointments).Error; err != nil {
		http.Error(w, "Failed to fetch reschedule requests", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"rescheduleRequests": appointments,
	})
}

// returns all upcoming appointments for a doctor
func GetUpcomingAppointmentsForDoctor(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var profile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	var appointments []models.Appointment
	if err := config.DB.
		Preload("Patient").
		Preload("DoctorProfile").
		Where("doctor_profile_id = ? AND scheduled_at > NOW() AND status IN ?",
			profile.ID,
			[]models.AppointmentStatus{models.ACCEPTED, models.RESCHEDULED_CONFIRMED}).
		Order("scheduled_at ASC").
		Find(&appointments).Error; err != nil {
		http.Error(w, "Failed to fetch appointments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"upcomingAppointments": appointments,
	})
}

// mark appointment as completed
func CompleteAppointment(w http.ResponseWriter, r *http.Request) {
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

	var profile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	var appointment models.Appointment
	if err := config.DB.
		Where("id = ? AND doctor_profile_id = ?", appointmentID, profile.ID).
		First(&appointment).Error; err != nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	if appointment.Status != models.ACCEPTED && appointment.Status != models.RESCHEDULED_CONFIRMED {
		http.Error(w, "Only accepted/rescheduled appointments can be marked as completed", http.StatusBadRequest)
		return
	}

	appointment.Status = models.COMPLETED
	if err := config.DB.Save(&appointment).Error; err != nil {
		http.Error(w, "Failed to mark appointment as completed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Appointment marked as completed",
	})
}

// Add Appointment Summary
func AddAppointmentSummary(w http.ResponseWriter, r *http.Request) {
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
		Summary string `json:"summary"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Summary == "" {
		http.Error(w, "Invalid or empty summary", http.StatusBadRequest)
		return
	}

	var profile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	var appointment models.Appointment
	if err := config.DB.Where("id = ? AND doctor_profile_id = ?", appointmentID, profile.ID).First(&appointment).Error; err != nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	if appointment.Status != models.COMPLETED {
		http.Error(w, "Summary can only be added to completed appointments", http.StatusBadRequest)
		return
	}

	appointment.Summary = &req.Summary
	if err := config.DB.Save(&appointment).Error; err != nil {
		http.Error(w, "Failed to save summary", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Summary added successfully",
	})
}

// Get All Unique Patients of a Doctor
func GetAllPatientsForDoctor(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var doctorProfile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&doctorProfile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	var patients []models.User
	if err := config.DB.
		Model(&models.Appointment{}).
		Select("DISTINCT users.*").
		Joins("JOIN users ON users.id = appointments.patient_id").
		Where("appointments.doctor_profile_id = ? AND appointments.status = ?", doctorProfile.ID, models.COMPLETED).
		Scan(&patients).Error; err != nil {
		http.Error(w, "Failed to fetch patients", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"patients": patients,
	})
}

// Download/Print Summary as PDF
func GenerateSummaryPDF(w http.ResponseWriter, r *http.Request) {
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

	var doctorProfile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&doctorProfile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	var appointment models.Appointment
	if err := config.DB.
		Preload("Patient").
		Preload("DoctorProfile").
		Where("id = ? AND doctor_profile_id = ?", appointmentID, doctorProfile.ID).
		First(&appointment).Error; err != nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	if *appointment.Summary == "" {
		http.Error(w, "No summary found for this appointment", http.StatusNotFound)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Appointment Summary")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Patient Name: "+appointment.Patient.Name)
	pdf.Ln(8)
	pdf.Cell(40, 10, "Doctor: "+appointment.DoctorProfile.ClinicName+" - "+appointment.DoctorProfile.Specialization)
	pdf.Ln(8)
	pdf.Cell(40, 10, "Date: "+appointment.ScheduledAt.Format("Jan 2, 2006 3:04 PM"))
	pdf.Ln(8)
	pdf.MultiCell(0, 10, "Summary:\n"+*appointment.Summary, "", "", false)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=appointment_summary.pdf")
	err := pdf.Output(w)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
	}
}

// Add Test from Doctor Side
func CreateMedicalCheck(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		AppointmentID string          `json:"appointmentId"`
		Type          models.TestType `json:"type"`
		Location      string          `json:"location"`
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var profile models.DoctorProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	var appointment models.Appointment
	if err := config.DB.Where("id = ? AND doctor_profile_id = ?", req.AppointmentID, profile.ID).First(&appointment).Error; err != nil {
		http.Error(w, "Appointment not found or not authorized", http.StatusNotFound)
		return
	}

	check := models.MedicalCheck{
		AppointmentID: req.AppointmentID,
		Type:          req.Type,
		Location:      req.Location,
	}
	if err := config.DB.Create(&check).Error; err != nil {
		http.Error(w, "Failed to create test request", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Test request created successfully",
		"testId":  check.ID,
	})
}

// Uploading Test Report
func UploadTestReport(w http.ResponseWriter, r *http.Request) {
	testID := chi.URLParam(r, "id")
	if testID == "" {
		http.Error(w, "Missing test ID", http.StatusBadRequest)
		return
	}
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ReportURL string `json:"reportUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ReportURL == "" {
		http.Error(w, "Invalid or missing report URL", http.StatusBadRequest)
		return
	}

	var test models.MedicalCheck
	if err := config.DB.
		Joins("JOIN appointments ON appointments.id = medical_checks.appointment_id").
		Joins("JOIN doctor_profiles ON appointments.doctor_profile_id = doctor_profiles.id").
		Where("medical_checks.id = ? AND doctor_profiles.user_id = ?", testID, userID).
		First(&test).Error; err != nil {
		http.Error(w, "Test not found or unauthorized", http.StatusNotFound)
		return
	}

	test.ReportUrl = &req.ReportURL
	test.ReportUploaded = true
	test.Status = models.REPORTED

	if err := config.DB.Save(&test).Error; err != nil {
		http.Error(w, "Failed to update test record", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Report URL saved successfully",
		"url":     req.ReportURL,
	})
}
