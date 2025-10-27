package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

func ServeSummary(videoUseCase usecases.VideoUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "id parameter is required",
			})
			return
		}

		subtitlePath, err := videoUseCase.GetSummary(id)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Summary file not found",
				Details: err.Error(),
			})
			return
		}

		c.Header("Content-Type", "text/vtt; charset=utf-8")
		c.Header("Cache-Control", "public, max-age=3600")
		c.Header("Access-Control-Allow-Origin", "*")
		c.File(subtitlePath)
	}
}
