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

func fallback(origValue, fallback string) string {
	if origValue == "" {
		return fallback
	}
	return origValue
}

func getSMTPConfig() (smtpServer, smtpPort, sender, password string) {
	smtpServer = fallback(os.Getenv("SMTP_SERVER"), "smtp.gmail.com")
	smtpPort = fallback(os.Getenv("SMTP_PORT"), "587")
	sender = os.Getenv("SMTP_EMAIL")
	password = os.Getenv("SMTP_PASSWORD")
	return
}

func SendResetCode(to string, code string) error {
	smtpServer, smtpPort, sender, password := getSMTPConfig()
	emailBody := fmt.Sprintf(ResetPasswordTemplate, code)

	subject := "Food Delivery App Password Reset"
	body := emailBody
	msg := []byte("Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"To: " + to + "\r\n" +
		"From: " + sender + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", sender, password, smtpServer)
	addr := smtpServer + ":" + smtpPort
	return smtp.SendMail(addr, auth, sender, []string{to}, msg)
}

func SendSignUpForm(to, role, sigupURL string) error {
	smtpServer, smtpPort, sender, password := getSMTPConfig()
	emailBody := fmt.Sprintf(SignUpFormTemplate, role, sigupURL, sigupURL)

	subject := "Food Delivery App Sign Up Invitation"
	body := emailBody
	msg := []byte("Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"To: " + to + "\r\n" +
		"From: " + sender + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", sender, password, smtpServer)
	addr := smtpServer + ":" + smtpPort
	return smtp.SendMail(addr, auth, sender, []string{to}, msg)
}
