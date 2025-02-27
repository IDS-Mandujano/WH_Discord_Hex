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

func SendDiscordMessage(message string) {
	webhookURL := os.Getenv("DISCORD_wH_URL")

	if webhookURL == "" {
		log.Println("No se encontró la URL del webhook de Discord en .env")
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

	log.Println("Acción recibida del Push:", action)

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

		log.Println(message)

		return sendDiscordMessage(webhookURL, message)
	} else {
		log.Printf("Pull Request Action no es Closed ni ready_for_review: %s", action)
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

//asdasdasd