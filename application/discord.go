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

// Enviar mensaje a Discord
func SendDiscordMessage(message string) {
	webhookURL := os.Getenv("DISCORD_wH_URL")

	if webhookURL == "" {
		log.Println("⚠️ No se encontró la URL del webhook de Discord en .env")
		return
	}

	payload := DiscordMessage{Content: message}
	jsonValue, _ := json.Marshal(payload)

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println("Error al enviar mensaje a Discord:", err)
		return
	}
	defer resp.Body.Close()
	log.Println("Mensaje enviado con estado:", resp.StatusCode)
}

// Manejar evento de Push
func ProcessPush(payload []byte) int {
	webhookURL := os.Getenv("DISCORD_wH_URL")
	if webhookURL == "" {
		log.Println("Error: Webhook de Discord no configurado en .env")
		return 500
	}

	// Parsear JSON del payload
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		log.Println("Error al parsear el payload de push:", err)
		return 500
	}

	// Verificar si hay commits
	commitsData, exists := data["commits"].([]interface{})
	if !exists || len(commitsData) == 0 {
		log.Println("Error: No se encontraron commits en el payload")
		return 400
	}

	// Construir mensaje para Discord
	message := "📌 **Nuevo Push Detectado:**\n"
	for _, c := range commitsData {
		commit, ok := c.(map[string]interface{})
		if !ok {
			continue
		}

		// Obtener autor y mensaje del commit
		authorName := "Desconocido"
		if authorData, ok := commit["author"].(map[string]interface{}); ok {
			if name, ok := authorData["name"].(string); ok {
				authorName = name
			}
		}

		commitMessage := "Sin mensaje"
		if msg, ok := commit["message"].(string); ok {
			commitMessage = msg
		}

		message += "- ✏️ `" + authorName + "`: " + commitMessage + "\n"
	}

	
	return sendDiscordMessage(webhookURL, message)
}

func sendDiscordMessage(url, message string) int {
	body, _ := json.Marshal(DiscordMessage{Content: message})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Println("Error al enviar mensaje a Discord:", err)
		return 500
	}
	defer resp.Body.Close()

	log.Println("Mensaje enviado a Discord con estado:", resp.StatusCode)
	return resp.StatusCode
}