package application

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type DiscordMessage struct {
	Content string `json:"content"`
}

func SendDiscordMessage(message string) {
	godotenv.Load()
	webhookURL := os.Getenv("DISCORD_wH_URL")

	if webhookURL == "" {
		log.Println("⚠️ No se encontró la URL del webhook de Discord en .env")
		return
	}

	payload := DiscordMessage{Content: message}
	jsonValue, _ := json.Marshal(payload)

	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println("❌ Error al enviar mensaje a Discord:", err)
	}
}

//pureba de wenhook