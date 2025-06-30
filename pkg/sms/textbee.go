package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	appErr "food-delivery-app-server/pkg/errors"
)

type textBeeRequest struct {
	Recipients []string `json:"recipients"`
	Message    string   `json:"message"`
}

func SendOTPTextBee(phone, otp string) error {
	apiKey := os.Getenv("TEXTBEE_API_KEY")
	deviceID := os.Getenv("TEXTBEE_DEVICE_ID")

	if apiKey == "" {
		return appErr.NewInternal("TextBee API key not set", nil)
	}
	if deviceID == "" {
		return appErr.NewInternal("TextBee Device ID not set", nil)
	}

	textBeeApi := fmt.Sprintf("https://api.textbee.dev/api/v1/gateway/devices/%s/send-sms", deviceID)

	reqBody := textBeeRequest{
		Recipients: []string{phone},
		Message:    fmt.Sprintf("Good day! Your Food Delivery App OTP code for signing up is: %s", otp),
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return appErr.NewInternal("Failed to marshal SMS request body", err)
	}

	req, err := http.NewRequest("POST", textBeeApi, bytes.NewBuffer(body))
	if err != nil {
		return appErr.NewInternal("Failed to create SMS request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return appErr.NewInternal("Failed to send SMS request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return appErr.NewInternal(fmt.Sprintf("Failed to send SMS: status %d", resp.StatusCode), nil)
	}

	return nil
}
