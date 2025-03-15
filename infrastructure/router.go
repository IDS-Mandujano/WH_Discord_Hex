package infrastructure

import (
	"github_wb/infrastructure/handlers"

	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {
	routes := engine.Group("/")

	{
		routes.POST("development", handlers.HandleDevelopmentEvent)
		routes.POST("workflows", handlers.HandleWorkflowsEvent)
	}
}