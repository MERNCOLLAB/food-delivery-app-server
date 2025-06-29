package email

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"

	"math/rand"
	"time"
)

func GenerateResetCode() (string, time.Time) {
	code := rand.Intn(90000) + 10000
	codeStr := strconv.Itoa(code)
	expiresAt := time.Now().Add(5 * time.Minute)

	return codeStr, expiresAt
}

func SendResetCode(to string, code string) error {

	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	sender := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	subject := "Your Password Reset Code"
	body := fmt.Sprintf("Your password reset code is: %s . It will expire in 5 minutes.", code)
	msg := []byte("Subject: " + subject + "\r\n" +
		"To: " + to + "\r\n" +
		"From: " + sender + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", sender, password, smtpServer)
	addr := smtpServer + ":" + smtpPort
	return smtp.SendMail(addr, auth, sender, []string{to}, msg)
}
