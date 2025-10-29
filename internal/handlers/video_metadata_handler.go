package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/log"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

// VideoMetadata representa los metadatos de un video
type VideoMetadata struct {
	ID           string `json:"id"`
	Size         int64  `json:"size"`
	Duration     string `json:"duration,omitempty"`
	Width        int    `json:"width,omitempty"`
	Height       int    `json:"height,omitempty"`
	Bitrate      int    `json:"bitrate,omitempty"`
	Format       string `json:"format"`
	LastModified string `json:"last_modified"`
	HasThumbnail bool   `json:"has_thumbnail"`
	HasSubtitles bool   `json:"has_subtitles"`
	HasSummary   bool   `json:"has_summary"`
}

// GetVideoMetadata retorna metadatos del video sin cargar el archivo completo
func GetVideoMetadata(videoUseCase usecases.VideoUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := log.With(c.Request.Context(), log.UseCase("get_video_metadata"))

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "id parameter is required",
			})
			return
		}

		// Obtener la ruta del video
		videoPath, err := videoUseCase.GetVideo(ctx, id)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Video not found",
				Details: err.Error(),
			})
			return
		}

		// Obtener información del archivo
		fileInfo, err := os.Stat(videoPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Error getting file info",
				Details: err.Error(),
			})
			return
		}

		// Obtener el directorio del video
		videoDir := filepath.Dir(videoPath)

		// Verificar archivos relacionados
		hasThumbnail := false
		hasSubtitles := false
		hasSummary := false

		if _, err := os.Stat(filepath.Join(videoDir, "thumbnail.jpg")); err == nil {
			hasThumbnail = true
		}

		if _, err := os.Stat(filepath.Join(videoDir, "subtitles.vtt")); err == nil {
			hasSubtitles = true
		}

		if _, err := os.Stat(filepath.Join(videoDir, "summary.txt")); err == nil {
			hasSummary = true
		}

		metadata := VideoMetadata{
			ID:           id,
			Size:         fileInfo.Size(),
			Format:       "mp4",
			LastModified: fileInfo.ModTime().Format("2006-01-02T15:04:05Z07:00"),
			HasThumbnail: hasThumbnail,
			HasSubtitles: hasSubtitles,
			HasSummary:   hasSummary,
		}

		// Headers para cache agresivo de metadatos
		c.Header("Cache-Control", "public, max-age=3600") // 1 hora
		c.Header("Content-Type", "application/json")

		c.JSON(http.StatusOK, metadata)
	}
}
