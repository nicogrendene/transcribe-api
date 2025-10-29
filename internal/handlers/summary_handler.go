package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/log"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

func ServeSummary(videoUseCase usecases.VideoUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := log.With(c.Request.Context(), log.UseCase("get_summary"))

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "id parameter is required",
			})
			return
		}

		summary, err := videoUseCase.GetSummary(ctx, id)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Summary file not found",
				Details: err.Error(),
			})
			return
		}

		// Headers para evitar cache
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.String(http.StatusOK, summary)
	}
}
