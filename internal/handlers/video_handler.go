package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

// ServeVideo retorna un handler para servir videos
func ServeVideo(videoUseCase usecases.VideoUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		if filename == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "filename parameter is required",
			})
			return
		}

		videoPath, err := videoUseCase.GetVideo(filename)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Video file not found",
				Details: err.Error(),
			})
			return
		}

		c.Header("Content-Type", "video/mp4")
		c.Header("Accept-Ranges", "bytes")
		c.Header("Cache-Control", "public, max-age=3600")
		c.File(videoPath)
	}
}
