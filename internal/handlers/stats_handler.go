package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/log"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

func GetStats(statsUseCase usecases.StatsUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := log.With(c.Request.Context(), log.UseCase("stats"))

		stats, err := statsUseCase.GetStats(ctx)
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
