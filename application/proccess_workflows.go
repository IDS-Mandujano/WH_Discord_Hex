package application

import (
	"encoding/json"
	"fmt"
	"log"
	"github_wb/domain/value_objects"
)

func ProcessWorkflows(payload []byte) int {
	var eventPayload domain.WorkflowEventPayload

	if err := json.Unmarshal(payload, &eventPayload); err != nil {
		log.Println("Error al parsear JSON:", err)
		return 500
	}

	message := fmt.Sprintf(
		"**Workflow Event**\n- Workflow: %s\n- Estado: %s\n- URL: %s",
		eventPayload.WorkflowJob.WorkflowName,
		eventPayload.WorkflowJob.Status,
		eventPayload.WorkflowJob.RunURL,
	)

	return SendWorkflowsMessage(message)
}