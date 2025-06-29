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

func SendResetCode(to string, code string) error {

	smtpServer := fallback(os.Getenv("SMTP_SERVER"), "smtp.gmail.com")
	smtpPort := fallback(os.Getenv("SMTP_PORT"), "587")
	sender := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	emailBody := fmt.Sprintf(`
    	<div style="font-family: Arial, sans-serif; line-height:1.6; color: #545454;">
        	<h1 style="font-size: 24px; font-weight: bold; color:#d417ff">
            	Food Delivery App Password Reset
        	</h1>
        	<p style="font-size:16px; margin-bottom:10px">
            	Your verification code is
            	<strong style="font-size: 18px; color:#d417ff;">%s</strong>
        	</p>
        	<p style="font-size: 14px; margin-bottom: 20px;">
            	The code expires in <strong>5 minutes</strong> after this email was sent.
        	</p>
        	<p style="font-size: 14px;">
            	Enter the code in the reset password section of the app to reset your password.
        	</p>
        	<hr style="border: 0; border-top: 1px solid #ccc; margin: 20px 0;">
        	<p style="font-size: 12px; color: #999;">
            	If you did not request a password reset, please ignore this email.
        	</p>
    	</div>`, code)

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
