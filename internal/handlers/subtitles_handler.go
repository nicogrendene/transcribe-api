package handlers

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

// ServeSubtitles retorna un handler para servir subtítulos
func ServeSubtitles(videoUseCase usecases.VideoUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "id parameter is required",
			})
			return
		}

		subtitlePath, err := videoUseCase.GetSubtitles(id)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Subtitle file not found",
				Details: err.Error(),
			})
			return
		}

		// Leer el contenido del archivo para evitar problemas de cache
		file, err := os.Open(subtitlePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Error reading subtitle file",
				Details: err.Error(),
			})
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Error reading subtitle content",
				Details: err.Error(),
			})
			return
		}

		// Headers para evitar cache
		c.Header("Content-Type", "text/vtt; charset=utf-8")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.String(http.StatusOK, string(content))
	}
}
