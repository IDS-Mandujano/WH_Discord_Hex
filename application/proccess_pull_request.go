package application

import (
	"encoding/json"
	"fmt"
	"log"

	"github_wb/domain/value_objects"
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

	message := fmt.Sprintf(
		"**Nuevo Pull Request**\n- Usuario: %s\n- Base: %s\n- Head: %s\n- Estado: %s\n- URL: %s",
		eventPayload.PullRequest.User.Login,
		eventPayload.PullRequest.Base.Ref,
		eventPayload.PullRequest.Head.Ref,
		eventPayload.Action,
		eventPayload.PullRequest.HTMLURL,
	)

	return SendDevelopmentMessage(message)
}