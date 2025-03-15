package handlers

import (
	"github_wb/application"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleDevelopmentEvent(ctx *gin.Context) {
	payload, err := ctx.GetRawData()
	if err != nil {
		log.Println("Error al leer el cuerpo de la solicitud:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el cuerpo"})
		return
	}

	statusCode := application.ProcessPullRequest(payload)
	ctx.JSON(statusCode, gin.H{"status": "Evento Development procesado"})
}

func HandleWorkflowsEvent(ctx *gin.Context) {
	payload, err := ctx.GetRawData()
	if err != nil {
		log.Println("Error al leer el cuerpo de la solicitud:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el cuerpo"})
		return
	}

	statusCode := application.ProcessWorkflows(payload)
	ctx.JSON(statusCode, gin.H{"status": "Evento Workflow procesado"})
}