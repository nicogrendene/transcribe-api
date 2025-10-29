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

// VideoInfo representa la información de un video
type VideoInfo struct {
	ID        string    `json:"id"`
	Qualities []Quality `json:"qualities"`
}

// Quality representa una calidad de video disponible
type Quality struct {
	Name      string `json:"name"`
	Label     string `json:"label"`
	File      string `json:"file"`
	Available bool   `json:"available"`
}

// GetVideoInfo retorna información sobre las calidades disponibles de un video
func GetVideoInfo(videoUseCase usecases.VideoUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := log.With(c.Request.Context(), log.UseCase("get_video_info"))

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "id parameter is required",
			})
			return
		}

		// Obtener la ruta base del video
		videoPath, err := videoUseCase.GetVideo(ctx, id)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Video not found",
				Details: err.Error(),
			})
			return
		}

		// Obtener el directorio del video
		videoDir := filepath.Dir(videoPath)

		// Definir las calidades disponibles
		qualities := []Quality{
			{Name: "360p", Label: "Mobile (360p)", File: "video_360p.mp4"},
			{Name: "480p", Label: "Low (480p)", File: "video_480p.mp4"},
			{Name: "720p", Label: "Medium (720p)", File: "video_720p.mp4"},
			{Name: "1080p", Label: "High (1080p)", File: "video_1080p.mp4"},
			{Name: "default", Label: "Default", File: "video.mp4"},
		}

		// Verificar qué calidades están disponibles
		for i := range qualities {
			filePath := filepath.Join(videoDir, qualities[i].File)
			if _, err := os.Stat(filePath); err == nil {
				qualities[i].Available = true
			}
		}

		videoInfo := VideoInfo{
			ID:        id,
			Qualities: qualities,
		}

		c.JSON(http.StatusOK, videoInfo)
	}
}
