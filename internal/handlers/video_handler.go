package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/log"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

func ServeVideo(videoUseCase usecases.VideoUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := log.With(c.Request.Context(), log.UseCase("serve_video"))

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "id parameter is required",
			})
			return
		}

		videoPath, err := videoUseCase.GetVideo(ctx, id)
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
