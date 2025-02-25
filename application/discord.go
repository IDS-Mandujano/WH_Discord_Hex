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

// Manejar evento de Push (nombre cambiado)
func ProcessPush(payload []byte) int {
	webhookURL := os.Getenv("DISCORD_wH_URL")
	if webhookURL == "" {
		log.Println("Error: Webhook de Discord no configurado en .env")
		return 500
	}

	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		log.Println("Error al parsear el payload de Push:", err)
		return 500
	}

	action, ok := data["action"].(string)
	if !ok {
		log.Println("Error: No se encontró la acción en el payload")
		return 400
	}

	// Imprimir la acción recibida para debug
	log.Println("Acción recibida del Push:", action)

	// Verificar si el PR está listo para revisión o fue mergeado
	if action == "ready_for_review" || action == "closed" {
		pr, exists := data["pull_request"].(map[string]interface{})
		if !exists {
			log.Println("Error: No se encontró la información del Pull Request en el payload")
			return 400
		}

		user := pr["user"].(map[string]interface{})["login"].(string)
		title := pr["title"].(string)
		baseBranch := pr["base"].(map[string]interface{})["ref"].(string)
		headBranch := pr["head"].(map[string]interface{})["ref"].(string)
		prURL := pr["html_url"].(string)

		message := "Nuevo Pull Request:\n" +
			"- Usuario: " + user + "\n" +
			"- Título: " + title + "\n" +
			"- De: " + headBranch + " -> " + baseBranch + "\n" +
			"- URL: " + prURL

		// Imprimir en consola
		log.Println(message)

		// Enviar a Discord
		return sendDiscordMessage(webhookURL, message)
	} else {
		// Si la acción no es ni "ready_for_review" ni "closed", mostramos un mensaje
		log.Println("Push Action no es relevante:", action)
	}

	return 200
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