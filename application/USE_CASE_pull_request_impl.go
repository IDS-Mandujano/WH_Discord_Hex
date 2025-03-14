package application

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	domain "github_wb/domain/value_objects"
)

func ProcessPullRequest(payload []byte) int {
	var eventPayload domain.PullRequestEventPayload

	if err := json.Unmarshal(payload, &eventPayload); err != nil {
		log.Println("Error al parsear JSON:", err)
		return 500
	}

	if eventPayload.Action == "opened" && eventPayload.PullRequest.Draft {
		log.Println("El Pull Request está en Draft, no se envía mensaje.")
		return 200
	}

	log.Printf("Pull Request Recibido: Action=%s", eventPayload.Action)

	if eventPayload.Action == "closed" || eventPayload.Action == "ready_for_review" || eventPayload.Action == "opened" {
		base := eventPayload.PullRequest.Base.Ref
		branch := eventPayload.PullRequest.Head.Ref
		
		user := eventPayload.PullRequest.User.Login
		pRID := eventPayload.PullRequest.ID
		url := eventPayload.PullRequest.HTMLURL

		log.Printf("Pull Request Recibido:\nID:%d\nBase:%s\nHead:%s\nUser:%s\nURL:%s", pRID, base, branch, user, url)

		message := fmt.Sprintf(
			"**Nuevo Pull Request**\n- Usuario: %s\n- Base: %s\n- Head: %s\n- Estado: %s\n- URL: %s",
			user, base, branch, eventPayload.Action, url,
		)

		return sendDiscordMessage(os.Getenv("DISCORD_wH_URL"), message)
	}

	return 200
}