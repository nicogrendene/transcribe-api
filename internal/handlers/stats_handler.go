package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

// GetStats retorna un handler para estadísticas
func GetStats(statsUseCase usecases.StatsUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := statsUseCase.GetStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Error obteniendo estadísticas",
				Details: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}
