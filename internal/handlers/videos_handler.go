package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

func GetVideos(videoUseCase usecases.VideoUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		jsonData, err := videoUseCase.GetVideos()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Error obteniendo videos",
				Details: err.Error(),
			})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Header("Cache-Control", "public, max-age=300")
		c.Data(http.StatusOK, "application/json", jsonData)
	}
}
