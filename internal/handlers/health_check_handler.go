package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

// HealthCheck retorna un handler para health check
func HealthCheck(healthUseCase usecases.HealthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		response, err := healthUseCase.CheckHealth()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Error verificando salud del sistema",
				Details: err.Error(),
			})
			return
		}

		statusCode := http.StatusOK
		if response.Status == "error" {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, response)
	}
}
