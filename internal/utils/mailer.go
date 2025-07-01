package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmailOTP(to, otp string) error {
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")
	fromName := os.Getenv("EMAIL_FROM_NAME")

	msg := fmt.Sprintf("Subject: Your Wello OTP\n\nHi,\nYour OTP is: %s\n\nThis OTP is valid for 10 minutes.\n\nDo not share it with anyone.\n\nThanks,\n%s", otp, fromName)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg),
	)

	if err != nil {
		fmt.Println("Failed to send email:", err)
	}
	return err
}

// SendEmail - (for appointment/reschedule notifications)
func SendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")

	msg := "From: " + os.Getenv("EMAIL_FROM_NAME") + " <" + from + ">\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" + body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg),
	)

	if err != nil {
		fmt.Println("Failed to send email:", err)
	}
	return err
}
