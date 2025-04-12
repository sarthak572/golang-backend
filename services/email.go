package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type EmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
}

func SendEmail(to string) error {
	url := "https://api.resend.com/emails"
	apiKey := os.Getenv("RESEND_API_KEY")
	fromEmail := os.Getenv("RESEND_FROM_EMAIL")

	if apiKey == "" || fromEmail == "" {
		log.Println("Missing RESEND_API_KEY or RESEND_FROM_EMAIL in .env")
		return fmt.Errorf("missing email configuration")
	}

	payload := EmailRequest{
		From:    fromEmail, // Now using your verified domain!
		To:      []string{to},
		Subject: "Welcome to General Shop!",
		Text:    "Thank you for registering with us!",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var errBody map[string]interface{}
		json.NewDecoder(res.Body).Decode(&errBody)
		log.Printf("Error sending email: status: %d, error: %v\n", res.StatusCode, errBody)
		return fmt.Errorf("email failed with status %d", res.StatusCode)
	}

	log.Println("✅ Email sent successfully to", to)
	return nil
}
