package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

// ServeSubtitles retorna un handler para servir subtítulos
func ServeSubtitles(videoUseCase usecases.VideoUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		if filename == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "filename parameter is required",
			})
			return
		}

		subtitlePath, err := videoUseCase.GetSubtitles(filename)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Subtitle file not found",
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
