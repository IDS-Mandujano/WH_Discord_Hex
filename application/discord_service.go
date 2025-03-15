package application

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type DiscordMessage struct {
	Content string `json:"content"`
}

func sendDiscordMessage(webhookURL, message string) int {
	if webhookURL == "" {
		log.Println("Error: Webhook de Discord no configurado en .env")
		return 500
	}

	body, _ := json.Marshal(DiscordMessage{Content: message})
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Println("Error al enviar mensaje a Discord:", err)
		return 500
	}
	defer resp.Body.Close()

	log.Println("Mensaje enviado a Discord con estado:", resp.StatusCode)
	return resp.StatusCode
}

func SendDevelopmentMessage(message string) int {
	webhookURL := os.Getenv("DISCORD_DEVELOPMENT_URL")
	return sendDiscordMessage(webhookURL, message)
}

func SendWorkflowsMessage(message string) int {
	webhookURL := os.Getenv("DISCORD_WORKFLOWS_URL")
	return sendDiscordMessage(webhookURL, message)
}