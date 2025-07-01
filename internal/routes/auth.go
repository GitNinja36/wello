package routes

import (
	"github.com/GitNinja36/wello-backend/internal/controllers"
	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router) {
	r.Post("/send/phone/otp", controllers.SendOTPPhone)
	r.Post("/verify/phone/otp", controllers.VerifyOTPPhone)

	r.Post("/send/email/otp", controllers.SendOTPEmail)
	r.Post("/verify/email/otp", controllers.VerifyOTPEmail)
}
